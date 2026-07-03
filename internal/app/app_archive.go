package app

import (
	"path/filepath"

	"github.com/openqt/tcmd/internal/filesrc"
	"github.com/openqt/tcmd/internal/filesrc/archive"
)

func (m *Model) EnterArchive(path string) error {
	m.activePanel().Path = filesrc.EnterZip(path)
	return m.activePanel().Load(m.source)
}

func (m *Model) PackSelected() error {
	entries := m.activePanel().SelectedEntries()
	names := make([]string, 0, len(entries))
	for _, e := range entries {
		names = append(names, e.Name)
	}
	if len(names) == 0 {
		return nil
	}
	target := filepath.Join(m.activePanel().Path, "packed.zip")
	return archive.PackZip(target, m.activePanel().Path, names)
}

func (m *Model) UnpackCursor() error {
	cur := m.activePanel().Current()
	if cur == nil || !archive.IsZipPath(cur.Path) {
		return nil
	}
	return archive.ExtractAll(cur.Path, m.inactivePanel().Path)
}
