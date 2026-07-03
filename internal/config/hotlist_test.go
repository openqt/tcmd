package config_test

import (
	"testing"

	"github.com/openqt/tcmd/internal/config"
)

func TestHotlistSaveLoad(t *testing.T) {
	base := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", base)
	h := &config.Hotlist{Paths: []string{"/tmp", "/var"}}
	if err := h.Save(); err != nil {
		t.Fatal(err)
	}
	loaded, err := config.LoadHotlist()
	if err != nil {
		t.Fatal(err)
	}
	if len(loaded.Paths) != 2 {
		t.Fatalf("got %v", loaded.Paths)
	}
}
