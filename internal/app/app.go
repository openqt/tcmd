package app

import (
	"os"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/openqt/tcmd/internal/commands"
	"github.com/openqt/tcmd/internal/commands/builtin"
	"github.com/openqt/tcmd/internal/filesrc"
	"github.com/openqt/tcmd/internal/filesrc/local"
	"github.com/openqt/tcmd/internal/hotkeys"
	"github.com/openqt/tcmd/internal/operations"
	"github.com/openqt/tcmd/internal/panel"
	"github.com/openqt/tcmd/internal/platform"
	"github.com/openqt/tcmd/internal/ui/mainview"
)

// Model is the root Bubble Tea model.
type Model struct {
	width    int
	height   int
	active   panel.Side
	leftNB   *panel.Notebook
	rightNB  *panel.Notebook
	status   string
	prompt   string
	promptFn string
	cmdLine  string
	cmdFocus bool
	cmdHist  []string
	cmdIdx   int
	viewer   string
	quickView string
	quitting bool
	driveIdx  int
	driveFor  panel.Side
	driveMode bool
	dlg       dialogState

	source     filesrc.Source
	registry   *commands.Registry
	dispatcher *hotkeys.Dispatcher
}

// NewModel constructs the application model.
func NewModel(leftPath, rightPath string) *Model {
	source := local.New()
	registry := commands.NewRegistry()
	m := &Model{
		active:     panel.Left,
		leftNB:     panel.NewNotebook(panel.Left, leftPath),
		rightNB:    panel.NewNotebook(panel.Right, rightPath),
		status:     "dc-tui Phase 2",
		source:     source,
		registry:   registry,
		dispatcher: hotkeys.NewDispatcher(nil),
	}
	builtin.RegisterPhase0(registry)
	builtin.RegisterPhase1(registry, m)
	builtin.RegisterPhase2(registry, m)
	builtin.RegisterPhase3(registry, m)
	_ = m.leftNB.Current().Load(source)
	_ = m.rightNB.Current().Load(source)
	return m
}

// Init implements tea.Model.
func (m Model) Init() tea.Cmd { return nil }

// Update implements tea.Model.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		if m.handleInput(msg) {
			return m, nil
		}
	}
	return m, nil
}

func (m *Model) handleInput(msg tea.KeyMsg) bool {
	key := hotkeys.NormalizeKey(msg)

	if m.driveMode {
		m.handleDriveKey(key)
		return true
	}
	if m.dlg.kind != "" {
		return m.handleDialogKey(msg, key)
	}
	if m.cmdFocus {
		m.handleCmdLineKey(msg, key)
		return true
	}
	if m.prompt != "" {
		m.handlePromptKey(msg, key)
		return true
	}
	if m.viewer != "" && (key == "Esc" || key == "Q") {
		m.viewer = ""
		m.status = "Viewer closed"
		return true
	}
	if m.activePanel().Renaming {
		m.handleRenameKey(msg, key)
		return true
	}

	// Quick search: printable keys without modifiers
	if msg.Type == tea.KeyRunes && len(msg.Runes) == 1 && !msg.Alt && key != "Tab" {
		if key == " " {
			// fall through to command dispatch for space
		} else if key[0] >= 'A' && key[0] <= 'Z' || (key[0] >= '0' && key[0] <= '9') {
			m.activePanel().QuickSearch += strings.ToLower(key)
			m.activePanel().QuickFilter = m.activePanel().QuickSearch
			m.status = "Search: " + m.activePanel().QuickSearch
			return true
		}
	}

	if match, ok := m.dispatcher.LookupMsg(hotkeys.ContextMain, msg); ok {
		m.execute(match.Command, match.Params)
		return true
	}
	return false
}

func (m *Model) handleCmdLineKey(msg tea.KeyMsg, key string) {
	switch key {
	case "Esc":
		m.cmdFocus = false
		m.cmdLine = ""
	case "Enter":
		m.cmdHist = append(m.cmdHist, m.cmdLine)
		m.cmdIdx = len(m.cmdHist)
		if err := m.ExecuteCmdLine(); err != nil {
			m.status = err.Error()
		}
		m.cmdFocus = false
	case "Backspace":
		if len(m.cmdLine) > 0 {
			m.cmdLine = m.cmdLine[:len(m.cmdLine)-1]
		}
	default:
		if msg.Type == tea.KeyRunes {
			m.cmdLine += string(msg.Runes)
		}
	}
}

