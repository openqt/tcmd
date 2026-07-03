# dc-tui

Go 语言实现的终端双面板文件管理器，设计目标与 [Double Commander](https://github.com/doublecmd/doublecmd) 的快捷键与内部命令（`cm_*`）保持一致。

## 状态

**全部 6 个 Phase 已完成** — v1.0.0

| Phase | 内容 | 版本 |
|-------|------|------|
| 0 | 基础设施、命令/热键骨架、双面板空壳 | v0.0.1 |
| 1 | 本地文件管理、F 键操作链 | v0.1.0 |
| 2 | 标签页、历史、快速过滤/查看、命令行、剪贴板 | v0.2.0 |
| 3 | 工具对话框（查找、重命名、同步、收藏夹） | v0.3.0 |
| 4 | 内部查看器/编辑器/比较器 | v0.4.0 |
| 5 | ZIP 归档虚拟目录、打包/解包 | v0.5.0 |
| 6 | 配置、校验和、快捷键导入、命令浏览器 | v1.0.0 |

## 构建与运行

```bash
go build -o dctui ./cmd/dctui
./dctui [left-path] [right-path]
```

## 测试

```bash
go test ./...
```

## 配置目录

`$XDG_CONFIG_HOME/dc-tui/` 或 `~/.config/dc-tui/`：

| 文件 | 说明 |
|------|------|
| `settings.yaml` | 显示/排序/编辑器等偏好 |
| `dirhotlist.txt` | 目录收藏夹 |
| `shortcuts.yaml` | 自定义快捷键（可导入） |

## 核心快捷键（主窗口）

| 键 | 功能 |
|----|------|
| Tab | 切换面板 |
| F3/F4 | 查看/编辑 |
| F5/F6 | 复制/移动 |
| F7/F8 | 建目录/删除 |
| Ctrl+T/W | 新建/关闭标签 |
| Ctrl+U | 交换面板 |
| Alt+F7 | 查找文件 |
| Ctrl+M | 批量重命名 |
| Alt+F4 | 退出 |

完整列表见 [docs/SHORTCUTS.md](docs/SHORTCUTS.md)。

## 文档

- [实现方案](docs/PLAN.md)
- [TUI 限制说明](docs/TUI_LIMITATIONS.md)
- [默认快捷键](docs/SHORTCUTS.md)
- [命令实现状态](docs/COMMANDS.md)
- [终端配置指南](docs/TERMINAL_SETUP.md)

## 许可证

MIT License（见 [LICENSE](LICENSE)）
