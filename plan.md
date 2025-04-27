# MCPサーバ実装計画（更新版）

## 1. 設計概要

### 1.1 機能要件
- ローカルディレクトリのGoモジュールのgodocをMCP経由でアクセス可能にする
- 起動時のオプション
  - ディレクトリルートパス（必須）
  - 特定のパッケージディレクトリ（オプション）
- 起動時にGoファイルを再帰的に読み込み、ドキュメント情報を保持

### 1.2 使用ライブラリ
- go-mcp
  - MCPサーバーの実装に使用
  - 型安全性を確保
  - コード生成による効率的な開発
  - パッケージ名: `github.com/ktr0731/go-mcp`
  - 主要なサブパッケージ:
    - `github.com/ktr0731/go-mcp/codegen`: コード生成用
    - `github.com/ktr0731/go-mcp/mcp`: MCPサーバー実装用
- golang.org/x/tools/go/packages
  - Goパッケージの解析に使用
  - 高レベルAPIによる効率的な実装
  - 依存関係の自動解決
  - 型情報の完全な解析
  - Go Modulesとの互換性
- golang.org/x/tools/godoc
  - Goドキュメントの解析と表示に使用
  - コメントからドキュメント情報を抽出

### 1.3 MCPツール
1. `list_packages`
   - 読み込んだパッケージの一覧表示
   - パッケージコメントの出力

2. `inspect_package`
   - パッケージ内の公開構造体、メソッド、関数のリスト表示
   - 引数: パッケージ名、コメント有無（オプション）

3. `get_doc_struct`
   - 構造体情報の取得
   - 引数: パッケージ名、構造体名

4. `get_doc_func`
   - 関数情報の取得
   - 引数: パッケージ名、関数名

5. `get_doc_method`
   - 構造体メソッド情報の取得
   - 引数: パッケージ名、構造体名、メソッド名

6. `get_doc_const_and_var`
   - 定数・変数情報の取得
   - 引数: パッケージ名

## 2. 実装ステップ

### 2.1 プロジェクト構造の整理
- [x] プロジェクト構造の作成
  - [x] `cmd/mcpgen/`ディレクトリの作成
  - [x] `cmd/server/`ディレクトリの作成
  - [x] `internal/handler/`ディレクトリの作成
  - [x] `internal/model/`ディレクトリの作成
  - [x] `internal/parser/`ディレクトリの作成（Goファイル解析用）

### 2.2 コード生成の設定
- [x] `cmd/mcpgen/main.go`の実装
  - [x] サーバー定義の作成
    ```go
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
                Description: "読み込んだパッケージの一覧とパッケージコメントを表示します",
                InputSchema: struct{}{},
            },
            {
                Name:        "inspect_package",
                Description: "指定されたパッケージ内の公開されている構造体、メソッド、関数をリストします",
                InputSchema: struct {
                    PackageName    string `json:"package_name" jsonschema:"description=パッケージ名"`
                    IncludeComments bool   `json:"include_comments,omitempty" jsonschema:"description=コメントを含めるかどうか,default=true"`
                }{},
            },
            {
                Name:        "get_doc_struct",
                Description: "指定された構造体の情報を返します",
                InputSchema: struct {
                    PackageName string `json:"package_name" jsonschema:"description=構造体のパッケージ名"`
                    StructName  string `json:"struct_name" jsonschema:"description=構造体の名前"`
                }{},
            },
            {
                Name:        "get_doc_func",
                Description: "指定された関数の情報を返します",
                InputSchema: struct {
                    PackageName string `json:"package_name" jsonschema:"description=関数のパッケージ名"`
                    FuncName    string `json:"func_name" jsonschema:"description=関数の名前"`
                }{},
            },
            {
                Name:        "get_doc_method",
                Description: "指定された構造体のメソッドの情報を返します",
                InputSchema: struct {
                    PackageName string `json:"package_name" jsonschema:"description=構造体のパッケージ名"`
                    StructName  string `json:"struct_name" jsonschema:"description=構造体の名前"`
                    MethodName  string `json:"method_name" jsonschema:"description=メソッド名"`
                }{},
            },
            {
                Name:        "get_doc_const_and_var",
                Description: "指定されたパッケージの定数や変数情報を返します",
                InputSchema: struct {
                    PackageName string `json:"package_name" jsonschema:"description=パッケージ名"`
                }{},
            },
        },
    }
    ```
  - [x] コード生成の実行
    ```go
    if err := codegen.Generate(f, def, "godoc"); err != nil {
        log.Fatalf("failed to generate code: %v", err)
    }
    ```

