package tools_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/openqt/tcmd/internal/tools"
)

func TestFindByName(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "hello.txt"), []byte("world"), 0o644); err != nil {
		t.Fatal(err)
	}
	res, err := tools.Find(tools.FindOptions{Root: dir, NameMask: "*.txt"})
	if err != nil {
		t.Fatal(err)
	}
	if len(res) != 1 {
		t.Fatalf("got %d results", len(res))
	}
}

func TestApplyRename(t *testing.T) {
	got := tools.ApplyRename("file.txt", tools.RenameRule{Prefix: "pre_"})
	if got != "pre_file.txt" {
		t.Fatalf("got %q", got)
	}
}

func TestCompareDirs(t *testing.T) {
	left := t.TempDir()
	right := t.TempDir()
	if err := os.WriteFile(filepath.Join(left, "a.txt"), []byte("1"), 0o644); err != nil {
		t.Fatal(err)
	}
	items, err := tools.CompareDirs(left, right)
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 1 || items[0].Status != "left-only" {
		t.Fatalf("got %+v", items)
	}
}
