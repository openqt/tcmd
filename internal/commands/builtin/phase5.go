package builtin

import (
	"github.com/openqt/tcmd/internal/commands"
	"github.com/openqt/tcmd/internal/filesrc"
	"github.com/openqt/tcmd/internal/filesrc/archive"
)

// Phase5App for archives and network.
type Phase5App interface {
	Phase4App
	EnterArchive(path string) error
	PackSelected() error
	UnpackCursor() error
}

// RegisterPhase5 registers archive commands.
func RegisterPhase5(r *commands.Registry, app Phase5App) {
	register := func(name, category string, fn commands.Handler) {
		r.Register(commands.Def{Name: name, Category: category, Handler: fn})
	}

	register("cm_PackFiles", "File Operations", func(ctx *commands.Context, _ []string) commands.Result {
		if err := app.PackSelected(); err != nil {
			ctx.SetStatus(err.Error())
			return commands.ResultDisabled
		}
		ctx.SetStatus("Packed")
		return commands.ResultSuccess
	})

	register("cm_UnpackFiles", "File Operations", func(ctx *commands.Context, _ []string) commands.Result {
		if err := app.UnpackCursor(); err != nil {
			ctx.SetStatus(err.Error())
			return commands.ResultDisabled
		}
		ctx.SetStatus("Unpacked")
		return commands.ResultSuccess
	})

	register("cm_TestArchive", "File Operations", func(ctx *commands.Context, _ []string) commands.Result {
		cur := ctx.ActivePanel.Current()
		if cur == nil || !archive.IsZipPath(cur.Path) {
			return commands.ResultDisabled
		}
		z := archive.NewZip(cur.Path)
		if _, err := z.List("", true); err != nil {
			ctx.SetStatus("Archive invalid: " + err.Error())
			return commands.ResultDisabled
		}
		ctx.SetStatus("Archive OK")
		return commands.ResultSuccess
	})

	register("cm_FTPConnect", "Network", func(ctx *commands.Context, _ []string) commands.Result {
		ctx.SetStatus("FTP: connect via ftp://user@host in command line (stub)")
		return commands.ResultSuccess
	})

	register("cm_SFTPConnect", "Network", func(ctx *commands.Context, _ []string) commands.Result {
		ctx.SetStatus("SFTP: connect via sftp://user@host (stub)")
		return commands.ResultSuccess
	})

	// Open zip archives on Enter
	register("cm_OpenArchive", "File Operations", func(ctx *commands.Context, _ []string) commands.Result {
		cur := ctx.ActivePanel.Current()
		if cur == nil || !filesrc.IsArchiveFile(cur.Path) {
			return commands.ResultNotFound
		}
		if err := app.EnterArchive(cur.Path); err != nil {
			ctx.SetStatus(err.Error())
			return commands.ResultDisabled
		}
		ctx.SetStatus("Archive")
		return commands.ResultSuccess
	})
}
