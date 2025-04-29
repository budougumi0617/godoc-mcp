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
				Name:        "golang_list_packages",
				Description: "Display a list of Go packages and their package comments. You can check the description and purpose of each package.",
				InputSchema: struct{}{},
			},
			{
				Name:        "golang_inspect_package",
				Description: "List publicly available structs, methods, and functions in the specified Go package. You can check comments for each element.",
				InputSchema: struct {
					PackageName     string `json:"package_name" jsonschema:"description=Package name"`
					IncludeComments bool   `json:"include_comments,omitempty" jsonschema:"description=Whether to include comments,default=true"`
				}{},
			},
			{
				Name:        "golang_get_struct_doc",
				Description: "Display detailed information about the specified Go struct. You can check the struct's comments, fields, methods, and their comments.",
				InputSchema: struct {
					PackageName string `json:"package_name" jsonschema:"description=Package name where the struct is defined"`
					StructName  string `json:"struct_name" jsonschema:"description=Name of the struct"`
				}{},
			},
			{
				Name:        "golang_get_func_doc",
				Description: "Display detailed information about the specified Go function. You can check the function's signature, comments, and usage examples.",
				InputSchema: struct {
					PackageName string `json:"package_name" jsonschema:"description=Package name where the function is defined"`
					FuncName    string `json:"func_name" jsonschema:"description=Name of the function"`
				}{},
			},
			{
				Name:        "golang_get_method_doc",
				Description: "Display detailed information about the specified Go struct method. You can check the method's signature, comments, and usage examples.",
				InputSchema: struct {
					PackageName string `json:"package_name" jsonschema:"description=Package name where the method is defined"`
					StructName  string `json:"struct_name" jsonschema:"description=Name of the struct that owns the method"`
					MethodName  string `json:"method_name" jsonschema:"description=Name of the method"`
				}{},
			},
			{
				Name:        "golang_get_const_and_var_doc",
				Description: "Display detailed information about constants and variables in the specified Go package. You can check the type, value, and comments for each constant and variable.",
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
