package handler

import (
	"context"
	"fmt"
	"go/types"
	"strings"

	godoc "github.com/budougumi0617/godoc-mcp"
	"github.com/budougumi0617/godoc-mcp/internal/model"
	"github.com/budougumi0617/godoc-mcp/internal/parser"
	mcp "github.com/ktr0731/go-mcp"
)

// ToolHandler は、MCPツールリクエストを処理するハンドラー構造体です。
type ToolHandler struct {
	parser *parser.Parser
}

// NewToolHandler は、新しいToolHandlerインスタンスを作成します。
func NewToolHandler(p *parser.Parser) *ToolHandler {
	return &ToolHandler{
		parser: p,
	}
}

// HandleToolListPackages は、ロードされたすべてのパッケージのリストを返します。
func (h *ToolHandler) HandleToolListPackages(ctx context.Context, req *godoc.ToolListPackagesRequest) (*mcp.CallToolResult, error) {
	pkgs := h.parser.GetAllPackages()
	if len(pkgs) == 0 {
		return &mcp.CallToolResult{
			Content: []mcp.CallToolContent{
				mcp.TextContent{Text: "パッケージがロードされていません。"},
			},
		}, nil
	}

	var packages []model.PackageInfo
	for _, p := range pkgs {
		// パッケージのコメントを取得
		packages = append(packages, model.PackageInfo{
			Name:       p.Name,
			ImportPath: p.PkgPath,
			Comment:    parser.GetPackageComment(p),
		})
	}

	// マークダウン形式でフォーマット
	mdContent := model.FormatPackageListMarkdown(packages)

	return &mcp.CallToolResult{
		Content: []mcp.CallToolContent{
			mcp.TextContent{Text: mdContent},
		},
	}, nil
}

// HandleToolInspectPackage は、指定されたパッケージの公開されている構造体、メソッド、関数をリストします。
func (h *ToolHandler) HandleToolInspectPackage(ctx context.Context, req *godoc.ToolInspectPackageRequest) (*mcp.CallToolResult, error) {
	pkg, err := h.parser.GetPackage(req.PackageName)
	if err != nil {
		return nil, fmt.Errorf("パッケージの取得に失敗しました: %w", err)
	}

	// パッケージ情報を作成
	pkgInfo := model.PackageInfo{
		Name:       pkg.Name,
		ImportPath: pkg.PkgPath,
		Comment:    parser.GetPackageComment(pkg),
	}

	// 構造体、関数、メソッドの情報を収集
	scope := pkg.Types.Scope()
	var structs []model.StructSummary
	var funcs []model.FuncSummary
	var methods []model.MethodSummary

	for _, name := range scope.Names() {
		obj := scope.Lookup(name)
		if !obj.Exported() {
			continue
		}

		switch obj := obj.(type) {
		case *types.TypeName:
			// 型定義を取得
			if _, ok := obj.Type().Underlying().(*types.Struct); ok {
				structs = append(structs, model.StructSummary{
					Name:    obj.Name(),
					Comment: parser.GetComment(pkg, obj), // parserパッケージの公開された関数を使用
				})
			}
		case *types.Func:
			sig, ok := obj.Type().(*types.Signature)
			if !ok {
				continue
			}

			// メソッドの場合（レシーバーがある）
			if sig.Recv() != nil {
				recvType := sig.Recv().Type().String()
				// ポインタ型の場合は*を削除
				if strings.HasPrefix(recvType, "*") {
					recvType = recvType[1:]
				}
				methods = append(methods, model.MethodSummary{
					ReceiverType: recvType,
					Name:         obj.Name(),
					Comment:      parser.GetComment(pkg, obj), // parserパッケージの公開された関数を使用
				})
			} else {
				// 関数の場合
				funcs = append(funcs, model.FuncSummary{
					Name:    obj.Name(),
					Comment: parser.GetComment(pkg, obj), // parserパッケージの公開された関数を使用
				})
			}
		}
	}

	// マークダウン形式でフォーマット
	mdContent := model.FormatPackageInspectionMarkdown(pkgInfo, structs, funcs, methods, req.IncludeComments)

	return &mcp.CallToolResult{
		Content: []mcp.CallToolContent{
			mcp.TextContent{Text: mdContent},
		},
	}, nil
}

