package backup

import (
	"os"
	"path/filepath"
	"testing"
)

func TestResolveFormat(t *testing.T) {
	d := newPostgresDumper()

	tests := []struct {
		name       string
		fileExt    string
		formatFlag string
		want       string
		wantErr    bool
	}{
		{name: "both set errors", fileExt: ".sql", formatFlag: "custom", wantErr: true},
		{name: "extension only", fileExt: ".sql", want: "plain"},
		{name: "extension only dump", fileExt: ".dump", want: "custom"},
		{name: "unknown extension errors", fileExt: ".txt", wantErr: true},
		{name: "flag only", formatFlag: "tar", want: "tar"},
		{name: "unknown flag errors", formatFlag: "bogus", wantErr: true},
		{name: "neither set uses default", want: "custom"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ResolveFormat(tt.fileExt, tt.formatFlag, d)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Fatalf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestResolvePath(t *testing.T) {
	dir := t.TempDir()

	t.Run("existing directory", func(t *testing.T) {
		got := ResolvePath(dir, "mydb", "dump")
		gotDir := filepath.Dir(got)
		if gotDir != dir {
			t.Fatalf("got dir %q, want %q", gotDir, dir)
		}
		if filepath.Ext(got) != ".dump" {
			t.Fatalf("expected .dump extension, got %q", got)
		}
	})

	t.Run("file with extension used as-is", func(t *testing.T) {
		path := filepath.Join(dir, "out.sql")
		got := ResolvePath(path, "mydb", "sql")
		if got != path {
			t.Fatalf("got %q, want %q", got, path)
		}
	})

	t.Run("file without extension gets extension appended", func(t *testing.T) {
		path := filepath.Join(dir, "out")
		got := ResolvePath(path, "mydb", "dump")
		want := path + ".dump"
		if got != want {
			t.Fatalf("got %q, want %q", got, want)
		}
	})

	t.Run("empty path uses cwd", func(t *testing.T) {
		cwd, err := os.Getwd()
		if err != nil {
			t.Fatal(err)
		}
		got := ResolvePath("", "mydb", "dump")
		if filepath.Dir(got) != cwd {
			t.Fatalf("got dir %q, want cwd %q", filepath.Dir(got), cwd)
		}
	})
}
