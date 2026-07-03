# dc-tui：Go TUI 版 Double Commander 实现方案

> 参考项目：[doublecmd/doublecmd](https://github.com/doublecmd/doublecmd)  
> 目标：在终端中实现与 Double Commander（以下简称 DC）行为一致的双面板文件管理器，快捷键与内部命令（`cm_*`）完全对齐 DC 设计。

---

## 1. 项目定位

| 维度 | 说明 |
|------|------|
| 名称 | `dc-tui`（工作名，可调整） |
| 语言 | Go 1.22+ |
| 界面 | 全屏 TUI（Alt Screen），支持窄终端自适应 |
| 兼容性目标 | DC 默认快捷键、内部命令语义、面板交互逻辑 |
| 非目标 | 复刻 DC 的 Lazarus/GTK/Qt GUI；运行 Windows TC 插件（WCX/WDX/WFX/WLX） |

DC 源码规模：约 500+ Pascal 单元、723 个 `cm_*` 内部命令标识符、11 个快捷键上下文（Main / Copy-Move / Viewer / Editor / …）。本方案采用**分阶段交付**，优先保证主窗口核心体验与 F 键操作链，再逐步覆盖工具对话框与高级功能。

---

## 2. DC 架构分析（对标依据）

### 2.1 核心子系统

```
┌─────────────────────────────────────────────────────────────┐
│                      Main Window (fmain)                     │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────────┐  │
│  │ Left Panel   │  │ Right Panel  │  │ Quick View       │  │
│  │ (FileView +  │  │ (FileView +  │  │ (optional)       │  │
│  │  Tabs)       │  │  Tabs)       │  │                  │  │
│  └──────┬───────┘  └──────┬───────┘  └──────────────────┘  │
│         │                  │                                 │
│  ┌──────┴──────────────────┴──────────────────────────────┐ │
│  │ Command Line + Function Key Bar + Drive Buttons        │ │
│  └────────────────────────────────────────────────────────┘ │
└───────────────────────────┬─────────────────────────────────┘
                            │
        ┌───────────────────┼───────────────────┐
        ▼                   ▼                   ▼
  TMainCommands      TFormCommands         THotkeyManager
  (umaincommands)    (uformcommands)       (uhotkeymanager)
        │                   │                   │
        └───────────────────┼───────────────────┘
                            ▼
                   TFileSourceManager
              (local / archive / FTP / SFTP / VFS / …)
                            │
                            ▼
                 TOperationsManager (后台文件操作队列)
```

### 2.2 必须复刻的设计模式

1. **命令总线**：所有操作通过 `cm_<Name>` 字符串命令触发，`ExecuteCommand(cmd, params[])` 分发。
2. **快捷键上下文**：同一按键在不同界面（主窗口、查看器、编辑器等）映射不同命令。
3. **文件源抽象**：`IFileSource` 统一本地目录、压缩包虚拟目录、网络协议。
4. **双面板语义**：Active/Inactive Panel、Source/Target 概念贯穿复制/移动/快速查看。
5. **可配置热键**：支持 `shortcuts.scf` 风格的快捷键配置文件（见 DC 文档 Configuration > Keys）。

### 2.3 DC 功能清单（按优先级）

| 优先级 | 功能域 | DC 代表能力 |
|--------|--------|-------------|
| P0 | 双面板浏览 | 目录列表、排序、隐藏文件、驱动器切换、面板焦点切换 |
| P0 | 文件操作 | F3 查看、F4 编辑、F5 复制、F6 移动/重命名、F7 建目录、F8 删除 |
| P0 | 选择与标记 | Space/Insert 选择、Ctrl+A 全选、扩展/收缩选择、同扩展名选择 |
| P0 | 导航 | Enter/Backspace、Alt+←/→ 历史、Ctrl+←/→ 根目录/主目录 |
| P1 | 标签页 | Ctrl+T 新建、Ctrl+W 关闭、Ctrl+Tab 切换、Alt+1..9 激活 |
| P1 | 工具对话框 | 查找文件、批量重命名、目录同步、目录收藏夹 |
| P1 | 命令行 | 底部命令行、历史、执行外部命令 |
| P2 | 内部查看器/编辑器 | 文本/十六进制/二进制模式 |
| P2 | 压缩包 | ZIP/TAR/GZ/BZ2/XZ/7Z/RAR 作为子目录浏览 |
| P2 | 网络 | FTP、SFTP、SSH+SCP |
| P3 | 高级 | 自定义列、文件注释、校验和、分割/合并、日志、配置 UI |
| P3 | 插件 | TC 插件接口（见限制文档） |

---

## 3. 技术选型

### 3.1 TUI 框架：**Bubble Tea v2** + **Lip Gloss** + **Bubbles**

| 库 | 用途 |
|----|------|
| [bubbletea](https://github.com/charmbracelet/bubbletea) | Elm 架构状态机，适合多子界面（主窗口/对话框/查看器）切换 |
| [lipgloss](https://github.com/charmbracelet/lipgloss) | 双面板水平布局、边框、颜色主题 |
| [bubbles](https://github.com/charmbracelet/bubbles) | list、textinput、viewport、progress 等复用组件 |

**选择理由**：DC 是状态复杂的多窗口应用；Bubble Tea 的 `Model/Update/View` 与 DC 的「命令 → 状态变更 → 重绘」模型天然契合。子界面可作为嵌套 Model 路由消息（类似 SprintOS、Superfile 的实践）。

备选 `tview`：控件更丰富但命令式回调在 700+ 命令场景下难以维护，不作为首选。

### 3.2 支撑库

| 领域 | 库 |
|------|-----|
| 配置 | `gopkg.in/yaml.v3` 或兼容 DC 的 XML（`dcxmlconfig` 格式调研后决定） |
| 归档 | `archive/zip`、`github.com/nwaples/rardecode`、`github.com/ulikunitz/xz`、`github.com/dsnet/compress/bzip2`、`github.com/klauspost/compress/gzip/zstd` |
| SFTP/FTP | `github.com/pkg/sftp`、`github.com/jlaffaye/ftp` |
| 语法高亮（编辑器） | `github.com/alecthomas/chroma`（输出 ANSI） |
| 文件监控 | `github.com/fsnotify/fsnotify` |
| 校验和 | `crypto/md5`、`crypto/sha256` 等标准库 |
| 测试 | `testing` + `github.com/charmbracelet/x/exp/golden` |

---

## 4. 项目结构

```
dc-tui/
├── cmd/dctui/                 # 入口
├── internal/
│   ├── app/                   # 根 Tea Model，界面路由
│   ├── commands/              # cm_* 注册表、执行器、参数解析
│   │   ├── registry.go
│   │   ├── categories.go      # 对标 DC 17 个命令分类
│   │   └── builtin/           # 各分类命令实现
│   ├── hotkeys/
│   │   ├── context.go         # 11 个快捷键上下文
│   │   ├── binding.go
│   │   └── scf/               # shortcuts.scf 导入/导出
│   ├── panel/
│   │   ├── panel.go           # 单面板状态
│   │   ├── filelist.go        # 文件列表、排序、过滤
│   │   ├── selection.go       # 标记逻辑
│   │   └── tabs.go              # 标签页
│   ├── filesrc/
│   │   ├── source.go          # IFileSource 接口
│   │   ├── local/
│   │   ├── archive/
│   │   ├── ftp/
│   │   └── sftp/
│   ├── operations/
│   │   ├── queue.go           # 后台操作队列（对标 OperationsManager）
│   │   ├── copy.go
│   │   ├── move.go
│   │   └── delete.go
│   ├── ui/
│   │   ├── mainview/          # 主窗口布局
│   │   ├── dialogs/           # 各类模态对话框
│   │   ├── viewer/
│   │   ├── editor/
│   │   └── widgets/           # 可复用 TUI 组件
│   ├── config/
│   └── platform/              # OS 相关（回收站、打开方式、终端）
├── docs/
│   ├── PLAN.md                # 本文档
│   ├── TUI_LIMITATIONS.md     # TUI 不兼容项
│   ├── COMMANDS.md            # cm_* 命令实现状态追踪
│   └── SHORTCUTS.md           # 默认快捷键对照表（摘自 DC 文档）
├── test/
│   └── integration/
├── go.mod
└── README.md
```

---

## 5. 核心设计

### 5.1 命令系统

```go
// 对标 DC TFormCommands.ExecuteCommand
type CommandFunc func(ctx *CommandContext, params []string) Result

type CommandContext struct {
    App           *App
    ActivePanel   *panel.Panel
    InactivePanel *panel.Panel
    SourcePanel   *panel.Panel  // 复制操作源
    TargetPanel   *panel.Panel  // 复制操作目标
}

type Registry struct {
    commands map[string]CommandDef
}

type CommandDef struct {
    Name        string   // cm_Copy
    Category    string   // File Operations
    Handler     CommandFunc
    Enabled     func(*CommandContext) bool
}
```

**实现策略**：

1. 从 DC 的 `uglobs.pas`（默认热键）、`cmds.html` 文档提取命令清单，生成 `docs/COMMANDS.md` 追踪表。
2. 未实现的命令返回 `cfrNotFound` 等价行为，状态栏提示「未实现：cm_XXX」。
3. 命令参数格式与 DC 一致（`key=value` 多行参数，如 `cm_Delete` 的 `trashcan=reversesetting`）。

### 5.2 快捷键系统

```go
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

type Binding struct {
    Keys       []string      // 支持多快捷键，如 F8 与 Del
    Command    string
    Params     []string
    Controls   []string      // OnlyForControls: files, cmdline, quicksearch
}
```

**终端适配层**：

- 建立 `KeyEvent → DC KeyChord` 归一化（处理 Kitty/legacy 终端差异）。
- 数字键盘、F1–F12、Alt+字母 在部分终端需配置；提供 `docs/TERMINAL_SETUP.md`。
- 冲突检测逻辑复刻 DC：新绑定若与已有冲突则警告。

### 5.3 主窗口布局（TUI 映射）

```
┌─ Function Keys: F1 F2 F3 F4 F5 F6 F7 F8 F9 F10 ─────────────────┐
├─ Drives: [/] [home] [tmp] ... ────────────────────────────────────┤
├─ Tabs L: [~/proj] [~/doc] │ Tabs R: [/tmp] [~/dl] ───────────────┤
├──────────────────────┬──────────────────────┬─────────────────────┤
│ LEFT PANEL (active)  │ RIGHT PANEL          │ QUICK VIEW (opt)  │
│ > ..                 │   ..                 │ preview text...   │
│   src/               │   bin/               │                   │
│   main.go            │   app                │                   │
├──────────────────────┴──────────────────────┴─────────────────────┤
│ > command line___________________________________________          │
├─ Status: 3 files | 1 sel | 12.4 MB ──────────────────────────────┤
└─ Operations: [████████░░] copying file.zip ────────────────────────┘
```

- **Brief View**（Ctrl+F1）：仅文件名单列。
- **Columns View**（Ctrl+F2）：Name | Size | Date | Attr（TUI 宽度自适应裁列）。
- **Thumbnails View**：见 `TUI_LIMITATIONS.md`（降级为图标字符或禁用）。
- **Flat View**（Ctrl+B）：递归展平子目录文件列表。

### 5.4 文件源接口

```go
type FileSource interface {
    Scheme() string
    List(path string) ([]FileEntry, error)
    Stat(path string) (FileInfo, error)
    OpenRead(path string) (io.ReadCloser, error)
    // 操作接口拆分：CopyTo, MoveTo, Delete, Mkdir, ...
    Capabilities() FileSourceCaps
}
```

首批实现：`local://`、`file://`、`.zip`/`.tar.*`/`.7z` 虚拟路径（`archive://path.zip/inner/file`）。

### 5.5 后台操作队列

对标 DC `OperationsManager`：

- F5/F6 弹出 Copy/Move 对话框，确认后加入队列。
- 底部操作进度条，Ctrl+F12 / Alt+V 查看详情。
- 支持暂停、跳过、错误策略（询问/全部跳过/中止）。

---

## 6. 分阶段实施计划

### Phase 0 — 基础设施

- [x] 初始化 Go module、CI（`go test ./...`、`golangci-lint`）
- [x] 命令注册表骨架 + 热键分发器
- [x] 从 DC `uglobs.pas` 提取默认快捷键，生成 `docs/SHORTCUTS.md`
- [x] 主窗口空壳：双面板 + 状态栏 + F 键栏
- [x] 配置目录：`~/.config/dc-tui/`（或 `$XDG_CONFIG_HOME`）

**交付物**：可启动的空双面板应用，Tab 切换焦点，Quit（Alt+F4 / Alt+X）。

### Phase 1 — 核心文件管理（MVP）

- [x] 本地 `FileSource`：列表、排序（Ctrl+F3–F6）、隐藏文件（Ctrl+H）
- [x] 导航：Enter、Backspace、Ctrl+\、Ctrl+←/→、Space 计算目录大小
- [x] 选择：Insert、Ctrl+A、小键盘 *、扩展/收缩选择
- [x] 文件操作：F5 复制、F6 移动/重命名、F7 建目录、F8 删除（含 Shift 永久删除）
- [x] F2 重命名、F3 简易查看（文本）、F4 外部编辑器调用
- [x] 面板交换 Ctrl+U、目标=源 Ctrl+\
- [x] 驱动器列表 Alt+F1/F2

**验收**：完成日常双面板复制/移动/删除，快捷键与 DC 主窗口表一致。

### Phase 2 — 导航增强与标签页

- [x] 目录历史 Alt+↓、Alt+←/→
- [x] 标签页：Ctrl+T/W、Ctrl+Tab、Alt+1..9/0
- [x] Quick Filter（Ctrl+S）、Quick Search（输入即搜）
- [x] Flat View（Ctrl+B / Ctrl+Shift+B）
- [x] Quick View（Ctrl+Q）在第三栏显示文件头内容
- [x] 命令行：Ctrl+L 聚焦、历史 Alt+F8/Ctrl+↓、Enter 执行
- [x] 剪贴板：Ctrl+C/X/V 文件路径

### Phase 3 — 工具对话框

- [x] **Find Files**（Alt+F7）：标准/高级/结果页，F9 开始搜索
- [x] **Multi-Rename Tool**（Ctrl+M）：规则预设、占位符、F9 执行
- [x] **Synchronize Directories**（Shift+F2）：双目录比较
- [x] **Directory Hotlist**（Ctrl+D）：收藏目录
- [x] **Copy/Move Dialog** 完整字段循环（Shift+F5）
- [x] 文件注释（Ctrl+Z）、属性（Alt+Enter 精简版）

### Phase 4 — 内部查看器与编辑器

- [x] Viewer 上下文全部快捷键（hex/text/binary 模式切换）
- [x] Editor：文本编辑（语法高亮降级为纯文本，见 TUI_LIMITATIONS.md）
- [x] Differ（简易双文件 diff 视图）

### Phase 5 — 归档与网络

- [ ] 归档文件源：ZIP、TAR、GZ、BZ2、XZ、7Z、RAR
- [ ] 打包/解包 Alt+F5/F9
- [ ] SFTP/FTP 浏览（对标 DC Network 分类命令）
- [ ] 归档完整性测试 Alt+Shift+F9

### Phase 6 — 配置与高级功能

- [ ] 配置对话框（cm_Config 各 section）
- [ ] 自定义列（Columns View 扩展）
- [ ] 校验和计算/验证
- [ ] 文件分割/合并
- [ ] 操作日志
- [ ] `shortcuts.scf` 导入导出
- [ ] 命令浏览器（Shift+F12）

---

## 7. 测试策略

| 类型 | 内容 |
|------|------|
| 单元测试 | 命令参数解析、路径处理、选择逻辑、排序 |
| 黄金测试 | 主界面渲染快照（lipgloss 输出） |
| 集成测试 | 临时目录中 F5/F6/F7/F8 操作链 |
| 兼容性测试 | 每个 `cm_*` 是否有 handler；快捷键表与 DC 默认一致 |
| 终端矩阵 | xterm、kitty、alacritty、tmux 下 F 键与 Alt 组合 |

---

## 8. 配置与 DC 互操作

| 文件 | 策略 |
|------|------|
| `shortcuts.scf` | 优先支持导入；原生保存为 YAML 超集 |
| `dirhotlist.txt` | 尝试兼容 DC 格式 |
| `multiarc.ini` | 归档关联参考 `default/multiarc.ini` |
| `doublecmd.xml` | 分阶段映射关键配置项（显示隐藏文件、排序方式等） |

---

## 9. 风险与缓解

| 风险 | 缓解 |
|------|------|
| 723 个命令无法一次实现 | `COMMANDS.md` 追踪；未实现命令明确提示 |
| 终端快捷键差异 | 提供默认 + 可配置映射；文档说明 |
| 归档格式专利/许可 | RAR 只读解压；7z 调用外部 `7z` 可选 |
| 大目录性能 | 虚拟滚动（viewport 仅渲染可见行）；异步目录统计 |
| 双击/鼠标 | 可选鼠标支持（Bubble Tea 已支持），非核心 |

---

## 10. 文档交付物

| 文档 | 说明 |
|------|------|
| `docs/PLAN.md` | 本方案 |
| `docs/TUI_LIMITATIONS.md` | 无法实现或与 TUI 不兼容的功能 |
| `docs/SHORTCUTS.md` | DC 默认快捷键完整对照 |
| `docs/COMMANDS.md` | `cm_*` 实现状态矩阵 |
| `docs/TERMINAL_SETUP.md` | 各终端 F 键/Alt 键配置指南 |
| `README.md` | 构建、运行、快速上手 |

---

## 11. 建议的首个里程碑（v0.1.0）

> 双面板本地文件管理 + DC 默认主窗口快捷键 + F 键操作链

预估涉及组件：Phase 0 + Phase 1 全部，约 40–60 个核心 `cm_*` 命令。

---

## 12. 参考链接

- [DC 官方文档](https://doublecmd.github.io/doc/en/)
- [DC 快捷键](https://doublecmd.github.io/doc/en/shortcuts.html)
- [DC 内部命令](https://doublecmd.github.io/doc/en/cmds.html)
- [DC 配置](https://doublecmd.github.io/doc/en/configuration.html)
- [DC 源码](https://github.com/doublecmd/doublecmd)
