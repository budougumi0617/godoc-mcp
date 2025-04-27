package parser

import (
	"fmt"
	"go/ast"
	"go/format"
	"go/token"
	"go/types"
	"strings"

	"golang.org/x/tools/go/packages"
)

// Parser は、ロードしたパッケージ情報を保持する構造体です。
type Parser struct {
	pkgs map[string]*packages.Package
}

// New は、指定されたディレクトリからGoパッケージをロードし、Parserインスタンスを作成します。
// rootDirは基点となるディレクトリ、pkgDirはロード対象のパターンを指定します。
// pkgDirが空の場合は、"./..."がパターンとして使用されます。
func New(rootDir string, pkgDir string) (*Parser, error) {
	patterns := []string{"./..."}
	if pkgDir != "" {
		patterns = []string{pkgDir}
	}

	cfg := &packages.Config{
		Mode: packages.NeedName |
			packages.NeedFiles |
			packages.NeedCompiledGoFiles |
			packages.NeedImports |
			packages.NeedDeps |
			packages.NeedTypes |
			packages.NeedSyntax |
			packages.NeedTypesInfo,
		Dir:   rootDir,
		Tests: false,
	}

	pkgs, err := packages.Load(cfg, patterns...)
	if err != nil {
		return nil, fmt.Errorf("failed to load packages: %w", err)
	}

	parser := &Parser{
		pkgs: make(map[string]*packages.Package),
	}

	// パッケージをマップに格納
	for _, pkg := range pkgs {
		parser.pkgs[pkg.PkgPath] = pkg
	}

	return parser, nil
}

// GetAllPackages は、ロードされたすべてのパッケージを返します。
func (p *Parser) GetAllPackages() []*packages.Package {
	result := make([]*packages.Package, 0, len(p.pkgs))
	for _, pkg := range p.pkgs {
		result = append(result, pkg)
	}
	return result
}

// GetPackage は、指定されたパッケージパスのパッケージを返します。
// パッケージが見つからない場合はエラーを返します。
func (p *Parser) GetPackage(pkgPath string) (*packages.Package, error) {
	pkg, ok := p.pkgs[pkgPath]
	if !ok {
		return nil, fmt.Errorf("package not found: %s", pkgPath)
	}
	return pkg, nil
}

// StructInfo は、構造体の情報を表す構造体です。
type StructInfo struct {
	Name    string   // 構造体名
	Comment string   // 構造体のコメント
	Fields  []Field  // フィールドのリスト
	Methods []Method // メソッドのリスト
}

// Field は、構造体のフィールド情報を表す構造体です。
type Field struct {
	Name       string // フィールド名
	Type       string // フィールドの型
	Comment    string // フィールドのコメント
	IsExported bool   // エクスポートされているかどうか
}

// Method は、構造体のメソッド情報を表す構造体です。
type Method struct {
	Name      string    // メソッド名
	Signature string    // メソッドのシグネチャ
	Comment   string    // メソッドのコメント
	Examples  []Example // メソッドの例
}

// GetStructInfo は、指定されたパッケージ内の構造体情報を返します。
func (p *Parser) GetStructInfo(pkgPath, structName string) (*StructInfo, error) {
	pkg, err := p.GetPackage(pkgPath)
	if err != nil {
		return nil, err
	}

	// パッケージ内の型情報を取得
	scope := pkg.Types.Scope()
	obj := scope.Lookup(structName)
	if obj == nil {
		return nil, fmt.Errorf("struct not found: %s in package %s", structName, pkgPath)
	}

	// 型が構造体かどうかを確認
	named, ok := obj.Type().(*types.Named)
	if !ok {
		return nil, fmt.Errorf("not a struct: %s in package %s", structName, pkgPath)
	}

	structType, ok := named.Underlying().(*types.Struct)
	if !ok {
		return nil, fmt.Errorf("not a struct: %s in package %s", structName, pkgPath)
	}

	// 構造体情報を構築
	info := &StructInfo{
		Name:    structName,
		Comment: GetComment(pkg, obj),
		Fields:  make([]Field, 0, structType.NumFields()),
		Methods: make([]Method, 0),
	}

	// フィールド情報を取得
	for i := 0; i < structType.NumFields(); i++ {
		field := structType.Field(i)
		info.Fields = append(info.Fields, Field{
			Name:       field.Name(),
			Type:       field.Type().String(),
			Comment:    GetComment(pkg, field),
			IsExported: field.Exported(),
		})
	}

	// メソッド情報を取得
	for i := 0; i < named.NumMethods(); i++ {
		method := named.Method(i)
		info.Methods = append(info.Methods, Method{
			Name:      method.Name(),
			Signature: method.Type().String(),
			Comment:   GetComment(pkg, method),
		})
	}

	return info, nil
}

