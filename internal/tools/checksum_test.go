package tools_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/openqt/tcmd/internal/tools"
)

func TestChecksum(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "a.txt")
	if err := os.WriteFile(path, []byte("test"), 0o644); err != nil {
		t.Fatal(err)
	}
	sum, err := tools.CalcChecksum(path, tools.AlgoMD5)
	if err != nil || len(sum) != 32 {
		t.Fatalf("got %q err=%v", sum, err)
	}
}

func TestSplitCombine(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "big.bin")
	data := make([]byte, 2000)
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatal(err)
	}
	parts, err := tools.SplitFile(path, 1000)
	if err != nil || len(parts) != 2 {
		t.Fatalf("parts=%v err=%v", parts, err)
	}
	target := filepath.Join(dir, "out.bin")
	if err := tools.CombineFiles(parts, target); err != nil {
		t.Fatal(err)
	}
}
