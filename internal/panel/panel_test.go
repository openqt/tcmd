package panel_test

import (
	"testing"

	"github.com/openqt/tcmd/internal/panel"
)

func TestNotebookTabs(t *testing.T) {
	nb := panel.NewNotebook(panel.Left, "/a")
	nb.NewTab("/b")
	if len(nb.Tabs) != 2 {
		t.Fatalf("got %d tabs", len(nb.Tabs))
	}
	nb.NextTab()
	if nb.Current().Path != "/a" {
		t.Fatalf("next tab got %s", nb.Current().Path)
	}
}

func TestHistoryBackForward(t *testing.T) {
	h := panel.NewHistory()
	h.Push("/a")
	h.Push("/b")
	if p, ok := h.Back(); !ok || p != "/a" {
		t.Fatalf("back got %v %v", p, ok)
	}
	if p, ok := h.Forward(); !ok || p != "/b" {
		t.Fatalf("forward got %v %v", p, ok)
	}
}