### 2.3 サーバー実装
- [ ] `cmd/server/main.go`の実装
  - [ ] コマンドラインオプションの実装
    ```go
    rootDir := flag.String("root", ".", "ディレクトリルートパス")
    pkgDir := flag.String("pkg", "", "特定のパッケージディレクトリ（オプション）")
    flag.Parse()
    ```
  - [ ] パーサーの初期化
    ```go
    p, err := parser.New(*rootDir, *pkgDir)
    if err != nil {
        log.Fatalf("failed to initialize parser: %v", err)
    }
    ```
  - [ ] MCPサーバーの起動
    ```go
    handler := NewHandler(&ToolHandler{parser: p})
    
    ctx, listener, binder := mcp.NewStdioTransport(context.Background(), handler, nil)
    srv, err := jsonrpc2.Serve(ctx, listener, binder)
    if err != nil {
        log.Fatalf("failed to serve: %v", err)
    }
    
    srv.Wait()
    ```
  - [ ] エラーハンドリング
    - [ ] パーサー初期化エラーの処理
    - [ ] サーバー起動エラーの処理

### 2.4 ツールハンドラーの実装
- [ ] `internal/handler/handler.go`の実装
  - [ ] ハンドラー構造体の定義
    ```go
    type ToolHandler struct {
        parser *parser.Parser
    }
    ```
  - [ ] `list_packages`ハンドラーの実装
    ```go
    func (h *ToolHandler) HandleToolListPackages(ctx context.Context, req *ToolListPackagesRequest) (*mcp.CallToolResult, error) {
        packages := h.parser.GetAllPackages()
        return &mcp.CallToolResult{
            Content: []mcp.CallToolContent{
                mcp.TextContent{Text: formatPackageList(packages)},
            },
        }, nil
    }
    ```
  - [ ] `inspect_package`ハンドラーの実装
    ```go
    func (h *ToolHandler) HandleToolInspectPackage(ctx context.Context, req *ToolInspectPackageRequest) (*mcp.CallToolResult, error) {
        pkg, err := h.parser.GetPackage(req.PackageName)
        if err != nil {
            return nil, fmt.Errorf("package not found: %s", req.PackageName)
        }
        return &mcp.CallToolResult{
            Content: []mcp.CallToolContent{
                mcp.TextContent{Text: formatPackageInspection(pkg, req.IncludeComments)},
            },
        }, nil
    }
    ```
  - [ ] `get_doc_struct`ハンドラーの実装
    ```go
    func (h *ToolHandler) HandleToolGetDocStruct(ctx context.Context, req *ToolGetDocStructRequest) (*mcp.CallToolResult, error) {
        pkg, err := h.parser.GetPackage(req.PackageName)
        if err != nil {
            return nil, fmt.Errorf("package not found: %s", req.PackageName)
        }
        
        structInfo, err := pkg.GetStruct(req.StructName)
        if err != nil {
            return nil, fmt.Errorf("struct not found: %s in package %s", req.StructName, req.PackageName)
        }
        
        return &mcp.CallToolResult{
            Content: []mcp.CallToolContent{
                mcp.TextContent{Text: formatStructDoc(structInfo)},
            },
        }, nil
    }
    ```
  - [ ] `get_doc_func`ハンドラーの実装
  - [ ] `get_doc_method`ハンドラーの実装
  - [ ] `get_doc_const_and_var`ハンドラーの実装
  - [ ] エラーハンドリング
    - [ ] パッケージが見つからない場合
    - [ ] 構造体や関数が見つからない場合
    - [ ] 解析エラーの場合
  - [ ] レスポンスの整形
    - [ ] マークダウン形式での出力
    - [ ] コードブロックの適切な使用

