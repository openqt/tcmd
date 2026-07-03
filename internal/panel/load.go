package panel

import (
	"path/filepath"

	"github.com/openqt/tcmd/internal/filesrc"
)

// Load reads directory entries into the panel.
func (p *Panel) Load(src filesrc.Source) error {
	entries, err := src.List(p.Path, p.ShowHidden)
	if err != nil {
		return err
	}
	p.Entries = entries
	p.Sort()
	if p.Cursor >= len(p.VisibleEntries()) {
		p.Cursor = max(0, len(p.VisibleEntries())-1)
	}
	p.ClearDirSizeHint()
	if p.History != nil {
		p.History.Push(p.Path)
	}
	return nil
}

// EnterDirectory changes path when cursor is on a directory.
func (p *Panel) EnterDirectory(src filesrc.Source) error {
	cur := p.Current()
	if cur == nil || !cur.IsDir {
		return nil
	}
	p.Path = cur.Path
	p.Cursor = 0
	p.ClearSelection()
	return p.Load(src)
}

// ParentDirectory moves to parent path.
func (p *Panel) ParentDirectory(src filesrc.Source) error {
	parent := filepath.Dir(p.Path)
	if parent == p.Path && p.Path != "/" {
		return nil
	}
	if parent == "" {
		parent = "/"
	}
	p.Path = parent
	p.Cursor = 0
	p.ClearSelection()
	return p.Load(src)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
