package builtin

import "github.com/openqt/tcmd/internal/commands"

// RegisterPhase0 wires commands required for the foundation milestone.
func RegisterPhase0(r *commands.Registry) {
	r.Register(commands.Def{
		Name:     "cm_Exit",
		Category: "Window",
		Handler: func(ctx *commands.Context, _ []string) commands.Result {
			if ctx != nil && ctx.Quit != nil {
				ctx.Quit()
			}
			return commands.ResultSuccess
		},
	})
	r.Register(commands.Def{
		Name:     "cm_Close",
		Category: "Window",
		Handler: func(ctx *commands.Context, _ []string) commands.Result {
			if ctx != nil && ctx.Quit != nil {
				ctx.Quit()
			}
			return commands.ResultSuccess
		},
	})
	r.Register(commands.Def{
		Name:     "cm_SwitchPnl",
		Category: "Window",
		Handler: func(ctx *commands.Context, _ []string) commands.Result {
			if ctx != nil && ctx.SwitchPanel != nil {
				ctx.SwitchPanel()
				return commands.ResultSuccess
			}
			return commands.ResultDisabled
		},
	})
	r.Register(commands.Def{
		Name:     "cm_FocusCmdLine",
		Category: "Command Line",
		Handler: func(ctx *commands.Context, _ []string) commands.Result {
			if ctx != nil && ctx.SetStatus != nil {
				ctx.SetStatus("Command line focus (Phase 2)")
			}
			return commands.ResultSuccess
		},
	})
	r.Register(commands.Def{
		Name:     "cm_HelpIndex",
		Category: "Help",
		Handler: func(ctx *commands.Context, _ []string) commands.Result {
			if ctx != nil && ctx.SetStatus != nil {
				ctx.SetStatus("Help: see docs/SHORTCUTS.md")
			}
			return commands.ResultSuccess
		},
	})
}