### 2.5 Goファイル解析の実装
- [ ] `internal/parser/parser.go`の実装
  - [ ] `golang.org/x/tools/go/packages`を使用したパッケージロード
    ```go
    func LoadPackages(dir string, patterns ...string) ([]*packages.Package, error) {
        if len(patterns) == 0 {
            patterns = []string{"./..."}
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
            Dir: dir,
            Tests: false,
        }
        
        return packages.Load(cfg, patterns...)
    }
    ```
  - [ ] パッケージ情報の抽出
    ```go
    func extractPackageInfo(pkg *packages.Package) *model.Package {
        info := &model.Package{
            Name:       pkg.Name,
            ImportPath: pkg.PkgPath,
        }
        
        // パッケージドキュメントの抽出
        docPkg := extractDocPackage(pkg)
        if docPkg != nil {
            info.Doc = docPkg.Doc
        }
        
        // 型情報の抽出
        scope := pkg.Types.Scope()
        for _, name := range scope.Names() {
            obj := scope.Lookup(name)
            if !obj.Exported() {
                continue
            }
            
            switch obj := obj.(type) {
            case *types.TypeName:
                if named, ok := obj.Type().(*types.Named); ok {
                    if struct_, ok := named.Underlying().(*types.Struct); ok {
                        info.Structs = append(info.Structs, extractStructInfo(obj, named, struct_, pkg, docPkg))
                    }
                }
            case *types.Func:
                info.Functions = append(info.Functions, extractFunctionInfo(obj, pkg, docPkg))
            case *types.Const:
                info.Constants = append(info.Constants, extractConstantInfo(obj, pkg, docPkg))
            case *types.Var:
                info.Variables = append(info.Variables, extractVariableInfo(obj, pkg, docPkg))
            }
        }
        
        return info
    }
    ```
  - [ ] 構造体情報の抽出
  - [ ] 関数情報の抽出
  - [ ] メソッド情報の抽出
  - [ ] 定数・変数情報の抽出
  - [ ] データモデルへの変換
  - [ ] パフォーマンス最適化
    - [ ] 並行処理による解析の高速化
    - [ ] メモリ使用量の最適化
    - [ ] インクリメンタル解析の実装（オプション）

### 2.6 データモデルの設計
- [ ] `internal/model/model.go`の実装
  - [ ] パッケージ情報の構造体
    ```go
    type Package struct {
        Name       string
        ImportPath string
        Doc        string
        Structs    []*Struct
        Functions  []*Function
        Constants  []*Constant
        Variables  []*Variable
        
        GetStruct(name string) (*Struct, error)
        GetFunction(name string) (*Function, error)
    }
    ```
  - [ ] 構造体情報の構造体
    ```go
    type Struct struct {
        Name    string
        Doc     string
        Fields  []*Field
        Methods []*Method
        
        GetMethod(name string) (*Method, error)
    }
    ```
  - [ ] 関数情報の構造体
    ```go
    type Function struct {
        Name      string
        Doc       string
        Signature string
        Params    []*Parameter
        Results   []*Parameter
    }
    ```
  - [ ] メソッド情報の構造体
    ```go
    type Method struct {
        Name         string
        Doc          string
        Signature    string
        ReceiverType string
        Params       []*Parameter
        Results      []*Parameter
    }
    ```
  - [ ] フィールド情報の構造体
    ```go
    type Field struct {
        Name       string
        Type       string
        Doc        string
        IsExported bool
    }
    ```
  - [ ] パラメータ情報の構造体
    ```go
    type Parameter struct {
        Name string
        Type string
    }
    ```
  - [ ] 定数情報の構造体
    ```go
    type Constant struct {
        Name  string
        Type  string
        Value string
        Doc   string
    }
    ```
  - [ ] 変数情報の構造体
    ```go
    type Variable struct {
        Name string
        Type string
        Doc  string
    }
    ```

### 2.7 レスポンス設計
- [ ] 各ツールのレスポンス形式の定義
  - [ ] 共通レスポンス構造の設計
    - [ ] `mcp.CallToolResult`構造体の活用
    - [ ] エラーハンドリングの統一
  - [ ] `list_packages`のレスポンス
    ```go
    return &mcp.CallToolResult{
        Content: []mcp.CallToolContent{
            mcp.TextContent{Text: formatPackageList(packages)},
        },
    }, nil
    ```
  - [ ] `inspect_package`のレスポンス
    ```go
    return &mcp.CallToolResult{
        Content: []mcp.CallToolContent{
            mcp.TextContent{Text: formatPackageInspection(pkg, includeComments)},
        },
    }, nil
    ```
  - [ ] `get_doc_struct`のレスポンス
    ```go
    return &mcp.CallToolResult{
        Content: []mcp.CallToolContent{
            mcp.TextContent{Text: formatStructDoc(structInfo)},
        },
    }, nil
    ```
  - [ ] `get_doc_func`のレスポンス
    ```go
    return &mcp.CallToolResult{
        Content: []mcp.CallToolContent{
            mcp.TextContent{Text: formatFuncDoc(funcInfo)},
        },
    }, nil
    ```
  - [ ] `get_doc_method`のレスポンス
    ```go
    return &mcp.CallToolResult{
        Content: []mcp.CallToolContent{
            mcp.TextContent{Text: formatMethodDoc(methodInfo)},
        },
    }, nil
    ```
  - [ ] `get_doc_const_and_var`のレスポンス
    ```go
    return &mcp.CallToolResult{
        Content: []mcp.CallToolContent{
            mcp.TextContent{Text: formatConstAndVarDoc(constVarInfo)},
        },
    }, nil
    ```

