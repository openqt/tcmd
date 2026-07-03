package app_test

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/openqt/tcmd/internal/app"
	"github.com/openqt/tcmd/internal/panel"
)

func TestTabSwitchesPanel(t *testing.T) {
	m := app.NewModel("/left", "/right")
	next, _ := m.Update(tea.KeyMsg{Type: tea.KeyTab})
	model := next.(app.Model)
	if model.ActivePanel() != panel.Right {
		t.Fatalf("got %v want right", model.ActivePanel())
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
