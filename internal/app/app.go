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
	left     *panel.Panel
	right    *panel.Panel
	status   string
	prompt   string
	promptFn string
	viewer   string
	quitting bool
	driveIdx int
	driveFor panel.Side
	driveMode bool

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
		left:       panel.New(panel.Left, leftPath),
		right:      panel.New(panel.Right, rightPath),
		status:     "dc-tui Phase 1",
		source:     source,
		registry:   registry,
		dispatcher: hotkeys.NewDispatcher(nil),
	}
	builtin.RegisterPhase0(registry)
	builtin.RegisterPhase1(registry, m)
	_ = m.left.Load(source)
	_ = m.right.Load(source)
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

	if match, ok := m.dispatcher.LookupMsg(hotkeys.ContextMain, msg); ok {
		m.execute(match.Command, match.Params)
		return true
	}
	return false
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
		p := m.left
		if m.driveFor == panel.Right {
			p = m.right
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
	if m.driveMode {
		drives := platform.Drives()
		if m.driveIdx < len(drives) {
			m.status = "Drive: " + drives[m.driveIdx]
		}
	}
	return mainview.Render(m.width, m.height, m.active, m.left, m.right, m.status, m.prompt, viewer)
}

func (m *Model) execute(command string, params []string) {
	ctx := m.commandContext()
	m.registry.Execute(ctx, command, params)
}

func (m *Model) commandContext() *commands.Context {
	active := m.left
	inactive := m.right
	if m.active == panel.Right {
		active = m.right
		inactive = m.left
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
		return m.right
	}
	return m.left
}

func (m *Model) inactivePanel() *panel.Panel {
	if m.active == panel.Right {
		return m.left
	}
	return m.right
}

// AppServices implementation

func (m *Model) ActivePanel() *panel.Panel   { return m.activePanel() }
func (m *Model) InactivePanel() *panel.Panel { return m.inactivePanel() }
func (m *Model) Source() filesrc.Source      { return m.source }

func (m *Model) ReloadActive() error {
	return m.activePanel().Load(m.source)
}

func (m *Model) ReloadBoth() error {
	if err := m.left.Load(m.source); err != nil {
		return err
	}
	return m.right.Load(m.source)
}

func (m *Model) SwapPanels() {
	m.left, m.right = m.right, m.left
	m.left.Side = panel.Left
	m.right.Side = panel.Right
}

func (m *Model) OpenOtherPanelPath() {
	other := m.inactivePanel()
	other.Path = m.activePanel().Path
	_ = other.Load(m.source)
}

func (m *Model) SetPrompt(text string) {
	m.prompt = strings.TrimPrefix(text, strings.SplitN(text, ":", 2)[0]+": ")
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
	m.promptFn = ""
}

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

// ActivePanelSide returns focused panel side (testing).
func (m Model) ActivePanelSide() panel.Side { return m.active }
func (m Model) Status() string              { return m.status }
func (m Model) Quitting() bool              { return m.quitting }
