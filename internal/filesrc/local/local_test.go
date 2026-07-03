package local_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/openqt/tcmd/internal/filesrc/local"
)

func TestListIncludesParent(t *testing.T) {
	dir := t.TempDir()
	sub := filepath.Join(dir, "child")
	if err := os.Mkdir(sub, 0o755); err != nil {
		t.Fatal(err)
	}
	src := local.New()
	entries, err := src.List(sub, false)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) == 0 || entries[0].Name != ".." {
		t.Fatalf("expected .. first, got %+v", entries)
	}
}

func TestCopyFile(t *testing.T) {
	dir := t.TempDir()
	srcPath := filepath.Join(dir, "a.txt")
	dstPath := filepath.Join(dir, "b.txt")
	if err := os.WriteFile(srcPath, []byte("hello"), 0o644); err != nil {
		t.Fatal(err)
	}
	src := local.New()
	if err := src.CopyFile(srcPath, dstPath); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(dstPath)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "hello" {
		t.Fatalf("got %q", data)
	}
}
