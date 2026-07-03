package hotkeys_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/openqt/tcmd/internal/hotkeys"
)

func TestSaveLoadYAML(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "shortcuts.yaml")
	files := []hotkeys.File{{Context: "Main", Bindings: []hotkeys.Binding{{Key: "F1", Command: "cm_HelpIndex"}}}}
	if err := hotkeys.SaveYAML(path, files); err != nil {
		t.Fatal(err)
	}
	loaded, err := hotkeys.LoadYAML(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(loaded) != 1 || loaded[0].Bindings[0].Command != "cm_HelpIndex" {
		t.Fatalf("got %+v", loaded)
	}
	_ = os.Remove(path)
}