func (m *Model) handleRenameKey(msg tea.KeyMsg, key string) {
	p := m.activePanel()
	switch key {
	case "Esc":
		p.Renaming = false
		p.RenameValue = ""
	case "Enter":
		cur := p.Current()
		if cur != nil && cur.Name == ".." && len(p.Selected) == 0 {
			p.Path = strings.TrimSpace(p.RenameValue)
			p.Renaming = false
			_ = p.Load(m.source)
			m.status = "Path: " + p.Path
			return
		}
		if cur != nil {
			if err := operations.Rename(m.source, *cur, strings.TrimSpace(p.RenameValue)); err != nil {
				m.status = err.Error()
			} else {
				m.status = "Renamed"
				_ = p.Load(m.source)
			}
		}
		p.Renaming = false
	case "Backspace":
		if len(p.RenameValue) > 0 {
			p.RenameValue = p.RenameValue[:len(p.RenameValue)-1]
		}
	default:
		if msg.Type == tea.KeyRunes && len(msg.Runes) == 1 {
			p.RenameValue += string(msg.Runes)
		}
	}
}

func (m *Model) handlePromptKey(msg tea.KeyMsg, key string) {
	switch key {
	case "Esc":
		m.prompt = ""
		m.promptFn = ""
	case "Enter":
		value := strings.TrimSpace(m.prompt)
		m.prompt = ""
		switch m.promptFn {
		case "mkdir":
			if _, err := operations.Mkdir(m.source, m.activePanel().Path, value); err != nil {
				m.status = err.Error()
			} else {
				m.status = "Created " + value
				_ = m.activePanel().Load(m.source)
			}
		case "move":
			entries := m.activePanel().SelectedEntries()
			if len(entries) == 1 {
				target := value
				if !strings.Contains(target, string(os.PathSeparator)) {
					target = filepath.Join(m.inactivePanel().Path, value)
				}
				if err := m.source.Rename(entries[0].Path, target); err != nil {
					m.status = err.Error()
				} else {
					m.status = "Moved"
					_ = m.ReloadBoth()
				}
			}
		case "filter":
			m.activePanel().QuickFilter = value
			m.activePanel().QuickSearch = value
			m.status = "Filter: " + value
		}
		m.promptFn = ""
	default:
		if msg.Type == tea.KeyRunes && len(msg.Runes) == 1 {
			m.prompt += string(msg.Runes)
		}
	}
}

func (m *Model) handleDriveKey(key string) {
	drives := platform.Drives()
	if len(drives) == 0 {
		m.driveMode = false
		return
	}
	switch key {
	case "Esc":
		m.driveMode = false
	case "Up":
		m.driveIdx--
		if m.driveIdx < 0 {
			m.driveIdx = len(drives) - 1
		}
	case "Down":
		m.driveIdx++
		if m.driveIdx >= len(drives) {
			m.driveIdx = 0
		}
	case "Enter":
		p := m.leftNB.Current()
		if m.driveFor == panel.Right {
			p = m.rightNB.Current()
		}
		p.Path = drives[m.driveIdx]
		_ = p.Load(m.source)
		m.driveMode = false
		m.status = "Drive: " + p.Path
	}
}

// View implements tea.Model.
func (m Model) View() string {
	if m.quitting {
		return ""
	}
	viewer := m.viewer
	if viewer != "" {
		viewer = mainview.RenderViewer("", viewer, 4)
	}
	qv := m.quickView
	if m.driveMode {
		drives := platform.Drives()
		if m.driveIdx < len(drives) {
			m.status = "Drive: " + drives[m.driveIdx]
		}
	}
	cmd := m.cmdLine
	if m.cmdFocus {
		cmd = m.cmdLine + "_"
	} else if m.prompt != "" {
		cmd = m.prompt
	}
	return mainview.Render(m.width, m.height, m.active, m.leftNB.Current(), m.rightNB.Current(), m.status, cmd, viewer, qv,
		m.leftNB.TabTitles(), m.rightNB.TabTitles(), m.renderDialog())
}

func (m *Model) execute(command string, params []string) {
	ctx := m.commandContext()
	m.registry.Execute(ctx, command, params)
}

func (m *Model) commandContext() *commands.Context {
	active := m.leftNB.Current()
	inactive := m.rightNB.Current()
	if m.active == panel.Right {
		active = m.rightNB.Current()
		inactive = m.leftNB.Current()
	}
	return &commands.Context{
		ActivePanel:   active,
		InactivePanel: inactive,
		SourcePanel:   active,
		TargetPanel:   inactive,
		SetStatus:     func(s string) { m.status = s },
		Quit:          func() { m.quitting = true },
		SwitchPanel:   m.switchPanel,
	}
}

