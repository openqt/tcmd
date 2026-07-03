package config_test

import (
	"testing"

	"github.com/openqt/tcmd/internal/config"
)

func TestSettingsRoundTrip(t *testing.T) {
	base := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", base)
	s := config.Settings{ShowHidden: true, SortBy: "size", Editor: "vim"}
	if err := config.SaveSettings(s); err != nil {
		t.Fatal(err)
	}
	loaded, err := config.LoadSettings()
	if err != nil {
		t.Fatal(err)
	}
	if !loaded.ShowHidden || loaded.SortBy != "size" {
		t.Fatalf("got %+v", loaded)
	}
}
