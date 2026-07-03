package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/openqt/tcmd/internal/config"
)

func TestDirUsesXDG(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", filepath.Join(t.TempDir(), "xdg"))
	dir, err := config.Dir()
	if err != nil {
		t.Fatal(err)
	}
	if filepath.Base(dir) != "dc-tui" {
		t.Fatalf("got %q", dir)
	}
}

func TestEnsureDirCreatesDirectory(t *testing.T) {
	base := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", base)
	dir, err := config.EnsureDir()
	if err != nil {
		t.Fatal(err)
	}
	info, err := os.Stat(dir)
	if err != nil {
		t.Fatal(err)
	}
	if !info.IsDir() {
		t.Fatal("expected directory")
	}
}