func (m *Model) switchPanel() {
	if m.active == panel.Left {
		m.active = panel.Right
		m.status = "Active panel: right"
		return
	}
	m.active = panel.Left
	m.status = "Active panel: left"
}

func (m *Model) activePanel() *panel.Panel {
	if m.active == panel.Right {
		return m.rightNB.Current()
	}
	return m.leftNB.Current()
}

func (m *Model) inactivePanel() *panel.Panel {
	if m.active == panel.Right {
		return m.leftNB.Current()
	}
	return m.rightNB.Current()
}

// AppServices + Phase2App

func (m *Model) ActivePanel() *panel.Panel     { return m.activePanel() }
func (m *Model) InactivePanel() *panel.Panel   { return m.inactivePanel() }
func (m *Model) Source() filesrc.Source        { return m.source }
func (m *Model) Notebook(side panel.Side) *panel.Notebook {
	if side == panel.Right {
		return m.rightNB
	}
	return m.leftNB
}
func (m *Model) ActiveNotebook() *panel.Notebook { return m.Notebook(m.active) }

func (m *Model) ReloadActive() error  { return m.activePanel().Load(m.source) }
func (m *Model) ReloadBoth() error {
	if err := m.leftNB.Current().Load(m.source); err != nil {
		return err
	}
	return m.rightNB.Current().Load(m.source)
}
func (m *Model) ReloadNotebook(nb *panel.Notebook) error {
	return nb.Current().Load(m.source)
}

func (m *Model) SwapPanels() {
	m.leftNB, m.rightNB = m.rightNB, m.leftNB
	m.leftNB.Side = panel.Left
	m.rightNB.Side = panel.Right
}

func (m *Model) OpenOtherPanelPath() {
	other := m.inactivePanel()
	other.Path = m.activePanel().Path
	_ = other.Load(m.source)
}

func (m *Model) SetPrompt(text string) {
	if strings.HasPrefix(text, "New directory") {
		m.promptFn = "mkdir"
		m.prompt = ""
		return
	}
	if strings.HasPrefix(text, "Move to") {
		m.promptFn = "move"
		m.prompt = ""
		return
	}
	m.promptFn = "filter"
	m.prompt = ""
}

func (m *Model) SetCmdLineFocus(on bool) { m.cmdFocus = on }
func (m *Model) AppendCmdLine(s string)  { m.cmdLine += s }

func (m *Model) ExecuteCmdLine() error {
	cmd := strings.TrimSpace(m.cmdLine)
	if cmd == "" {
		return nil
	}
	if strings.HasPrefix(cmd, "cd ") {
		m.activePanel().Path = strings.TrimSpace(strings.TrimPrefix(cmd, "cd "))
		return m.activePanel().Load(m.source)
	}
	return builtin.RunShell(cmd, m.activePanel().Path)
}

func (m *Model) CmdLineHistoryPrev() {
	if len(m.cmdHist) == 0 {
		return
	}
	if m.cmdIdx > 0 {
		m.cmdIdx--
	}
	m.cmdLine = m.cmdHist[m.cmdIdx]
	m.cmdFocus = true
}

func (m *Model) SetQuickView(path string) {
	data, err := os.ReadFile(path)
	if err != nil {
		m.status = err.Error()
		return
	}
	content := string(data)
	if len(content) > 800 {
		content = content[:800] + "..."
	}
	m.quickView = mainview.RenderViewer(path, content, 4)
}

func (m *Model) ClearQuickView() { m.quickView = "" }

func (m *Model) SetViewer(path string) {
	data, err := os.ReadFile(path)
	if err != nil {
		m.status = err.Error()
		return
	}
	content := string(data)
	if len(content) > 4000 {
		content = content[:4000] + "\n..."
	}
	m.viewer = mainview.RenderViewer(path, content, 6)
	m.status = "Viewing " + filepath.Base(path)
}

func (m *Model) DriveMenu(side panel.Side) {
	m.driveMode = true
	m.driveFor = side
	m.driveIdx = 0
	m.status = "Select drive"
}

func (m Model) ActivePanelSide() panel.Side { return m.active }
func (m Model) Status() string               { return m.status }
func (m Model) Quitting() bool               { return m.quitting }
