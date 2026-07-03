package hotkeys_test

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/openqt/tcmd/internal/hotkeys"
)

func TestDispatcherLookupTab(t *testing.T) {
	t.Parallel()
	d := hotkeys.NewDispatcher(nil)
	match, ok := d.Lookup(hotkeys.ContextMain, "Tab")
	if !ok {
		t.Fatal("expected Tab binding")
	}
	if match.Command != "cm_SwitchPnl" {
		t.Fatalf("got %q want cm_SwitchPnl", match.Command)
	}
}

func TestDispatcherLookupExit(t *testing.T) {
	t.Parallel()
	d := hotkeys.NewDispatcher(nil)
	for _, key := range []string{"Alt+F4", "Alt+X"} {
		match, ok := d.Lookup(hotkeys.ContextMain, key)
		if !ok {
			t.Fatalf("expected binding for %s", key)
		}
		if match.Command != "cm_Exit" {
			t.Fatalf("%s got %q", key, match.Command)
		}
	}
}

func TestNormalizeKeyCtrlL(t *testing.T) {
	t.Parallel()
	msg := tea.KeyMsg{Type: tea.KeyCtrlL}
	if got := hotkeys.NormalizeKey(msg); got != "Ctrl+L" {
		t.Fatalf("got %q want Ctrl+L", got)
	}
}
