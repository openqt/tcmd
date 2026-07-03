package viewer_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/openqt/tcmd/internal/ui/viewer"
)

func TestRenderHex(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "a.bin")
	if err := os.WriteFile(path, []byte{0x01, 0x02}, 0o644); err != nil {
		t.Fatal(err)
	}
	out, err := viewer.Render(path, viewer.ModeHex, 100)
	if err != nil {
		t.Fatal(err)
	}
	if out == "" {
		t.Fatal("empty output")
	}
}
