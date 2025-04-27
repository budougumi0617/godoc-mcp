package handler

import (
	"context"
	"fmt"
	"strings"

	godoc "github.com/budougumi0617/godoc-mcp"
	"github.com/budougumi0617/godoc-mcp/internal/parser"
	mcp "github.com/ktr0731/go-mcp"

	// 型情報のためにインポート

	_ "golang.org/x/tools/go/packages"
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

	var packageNames []string
	for _, p := range pkgs {
		packageNames = append(packageNames, p.PkgPath)
	}

	// TODO: 適切なフォーマットを実装する
	responseText := fmt.Sprintf("ロードされたパッケージ:\n- %s", strings.Join(packageNames, "\n- "))

	return &mcp.CallToolResult{
		Content: []mcp.CallToolContent{
			mcp.TextContent{Text: responseText},
		},
	}, nil
}

// HandleToolInspectPackage は、指定されたパッケージの公開されている構造体、メソッド、関数をリストします。
func (h *ToolHandler) HandleToolInspectPackage(ctx context.Context, req *godoc.ToolInspectPackageRequest) (*mcp.CallToolResult, error) {
	pkg, err := h.parser.GetPackage(req.PackageName)
	if err != nil {
		return nil, fmt.Errorf("パッケージの取得に失敗しました: %w", err)
	}

	// TODO: パッケージの構造体、メソッド、関数の情報を抽出する
	responseText := fmt.Sprintf("パッケージ '%s' の情報:\n", pkg.PkgPath)
	responseText += fmt.Sprintf("パッケージ名: %s\n", pkg.Name)
	responseText += "構造体、メソッド、関数のリスト（プレースホルダー）"

	return &mcp.CallToolResult{
		Content: []mcp.CallToolContent{
			mcp.TextContent{Text: responseText},
		},
	}, nil
}

// HandleToolGetDocStruct は、指定された構造体に関する情報を返します。
func (h *ToolHandler) HandleToolGetDocStruct(ctx context.Context, req *godoc.ToolGetDocStructRequest) (*mcp.CallToolResult, error) {
	pkg, err := h.parser.GetPackage(req.PackageName)
	if err != nil {
		return nil, fmt.Errorf("パッケージの取得に失敗しました: %w", err)
	}

	// TODO: 指定された構造体の情報を抽出する
	responseText := fmt.Sprintf("パッケージ '%s' の構造体 '%s' の情報:\n", pkg.PkgPath, req.StructName)
	responseText += "構造体の詳細情報（プレースホルダー）"

	return &mcp.CallToolResult{
		Content: []mcp.CallToolContent{
			mcp.TextContent{Text: responseText},
		},
	}, nil
}

// HandleToolGetDocFunc は、指定された関数に関する情報を返します。
func (h *ToolHandler) HandleToolGetDocFunc(ctx context.Context, req *godoc.ToolGetDocFuncRequest) (*mcp.CallToolResult, error) {
	pkg, err := h.parser.GetPackage(req.PackageName)
	if err != nil {
		return nil, fmt.Errorf("パッケージの取得に失敗しました: %w", err)
	}

	// TODO: 指定された関数の情報を抽出する
	responseText := fmt.Sprintf("パッケージ '%s' の関数 '%s' の情報:\n", pkg.PkgPath, req.FuncName)
	responseText += "関数の詳細情報（プレースホルダー）"

	return &mcp.CallToolResult{
		Content: []mcp.CallToolContent{
			mcp.TextContent{Text: responseText},
		},
	}, nil
}

// HandleToolGetDocMethod は、指定された構造体のメソッドに関する情報を返します。
func (h *ToolHandler) HandleToolGetDocMethod(ctx context.Context, req *godoc.ToolGetDocMethodRequest) (*mcp.CallToolResult, error) {
	pkg, err := h.parser.GetPackage(req.PackageName)
	if err != nil {
		return nil, fmt.Errorf("パッケージの取得に失敗しました: %w", err)
	}

	// TODO: 指定された構造体のメソッドの情報を抽出する
	responseText := fmt.Sprintf("パッケージ '%s' の構造体 '%s' のメソッド '%s' の情報:\n",
		pkg.PkgPath, req.StructName, req.MethodName)
	responseText += "メソッドの詳細情報（プレースホルダー）"

	return &mcp.CallToolResult{
		Content: []mcp.CallToolContent{
			mcp.TextContent{Text: responseText},
		},
	}, nil
}

// HandleToolGetDocConstAndVar は、指定されたパッケージの定数と変数に関する情報を返します。
func (h *ToolHandler) HandleToolGetDocConstAndVar(ctx context.Context, req *godoc.ToolGetDocConstAndVarRequest) (*mcp.CallToolResult, error) {
	pkg, err := h.parser.GetPackage(req.PackageName)
	if err != nil {
		return nil, fmt.Errorf("パッケージの取得に失敗しました: %w", err)
	}

	// TODO: パッケージの定数と変数の情報を抽出する
	responseText := fmt.Sprintf("パッケージ '%s' の定数と変数の情報:\n", pkg.PkgPath)
	responseText += "定数と変数のリスト（プレースホルダー）"

	return &mcp.CallToolResult{
		Content: []mcp.CallToolContent{
			mcp.TextContent{Text: responseText},
		},
	}, nil
}
