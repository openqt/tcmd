package app_test

import (
	"os"
	"path/filepath"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/openqt/tcmd/internal/app"
	"github.com/openqt/tcmd/internal/panel"
)

func TestTabSwitchesPanel(t *testing.T) {
	m := app.NewModel("/left", "/right")
	next, _ := m.Update(tea.KeyMsg{Type: tea.KeyTab})
	model := next.(app.Model)
	if model.ActivePanelSide() != panel.Right {
		t.Fatalf("got %v want right", model.ActivePanelSide())
	}
}

func TestAltF4Quits(t *testing.T) {
	m := app.NewModel("/left", "/right")
	next, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'x'}, Alt: true})
	model := next.(app.Model)
	if !model.Quitting() {
		t.Fatal("expected quitting after Alt+X")
	}
}

func TestRenderNotEmpty(t *testing.T) {
	m := app.NewModel("/left", "/right")
	view := m.View()
	if view == "" {
		t.Fatal("expected non-empty view")
	}
}

func TestCopyCreatesFileInInactivePanel(t *testing.T) {
	dir := t.TempDir()
	left := filepath.Join(dir, "left")
	right := filepath.Join(dir, "right")
	if err := os.MkdirAll(left, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(right, 0o755); err != nil {
		t.Fatal(err)
	}
	file := filepath.Join(left, "a.txt")
	if err := os.WriteFile(file, []byte("data"), 0o644); err != nil {
		t.Fatal(err)
	}

	m := app.NewModel(left, right)
	// move cursor to a.txt (index 1 after ..)
	next, _ := m.Update(tea.KeyMsg{Type: tea.KeyDown})
	next, _ = next.Update(tea.KeyMsg{Type: tea.KeyDown})
	next, _ = next.Update(tea.KeyMsg{Type: tea.KeyInsert})
	next, _ = next.Update(tea.KeyMsg{Type: tea.KeyF5})
	model := next.(app.Model)

	dst := filepath.Join(right, "a.txt")
	if _, err := os.Stat(dst); err != nil {
		t.Fatalf("copy failed: %v status=%s", err, model.Status())
	}
}
