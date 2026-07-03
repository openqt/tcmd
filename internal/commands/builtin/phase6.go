package builtin

import (
	"fmt"
	"path/filepath"

	"github.com/openqt/tcmd/internal/commands"
	"github.com/openqt/tcmd/internal/config"
	"github.com/openqt/tcmd/internal/hotkeys"
	"github.com/openqt/tcmd/internal/operations"
	"github.com/openqt/tcmd/internal/tools"
)

// Phase6App for configuration and advanced features.
type Phase6App interface {
	Phase5App
	OpenConfig()
	OpenCommandBrowser()
	ImportShortcuts(path string) error
	ShowOperationLog()
}

// RegisterPhase6 registers advanced commands.
func RegisterPhase6(r *commands.Registry, app Phase6App) {
	register := func(name, category string, fn commands.Handler) {
		r.Register(commands.Def{Name: name, Category: category, Handler: fn})
	}

	register("cm_Options", "Configuration", func(ctx *commands.Context, _ []string) commands.Result {
		app.OpenConfig()
		return commands.ResultSuccess
	})
	register("cm_Config", "Configuration", func(ctx *commands.Context, _ []string) commands.Result {
		app.OpenConfig()
		return commands.ResultSuccess
	})
	register("cm_CommandBrowser", "Tools", func(ctx *commands.Context, _ []string) commands.Result {
		app.OpenCommandBrowser()
		return commands.ResultSuccess
	})
	register("cm_CheckSumCalc", "File Operations", func(ctx *commands.Context, _ []string) commands.Result {
		cur := ctx.ActivePanel.Current()
		if cur == nil || cur.IsDir {
			return commands.ResultDisabled
		}
		sum, err := tools.CalcChecksum(cur.Path, tools.AlgoMD5)
		if err != nil {
			ctx.SetStatus(err.Error())
			return commands.ResultDisabled
		}
		ctx.SetStatus("MD5: " + sum)
		return commands.ResultSuccess
	})
	register("cm_SplitFile", "File Operations", func(ctx *commands.Context, _ []string) commands.Result {
		cur := ctx.ActivePanel.Current()
		if cur == nil || cur.IsDir {
			return commands.ResultDisabled
		}
		parts, err := tools.SplitFile(cur.Path, 1024*1024)
		if err != nil {
			ctx.SetStatus(err.Error())
			return commands.ResultDisabled
		}
		operations.DefaultLog.Add(fmt.Sprintf("split %s into %d parts", cur.Name, len(parts)))
		ctx.SetStatus(fmt.Sprintf("Split into %d parts", len(parts)))
		return commands.ResultSuccess
	})
	register("cm_CombineFiles", "File Operations", func(ctx *commands.Context, _ []string) commands.Result {
		ctx.SetStatus("Combine: select .001 part and use command (manual)")
		return commands.ResultSuccess
	})
	register("cm_ShowOperations", "Logs", func(ctx *commands.Context, _ []string) commands.Result {
		app.ShowOperationLog()
		return commands.ResultSuccess
	})
	register("cm_ImportShortcuts", "Configuration", func(ctx *commands.Context, _ []string) commands.Result {
		dir, err := config.EnsureDir()
		if err != nil {
			ctx.SetStatus(err.Error())
			return commands.ResultDisabled
		}
		path := filepath.Join(dir, "shortcuts.yaml")
		if err := app.ImportShortcuts(path); err != nil {
			ctx.SetStatus(err.Error())
			return commands.ResultDisabled
		}
		ctx.SetStatus("Imported " + path)
		return commands.ResultSuccess
	})
}

// ExportDefaultShortcuts writes default bindings to path.
func ExportDefaultShortcuts(path string) error {
	return hotkeys.SaveYAML(path, []hotkeys.File{{
		Context:  "Main",
		Bindings: hotkeys.DefaultMainBindings(),
	}})
}
