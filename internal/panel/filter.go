package panel

import (
	"strings"

	"github.com/openqt/tcmd/internal/filesrc"
)

// FilterEntries returns entries matching quick filter text.
func (p *Panel) FilterEntries() []filesrc.Entry {
	if p == nil || p.QuickFilter == "" {
		return p.Entries
	}
	q := strings.ToLower(p.QuickFilter)
	out := make([]filesrc.Entry, 0, len(p.Entries))
	for _, e := range p.Entries {
		if e.Name == ".." {
			out = append(out, e)
			continue
		}
		if strings.Contains(strings.ToLower(e.Name), q) {
			out = append(out, e)
		}
	}
	return out
}

// VisibleEntries returns entries after filter or flat expansion.
func (p *Panel) VisibleEntries() []filesrc.Entry {
	if p.FlatView {
		return p.FlatEntries
	}
	return p.FilterEntries()
}
