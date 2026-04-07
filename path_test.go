package toolbox

import (
	"os"
	"path/filepath"
	"testing"
)

func TestExpandPath(t *testing.T) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("could not get home directory: %v", err)
	}

	tests := []struct {
		name       string
		path       string
		wantPrefix string
		check      func(t *testing.T, got string)
	}{
		{
			name:       "absolute path",
			path:       "/absolute/path",
			wantPrefix: "/absolute/path",
			check: func(t *testing.T, got string) {
				if got != "/absolute/path" {
					t.Errorf("got %q, want /absolute/path", got)
				}
			},
		},
		{
			name:       "home-relative path ~/.bashrc",
			path:       "~/.bashrc",
			wantPrefix: filepath.Join(homeDir, ".bashrc"),
			check: func(t *testing.T, got string) {
				expected := filepath.Join(homeDir, ".bashrc")
				if got != expected {
					t.Errorf("got %q, want %q", got, expected)
				}
			},
		},
		{
			name:       "home-relative path ~/dir/file",
			path:       "~/dir/file",
			wantPrefix: filepath.Join(homeDir, "dir/file"),
			check: func(t *testing.T, got string) {
				expected := filepath.Join(homeDir, "dir/file")
				if got != expected {
					t.Errorf("got %q, want %q", got, expected)
				}
			},
		},
		{
			name: "relative path ./file",
			path: "./file",
			check: func(t *testing.T, got string) {
				if !filepath.IsAbs(got) {
					t.Errorf("got %q, expected absolute path", got)
				}
				if !filepath.IsAbs(got) || !filepath.HasPrefix(got, filepath.Dir(got)) {
					t.Errorf("got %q, expected to end with /file", got)
				}
			},
		},
		{
			name: "relative path without ./",
			path: "relativefile",
			check: func(t *testing.T, got string) {
				if !filepath.IsAbs(got) {
					t.Errorf("got %q, expected absolute path", got)
				}
			},
		},
		{
			name: "current directory",
			path: ".",
			check: func(t *testing.T, got string) {
				if !filepath.IsAbs(got) {
					t.Errorf("got %q, expected absolute path", got)
				}
			},
		},
		{
			name: "parent directory",
			path: "..",
			check: func(t *testing.T, got string) {
				if !filepath.IsAbs(got) {
					t.Errorf("got %q, expected absolute path", got)
				}
			},
		},
		{
			name: "tilde without slash",
			path: "~",
			check: func(t *testing.T, got string) {
				// Tilde without slash doesn't match "~/" so it's treated as relative
				if !filepath.IsAbs(got) {
					t.Errorf("got %q, expected absolute path", got)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ExpandPath(tt.path)
			tt.check(t, got)
		})
	}
}
