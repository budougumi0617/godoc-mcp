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
		"Command line argument takes precedence": {
			cmdRootDir: "/path/to/root",
			envRootDir: "/path/to/env",
			want:       "/path/to/root",
		},
		"Environment variable is used": {
			cmdRootDir: "",
			envRootDir: "/path/to/env",
			want:       "/path/to/env",
		},
		"Default value is used": {
			cmdRootDir: "",
			envRootDir: "",
			want:       ".",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			// Set environment variable
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

func TestGetAbsPath(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		path string
		want string
	}{
		"Absolute path is returned as is": {
			path: "/absolute/path",
			want: "/absolute/path",
		},
		"Relative path is converted to absolute path": {
			path: "relative/path",
			want: filepath.Join(t.TempDir(), "relative/path"),
		},
	}

	for name, tt := range tests {
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
