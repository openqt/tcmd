package archive_test

import (
	"archive/zip"
	"os"
	"path/filepath"
	"testing"

	"github.com/openqt/tcmd/internal/filesrc/archive"
)

func TestZipList(t *testing.T) {
	dir := t.TempDir()
	zipPath := filepath.Join(dir, "a.zip")
	f, err := os.Create(zipPath)
	if err != nil {
		t.Fatal(err)
	}
	w := zip.NewWriter(f)
	_, _ = w.Create("hello.txt")
	w.Close()
	f.Close()

	z := archive.NewZip(zipPath)
	entries, err := z.List("", true)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 || entries[0].Name != "hello.txt" {
		t.Fatalf("got %+v", entries)
	}
}

func TestPackExtract(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "f.txt")
	if err := os.WriteFile(src, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	zipPath := filepath.Join(dir, "out.zip")
	if err := archive.PackZip(zipPath, dir, []string{"f.txt"}); err != nil {
		t.Fatal(err)
	}
	dst := filepath.Join(dir, "extracted")
	if err := os.Mkdir(dst, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := archive.ExtractAll(zipPath, dst); err != nil {
		t.Fatal(err)
	}
}
