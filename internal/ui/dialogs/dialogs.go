package dialogs

import "github.com/charmbracelet/lipgloss"

var boxStyle = lipgloss.NewStyle().
	Border(lipgloss.RoundedBorder()).
	BorderForeground(lipgloss.Color("81")).
	Padding(1, 2)

// Find renders the find files dialog.
func Find(nameMask, content, results string, page int) string {
	body := "Find Files (Alt+F7)\n"
	body += "Page: Standard | Advanced | Results\n"
	body += "Name mask: " + nameMask + "_\n"
	body += "Content:   " + content + "\n"
	body += "Results:\n" + results + "\n"
	body += "F9: Start | Esc: Close"
	return boxStyle.Render(body)
}

// MultiRename renders batch rename dialog.
func MultiRename(rule, preview string) string {
	body := "Multi-Rename Tool (Ctrl+M)\n"
	body += "Rule: " + rule + "_\n"
	body += preview + "\nF9: Rename | Esc: Close"
	return boxStyle.Render(body)
}

// Sync renders directory sync dialog.
func Sync(left, right, listing string) string {
	body := "Synchronize Directories\n"
	body += "Left:  " + left + "\nRight: " + right + "\n"
	body += listing + "\nF3: view left | Shift+F3: view right | Esc: close"
	return boxStyle.Render(body)
}

// Hotlist renders directory hotlist.
func Hotlist(paths []string, cursor int) string {
	body := "Directory Hotlist (Ctrl+D)\n"
	for i, p := range paths {
		prefix := "  "
		if i == cursor {
			prefix = "> "
		}
		body += prefix + p + "\n"
	}
	body += "Enter: go | Esc: close"
	return boxStyle.Render(body)
}

// Properties renders file properties.
func Properties(name, path, size, mode string) string {
	body := "Properties (Alt+Enter)\n"
	body += "Name: " + name + "\nPath: " + path + "\nSize: " + size + "\nMode: " + mode + "\nEsc: close"
	return boxStyle.Render(body)
}

// CopyMove renders copy/move confirmation dialog.
func CopyMove(op, target, field string) string {
	body := "Copy/Move Dialog\n"
	body += "Operation: " + op + "\nTarget: " + target + "\nField: " + field + "\nF2: queue | F5/F6: cycle field | Esc: cancel"
	return boxStyle.Render(body)
}