// FuncInfo は、関数の情報を表す構造体です。
type FuncInfo struct {
	Name      string    // 関数名
	Signature string    // 関数のシグネチャ
	Comment   string    // 関数のコメント
	Examples  []Example // 関数の例
}

// Example は、関数の例を表す構造体です。
type Example struct {
	Name   string // 例の名前
	Code   string // 例のコード
	Output string // 例の出力
}

// GetFuncInfo は、指定されたパッケージ内の関数情報を返します。
func (p *Parser) GetFuncInfo(pkgPath, funcName string) (*FuncInfo, error) {
	pkg, err := p.GetPackage(pkgPath)
	if err != nil {
		return nil, err
	}

	// パッケージ内の型情報を取得
	scope := pkg.Types.Scope()
	obj := scope.Lookup(funcName)
	if obj == nil {
		return nil, fmt.Errorf("function not found: %s in package %s", funcName, pkgPath)
	}

	// 型が関数かどうかを確認
	fn, ok := obj.(*types.Func)
	if !ok {
		return nil, fmt.Errorf("not a function: %s in package %s", funcName, pkgPath)
	}

	// 関数情報を構築
	info := &FuncInfo{
		Name:      funcName,
		Signature: fn.Type().String(),
		Comment:   GetComment(pkg, obj),
		Examples:  make([]Example, 0),
	}

	// 例を取得
	info.Examples = GetExamples(pkg, funcName)

	return info, nil
}

// GetMethodInfo は、指定されたパッケージ内のメソッド情報を返します。
func (p *Parser) GetMethodInfo(pkgPath, structName, methodName string) (*Method, error) {
	pkg, err := p.GetPackage(pkgPath)
	if err != nil {
		return nil, err
	}

	// パッケージ内の型情報を取得
	scope := pkg.Types.Scope()
	obj := scope.Lookup(structName)
	if obj == nil {
		return nil, fmt.Errorf("struct not found: %s in package %s", structName, pkgPath)
	}

	// 型が構造体かどうかを確認
	named, ok := obj.Type().(*types.Named)
	if !ok {
		return nil, fmt.Errorf("not a struct: %s in package %s", structName, pkgPath)
	}

	// メソッドを探す
	var method *types.Func
	for i := 0; i < named.NumMethods(); i++ {
		m := named.Method(i)
		if m.Name() == methodName {
			method = m
			break
		}
	}

	if method == nil {
		return nil, fmt.Errorf("method not found: %s.%s in package %s", structName, methodName, pkgPath)
	}

	// メソッド情報を構築
	info := &Method{
		Name:      methodName,
		Signature: method.Type().String(),
		Comment:   GetComment(pkg, method),
		Examples:  make([]Example, 0),
	}

	// 例を取得
	info.Examples = GetExamples(pkg, methodName)

	return info, nil
}

// ConstInfo は、定数の情報を表す構造体です。
type ConstInfo struct {
	Name    string // 定数名
	Type    string // 定数の型
	Value   string // 定数の値
	Comment string // 定数のコメント
}

// VarInfo は、変数の情報を表す構造体です。
type VarInfo struct {
	Name    string // 変数名
	Type    string // 変数の型
	Comment string // 変数のコメント
}

