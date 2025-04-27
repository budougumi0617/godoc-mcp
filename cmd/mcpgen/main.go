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

	// サーバ定義
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
		Tools:             []codegen.Tool{},
		Prompts:           []codegen.Prompt{},
		ResourceTemplates: []codegen.ResourceTemplate{},
	}

	// コード生成
	if err := codegen.Generate(f, def, "godoc"); err != nil {
		log.Fatalf("failed to generate code: %v", err)
	}
}
