package mainview

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/openqt/tcmd/internal/panel"
)

var (
	fnKeyStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("81")).Bold(true)
	driveStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))
	tabStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("252"))
	activeBorder = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("81")).
			Padding(0, 1)
	inactiveBorder = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("238")).
			Padding(0, 1)
	statusStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))
	cmdLineStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("252"))
)

// Render draws the main window shell.
func Render(width, height int, active panel.Side, left, right *panel.Panel, status string) string {
	if width < 20 {
		width = 20
	}
	if height < 8 {
		height = 8
	}

	fnBar := fnKeyStyle.Render("F1 Help  F2 Rename  F3 View  F4 Edit  F5 Copy  F6 Move  F7 Mkdir  F8 Delete  F9 Term  F10 Menu")
	driveBar := driveStyle.Render("Drives: [/] [home] [tmp] ...")
	leftTab := tabStyle.Render("L: " + left.Title())
	rightTab := tabStyle.Render("R: " + right.Title())
	tabBar := lipgloss.JoinHorizontal(lipgloss.Top, leftTab, "  |  ", rightTab)

	chromeHeight := 6
	panelHeight := height - chromeHeight
	if panelHeight < 3 {
		panelHeight = 3
	}
	panelWidth := (width - 3) / 2
	if panelWidth < 10 {
		panelWidth = 10
	}

	leftBody := renderPanelBody(left, panelHeight-2)
	rightBody := renderPanelBody(right, panelHeight-2)

	leftBox := activeBorder.Width(panelWidth).Height(panelHeight).Render(leftBody)
	if active != panel.Left {
		leftBox = inactiveBorder.Width(panelWidth).Height(panelHeight).Render(leftBody)
	}

	rightBox := activeBorder.Width(panelWidth).Height(panelHeight).Render(rightBody)
	if active != panel.Right {
		rightBox = inactiveBorder.Width(panelWidth).Height(panelHeight).Render(rightBody)
	}

	panels := lipgloss.JoinHorizontal(lipgloss.Top, leftBox, " ", rightBox)
	cmdLine := cmdLineStyle.Render("> ")
	if status == "" {
		status = "Tab: switch panel | Alt+F4: quit"
	}
	statusLine := statusStyle.Render(status)

	return strings.Join([]string{fnBar, driveBar, tabBar, panels, cmdLine, statusLine}, "\n")
}

func renderPanelBody(p *panel.Panel, rows int) string {
	if rows < 1 {
		rows = 1
	}
	lines := []string{
		"> ..",
		"  (empty panel)",
		fmt.Sprintf("  path: %s", p.Title()),
	}
	if len(lines) > rows {
		lines = lines[:rows]
	}
	for len(lines) < rows {
		lines = append(lines, "")
	}
	return strings.Join(lines, "\n")
}