// GetConstAndVarInfo は、指定されたパッケージ内の定数・変数情報を返します。
func (p *Parser) GetConstAndVarInfo(pkgPath string) ([]ConstInfo, []VarInfo, error) {
	pkg, err := p.GetPackage(pkgPath)
	if err != nil {
		return nil, nil, err
	}

	// パッケージ内の型情報を取得
	scope := pkg.Types.Scope()

	// 定数情報を取得
	constants := make([]ConstInfo, 0)
	for _, name := range scope.Names() {
		obj := scope.Lookup(name)
		if obj == nil {
			continue
		}

		// 定数かどうかを確認
		if constObj, ok := obj.(*types.Const); ok {
			constants = append(constants, ConstInfo{
				Name:    name,
				Type:    constObj.Type().String(),
				Value:   constObj.Val().String(),
				Comment: GetComment(pkg, obj),
			})
		}
	}

	// 変数情報を取得
	variables := make([]VarInfo, 0)
	for _, name := range scope.Names() {
		obj := scope.Lookup(name)
		if obj == nil {
			continue
		}

		// 変数かどうかを確認
		if varObj, ok := obj.(*types.Var); ok {
			variables = append(variables, VarInfo{
				Name:    name,
				Type:    varObj.Type().String(),
				Comment: GetComment(pkg, obj),
			})
		}
	}

	return constants, variables, nil
}

// GetExamples は、指定された関数の例を取得します。
func GetExamples(pkg *packages.Package, funcName string) []Example {
	var examples []Example

	// パッケージ内のASTノードを検索
	for _, file := range pkg.Syntax {
		ast.Inspect(file, func(n ast.Node) bool {
			// 関数宣言を探す
			funcDecl, ok := n.(*ast.FuncDecl)
			if !ok {
				return true
			}

			// 例の関数かどうかを確認
			if !strings.HasPrefix(funcDecl.Name.Name, "Example") {
				return true
			}

			// 対象の関数の例かどうかを確認
			if !strings.HasSuffix(funcDecl.Name.Name, funcName) {
				return true
			}

			// 例の情報を取得
			example := Example{
				Name: funcDecl.Name.Name,
			}

			// 例のコードを取得
			if funcDecl.Body != nil {
				example.Code = GetNodeString(pkg.Fset, funcDecl.Body)
			}

			// 例の出力を取得
			if funcDecl.Doc != nil {
				for _, comment := range funcDecl.Doc.List {
					if strings.HasPrefix(comment.Text, "// Output:") {
						example.Output = strings.TrimPrefix(comment.Text, "// Output:")
						example.Output = strings.TrimSpace(example.Output)
					}
				}
			}

			examples = append(examples, example)
			return true
		})
	}

	return examples
}

// getNodeString は、指定されたノードの文字列表現を取得します。
func getNodeString(fset *token.FileSet, node ast.Node) string {
	var buf strings.Builder
	if err := format.Node(&buf, fset, node); err != nil {
		return ""
	}
	return buf.String()
}

// getComment は、指定されたオブジェクトのコメントを取得します。
func GetComment(pkg *packages.Package, obj types.Object) string {
	var comment string

	// パッケージ内のASTノードを検索
	for _, file := range pkg.Syntax {
		ast.Inspect(file, func(n ast.Node) bool {
			switch node := n.(type) {
			case *ast.TypeSpec:
				// 型宣言の場合
				if node.Name.Name == obj.Name() {
					if node.Doc != nil {
						comment = node.Doc.Text()
					}
					return false
				}
			case *ast.Field:
				// フィールドの場合
				for _, name := range node.Names {
					if name.Name == obj.Name() {
						if node.Doc != nil {
							comment = node.Doc.Text()
						} else if node.Comment != nil {
							comment = node.Comment.Text()
						}
						return false
					}
				}
			case *ast.FuncDecl:
				// メソッドの場合
				if node.Name.Name == obj.Name() {
					if node.Doc != nil {
						comment = node.Doc.Text()
					}
					return false
				}
			case *ast.ValueSpec:
				// 定数・変数の場合
				for _, name := range node.Names {
					if name.Name == obj.Name() {
						if node.Doc != nil {
							comment = node.Doc.Text()
						} else if node.Comment != nil {
							comment = node.Comment.Text()
						}
						return false
					}
				}
			}
			return true
		})
	}

	// コメントの前後の空白を削除
	return strings.TrimSpace(comment)
}

// GetNodeString は、指定されたノードの文字列表現を取得します。
func GetNodeString(fset *token.FileSet, node ast.Node) string {
	var buf strings.Builder
	if err := format.Node(&buf, fset, node); err != nil {
		return ""
	}
	return buf.String()
}
