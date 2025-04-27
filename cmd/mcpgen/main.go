package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/ktr0731/go-mcp/codegen"
)

func main() {
	// 出力ディレクトリの作成
	outDir := "."
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		log.Fatalf("failed to create output directory: %v", err)
	}

	// 出力ファイルの作成
	f, err := os.Create(filepath.Join(outDir, "mcp.gen.go"))
	if err != nil {
		log.Fatalf("failed to create file: %v", err)
	}
	defer f.Close()

	// サーバー定義の作成
	def := &codegen.ServerDefinition{
		Capabilities: codegen.ServerCapabilities{
			Tools:   &codegen.ToolCapability{},
			Logging: &codegen.LoggingCapability{},
		},
		Implementation: codegen.Implementation{
			Name:    "GoDoc MCP Server",
			Version: "0.0.1",
		},
		// ツール定義
		Tools: []codegen.Tool{
			{
				Name:        "list_packages",
				Description: "Display a list of loaded packages and their package comments",
				InputSchema: struct{}{},
			},
			{
				Name:        "inspect_package",
				Description: "List publicly available structs, methods, and functions in the specified package",
				InputSchema: struct {
					PackageName     string `json:"package_name" jsonschema:"description=Package name"`
					IncludeComments bool   `json:"include_comments,omitempty" jsonschema:"description=Whether to include comments,default=true"`
				}{},
			},
			{
				Name:        "get_doc_struct",
				Description: "Return information about the specified struct",
				InputSchema: struct {
					PackageName string `json:"package_name" jsonschema:"description=Package name of the struct"`
					StructName  string `json:"struct_name" jsonschema:"description=Name of the struct"`
				}{},
			},
			{
				Name:        "get_doc_func",
				Description: "Return information about the specified function",
				InputSchema: struct {
					PackageName string `json:"package_name" jsonschema:"description=Package name of the function"`
					FuncName    string `json:"func_name" jsonschema:"description=Name of the function"`
				}{},
			},
			{
				Name:        "get_doc_method",
				Description: "Return information about the specified method of a struct",
				InputSchema: struct {
					PackageName string `json:"package_name" jsonschema:"description=Package name of the struct"`
					StructName  string `json:"struct_name" jsonschema:"description=Name of the struct"`
					MethodName  string `json:"method_name" jsonschema:"description=Name of the method"`
				}{},
			},
			{
				Name:        "get_doc_const_and_var",
				Description: "Return information about constants and variables in the specified package",
				InputSchema: struct {
					PackageName string `json:"package_name" jsonschema:"description=Package name"`
				}{},
			},
		},
	}

	// コード生成の実行
	if err := codegen.Generate(f, def, "godoc"); err != nil {
		log.Fatalf("failed to generate code: %v", err)
	}
}
