package panel

import (
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/openqt/tcmd/internal/filesrc"
)

// Side identifies the left or right panel.
type Side int

const (
	Left Side = iota
	Right
)

// SortMode matches DC sort commands.
type SortMode int

const (
	SortByName SortMode = iota
	SortByExt
	SortBySize
	SortByDate
)

// Panel holds per-panel file browser state.
type Panel struct {
	Side        Side
	Path        string
	Entries     []filesrc.Entry
	Cursor      int
	Selected    map[string]bool
	ShowHidden  bool
	SortMode    SortMode
	Reverse     bool
	DirSizeHint string
	Renaming    bool
	RenameValue string
}

// New creates a panel at the given path.
func New(side Side, path string) *Panel {
	return &Panel{
		Side:     side,
		Path:     path,
		Selected: make(map[string]bool),
		SortMode: SortByName,
	}
}

// Title returns a short path label.
func (p *Panel) Title() string {
	if p == nil || p.Path == "" {
		return "/"
	}
	return p.Path
}

// Current returns the entry under the cursor.
func (p *Panel) Current() *filesrc.Entry {
	if p == nil || p.Cursor < 0 || p.Cursor >= len(p.Entries) {
		return nil
	}
	entry := p.Entries[p.Cursor]
	return &entry
}

// SelectedEntries returns all marked entries (excluding ..).
func (p *Panel) SelectedEntries() []filesrc.Entry {
	if p == nil {
		return nil
	}
	if len(p.Selected) == 0 {
		if cur := p.Current(); cur != nil && cur.Name != ".." {
			return []filesrc.Entry{*cur}
		}
		return nil
	}
	out := make([]filesrc.Entry, 0, len(p.Selected))
	for _, e := range p.Entries {
		if e.Name == ".." {
			continue
		}
		if p.Selected[e.Path] {
			out = append(out, e)
		}
	}
	return out
}

// ToggleSelect toggles selection for cursor entry.
func (p *Panel) ToggleSelect() {
	cur := p.Current()
	if cur == nil || cur.Name == ".." {
		return
	}
	if p.Selected[cur.Path] {
		delete(p.Selected, cur.Path)
		return
	}
	p.Selected[cur.Path] = true
}

// SelectAll marks all entries except ...
func (p *Panel) SelectAll() {
	for _, e := range p.Entries {
		if e.Name != ".." {
			p.Selected[e.Path] = true
		}
	}
}

// ClearSelection clears all marks.
func (p *Panel) ClearSelection() {
	p.Selected = make(map[string]bool)
}

// InvertSelection inverts marks.
func (p *Panel) InvertSelection() {
	for _, e := range p.Entries {
		if e.Name == ".." {
			continue
		}
		if p.Selected[e.Path] {
			delete(p.Selected, e.Path)
		} else {
			p.Selected[e.Path] = true
		}
	}
}

// Sort applies the panel sort mode.
func (p *Panel) Sort() {
	if len(p.Entries) <= 1 {
		return
	}
	start := 0
	if p.Entries[0].Name == ".." {
		start = 1
	}
	slice := p.Entries[start:]
	sort.SliceStable(slice, func(i, j int) bool {
		a, b := slice[i], slice[j]
		if a.IsDir != b.IsDir {
			if p.Reverse {
				return !a.IsDir
			}
			return a.IsDir
		}
		less := false
		switch p.SortMode {
		case SortByExt:
			less = strings.ToLower(filepath.Ext(a.Name)) < strings.ToLower(filepath.Ext(b.Name))
		case SortBySize:
			less = a.Size < b.Size
		case SortByDate:
			less = a.ModTime.Before(b.ModTime)
		default:
			less = strings.ToLower(a.Name) < strings.ToLower(b.Name)
		}
		if p.Reverse {
			return !less
		}
		return less
	})
}

// MoveCursor moves the cursor by delta clamped to list bounds.
func (p *Panel) MoveCursor(delta int) {
	if len(p.Entries) == 0 {
		p.Cursor = 0
		return
	}
	p.Cursor += delta
	if p.Cursor < 0 {
		p.Cursor = 0
	}
	if p.Cursor >= len(p.Entries) {
		p.Cursor = len(p.Entries) - 1
	}
}

// FormatEntry returns a display line for an entry.
func (p *Panel) FormatEntry(i int, width int) string {
	if i < 0 || i >= len(p.Entries) {
		return ""
	}
	e := p.Entries[i]
	prefix := "  "
	if i == p.Cursor {
		prefix = "> "
	}
	mark := " "
	if p.Selected[e.Path] {
		mark = "*"
	}
	name := e.Name
	if e.IsDir && e.Name != ".." {
		name += "/"
	}
	size := ""
	if !e.IsDir && e.Name != ".." {
		size = formatSize(e.Size)
	}
	line := prefix + mark + " " + padRight(name, width-6) + " " + padLeft(size, 10)
	if i == p.Cursor && p.DirSizeHint != "" && e.IsDir {
		line += " (" + p.DirSizeHint + ")"
	}
	return line
}

func formatSize(n int64) string {
	const unit = 1024
	if n < unit {
		return strings.TrimSpace(padLeft(strconv.FormatInt(n, 10), 6)) + " B"
	}
	div, exp := int64(unit), 0
	for n/div >= unit && exp < 4 {
		div *= unit
		exp++
	}
	val := float64(n) / float64(div)
	suffix := []string{"KiB", "MiB", "GiB", "TiB"}[exp]
	return strings.TrimSpace(padLeft(strconv.Itoa(int(val)), 6)) + " " + suffix
}

func padRight(s string, w int) string {
	if w <= 0 {
		return s
	}
	if len(s) >= w {
		return s[:w]
	}
	return s + strings.Repeat(" ", w-len(s))
}

func padLeft(s string, w int) string {
	if len(s) >= w {
		return s
	}
	return strings.Repeat(" ", w-len(s)) + s
}

// SetDirSizeHint stores temporary directory size text.
func (p *Panel) SetDirSizeHint(hint string) {
	p.DirSizeHint = hint
}

// ClearDirSizeHint clears directory size hint.
func (p *Panel) ClearDirSizeHint() {
	p.DirSizeHint = ""
}
