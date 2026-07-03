package builtin

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/openqt/tcmd/internal/commands"
	"github.com/openqt/tcmd/internal/panel"
	"github.com/openqt/tcmd/internal/platform"
)

// Phase2App extends AppServices for navigation features.
type Phase2App interface {
	AppServices
	Notebook(side panel.Side) *panel.Notebook
	ActiveNotebook() *panel.Notebook
	ReloadNotebook(nb *panel.Notebook) error
	SetCmdLineFocus(bool)
	AppendCmdLine(string)
	ExecuteCmdLine() error
	CmdLineHistoryPrev()
	SetQuickView(string)
	ClearQuickView()
}

// RegisterPhase2 registers navigation and tab commands.
func RegisterPhase2(r *commands.Registry, app Phase2App) {
	register := func(name, category string, fn commands.Handler) {
		r.Register(commands.Def{Name: name, Category: category, Handler: fn})
	}

	register("cm_NewTab", "Tabs", func(ctx *commands.Context, _ []string) commands.Result {
		nb := app.ActiveNotebook()
		nb.NewTab(ctx.ActivePanel.Path)
		_ = app.ReloadNotebook(nb)
		ctx.SetStatus("New tab")
		return commands.ResultSuccess
	})

	register("cm_CloseTab", "Tabs", func(ctx *commands.Context, _ []string) commands.Result {
		nb := app.ActiveNotebook()
		if !nb.CloseTab() {
			ctx.SetStatus("Cannot close last tab")
			return commands.ResultDisabled
		}
		_ = app.ReloadNotebook(nb)
		ctx.SetStatus("Tab closed")
		return commands.ResultSuccess
	})

	register("cm_NextTab", "Tabs", func(ctx *commands.Context, _ []string) commands.Result {
		nb := app.ActiveNotebook()
		nb.NextTab()
		ctx.SetStatus("Next tab")
		return commands.ResultSuccess
	})

	register("cm_PrevTab", "Tabs", func(ctx *commands.Context, _ []string) commands.Result {
		nb := app.ActiveNotebook()
		nb.PrevTab()
		ctx.SetStatus("Previous tab")
		return commands.ResultSuccess
	})

	register("cm_ActivateTabByIndex", "Tabs", func(ctx *commands.Context, params []string) commands.Result {
		idx := 0
		if len(params) > 0 {
			if params[0] == "index=0" || params[0] == "0" {
				idx = 9
			} else if len(params[0]) == 1 && params[0][0] >= '1' && params[0][0] <= '9' {
				idx = int(params[0][0] - '1')
			}
		}
		app.ActiveNotebook().ActivateTab(idx)
		ctx.SetStatus("Tab activated")
		return commands.ResultSuccess
	})

	register("cm_ViewHistoryPrev", "Navigation", func(ctx *commands.Context, _ []string) commands.Result {
		p := ctx.ActivePanel
		if path, ok := p.History.Back(); ok {
			p.Path = path
			_ = p.Load(app.Source())
			ctx.SetStatus(path)
		}
		return commands.ResultSuccess
	})

	register("cm_ViewHistoryNext", "Navigation", func(ctx *commands.Context, _ []string) commands.Result {
		p := ctx.ActivePanel
		if path, ok := p.History.Forward(); ok {
			p.Path = path
			_ = p.Load(app.Source())
			ctx.SetStatus(path)
		}
		return commands.ResultSuccess
	})

	register("cm_DirHistory", "Navigation", func(ctx *commands.Context, _ []string) commands.Result {
		entries := ctx.ActivePanel.History.Entries()
		if len(entries) == 0 {
			ctx.SetStatus("History empty")
			return commands.ResultSuccess
		}
		ctx.SetStatus("History: " + strings.Join(entries, " | "))
		return commands.ResultSuccess
	})

	register("cm_QuickFilter", "Active Panel", func(ctx *commands.Context, _ []string) commands.Result {
		app.SetCmdLineFocus(true)
		ctx.SetStatus("Quick filter: type and press Enter")
		return commands.ResultSuccess
	})

	register("cm_FlatView", "Active Panel", func(ctx *commands.Context, _ []string) commands.Result {
		p := ctx.ActivePanel
		if p.FlatView {
			_ = p.DisableFlat(app.Source())
			ctx.SetStatus("Flat view off")
			return commands.ResultSuccess
		}
		if err := p.LoadFlat(app.Source()); err != nil {
			ctx.SetStatus(err.Error())
			return commands.ResultDisabled
		}
		ctx.SetStatus("Flat view on")
		return commands.ResultSuccess
	})

	register("cm_FlatViewSel", "Active Panel", func(ctx *commands.Context, _ []string) commands.Result {
		ctx.SetStatus("Flat view selection (same as flat view in TUI)")
		return registerFlatSel(app, ctx)
	})

	register("cm_QuickView", "Active Panel", func(ctx *commands.Context, _ []string) commands.Result {
		cur := ctx.ActivePanel.Current()
		if cur == nil || cur.IsDir {
			app.ClearQuickView()
			return commands.ResultDisabled
		}
		app.SetQuickView(cur.Path)
		ctx.SetStatus("Quick view")
		return commands.ResultSuccess
	})

	register("cm_FocusCmdLine", "Command Line", func(ctx *commands.Context, _ []string) commands.Result {
		app.SetCmdLineFocus(true)
		ctx.SetStatus("Command line")
		return commands.ResultSuccess
	})

	register("cm_ShowCmdLineHistory", "Command Line", func(ctx *commands.Context, _ []string) commands.Result {
		app.CmdLineHistoryPrev()
		return commands.ResultSuccess
	})

	register("cm_CopyNamesToClip", "Clipboard", func(ctx *commands.Context, _ []string) commands.Result {
		return copyToClipboard(ctx, ctx.ActivePanel, false)
	})

	register("cm_CopyFullNamesToClip", "Clipboard", func(ctx *commands.Context, _ []string) commands.Result {
		return copyToClipboard(ctx, ctx.ActivePanel, true)
	})

	register("cm_CopyToClipboard", "Clipboard", func(ctx *commands.Context, _ []string) commands.Result {
		return copyToClipboard(ctx, ctx.ActivePanel, true)
	})

	register("cm_CutToClipboard", "Clipboard", func(ctx *commands.Context, _ []string) commands.Result {
		res := copyToClipboard(ctx, ctx.ActivePanel, true)
		if res == commands.ResultSuccess {
			ctx.SetStatus("Cut paths to clipboard (paste with Ctrl+V)")
		}
		return res
	})

	register("cm_PasteFromClipboard", "Clipboard", func(ctx *commands.Context, _ []string) commands.Result {
		text, err := platform.GetClipboard()
		if err != nil {
			ctx.SetStatus(err.Error())
			return commands.ResultDisabled
		}
		target := filepath.Join(ctx.ActivePanel.Path, filepath.Base(strings.TrimSpace(text)))
		if err := app.Source().Rename(strings.TrimSpace(text), target); err != nil {
			ctx.SetStatus("Clipboard paste: " + text)
			return commands.ResultSuccess
		}
		_ = app.ReloadActive()
		ctx.SetStatus("Pasted")
		return commands.ResultSuccess
	})

	register("cm_ExecuteFile", "Command Line", func(ctx *commands.Context, _ []string) commands.Result {
		if err := app.ExecuteCmdLine(); err != nil {
			ctx.SetStatus(err.Error())
			return commands.ResultDisabled
		}
		return commands.ResultSuccess
	})
}

func registerFlatSel(app Phase2App, ctx *commands.Context) commands.Result {
	p := ctx.ActivePanel
	if len(p.SelectedEntries()) == 0 {
		return commands.ResultDisabled
	}
	if err := p.LoadFlat(app.Source()); err != nil {
		ctx.SetStatus(err.Error())
		return commands.ResultDisabled
	}
	ctx.SetStatus("Flat view (selection)")
	return commands.ResultSuccess
}

func copyToClipboard(ctx *commands.Context, p *panel.Panel, full bool) commands.Result {
	entries := p.SelectedEntries()
	if len(entries) == 0 {
		return commands.ResultDisabled
	}
	lines := make([]string, len(entries))
	for i, e := range entries {
		if full {
			lines[i] = e.Path
		} else {
			lines[i] = e.Name
		}
	}
	if err := platform.SetClipboard(strings.Join(lines, "\n")); err != nil {
		ctx.SetStatus(err.Error())
		return commands.ResultDisabled
	}
	ctx.SetStatus("Copied to clipboard")
	return commands.ResultSuccess
}

// RunShell executes a shell command in the active directory.
func RunShell(command, dir string) error {
	cmd := exec.Command("sh", "-c", command)
	cmd.Dir = dir
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
