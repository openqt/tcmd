# dc-tui cm_* Command Status (v1.0.0)

| Command | Phase | Status |
|---------|-------|--------|
| cm_Close / cm_Exit | 0 | done |
| cm_SwitchPnl | 0 | done |
| cm_HelpIndex | 0 | done |
| cm_FocusCmdLine | 2 | done |
| cm_Open / cm_ChangeDirToParent / cm_ChangeDirToRoot / cm_ChangeDirToHome | 1 | done |
| cm_Copy / cm_Rename / cm_RenameOnly / cm_MakeDir / cm_Delete | 1 | done |
| cm_View / cm_Edit | 1 | done (edit: external `$EDITOR`) |
| cm_Select / cm_SelectAll / cm_UnselectAll / cm_InvertSelection | 1 | done |
| cm_SortByName/Ext/Size/Date | 1 | done |
| cm_ShowSysFiles / cm_Refresh / cm_Exchange / cm_TransferPath | 1 | done |
| cm_LeftOpenDrives / cm_RightOpenDrives | 1 | done |
| cm_CalculateSpace | 1 | done |
| cm_NewTab / cm_CloseTab / cm_NextTab / cm_PrevTab | 2 | done |
| cm_ViewHistoryPrev / cm_ViewHistoryNext / cm_DirHistory | 2 | done |
| cm_QuickFilter / cm_QuickView / cm_FlatView / cm_FlatViewSel | 2 | done |
| cm_CopyFullNamesToClip / cm_CutToClipboard / cm_PasteFromClipboard | 2 | done |
| cm_Find / cm_MultiRename / cm_CompareDirectories / cm_DirHotList | 3 | done |
| cm_Properties / cm_EditDescr / cm_ShowCopyMoveDialog | 3 | done |
| cm_ShowAsText / cm_ShowAsHex / cm_ShowAsBin | 4 | done |
| cm_EditInternal / cm_CompareFilesByContent | 4 | done |
| cm_PackFiles / cm_UnpackFiles / cm_TestArchive | 5 | done |
| cm_FTPConnect / cm_SFTPConnect | 5 | external (stub) |
| cm_Options / cm_Config / cm_CommandBrowser | 6 | done |
| cm_CheckSumCalc / cm_SplitFile | 6 | done |
| cm_ShowOperations / cm_ImportShortcuts | 6 | done |
| cm_BriefView / cm_ColumnsView | 1 | degraded |

**未列出命令**：执行时显示 `[dc-tui] cm_XXX: not implemented`。完整 DC 命令集见 [官方文档](https://doublecmd.github.io/doc/en/cmds.html)。

Status: `done` | `degraded` | `external` | `unsupported` | `planned`
