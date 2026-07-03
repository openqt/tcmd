package hotkeys

// DefaultMainBindings returns Phase 0 main-window shortcuts from DC defaults.
func DefaultMainBindings() []Binding {
	return []Binding{
		{Key: "F1", Command: "cm_HelpIndex"},
		{Key: "Tab", Command: "cm_SwitchPnl"},
		{Key: "Ctrl+L", Command: "cm_FocusCmdLine"},
		{Key: "Alt+F4", Command: "cm_Exit"},
		{Key: "Alt+X", Command: "cm_Exit"},
	}
}
