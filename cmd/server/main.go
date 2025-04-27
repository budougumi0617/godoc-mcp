package main

import (
	"context"
	"flag"
	"log"

	godoc "github.com/budougumi0617/godoc-mcp"
	"github.com/budougumi0617/godoc-mcp/internal/handler"
	"github.com/budougumi0617/godoc-mcp/internal/parser"
	mcp "github.com/ktr0731/go-mcp"
	jsonrpc2 "golang.org/x/exp/jsonrpc2"
)

func main() {
	// コマンドライン引数のパース
	rootDir := flag.String("root", ".", "ディレクトリルートパス")
	pkgDir := flag.String("pkg", "", "特定のパッケージディレクトリ（オプション）")
	flag.Parse()

	// パーサーの初期化
	p, err := parser.New(*rootDir, *pkgDir)
	if err != nil {
		log.Fatalf("パーサーの初期化に失敗しました: %v", err)
	}

	// ツールハンドラーの初期化
	toolHandler := handler.NewToolHandler(p)

	// MCPハンドラーの作成
	mcpHandler := godoc.NewHandler(toolHandler)

	// MCPサーバーの起動
	ctx, listener, binder := mcp.NewStdioTransport(context.Background(), mcpHandler, nil)
	srv, err := jsonrpc2.Serve(ctx, listener, binder)
	if err != nil {
		log.Fatalf("サーバーの起動に失敗しました: %v", err)
	}

	// サーバー待機
	srv.Wait()
}
