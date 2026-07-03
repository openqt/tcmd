# dc-tui

Go 语言实现的终端双面板文件管理器，设计目标与 [Double Commander](https://github.com/doublecmd/doublecmd) 的快捷键与内部命令（`cm_*`）保持一致。

## 状态

| Phase | 状态 | 版本 |
|-------|------|------|
| 0 — 基础设施 | ✅ 完成 | v0.0.1 |
| 1 — 核心文件管理 | ✅ 完成 | v0.1.0 |
| 2 — 导航增强 | ✅ 完成 | v0.2.0 |
| 3 — 工具对话框 | 🚧 进行中 | — |

## 构建与运行

```bash
go build -o dctui ./cmd/dctui
./dctui [left-path] [right-path]
```

## 测试

```bash
go test ./...
```

## Phase 0 功能

- 双面板主窗口（F 键栏、驱动器栏、标签栏、状态栏）
- `cm_*` 命令注册表与热键分发器（11 个上下文骨架）
- Tab 切换面板（`cm_SwitchPnl`）
- Alt+F4 / Alt+X 退出（`cm_Exit`）
- 配置目录：`$XDG_CONFIG_HOME/dc-tui` 或 `~/.config/dc-tui`

## 文档

- [实现方案](docs/PLAN.md)
- [TUI 限制说明](docs/TUI_LIMITATIONS.md)
- [默认快捷键对照](docs/SHORTCUTS.md)
- [命令实现状态](docs/COMMANDS.md)

## 许可证

MIT License（见 [LICENSE](LICENSE)）
