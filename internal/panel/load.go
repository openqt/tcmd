package panel

import (
	"path/filepath"

	"github.com/openqt/tcmd/internal/filesrc"
	"github.com/openqt/tcmd/internal/filesrc/archive"
)

// Load reads directory entries into the panel.
func (p *Panel) Load(src filesrc.Source) error {
	if arch, inner, ok := archive.ParseArchivePath(p.Path); ok {
		z := archive.NewZip(arch)
		zentries, err := z.List(inner, p.ShowHidden)
		if err != nil {
			return err
		}
		p.Entries = make([]filesrc.Entry, len(zentries))
		for i, e := range zentries {
			p.Entries[i] = filesrc.Entry{
				Name:    e.Name,
				Path:    archive.JoinArchivePath(arch, e.Path),
				IsDir:   e.IsDir,
				Size:    e.Size,
				ModTime: e.ModTime,
			}
		}
		if p.Cursor >= len(p.VisibleEntries()) {
			p.Cursor = max(0, len(p.VisibleEntries())-1)
		}
		p.ClearDirSizeHint()
		return nil
	}

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
	if arch, inner, ok := archive.ParseArchivePath(p.Path); ok {
		if inner == "" {
			p.Path = arch
		} else {
			parent := filepath.Dir(inner)
			if parent == "." {
				p.Path = archive.JoinArchivePath(arch, "")
			} else {
				p.Path = archive.JoinArchivePath(arch, parent)
			}
		}
		p.Cursor = 0
		p.ClearSelection()
		return p.Load(src)
	}
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
