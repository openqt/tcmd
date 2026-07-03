# 终端配置指南

dc-tui 依赖终端正确转发功能键与 Alt 组合键。以下为常见终端建议配置。

## 通用

- 使用 **UTF-8**  locale
- 推荐终端宽度 ≥ 120 列、高度 ≥ 30 行
- 设置 `COLORTERM=truecolor` 以获得更好配色

## xterm / 默认 Linux 终端

- F1–F12 通常无需额外配置
- Alt+方向键：确保终端发送 ESC 序列而非菜单快捷键

## tmux

```bash
# ~/.tmux.conf
set -g extended-keys on
set -g extended-keys-format csi-u
```

## Kitty

Kitty 默认支持完整键盘协议，推荐用于 dc-tui。

## 已知冲突

| 按键 | 冲突 | 处理 |
|------|------|------|
| Alt+F4 | 窗口管理器关闭 | 可用 `Alt+X` 退出 |
| Ctrl+S | XON 流控 | 在 `shortcuts.yaml` 中改绑 `cm_QuickFilter` |

## 快捷键自定义

编辑 `~/.config/dc-tui/shortcuts.yaml` 或在应用内执行 `cm_ImportShortcuts`（首次运行会自动生成模板）。

```yaml
- context: Main
  bindings:
    - key: F5
      command: cm_Copy
```

更多限制见 [TUI_LIMITATIONS.md](TUI_LIMITATIONS.md)。
