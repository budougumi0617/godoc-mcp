package main

import (
	"context"
	"flag"
	"log"

	godoc "github.com/budougumi0617/godoc-mcp"
	"github.com/budougumi0617/godoc-mcp/internal/config"
	"github.com/budougumi0617/godoc-mcp/internal/handler"
	"github.com/budougumi0617/godoc-mcp/internal/parser"
	mcp "github.com/ktr0731/go-mcp"
	jsonrpc2 "golang.org/x/exp/jsonrpc2"
)

func main() {
	// Parse command line arguments
	rootDir := flag.String("root", "", "Root directory path")
	pkgDir := flag.String("pkg", "", "Specific package directory (optional)")
	flag.Parse()

	// Get configuration values
	rootPath := config.GetRootDir(*rootDir)
	pkgPath := config.GetPkgDir(*pkgDir)

	// Initialize parser
	p, err := parser.New(rootPath, pkgPath)
	if err != nil {
		log.Fatalf("Failed to initialize parser: %v", err)
	}

	// Initialize tool handler
	toolHandler := handler.NewToolHandler(p)

	// Create MCP handler
	mcpHandler := godoc.NewHandler(toolHandler)

	// Start MCP server
	ctx, listener, binder := mcp.NewStdioTransport(context.Background(), mcpHandler, nil)
	srv, err := jsonrpc2.Serve(ctx, listener, binder)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

	// Wait for server
	srv.Wait()
}