// HandleToolGetDocStruct は、指定された構造体に関する情報を返します。
func (h *ToolHandler) HandleToolGetDocStruct(ctx context.Context, req *godoc.ToolGetDocStructRequest) (*mcp.CallToolResult, error) {
	structInfo, err := h.parser.GetStructInfo(req.PackageName, req.StructName)
	if err != nil {
		return nil, fmt.Errorf("構造体情報の取得に失敗しました: %w", err)
	}

	// フィールドとメソッドの情報を変換
	var fields []model.FieldDoc
	var methods []model.MethodDoc

	for _, f := range structInfo.Fields {
		fields = append(fields, model.FieldDoc{
			Name:       f.Name,
			Type:       f.Type,
			Comment:    f.Comment,
			IsExported: f.IsExported,
		})
	}

	for _, m := range structInfo.Methods {
		methods = append(methods, model.MethodDoc{
			Name:      m.Name,
			Signature: m.Signature,
			Comment:   m.Comment,
		})
	}

	// マークダウン形式でフォーマット
	mdContent := model.FormatStructDocMarkdown(structInfo.Name, structInfo.Comment, fields, methods)

	return &mcp.CallToolResult{
		Content: []mcp.CallToolContent{
			mcp.TextContent{Text: mdContent},
		},
	}, nil
}

// HandleToolGetDocFunc は、指定された関数に関する情報を返します。
func (h *ToolHandler) HandleToolGetDocFunc(ctx context.Context, req *godoc.ToolGetDocFuncRequest) (*mcp.CallToolResult, error) {
	funcInfo, err := h.parser.GetFuncInfo(req.PackageName, req.FuncName)
	if err != nil {
		return nil, fmt.Errorf("関数情報の取得に失敗しました: %w", err)
	}

	// 例を変換
	var examples []model.Example
	for _, e := range funcInfo.Examples {
		examples = append(examples, model.Example{
			Name:   e.Name,
			Code:   e.Code,
			Output: e.Output,
		})
	}

	// マークダウン形式でフォーマット
	mdContent := model.FormatFuncDocMarkdown(funcInfo.Name, funcInfo.Signature, funcInfo.Comment, examples)

	return &mcp.CallToolResult{
		Content: []mcp.CallToolContent{
			mcp.TextContent{Text: mdContent},
		},
	}, nil
}

// HandleToolGetDocMethod は、指定された構造体のメソッドに関する情報を返します。
func (h *ToolHandler) HandleToolGetDocMethod(ctx context.Context, req *godoc.ToolGetDocMethodRequest) (*mcp.CallToolResult, error) {
	methodInfo, err := h.parser.GetMethodInfo(req.PackageName, req.StructName, req.MethodName)
	if err != nil {
		return nil, fmt.Errorf("メソッド情報の取得に失敗しました: %w", err)
	}

	// 例を変換
	var examples []model.Example
	for _, e := range methodInfo.Examples {
		examples = append(examples, model.Example{
			Name:   e.Name,
			Code:   e.Code,
			Output: e.Output,
		})
	}

	// マークダウン形式でフォーマット
	mdContent := model.FormatMethodDocMarkdown(req.StructName, methodInfo.Name, methodInfo.Signature, methodInfo.Comment, examples)

	return &mcp.CallToolResult{
		Content: []mcp.CallToolContent{
			mcp.TextContent{Text: mdContent},
		},
	}, nil
}

// HandleToolGetDocConstAndVar は、指定されたパッケージの定数と変数に関する情報を返します。
func (h *ToolHandler) HandleToolGetDocConstAndVar(ctx context.Context, req *godoc.ToolGetDocConstAndVarRequest) (*mcp.CallToolResult, error) {
	constInfos, varInfos, err := h.parser.GetConstAndVarInfo(req.PackageName)
	if err != nil {
		return nil, fmt.Errorf("定数と変数の情報の取得に失敗しました: %w", err)
	}

	// 定数と変数の情報を変換
	var constants []model.ConstDoc
	var variables []model.VarDoc

	for _, c := range constInfos {
		constants = append(constants, model.ConstDoc{
			Name:    c.Name,
			Type:    c.Type,
			Value:   c.Value,
			Comment: c.Comment,
		})
	}

	for _, v := range varInfos {
		variables = append(variables, model.VarDoc{
			Name:    v.Name,
			Type:    v.Type,
			Comment: v.Comment,
		})
	}

	// マークダウン形式でフォーマット
	mdContent := model.FormatConstAndVarDocMarkdown(constants, variables)

	return &mcp.CallToolResult{
		Content: []mcp.CallToolContent{
			mcp.TextContent{Text: mdContent},
		},
	}, nil
}
