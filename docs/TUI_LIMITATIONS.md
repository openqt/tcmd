# TUI 限制与无法实现功能说明

本文档列出 Go TUI 版 Double Commander（`dc-tui`）中**无法完全复刻**或**需要降级实现**的功能，并说明原因与替代方案。

> 参考： [doublecmd/doublecmd](https://github.com/doublecmd/doublecmd)  
> 原则：凡有差异，行为上尽量接近 DC；无法满足时，在状态栏或帮助中明确提示。

---

## 1. 完全不支持（架构性限制）

### 1.1 Total Commander 插件（WCX / WDX / WFX / WLX）

| 项目 | 说明 |
|------|------|
| DC 行为 | 加载 Windows 原生 DLL 插件扩展归档、查看器、字段、文件系统 |
| TUI 限制 | Go 无法加载 Windows TC 插件；跨平台无对应 ABI |
| 替代方案 | 用 Go 原生库实现常用归档；`docs/COMMANDS.md` 标记相关 `cm_*` 为 `unsupported` |
| 受影响命令示例 | `cm_ConfigPlugins`、WDX 自定义列、部分 WLX 查看器模式 |

### 1.2 原生 GUI 外壳集成

| 项目 | 说明 |
|------|------|
| DC 行为 | GTK/Qt/Win32 原生窗口、系统托盘、多窗口拖拽 |
| TUI 限制 | 单终端全屏/内联，无独立窗口 |
| 替代方案 | Alt Screen 全屏模式；对话框以模态 TUI 面板实现 |

### 1.3 拖放（Drag & Drop）

| 项目 | 说明 |
|------|------|
| DC 行为 | 鼠标拖放文件到面板或工具栏 |
| TUI 限制 | 终端内无操作系统级拖放 |
| 替代方案 | 保留 F5/F6 键盘流程；可选终端鼠标选中（非拖放） |

### 1.4 打印（cm_Print / Ctrl+P）

| 项目 | 说明 |
|------|------|
| DC 行为 | 内部查看器打印到系统打印机 |
| TUI 限制 | 终端无法直接驱动打印机 |
| 替代方案 | 导出到文件后调用 `lp`/`lpr` 或提示用户手动打印 |

### 1.5 Windows 特有文件系统功能

| 功能 | 说明 |
|------|------|
| 回收站深度集成 | Windows Shell 回收站 API；Linux 可部分实现（Freedesktop Trash），macOS 另议 |
| UNC 路径 / WinNet | `winnet` 文件源为 Windows 专用 |
| WSL 文件源 | DC 的 `uwslfilesource` 依赖 Windows 环境 |
| AirDrop、iCloud Drive | macOS 平台专用集成 |
| 快捷方式 `.lnk` 解析 | 可部分实现，但不如 DC 完整 |

---

## 2. 降级实现（功能缩减）

### 2.1 缩略图视图（Thumbnails View, Ctrl+Shift+F1）

| 项目 | 说明 |
|------|------|
| DC 行为 | 显示图片/GIF/视频/PDF 等缩略图网格 |
| TUI 限制 | 大多数终端无法渲染位图；Sixel/iTerm2 图片协议支持有限 |
| 降级方案 A | 用文件类型图标字符（📁 📄 🖼）+ 文件名网格 |
| 降级方案 B | 检测终端能力后可选 Sixel 预览（Phase 6+ 可选） |
| 默认 | 降级方案 A；状态栏提示「Thumbnails 已降级为图标视图」 |

### 2.2 内部查看器 — 富格式预览

| 模式 | DC 能力 | TUI 降级 |
|------|---------|----------|
| 图片 | 内嵌解码显示 | 不支持或 Sixel 可选 |
| GIF 动画 | `gifview` 组件播放 | 不支持 |
| Office 文档 | `cm_ShowOffice` | 不支持；提示用外部程序打开 |
| PDF | 内嵌预览 | 不支持；`cm_Open` 调外部查看器 |
| 二进制 | 完整 hex + 解码 | 支持 hex + 文本（与 DC 一致） |
| 语法高亮 | SynUniHighlighter | Chroma 输出 ANSI（色彩取决于终端） |

### 2.3 内部编辑器

| 项目 | 说明 |
|------|------|
| DC 行为 | SynEdit 完整 IDE 体验：多标签、断点、复杂撤销树 |
| TUI 降级 | 单文件编辑 + 基础撤销 + Chroma 高亮；复杂编辑建议 F4 调 `$EDITOR` |
| 受影响 | 部分 Editor 上下文高级快捷键保留绑定，但功能可能简化为「提示使用外部编辑器」 |

### 2.4 树形面板（Tree View / cm_ShowTreeView）

| 项目 | 说明 |
|------|------|
| DC 行为 | 左侧可展开目录树，与列表面板联动 |
| TUI 降级 | 用可折叠树形列表（缩进 + ▶/▼）占单列或弹出层；宽度受限时改为模态树 |
| 快捷键 | `Ctrl+Shift+F8` 等保留，行为为切换树形侧栏 |

### 2.5 内嵌终端（F9 / cm_RunTerm）

| 项目 | 说明 |
|------|------|
| DC 行为 | 在 GUI 底部/分屏嵌入 VT |
| TUI 限制 | 已在终端中运行，无法再嵌套完整 VT（除非 tmux 分屏） |
| 降级方案 | F9 在当前目录启动外部 shell 脚本；或提示 `Ctrl+Z` 挂起后手动开终端 |
| 可选 | 检测 tmux 时 `split-window` 水平分屏 |

### 2.6 上下文菜单（右键 / Shift+F10 / cm_ContextMenu）

| 项目 | 说明 |
|------|------|
| DC 行为 | 系统 Shell 右键菜单（Windows COM / GTK） |
| TUI 降级 | 弹出 TUI 菜单，包含常用内部命令子集 |
| 差异 | 无动态「发送到」、无完整文件关联子菜单（除非自行解析 `.desktop`/注册表） |

### 2.7 配置对话框（cm_Config*）

| 项目 | 说明 |
|------|------|
| DC 行为 | 大型分页 GUI 配置（数十页） |
| TUI 降级 | 分页 TUI 表单；复杂选项（颜色选择器、字体预览）简化 |
| 策略 | 优先支持：热键、显示、文件操作、归档；其余只读展示或 YAML 手工编辑 |

### 2.8 工具栏 / 按钮栏（可配置图标按钮）

| 项目 | 说明 |
|------|------|
| DC 行为 | 可拖拽图标工具栏，绑定 `cm_*` |
| TUI 降级 | 顶部 F1–F10 文字键栏 + 可选一行数字快捷键提示 |
| 配置 | `cm_ConfigToolbars` 简化为键位绑定，无图标拖拽 |

### 2.9 文件关联（cm_Open 调系统关联）

| 项目 | 说明 |
|------|------|
| DC 行为 | 调用 OS 默认程序打开 |
| TUI 降级 | Linux：`xdg-open`；macOS：`open`；Windows：`start`；无关联时提示 |
| 差异 | 无法列出完整关联菜单（同 2.6） |

### 2.10 多显示器 / 窗口布局

| 项目 | 说明 |
|------|------|
| DC 行为 | 自由调整窗口大小、多显示器、面板水平/垂直分割 |
| TUI 降级 | 响应终端尺寸；`cm_VerticalPanels` 在窄屏自动垂直堆叠双面板 |
| 全屏 | `cm_FullScreen` 映射为终端全屏（部分终端支持） |

---

## 3. 终端环境限制

### 3.1 快捷键冲突

| 按键 | 问题 | 处理 |
|------|------|------|
| Alt+F4 | 多数桌面环境关闭窗口 | 保留绑定；文档说明 WM 可能拦截 |
| Ctrl+S | 部分终端 XON/XOFF 冻结 | 提供配置项改键 |
| Alt+方向键 | 部分终端发送 ESC 序列不完整 | 提供 Kitty keyboard protocol 优先 |
| Fn 键 | 依赖终端转发 | `TERMINAL_SETUP.md` 说明 |
| 小键盘 | NumLock 状态影响 | 检测并回退到主键盘等价键 |

### 3.2 鼠标

| 项目 | 说明 |
|------|------|
| DC | 完整鼠标操作：点击、双击、右键、拖放 |
| TUI | Bubble Tea 支持基础鼠标；**不保证**双击、拖放 |
| 策略 | 鼠标为可选增强；所有功能必须可用键盘完成 |

### 3.3 颜色与 Unicode

| 项目 | 说明 |
|------|------|
| 真彩色 | 需 `COLORTERM=truecolor` |
| 中日韩宽字符 | 列对齐用 `runewidth` 计算 |
| Nerd Fonts 图标 | 可选；默认 ASCII 降级 |

---

## 4. 网络与协议差异

| 协议 | DC | dc-tui 计划 |
|------|-----|-------------|
| FTP/FTPS/FTPES | 完整 | Phase 5 实现 |
| SFTP/SCP | 完整 | Phase 5 实现 |
| WebDAV | 部分版本 | 待评估，可能不支持 |
| Google Drive | GIO 集成 | 不支持（OAuth GUI 复杂） |
| 管理员/提权 RPC | `rpc/` 模块 | 不支持；`sudo` 提示 |

---

## 5. 命令实现策略

对本文档涉及的功能，在 `docs/COMMANDS.md` 中使用以下状态：

| 状态 | 含义 |
|------|------|
| `done` | 与 DC 行为一致 |
| `degraded` | 已实现但功能缩减（见本文档章节） |
| `external` | 转交外部程序 |
| `unsupported` | 不计划实现 |
| `planned` | 后续阶段 |

用户执行 `unsupported` 命令时，状态栏显示：

```
[dc-tui] cm_ShowOffice: not available in TUI — use cm_Open with external viewer
```

---

## 6. 与 DC 配置文件兼容性

| 文件 | 兼容性 |
|------|--------|
| `shortcuts.scf` | 计划导入；不保证 100% 参数语义 |
| `doublecmd.xml` | 仅映射子集 |
| TC 插件配置 | 不兼容 |
| `wcx`/`wlx` 插件目录 | 忽略 |

---

## 7. 总结

`dc-tui` 的设计目标是：**在终端内尽可能完整地复刻 DC 的键盘驱动双面板工作流**，而非复制其 GUI 与 Windows 插件生态。本文档将随实现进展更新；每个 `degraded` 项应在 `--help` 或 `cm_Help` 中可查阅。
