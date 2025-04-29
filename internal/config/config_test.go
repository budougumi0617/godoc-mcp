package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGetRootDir(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		cmdRootDir string
		envRootDir string
		want       string
	}{
		"コマンドライン引数が優先": {
			cmdRootDir: "/path/to/root",
			envRootDir: "/path/to/env",
			want:       "/path/to/root",
		},
		"環境変数が使用される": {
			cmdRootDir: "",
			envRootDir: "/path/to/env",
			want:       "/path/to/env",
		},
		"デフォルト値が使用される": {
			cmdRootDir: "",
			envRootDir: "",
			want:       ".",
		},
	}

	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			// 環境変数の設定
			if tt.envRootDir != "" {
				os.Setenv(EnvRootDir, tt.envRootDir)
				defer os.Unsetenv(EnvRootDir)
			}

			got := GetRootDir(tt.cmdRootDir)
			if got != tt.want {
				t.Errorf("GetRootDir() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetPkgDir(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		cmdPkgDir string
		envPkgDir string
		want      string
	}{
		"コマンドライン引数が優先": {
			cmdPkgDir: "/path/to/pkg",
			envPkgDir: "/path/to/env/pkg",
			want:      "/path/to/pkg",
		},
		"環境変数が使用される": {
			cmdPkgDir: "",
			envPkgDir: "/path/to/env/pkg",
			want:      "/path/to/env/pkg",
		},
		"空文字列が返される": {
			cmdPkgDir: "",
			envPkgDir: "",
			want:      "",
		},
	}

	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			// 環境変数の設定
			if tt.envPkgDir != "" {
				os.Setenv(EnvPkgDir, tt.envPkgDir)
				defer os.Unsetenv(EnvPkgDir)
			}

			got := GetPkgDir(tt.cmdPkgDir)
			if got != tt.want {
				t.Errorf("GetPkgDir() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetAbsPath(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		path string
		want string
	}{
		"絶対パスはそのまま返される": {
			path: "/absolute/path",
			want: "/absolute/path",
		},
		"相対パスは絶対パスに変換される": {
			path: "relative/path",
			want: filepath.Join(t.TempDir(), "relative/path"),
		},
	}

	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got, err := GetAbsPath(tt.path)
			if err != nil {
				t.Fatalf("GetAbsPath() error = %v", err)
			}
			if got != tt.want {
				t.Errorf("GetAbsPath() = %v, want %v", got, tt.want)
			}
		})
	}
}
