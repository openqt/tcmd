package hotkeys

// Binding maps a normalized key chord to a DC internal command.
type Binding struct {
	Key     string
	Command string
	Params  []string
}

// Match describes a resolved hotkey hit.
type Match struct {
	Command string
	Params  []string
}
