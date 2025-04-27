package parser

import (
	"fmt"
	"go/ast"
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
	Name      string // メソッド名
	Signature string // メソッドのシグネチャ
	Comment   string // メソッドのコメント
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
		Comment: getComment(pkg, obj),
		Fields:  make([]Field, 0, structType.NumFields()),
		Methods: make([]Method, 0),
	}

	// フィールド情報を取得
	for i := 0; i < structType.NumFields(); i++ {
		field := structType.Field(i)
		info.Fields = append(info.Fields, Field{
			Name:       field.Name(),
			Type:       field.Type().String(),
			Comment:    getComment(pkg, field),
			IsExported: field.Exported(),
		})
	}

	// メソッド情報を取得
	for i := 0; i < named.NumMethods(); i++ {
		method := named.Method(i)
		info.Methods = append(info.Methods, Method{
			Name:      method.Name(),
			Signature: method.Type().String(),
			Comment:   getComment(pkg, method),
		})
	}

	return info, nil
}

// getComment は、指定されたオブジェクトのコメントを取得します。
func getComment(pkg *packages.Package, obj types.Object) string {
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
			}
			return true
		})
	}

	// コメントの前後の空白を削除
	return strings.TrimSpace(comment)
}
