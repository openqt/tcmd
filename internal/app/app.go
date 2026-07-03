package app

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/openqt/tcmd/internal/commands"
	"github.com/openqt/tcmd/internal/commands/builtin"
	"github.com/openqt/tcmd/internal/hotkeys"
	"github.com/openqt/tcmd/internal/panel"
	"github.com/openqt/tcmd/internal/ui/mainview"
)

// Model is the root Bubble Tea model.
type Model struct {
	width   int
	height  int
	active  panel.Side
	left    *panel.Panel
	right   *panel.Panel
	status  string
	quitting bool

	registry   *commands.Registry
	dispatcher *hotkeys.Dispatcher
}

// NewModel constructs the application model.
func NewModel(leftPath, rightPath string) *Model {
	registry := commands.NewRegistry()
	builtin.RegisterPhase0(registry)

	return &Model{
		active:     panel.Left,
		left:       panel.New(panel.Left, leftPath),
		right:      panel.New(panel.Right, rightPath),
		status:     "dc-tui Phase 0",
		registry:   registry,
		dispatcher: hotkeys.NewDispatcher(nil),
	}
}

// Init implements tea.Model.
func (m Model) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		if match, ok := m.dispatcher.LookupMsg(hotkeys.ContextMain, msg); ok {
			m.execute(match.Command, match.Params)
		}
	}
	return m, nil
}

// View implements tea.Model.
func (m Model) View() string {
	if m.quitting {
		return ""
	}
	return mainview.Render(m.width, m.height, m.active, m.left, m.right, m.status)
}

func (m *Model) execute(command string, params []string) {
	ctx := m.commandContext()
	switch m.registry.Execute(ctx, command, params) {
	case commands.ResultNotFound:
		// status already set by registry
	case commands.ResultDisabled:
		m.status = command + ": disabled"
	case commands.ResultSuccess:
		if m.status == "" || m.status == "dc-tui Phase 0" {
			m.status = command
		}
	}
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

// ActivePanel returns the focused panel side (testing helper).
func (m Model) ActivePanel() panel.Side {
	return m.active
}

// Status returns the status line (testing helper).
func (m Model) Status() string {
	return m.status
}

// Quitting reports whether the app requested exit (testing helper).
func (m Model) Quitting() bool {
	return m.quitting
}
