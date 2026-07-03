package builtin

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/openqt/tcmd/internal/commands"
	"github.com/openqt/tcmd/internal/filesrc"
	"github.com/openqt/tcmd/internal/operations"
	"github.com/openqt/tcmd/internal/panel"
)

// AppServices exposes callbacks needed by file commands.
type AppServices interface {
	ActivePanel() *panel.Panel
	InactivePanel() *panel.Panel
	Source() filesrc.Source
	ReloadActive() error
	ReloadBoth() error
	SwapPanels()
	OpenOtherPanelPath()
	SetPrompt(string)
	SetViewer(string)
	DriveMenu(panel.Side)
	EnterArchive(path string) error
}

// RegisterPhase1 registers core file management commands.
func RegisterPhase1(r *commands.Registry, app AppServices) {
	register := func(name, category string, fn commands.Handler) {
		r.Register(commands.Def{Name: name, Category: category, Handler: fn})
	}

	register("cm_Open", "Navigation", func(ctx *commands.Context, _ []string) commands.Result {
		p := ctx.ActivePanel
		if p == nil {
			return commands.ResultDisabled
		}
		cur := p.Current()
		if cur == nil {
			return commands.ResultDisabled
		}
		if cur.IsDir {
			if filesrc.IsArchiveFile(cur.Path) {
				if err := app.EnterArchive(cur.Path); err != nil {
					ctx.SetStatus(err.Error())
					return commands.ResultDisabled
				}
				ctx.SetStatus("Archive")
				return commands.ResultSuccess
			}
			if err := p.EnterDirectory(app.Source()); err != nil {
				ctx.SetStatus(err.Error())
				return commands.ResultDisabled
			}
			ctx.SetStatus("Opened " + p.Path)
			return commands.ResultSuccess
		}
		app.SetViewer(cur.Path)
		return commands.ResultSuccess
	})

	register("cm_ChangeDirToParent", "Navigation", func(ctx *commands.Context, _ []string) commands.Result {
		p := ctx.ActivePanel
		if p == nil {
			return commands.ResultDisabled
		}
		if err := p.ParentDirectory(app.Source()); err != nil {
			ctx.SetStatus(err.Error())
			return commands.ResultDisabled
		}
		ctx.SetStatus(p.Path)
		return commands.ResultSuccess
	})

	register("cm_ChangeDirToRoot", "Navigation", func(ctx *commands.Context, _ []string) commands.Result {
		p := ctx.ActivePanel
		if p == nil {
			return commands.ResultDisabled
		}
		p.Path = "/"
		p.Cursor = 0
		p.ClearSelection()
		_ = p.Load(app.Source())
		ctx.SetStatus(p.Path)
		return commands.ResultSuccess
	})

	register("cm_ChangeDirToHome", "Navigation", func(ctx *commands.Context, _ []string) commands.Result {
		p := ctx.ActivePanel
		home, err := os.UserHomeDir()
		if p == nil || err != nil {
			return commands.ResultDisabled
		}
		p.Path = home
		p.Cursor = 0
		p.ClearSelection()
		_ = p.Load(app.Source())
		ctx.SetStatus(p.Path)
		return commands.ResultSuccess
	})

	register("cm_Exchange", "Window", func(ctx *commands.Context, _ []string) commands.Result {
		app.SwapPanels()
		ctx.SetStatus("Panels swapped")
		return commands.ResultSuccess
	})

	register("cm_TransferPath", "Navigation", func(ctx *commands.Context, _ []string) commands.Result {
		app.OpenOtherPanelPath()
		ctx.SetStatus("Target = Source")
		return commands.ResultSuccess
	})

	register("cm_LeftOpenDrives", "Active Panel", func(ctx *commands.Context, _ []string) commands.Result {
		app.DriveMenu(panel.Left)
		return commands.ResultSuccess
	})
	register("cm_RightOpenDrives", "Active Panel", func(ctx *commands.Context, _ []string) commands.Result {
		app.DriveMenu(panel.Right)
		return commands.ResultSuccess
	})

	register("cm_ShowSysFiles", "View", func(ctx *commands.Context, _ []string) commands.Result {
		p := ctx.ActivePanel
		if p == nil {
			return commands.ResultDisabled
		}
		p.ShowHidden = !p.ShowHidden
		_ = p.Load(app.Source())
		ctx.SetStatus(fmt.Sprintf("Hidden files: %v", p.ShowHidden))
		return commands.ResultSuccess
	})

	register("cm_SortByName", "Active Panel", sortHandler(app, panel.SortByName))
	register("cm_SortByExt", "Active Panel", sortHandler(app, panel.SortByExt))
	register("cm_SortBySize", "Active Panel", sortHandler(app, panel.SortBySize))
	register("cm_SortByDate", "Active Panel", sortHandler(app, panel.SortByDate))

	register("cm_Refresh", "Active Panel", func(ctx *commands.Context, _ []string) commands.Result {
		if err := app.ReloadActive(); err != nil {
			ctx.SetStatus(err.Error())
			return commands.ResultDisabled
		}
		ctx.SetStatus("Refreshed")
		return commands.ResultSuccess
	})

	register("cm_Select", "Mark", func(ctx *commands.Context, _ []string) commands.Result {
		ctx.ActivePanel.ToggleSelect()
		return commands.ResultSuccess
	})
	register("cm_SelectAll", "Mark", func(ctx *commands.Context, _ []string) commands.Result {
		ctx.ActivePanel.SelectAll()
		return commands.ResultSuccess
	})
	register("cm_UnselectAll", "Mark", func(ctx *commands.Context, _ []string) commands.Result {
		ctx.ActivePanel.ClearSelection()
		return commands.ResultSuccess
	})
	register("cm_InvertSelection", "Mark", func(ctx *commands.Context, _ []string) commands.Result {
		ctx.ActivePanel.InvertSelection()
		return commands.ResultSuccess
	})

	register("cm_CalculateSpace", "File Operations", func(ctx *commands.Context, _ []string) commands.Result {
		cur := ctx.ActivePanel.Current()
		if cur == nil || !cur.IsDir {
			return commands.ResultDisabled
		}
		size, err := operations.DirSize(app.Source(), cur.Path)
		if err != nil {
			ctx.SetStatus(err.Error())
			return commands.ResultDisabled
		}
		ctx.ActivePanel.SetDirSizeHint(operations.FormatDirSize(size))
		ctx.SetStatus(fmt.Sprintf("%s: %s", cur.Name, operations.FormatDirSize(size)))
		return commands.ResultSuccess
	})

	register("cm_Copy", "File Operations", func(ctx *commands.Context, _ []string) commands.Result {
		entries := ctx.SourcePanel.SelectedEntries()
		if len(entries) == 0 {
			return commands.ResultDisabled
		}
		if err := operations.CopyEntries(app.Source(), entries, ctx.TargetPanel.Path); err != nil {
			ctx.SetStatus(err.Error())
			return commands.ResultDisabled
		}
		_ = app.ReloadBoth()
		ctx.SetStatus(fmt.Sprintf("Copied %d item(s)", len(entries)))
		return commands.ResultSuccess
	})

	register("cm_Rename", "File Operations", func(ctx *commands.Context, _ []string) commands.Result {
		entries := ctx.SourcePanel.SelectedEntries()
		if len(entries) != 1 {
			ctx.SetStatus("Select one item to move/rename")
			return commands.ResultDisabled
		}
		app.SetPrompt("Move to (directory or new name): " + filepath.Base(entries[0].Path))
		return commands.ResultSuccess
	})

	register("cm_RenameOnly", "File Operations", func(ctx *commands.Context, _ []string) commands.Result {
		cur := ctx.ActivePanel.Current()
		if cur == nil || cur.Name == ".." {
			if len(ctx.ActivePanel.Selected) == 0 {
				ctx.ActivePanel.Renaming = true
				ctx.ActivePanel.RenameValue = ctx.ActivePanel.Path
				ctx.SetStatus("Edit path")
				return commands.ResultSuccess
			}
			return commands.ResultDisabled
		}
		ctx.ActivePanel.Renaming = true
		ctx.ActivePanel.RenameValue = cur.Name
		ctx.SetStatus("Rename: " + cur.Name)
		return commands.ResultSuccess
	})

	register("cm_MakeDir", "File Operations", func(ctx *commands.Context, _ []string) commands.Result {
		app.SetPrompt("New directory name: ")
		return commands.ResultSuccess
	})

	register("cm_Delete", "File Operations", func(ctx *commands.Context, params []string) commands.Result {
		permanent := hasParam(params, "trashcan=reversesetting")
		entries := ctx.SourcePanel.SelectedEntries()
		if len(entries) == 0 {
			return commands.ResultDisabled
		}
		if err := operations.DeleteEntries(app.Source(), entries); err != nil {
			ctx.SetStatus(err.Error())
			return commands.ResultDisabled
		}
		ctx.SourcePanel.ClearSelection()
		_ = app.ReloadActive()
		mode := "recycle"
		if permanent {
			mode = "permanent"
		}
		ctx.SetStatus(fmt.Sprintf("Deleted %d item(s) [%s]", len(entries), mode))
		return commands.ResultSuccess
	})

	register("cm_View", "File Operations", func(ctx *commands.Context, _ []string) commands.Result {
		cur := ctx.ActivePanel.Current()
		if cur == nil || cur.IsDir {
			return commands.ResultDisabled
		}
		app.SetViewer(cur.Path)
		return commands.ResultSuccess
	})

	register("cm_Edit", "File Operations", func(ctx *commands.Context, _ []string) commands.Result {
		cur := ctx.ActivePanel.Current()
		if cur == nil || cur.IsDir {
			return commands.ResultDisabled
		}
		editor := os.Getenv("EDITOR")
		if editor == "" {
			editor = "vi"
		}
		cmd := exec.Command(editor, cur.Path)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			ctx.SetStatus(err.Error())
			return commands.ResultDisabled
		}
		_ = app.ReloadActive()
		ctx.SetStatus("Edited " + cur.Name)
		return commands.ResultSuccess
	})

	register("cm_MoveUp", "Navigation", func(ctx *commands.Context, _ []string) commands.Result {
		ctx.ActivePanel.MoveCursor(-1)
		return commands.ResultSuccess
	})
	register("cm_MoveDown", "Navigation", func(ctx *commands.Context, _ []string) commands.Result {
		ctx.ActivePanel.MoveCursor(1)
		return commands.ResultSuccess
	})

	register("cm_BriefView", "Active Panel", func(ctx *commands.Context, _ []string) commands.Result {
		ctx.SetStatus("Brief view")
		return commands.ResultSuccess
	})
	register("cm_ColumnsView", "Active Panel", func(ctx *commands.Context, _ []string) commands.Result {
		ctx.SetStatus("Columns view")
		return commands.ResultSuccess
	})
}

func sortHandler(app AppServices, mode panel.SortMode) commands.Handler {
	return func(ctx *commands.Context, _ []string) commands.Result {
		p := ctx.ActivePanel
		if p == nil {
			return commands.ResultDisabled
		}
		if p.SortMode == mode {
			p.Reverse = !p.Reverse
		} else {
			p.SortMode = mode
			p.Reverse = false
		}
		p.Sort()
		ctx.SetStatus(fmt.Sprintf("Sort: %v reverse=%v", mode, p.Reverse))
		return commands.ResultSuccess
	}
}

func hasParam(params []string, want string) bool {
	for _, p := range params {
		if strings.EqualFold(strings.TrimSpace(p), want) {
			return true
		}
	}
	return false
}
