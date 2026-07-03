package app

import (
	"fmt"
	"strings"

	"github.com/openqt/tcmd/internal/commands/builtin"
	"github.com/openqt/tcmd/internal/config"
	"github.com/openqt/tcmd/internal/hotkeys"
	"github.com/openqt/tcmd/internal/operations"
)

func (m *Model) OpenConfig() {
	s, _ := config.LoadSettings()
	m.dlg = dialogState{kind: "config"}
	m.status = fmt.Sprintf("Config: hidden=%v sort=%s editor=%s", s.ShowHidden, s.SortBy, s.Editor)
}

func (m *Model) OpenCommandBrowser() {
	names := m.registry.Names()
	m.viewer = "COMMANDS\n" + strings.Join(names, "\n")
	m.status = "Command browser (Shift+F12)"
}

func (m *Model) ImportShortcuts(path string) error {
	files, err := hotkeys.LoadYAML(path)
	if err != nil {
		// create default export if missing
		if err2 := builtin.ExportDefaultShortcuts(path); err2 != nil {
			return err
		}
		files, err = hotkeys.LoadYAML(path)
		if err != nil {
			return err
		}
	}
	for _, f := range files {
		if f.Context == "Main" || f.Context == "" {
			merged := hotkeys.MergeBindings(hotkeys.DefaultMainBindings(), f.Bindings)
			m.dispatcher = hotkeys.NewDispatcher(map[hotkeys.Context][]hotkeys.Binding{
				hotkeys.ContextMain: merged,
			})
		}
	}
	return nil
}

func (m *Model) ShowOperationLog() {
	m.viewer = "OPERATIONS LOG\n" + operations.DefaultLog.Format()
	m.status = "Operation log"
}
