package panel

import "os"

// Side identifies the left or right panel.
type Side int

const (
	Left Side = iota
	Right
)

// Panel holds per-panel UI state for Phase 0.
type Panel struct {
	Side Side
	Path string
}

// New creates a panel rooted at the given path.
func New(side Side, path string) *Panel {
	if path == "" {
		path, _ = os.Getwd()
	}
	return &Panel{Side: side, Path: path}
}

// Title returns a short label for tabs and headers.
func (p *Panel) Title() string {
	if p == nil || p.Path == "" {
		return "/"
	}
	return p.Path
}
