package mainview

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/openqt/tcmd/internal/panel"
)

var (
	fnKeyStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("81")).Bold(true)
	driveStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))
	tabStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("252"))
	activeBorder = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("81")).
			Padding(0, 1)
	inactiveBorder = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("238")).
			Padding(0, 1)
	statusStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))
	cmdLineStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("252"))
	viewerStyle  = lipgloss.NewStyle().Border(lipgloss.NormalBorder()).Padding(0, 1)
)

// Render draws the main window.
func Render(width, height int, active panel.Side, left, right *panel.Panel, status, prompt, viewer string) string {
	if width < 20 {
		width = 20
	}
	if height < 8 {
		height = 8
	}

	fnBar := fnKeyStyle.Render("F1 Help  F2 Rename  F3 View  F4 Edit  F5 Copy  F6 Move  F7 Mkdir  F8 Delete  F9 Term  F10 Menu")
	driveBar := driveStyle.Render("Drives: Alt+F1 left | Alt+F2 right")
	leftTab := tabStyle.Render("L: " + shorten(left.Title(), 24))
	rightTab := tabStyle.Render("R: " + shorten(right.Title(), 24))
	tabBar := lipgloss.JoinHorizontal(lipgloss.Top, leftTab, "  |  ", rightTab)

	chromeHeight := 6
	if viewer != "" {
		chromeHeight = 10
	}
	panelHeight := height - chromeHeight
	if panelHeight < 3 {
		panelHeight = 3
	}
	panelWidth := (width - 3) / 2
	if panelWidth < 10 {
		panelWidth = 10
	}

	leftBody := renderPanelBody(left, panelWidth-4, panelHeight-2, active == panel.Left)
	rightBody := renderPanelBody(right, panelWidth-4, panelHeight-2, active == panel.Right)

	leftBox := inactiveBorder.Width(panelWidth).Height(panelHeight).Render(leftBody)
	if active == panel.Left {
		leftBox = activeBorder.Width(panelWidth).Height(panelHeight).Render(leftBody)
	}
	rightBox := inactiveBorder.Width(panelWidth).Height(panelHeight).Render(rightBody)
	if active == panel.Right {
		rightBox = activeBorder.Width(panelWidth).Height(panelHeight).Render(rightBody)
	}

	panels := lipgloss.JoinHorizontal(lipgloss.Top, leftBox, " ", rightBox)

	cmdLine := cmdLineStyle.Render("> " + prompt)
	if status == "" {
		status = "Tab: switch panel | Alt+F4: quit"
	}
	statusLine := statusStyle.Render(status)

	lines := []string{fnBar, driveBar, tabBar, panels, cmdLine, statusLine}
	if viewer != "" {
		lines = append(lines[:4], viewerStyle.Width(width-2).Render(shorten(viewer, width*2)), lines[4], lines[5])
	}
	return strings.Join(lines, "\n")
}

func renderPanelBody(p *panel.Panel, width, rows int, _ bool) string {
	if rows < 1 {
		rows = 1
	}
	if p == nil || len(p.Entries) == 0 {
		return strings.Repeat("\n", rows-1) + "  (empty)"
	}
	start := 0
	if p.Cursor >= rows {
		start = p.Cursor - rows + 1
	}
	lines := make([]string, 0, rows)
	for i := start; i < len(p.Entries) && len(lines) < rows; i++ {
		if p.Renaming && i == p.Cursor {
			lines = append(lines, "> * "+p.RenameValue+"_")
			continue
		}
		lines = append(lines, p.FormatEntry(i, width))
	}
	for len(lines) < rows {
		lines = append(lines, "")
	}
	return strings.Join(lines, "\n")
}

func shorten(s string, max int) string {
	if len(s) <= max {
		return s
	}
	if max <= 3 {
		return s[:max]
	}
	return "..." + s[len(s)-(max-3):]
}

// RenderViewer formats viewer overlay text.
func RenderViewer(path, content string, maxLines int) string {
	lines := strings.Split(content, "\n")
	if len(lines) > maxLines {
		lines = lines[:maxLines]
	}
	return fmt.Sprintf("Viewer: %s\n%s", path, strings.Join(lines, "\n"))
}
