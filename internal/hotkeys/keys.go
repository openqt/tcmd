package hotkeys

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// NormalizeKey converts a Bubble Tea key message to a DC-style chord string.
func NormalizeKey(msg tea.KeyMsg) string {
	if msg.String() == "" {
		return ""
	}

	key := msg.String()
	parts := strings.Split(key, "+")
	modifiers := make([]string, 0, 3)
	main := parts[len(parts)-1]

	for _, part := range parts[:len(parts)-1] {
		switch strings.ToLower(part) {
		case "ctrl", "control":
			modifiers = append(modifiers, "Ctrl")
		case "alt", "opt", "option":
			modifiers = append(modifiers, "Alt")
		case "shift":
			modifiers = append(modifiers, "Shift")
		}
	}

	switch strings.ToLower(main) {
	case "esc":
		main = "Esc"
	case "enter":
		main = "Enter"
	case "tab":
		main = "Tab"
	case "backspace":
		main = "Backspace"
	case "delete":
		main = "Del"
	case "up", "down", "left", "right":
		main = strings.ToUpper(main[:1]) + main[1:]
	case "pgup":
		main = "PgUp"
	case "pgdown":
		main = "PgDn"
	default:
		if len(main) == 1 {
			if main[0] >= 'a' && main[0] <= 'z' {
				main = strings.ToUpper(main)
			}
		}
	}

	if strings.HasPrefix(strings.ToLower(main), "f") && len(main) <= 3 {
		main = strings.ToUpper(main)
	}

	if len(modifiers) == 0 {
		return main
	}
	return strings.Join(modifiers, "+") + "+" + main
}
