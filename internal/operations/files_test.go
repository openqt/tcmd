package operations_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/openqt/tcmd/internal/filesrc"
	"github.com/openqt/tcmd/internal/filesrc/local"
	"github.com/openqt/tcmd/internal/operations"
)

func TestCopyEntries(t *testing.T) {
	dir := t.TempDir()
	srcDir := filepath.Join(dir, "src")
	dstDir := filepath.Join(dir, "dst")
	if err := os.MkdirAll(srcDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(dstDir, 0o755); err != nil {
		t.Fatal(err)
	}
	file := filepath.Join(srcDir, "f.txt")
	if err := os.WriteFile(file, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	src := local.New()
	entries := []filesrc.Entry{{Name: "f.txt", Path: file, IsDir: false}}
	if err := operations.CopyEntries(src, entries, dstDir); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(dstDir, "f.txt")); err != nil {
		t.Fatal(err)
	}
}

func TestMkdir(t *testing.T) {
	dir := t.TempDir()
	src := local.New()
	path, err := operations.Mkdir(src, dir, "newdir")
	if err != nil {
		t.Fatal(err)
	}
	if st, err := os.Stat(path); err != nil || !st.IsDir() {
		t.Fatal("expected directory")
	}
}
