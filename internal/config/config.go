package config

import (
	"os"
	"path/filepath"
)

const (
	// Environment variable names
	EnvRootDir = "GODOC_MCP_ROOT_DIR"
	EnvPkgDir  = "GODOC_MCP_PKG_DIR"
)

// GetRootDir returns the root directory path.
// Priority order:
// 1. Command line argument
// 2. Environment variable
// 3. Current directory
func GetRootDir(cmdRootDir string) string {
	if cmdRootDir != "" {
		return cmdRootDir
	}
	if envRootDir := os.Getenv(EnvRootDir); envRootDir != "" {
		return envRootDir
	}
	// Get current directory
	wd, err := os.Getwd()
	if err != nil {
		return "."
	}
	return wd
}

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

// GetAbsPath converts the specified path to an absolute path.
func GetAbsPath(path string) (string, error) {
	if filepath.IsAbs(path) {
		return path, nil
	}
	return filepath.Abs(path)
}
