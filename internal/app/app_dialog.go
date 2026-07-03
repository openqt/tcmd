package app

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/openqt/tcmd/internal/commands/builtin"
	"github.com/openqt/tcmd/internal/config"
	"github.com/openqt/tcmd/internal/operations"
	"github.com/openqt/tcmd/internal/tools"
	"github.com/openqt/tcmd/internal/ui/dialogs"
)

// dialog state
type dialogState struct {
	kind       string
	input      string
	nameMask   string
	content    string
	results    []string
	renameRule tools.RenameRule
	syncItems  []tools.SyncItem
	hotlist    []string
	hotCursor  int
	copyField  int
	copyTarget string
}

func (m *Model) OpenDialog(kind string) {
	m.dlg = dialogState{kind: kind}
	switch kind {
	case "hotlist":
		h, _ := config.LoadHotlist()
		if h != nil {
			m.dlg.hotlist = h.Paths
		}
	case "sync":
		items, _ := builtin.RunSyncCompare(m.leftNB.Current().Path, m.rightNB.Current().Path)
		m.dlg.syncItems = items
	case "copymove":
		m.dlg.copyTarget = m.inactivePanel().Path
		m.dlg.copyField = 0
	case "properties":
		cur := m.activePanel().Current()
		if cur != nil {
			m.status = fmt.Sprintf("Properties: %s", cur.Name)
		}
	}
	m.status = "Dialog: " + kind
}

func (m *Model) DialogAction(action string) {
	switch action {
	case "find_start":
		res, err := builtin.RunFind(m.activePanel().Path, m.dlg.nameMask, m.dlg.content)
		if err != nil {
			m.status = err.Error()
			return
		}
		m.dlg.results = res
		m.status = fmt.Sprintf("Found %d files", len(res))
	case "rename_run":
		names := make([]string, 0)
		for _, e := range m.activePanel().SelectedEntries() {
			names = append(names, e.Name)
		}
		if len(names) == 0 {
			if c := m.activePanel().Current(); c != nil {
				names = []string{c.Name}
			}
		}
		if err := builtin.RunMultiRename(m.activePanel().Path, names, m.dlg.renameRule); err != nil {
			m.status = err.Error()
			return
		}
		_ = m.ReloadActive()
		m.dlg.kind = ""
		m.status = "Renamed"
	case "copy_confirm":
		entries := m.activePanel().SelectedEntries()
		if err := operations.CopyEntries(m.source, entries, m.dlg.copyTarget); err != nil {
			m.status = err.Error()
			return
		}
		_ = m.ReloadBoth()
		m.dlg.kind = ""
		m.status = "Copied"
	}
}

func (m *Model) handleDialogKey(msg tea.KeyMsg, key string) bool {
	if m.dlg.kind == "" {
		return false
	}
	switch key {
	case "Esc":
		m.dlg.kind = ""
		m.status = "Dialog closed"
		return true
	case "F9":
		if m.dlg.kind == "find" {
			m.DialogAction("find_start")
		}
		if m.dlg.kind == "rename" {
			m.DialogAction("rename_run")
		}
		return true
	case "Enter":
		if m.dlg.kind == "hotlist" && m.dlg.hotCursor < len(m.dlg.hotlist) {
			m.activePanel().Path = m.dlg.hotlist[m.dlg.hotCursor]
			_ = m.activePanel().Load(m.source)
			m.dlg.kind = ""
		}
		if m.dlg.kind == "copymove" {
			m.DialogAction("copy_confirm")
		}
		return true
	case "Up":
		if m.dlg.kind == "hotlist" {
			m.dlg.hotCursor--
			if m.dlg.hotCursor < 0 {
				m.dlg.hotCursor = 0
			}
		}
		return true
	case "Down":
		if m.dlg.kind == "hotlist" {
			m.dlg.hotCursor++
		}
		return true
	case "F5", "F6":
		if m.dlg.kind == "copymove" {
			m.dlg.copyField = (m.dlg.copyField + 1) % 4
		}
		return true
	case "Backspace":
		if len(m.dlg.input) > 0 {
			m.dlg.input = m.dlg.input[:len(m.dlg.input)-1]
		}
		if m.dlg.kind == "find" {
			m.dlg.nameMask = m.dlg.input
		}
		return true
	default:
		if msg.Type == tea.KeyRunes {
			m.dlg.input += string(msg.Runes)
			if m.dlg.kind == "find" {
				m.dlg.nameMask = m.dlg.input
			}
			if m.dlg.kind == "rename" {
				m.dlg.renameRule.Prefix = m.dlg.input
			}
		}
	}
	return true
}

func (m Model) renderDialog() string {
	switch m.dlg.kind {
	case "find":
		res := strings.Join(m.dlg.results, "\n")
		return dialogs.Find(m.dlg.nameMask, m.dlg.content, res, 0)
	case "rename":
		names := []string{"file.txt"}
		if c := m.activePanel().Current(); c != nil {
			names = []string{c.Name}
		}
		preview := tools.FormatRenameLog(tools.PreviewRenames(names, m.dlg.renameRule))
		return dialogs.MultiRename(m.dlg.input, preview)
	case "sync":
		var b strings.Builder
		for _, it := range m.dlg.syncItems {
			fmt.Fprintf(&b, "[%s] %s\n", it.Status, it.Path)
		}
		return dialogs.Sync(m.leftNB.Current().Path, m.rightNB.Current().Path, b.String())
	case "hotlist":
		return dialogs.Hotlist(m.dlg.hotlist, m.dlg.hotCursor)
	case "properties":
		cur := m.activePanel().Current()
		if cur == nil {
			return dialogs.Properties("-", "-", "-", "-")
		}
		info, _ := os.Stat(cur.Path)
		mode := ""
		size := fmt.Sprintf("%d", cur.Size)
		if info != nil {
			mode = info.Mode().String()
			size = fmt.Sprintf("%d", info.Size())
		}
		return dialogs.Properties(cur.Name, cur.Path, size, mode)
	case "copymove":
		fields := []string{"name", "ext", "path", "all"}
		field := fields[m.dlg.copyField%len(fields)]
		return dialogs.CopyMove("copy", m.dlg.copyTarget, field)
	}
	return ""
}
