package builtin

import (
	"os"

	"github.com/openqt/tcmd/internal/commands"
	"github.com/openqt/tcmd/internal/config"
	"github.com/openqt/tcmd/internal/tools"
)

// Phase3App extends app for tool dialogs.
type Phase3App interface {
	Phase2App
	OpenDialog(kind string)
	DialogAction(action string)
}

// RegisterPhase3 registers tool dialog commands.
func RegisterPhase3(r *commands.Registry, app Phase3App) {
	register := func(name, category string, fn commands.Handler) {
		r.Register(commands.Def{Name: name, Category: category, Handler: fn})
	}

	register("cm_Find", "Tools", func(_ *commands.Context, _ []string) commands.Result {
		app.OpenDialog("find")
		return commands.ResultSuccess
	})
	register("cm_MultiRename", "Tools", func(_ *commands.Context, _ []string) commands.Result {
		app.OpenDialog("rename")
		return commands.ResultSuccess
	})
	register("cm_CompareDirectories", "Tools", func(_ *commands.Context, _ []string) commands.Result {
		app.OpenDialog("sync")
		return commands.ResultSuccess
	})
	register("cm_DirHotList", "Tools", func(_ *commands.Context, _ []string) commands.Result {
		app.OpenDialog("hotlist")
		return commands.ResultSuccess
	})
	register("cm_ShowCopyMoveDialog", "Tools", func(_ *commands.Context, _ []string) commands.Result {
		app.OpenDialog("copymove")
		return commands.ResultSuccess
	})
	register("cm_EditDescr", "File Operations", func(ctx *commands.Context, _ []string) commands.Result {
		cur := ctx.ActivePanel.Current()
		if cur == nil {
			return commands.ResultDisabled
		}
		commentPath := cur.Path + ".desc"
		if _, err := os.Stat(commentPath); os.IsNotExist(err) {
			_ = os.WriteFile(commentPath, []byte(""), 0o644)
		}
		ctx.SetStatus("Comment file: " + commentPath)
		return commands.ResultSuccess
	})
	register("cm_Properties", "File Operations", func(_ *commands.Context, _ []string) commands.Result {
		app.OpenDialog("properties")
		return commands.ResultSuccess
	})
	register("cm_Start", "Find Files", func(_ *commands.Context, _ []string) commands.Result {
		app.DialogAction("find_start")
		return commands.ResultSuccess
	})
	register("cm_RunMultiRename", "Multi-Rename", func(_ *commands.Context, _ []string) commands.Result {
		app.DialogAction("rename_run")
		return commands.ResultSuccess
	})
}

// RunFind executes search.
func RunFind(root, mask, content string) ([]string, error) {
	return tools.Find(tools.FindOptions{Root: root, NameMask: mask, Content: content})
}

// RunMultiRename applies rename rule.
func RunMultiRename(dir string, names []string, rule tools.RenameRule) error {
	return tools.ApplyRenames(dir, names, rule)
}

// RunSyncCompare compares directories.
func RunSyncCompare(left, right string) ([]tools.SyncItem, error) {
	return tools.CompareDirs(left, right)
}

// SaveHotlistPath adds path to hotlist.
func SaveHotlistPath(path string) error {
	h, err := config.LoadHotlist()
	if err != nil {
		return err
	}
	h.Add(path)
	return h.Save()
}
