# pkgDir パラメータ削除リファクタリング計画

## 概要

現在の実装では、`rootDir`と`pkgDir`の両方を受け取っていますが、これを`rootDir`だけを受け取る実装に変更し、`pkgDir`を削除します。

## 変更対象ファイル

1. `internal/parser/parser.go`
2. `cmd/server/main.go`
3. `cmd/parser/main.go`
4. `internal/config/config.go`
5. `internal/config/config_test.go`

## 変更内容詳細

### 1. internal/parser/parser.go

#### 変更前

```go
// New creates a Parser instance by loading Go packages from the specified directory.
// rootDir is the base directory, and pkgDir specifies the pattern to load.
// If pkgDir is empty, "./..." is used as the pattern.
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

	// 以下省略
}
```

#### 変更後

```go
// New creates a Parser instance by loading Go packages from the specified directory.
// rootDir is the base directory where packages will be loaded from.
func New(rootDir string) (*Parser, error) {
	patterns := []string{"./..."}

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

	// 以下省略
}
```

### 2. cmd/server/main.go

#### 変更前

```go
// Parse command line arguments
rootDir := flag.String("root", "", "Root directory path")
pkgDir := flag.String("pkg", "", "Specific package directory (optional)")
flag.Parse()

// Get configuration values
rootPath := config.GetRootDir(*rootDir)
pkgPath := config.GetPkgDir(*pkgDir)

// Initialize parser
p, err := parser.New(rootPath, pkgPath)
```

#### 変更後

```go
// Parse command line arguments
rootDir := flag.String("root", "", "Root directory path")
flag.Parse()

// Get configuration values
rootPath := config.GetRootDir(*rootDir)

// Initialize parser
p, err := parser.New(rootPath)
```

### 3. cmd/parser/main.go

#### 変更前

```go
// Parse command line arguments
rootDir := flag.String("root", "", "Root directory path")
pkgDir := flag.String("pkg", "", "Specific package directory (optional)")
flag.Parse()

// Get configuration values
rootPath := config.GetRootDir(*rootDir)
pkgPath := config.GetPkgDir(*pkgDir)

// Initialize parser
p, err := parser.New(rootPath, pkgPath)
```

#### 変更後

```go
// Parse command line arguments
rootDir := flag.String("root", "", "Root directory path")
flag.Parse()

// Get configuration values
rootPath := config.GetRootDir(*rootDir)

// Initialize parser
p, err := parser.New(rootPath)
```

### 4. internal/config/config.go

#### 変更前

```go
const (
	// Environment variable names
	EnvRootDir = "GODOC_MCP_ROOT_DIR"
	EnvPkgDir  = "GODOC_MCP_PKG_DIR"
)

// GetPkgDir returns the package directory path.
// Priority order:
// 1. Command line argument
// 2. Environment variable
// 3. Empty string (no specification)
func GetPkgDir(cmdPkgDir string) string {
	if cmdPkgDir != "" {
		return cmdPkgDir
	}
	return os.Getenv(EnvPkgDir)
}
```

#### 変更後

```go
const (
	// Environment variable names
	EnvRootDir = "GODOC_MCP_ROOT_DIR"
)
```

`GetPkgDir`関数と`EnvPkgDir`定数を削除します。

### 5. internal/config/config_test.go

`TestGetPkgDir`テスト（52-93行目）を削除します。

## 変更理由

このリファクタリングにより、コードがシンプルになり、メンテナンス性が向上します。`pkgDir`パラメータは実際には常に`"./..."`パターンを使用するか、指定された場合はそのパターンを使用するだけでした。この変更により、常に`"./..."`パターンを使用するようになり、コードの複雑さが軽減されます。