- [ ] フォーマット関数の実装
  - [ ] `formatPackageList`: パッケージリストのフォーマット
    - [ ] パッケージ名と簡潔な説明を含む
    - [ ] マークダウン形式での出力
  - [ ] `formatPackageInspection`: パッケージ検査結果のフォーマット
    - [ ] 構造体、関数、メソッドのリスト
    - [ ] オプションでコメントを含める
  - [ ] `formatStructDoc`: 構造体ドキュメントのフォーマット
    - [ ] 構造体コメント
    - [ ] フィールドリストとコメント
    - [ ] メソッドリストとコメント
  - [ ] `formatFuncDoc`: 関数ドキュメントのフォーマット
    - [ ] 関数シグネチャ
    - [ ] 関数コメント
  - [ ] `formatMethodDoc`: メソッドドキュメントのフォーマット
    - [ ] メソッドシグネチャ
    - [ ] メソッドコメント
  - [ ] `formatConstAndVarDoc`: 定数・変数ドキュメントのフォーマット
    - [ ] 定数・変数の型情報
    - [ ] 定数・変数コメント

- [ ] 出力形式の統一
  - [ ] マークダウン形式での出力
  - [ ] コードブロックの適切な使用
  - [ ] 見出しレベルの一貫性
  - [ ] リストの適切な使用

### 2.8 テストの実装
- [ ] ユニットテストの作成
  - [ ] ハンドラーのテスト
    - [ ] 各ツールハンドラーの入出力テスト
    - [ ] エラーケースのテスト
  - [ ] パーサーのテスト
    - [ ] パッケージロードのテスト
    - [ ] 情報抽出のテスト
    - [ ] エラーケースのテスト
  - [ ] モデルのテスト
    - [ ] データモデルの操作テスト
    - [ ] 検索機能のテスト
- [ ] 統合テストの作成
  - [ ] サーバー全体のテスト
    - [ ] コマンドラインオプションのテスト
    - [ ] エンドツーエンドのテスト
  - [ ] テストデータの準備
    - [ ] サンプルGoコードの作成
    - [ ] 期待される出力の定義

## 3. 実装の優先順位

1. 基本的なMCPサーバーの構築
   - [ ] コード生成の設定
   - [ ] サーバーの起動
   - [ ] 基本的なツールの実装

2. Goファイル解析の実装
   - [ ] `golang.org/x/tools/go/packages`を使用したパーサーの実装
   - [ ] 情報抽出の実装
   - [ ] データモデルの実装

3. ツールハンドラーの実装
   - [ ] 各ツールのハンドラー実装
   - [ ] レスポンス形式の統一
   - [ ] エラーハンドリングの実装

4. テストとドキュメントの整備
   - [ ] ユニットテスト
   - [ ] 統合テスト
   - [ ] APIドキュメント
   - [ ] 使用例の作成

## 4. 注意点

- [ ] 型安全性の確保
  - [ ] コード生成による型安全なインターフェースの実現
  - [ ] 適切な型変換の実装

- [ ] エラーハンドリングの徹底
  - [ ] 具体的なエラーケースの特定と対応
  - [ ] ユーザーフレンドリーなエラーメッセージ
  - [ ] エラーの適切な伝播

- [ ] パフォーマンスの考慮
  - [ ] 大規模プロジェクトでの動作確認
  - [ ] メモリ使用量の最適化
  - [ ] 並行処理の活用

- [ ] Goファイル解析の正確性
  - [ ] コメント解析の正確性
  - [ ] 型情報の正確な抽出
  - [ ] エッジケースの考慮

## 5. 次のステップ

1. [x] プロジェクト構造の作成
2. [x] `cmd/mcpgen/main.go`の実装
3. [x] 基本的なツールの定義
4. [x] コード生成の実行
5. [ ] `internal/model/model.go`の実装
6. [ ] `internal/parser/parser.go`の実装
7. [ ] `internal/handler/handler.go`の実装
8. [ ] `cmd/server/main.go`の実装
9. [ ] テストの実装
10. [ ] ドキュメントの整備