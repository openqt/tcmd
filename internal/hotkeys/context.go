package hotkeys

// Context identifies the active Double Commander shortcut scope.
type Context int

const (
	ContextMain Context = iota
	ContextCopyMoveDialog
	ContextEditComment
	ContextFindFiles
	ContextMultiRename
	ContextSyncDirs
	ContextViewer
	ContextEditor
	ContextDiffer
	ContextConfig
	ContextDirHotlist
)

// String returns the DC context name.
func (c Context) String() string {
	switch c {
	case ContextMain:
		return "Main"
	case ContextCopyMoveDialog:
		return "CopyMoveDialog"
	case ContextEditComment:
		return "EditComment"
	case ContextFindFiles:
		return "FindFiles"
	case ContextMultiRename:
		return "MultiRename"
	case ContextSyncDirs:
		return "SyncDirs"
	case ContextViewer:
		return "Viewer"
	case ContextEditor:
		return "Editor"
	case ContextDiffer:
		return "Differ"
	case ContextConfig:
		return "Config"
	case ContextDirHotlist:
		return "DirHotlist"
	default:
		return "Unknown"
	}
}
