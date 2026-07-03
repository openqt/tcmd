package builtin

import (
	"github.com/openqt/tcmd/internal/commands"
	"github.com/openqt/tcmd/internal/ui/differ"
	"github.com/openqt/tcmd/internal/ui/editor"
	"github.com/openqt/tcmd/internal/ui/viewer"
)

// Phase4App extends app for viewer/editor/differ.
type Phase4App interface {
	Phase3App
	SetViewerMode(path string, mode viewer.Mode)
	OpenEditor(path string)
	OpenDiffer(left, right string)
}

// RegisterPhase4 registers viewer/editor/differ commands.
func RegisterPhase4(r *commands.Registry, app Phase4App) {
	register := func(name, category string, fn commands.Handler) {
		r.Register(commands.Def{Name: name, Category: category, Handler: fn})
	}

	register("cm_ShowAsText", "Viewer", func(ctx *commands.Context, _ []string) commands.Result {
		cur := ctx.ActivePanel.Current()
		if cur == nil || cur.IsDir {
			return commands.ResultDisabled
		}
		app.SetViewerMode(cur.Path, viewer.ModeText)
		return commands.ResultSuccess
	})
	register("cm_ShowAsHex", "Viewer", func(ctx *commands.Context, _ []string) commands.Result {
		cur := ctx.ActivePanel.Current()
		if cur == nil || cur.IsDir {
			return commands.ResultDisabled
		}
		app.SetViewerMode(cur.Path, viewer.ModeHex)
		return commands.ResultSuccess
	})
	register("cm_ShowAsBin", "Viewer", func(ctx *commands.Context, _ []string) commands.Result {
		cur := ctx.ActivePanel.Current()
		if cur == nil || cur.IsDir {
			return commands.ResultDisabled
		}
		app.SetViewerMode(cur.Path, viewer.ModeBinary)
		return commands.ResultSuccess
	})
	register("cm_EditInternal", "Editor", func(ctx *commands.Context, _ []string) commands.Result {
		cur := ctx.ActivePanel.Current()
		if cur == nil || cur.IsDir {
			return commands.ResultDisabled
		}
		app.OpenEditor(cur.Path)
		return commands.ResultSuccess
	})
	register("cm_CompareFilesByContent", "Tools", func(ctx *commands.Context, _ []string) commands.Result {
		left := ctx.ActivePanel.Current()
		right := ctx.InactivePanel.Current()
		if left == nil || right == nil || left.IsDir || right.IsDir {
			return commands.ResultDisabled
		}
		app.OpenDiffer(left.Path, right.Path)
		return commands.ResultSuccess
	})
}

// RenderEditor loads and highlights a file.
func RenderEditor(path string) (string, error) {
	content, err := editor.Load(path)
	if err != nil {
		return "", err
	}
	return editor.Highlight(path, content)
}

// RenderDiffer compares two paths.
func RenderDiffer(left, right string) (string, error) {
	return differ.Compare(left, right)
}
