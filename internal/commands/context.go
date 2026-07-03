package commands

import "github.com/openqt/tcmd/internal/panel"

// Context carries runtime state for command handlers.
type Context struct {
	ActivePanel   *panel.Panel
	InactivePanel *panel.Panel
	SourcePanel   *panel.Panel
	TargetPanel   *panel.Panel
	SetStatus     func(string)
	Quit          func()
	SwitchPanel   func()
}
