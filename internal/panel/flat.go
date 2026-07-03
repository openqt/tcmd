package panel

import (
	"github.com/openqt/tcmd/internal/filesrc"
)

// LoadFlat expands all files under current directory recursively.
func (p *Panel) LoadFlat(src filesrc.Source) error {
	flat, err := collectFlat(src, p.Path, p.ShowHidden)
	if err != nil {
		return err
	}
	p.FlatEntries = flat
	p.FlatView = true
	if p.Cursor >= len(p.FlatEntries) {
		p.Cursor = 0
	}
	return nil
}

// DisableFlat returns to normal directory listing.
func (p *Panel) DisableFlat(src filesrc.Source) error {
	p.FlatView = false
	p.FlatEntries = nil
	return p.Load(src)
}

func collectFlat(src filesrc.Source, root string, showHidden bool) ([]filesrc.Entry, error) {
	entries, err := src.List(root, showHidden)
	if err != nil {
		return nil, err
	}
	out := make([]filesrc.Entry, 0)
	for _, e := range entries {
		if e.Name == ".." {
			continue
		}
		out = append(out, e)
		if e.IsDir {
			sub, err := collectFlat(src, e.Path, showHidden)
			if err != nil {
				return nil, err
			}
			out = append(out, sub...)
		}
	}
	return out, nil
}
