package config

import (
	"os"
	"path/filepath"
)

const (
	// 環境変数名
	EnvRootDir = "GODOC_MCP_ROOT_DIR"
	EnvPkgDir  = "GODOC_MCP_PKG_DIR"
)

// GetRootDir はディレクトリルートパスを取得します。
// 優先順位:
// 1. コマンドライン引数
// 2. 環境変数
// 3. カレントディレクトリ
func GetRootDir(cmdRootDir string) string {
	if cmdRootDir != "" {
		return cmdRootDir
	}
	if envRootDir := os.Getenv(EnvRootDir); envRootDir != "" {
		return envRootDir
	}
	// カレントディレクトリを取得
	wd, err := os.Getwd()
	if err != nil {
		return "."
	}
	return wd
}

// GetPkgDir はパッケージディレクトリを取得します。
// 優先順位:
// 1. コマンドライン引数
// 2. 環境変数
// 3. 空文字列（指定なし）
func GetPkgDir(cmdPkgDir string) string {
	if cmdPkgDir != "" {
		return cmdPkgDir
	}
	return os.Getenv(EnvPkgDir)
}

// GetAbsPath は指定されたパスを絶対パスに変換します。
func GetAbsPath(path string) (string, error) {
	if filepath.IsAbs(path) {
		return path, nil
	}
	return filepath.Abs(path)
}
