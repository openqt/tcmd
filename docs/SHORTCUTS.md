# 默认快捷键对照表（摘自 Double Commander）

> 来源：[DC Shortcuts Documentation](https://doublecmd.github.io/doc/en/shortcuts.html)  
> 实现时以 `internal/hotkeys` 包为权威；本文档供对照与测试使用。

## 快捷键上下文

| # | 上下文 | DC 名称 |
|---|--------|---------|
| 1 | 主窗口 | Main window |
| 2 | 复制/移动对话框 | Copy/Move Dialog |
| 3 | 编辑注释 | Edit Comment Dialog |
| 4 | 查找文件 | Find Files |
| 5 | 批量重命名 | Multi-Rename Tool |
| 6 | 目录同步 | Synchronize Directories |
| 7 | 内部查看器 | Internal Viewer |
| 8 | 内部编辑器 | Internal Editor |
| 9 | 文件比较 | Differ |
| 10 | 配置 | Configuration |
| 11 | 目录收藏 | Directory Hotlist |

---

## 1. 主窗口（Main Window）

| 快捷键 | 命令/动作 |
|--------|-----------|
| F1 | 帮助 |
| F2, Shift+F6 | 重命名（光标在 `..` 且无选择时编辑路径） |
| F3 | 查看文件 / 进入目录 |
| Shift+F3 | 仅查看光标处文件 |
| F4 | 编辑 |
| Shift+F4 | 新建文本并编辑 |
| F5 | 复制到对面板 |
| Shift+F5 | 同目录复制 |
| F6 | 移动/重命名 |
| F7 | 新建目录 |
| F8, Del | 删除到回收站（Shift 反转） |
| Shift+F8, Shift+Del | 永久删除 |
| F9 | 启动终端 |
| Alt+F1 | 切换左驱动器 |
| Alt+F2 | 切换右驱动器 |
| Alt+F4, Alt+X | 退出 |
| Alt+F5 | 打包 |
| Alt+F7 | 查找文件 |
| Alt+F8 | 命令行历史菜单 |
| Alt+F9 | 解包 |
| Alt+1..9 | 按索引激活标签 |
| Alt+0 | 激活最后标签 |
| Alt+↓ | 目录历史 |
| Alt+← | 历史上一条 |
| Alt+→ | 历史下一条 |
| Alt+Shift+F9 | 校验归档 |
| Alt+Enter | 文件属性 |
| Alt+Shift+Enter | 计算所有目录大小 |
| Alt+Del | 擦除文件 |
| Alt+V | 操作进度窗口 |
| Alt+Z | 目录收藏夹配置 |
| Ctrl+F1 | 简要视图 |
| Ctrl+F2 | 列视图 |
| Ctrl+Shift+F1 | 缩略图视图 |
| Ctrl+F3 | 按名称排序 |
| Ctrl+F4 | 按扩展名排序 |
| Ctrl+F5 | 按日期排序 |
| Ctrl+F6 | 按大小排序 |
| Ctrl+1..9 | 按索引打开驱动器 |
| Ctrl+Alt+Enter | 用系统关联打开 |
| Ctrl+Tab | 下一标签 |
| Ctrl+Shift+Tab | 上一标签 |
| Ctrl+A | 全选 |
| Ctrl+B | 展平视图 |
| Ctrl+Shift+B | 展平视图（仅选中项） |
| Ctrl+C | 复制路径到剪贴板 |
| Ctrl+D | 目录收藏夹 |
| Ctrl+H | 显示/隐藏隐藏文件 |
| Ctrl+L | 聚焦命令行 |
| Ctrl+M | 批量重命名 |
| Ctrl+O | 在新标签打开目录 |
| Ctrl+P | 切换控制台全屏 |
| Ctrl+Q | 快速查看 |
| Ctrl+R | 刷新 |
| Ctrl+S | 快速搜索 |
| Ctrl+T | 新建标签 |
| Ctrl+U | 交换面板 |
| Ctrl+V | 粘贴 |
| Ctrl+W | 关闭标签 |
| Ctrl+X | 剪切到剪贴板 |
| Ctrl+Z | 编辑文件注释 |
| Ctrl+↑/↓ | 展开/折叠树（树视图） |
| Ctrl+← | 根目录 |
| Ctrl+→ | 主目录 |
| Ctrl+\\ | 对面板打开同目录 |
| Ctrl+. | 追加路径到命令行 |
| Ctrl+Enter | 追加选中项到命令行 |
| Ctrl+Shift+Enter | 追加路径+文件名到命令行 |
| Ctrl+Shift+F7 | 新建搜索实例 |
| Ctrl+Shift+F8 | 树形面板 |
| Ctrl+Shift+Home | 显示标签列表 |
| Ctrl+Shift+A/C/X | 复制全名/文件名到剪贴板 |
| Ctrl+Shift+D | 目录收藏夹配置 |
| Ctrl+Shift+H | 目录历史下拉 |
| Ctrl+PgDn/PgUp | 对面板打开同目录 |
| Ctrl+Num+/Num- | 全选/取消全选 |
| Num+ / Num- | 选择/取消选择 |
| Num* | 反选 |
| Shift+Num+/Num- | 扩展/收缩选择 |
| Shift+F2 | 比较目录 |
| Shift+F10 | 上下文菜单 |
| Shift+F12 | 命令浏览器 |
| Shift+Tab | 切换树视图焦点 |
| Shift+Enter | 执行/打开（同 Enter 扩展语义） |
| Tab | 切换面板 |
| Enter | 打开/执行/确认重命名 |
| Insert | 选择文件 |
| Backspace | 上级目录 |
| Space | 选择/取消；目录则计算大小 |
| 字母数字键 | 快速搜索（依配置） |
| ←/→ | 依模式移动光标/列 |
| 右键 | 上下文菜单 |

---

## 2. 复制/移动对话框

| 快捷键 | 动作 |
|--------|------|
| F2 | 加入操作队列 |
| F5, F6 | 切换目标字段选择（循环） |

---

## 3–11. 其他上下文

完整快捷键表见 DC 官方文档对应章节。实现各子界面时，从 Phase 3/4 起逐节补全本文件。

- [Copy/Move Dialog](https://doublecmd.github.io/doc/en/shortcuts.html#3-copymove-dialog)
- [Find Files](https://doublecmd.github.io/doc/en/shortcuts.html#5-find-files)
- [Multi-Rename Tool](https://doublecmd.github.io/doc/en/shortcuts.html#6-multi-rename-tool)
- [Internal Viewer](https://doublecmd.github.io/doc/en/shortcuts.html#8-internal-viewer)
- [Internal Editor](https://doublecmd.github.io/doc/en/shortcuts.html#9-internal-editor)

---

## 实现备注

1. 所有快捷键必须可通过 `cm_ConfigHotKeys` 配置（Phase 6）。
2. 同一命令允许多个快捷键（如 F8 与 Del）。
3. 命令可带参数（如 `cm_Delete` + `trashcan=reversesetting`）。
4. 「Only for these controls」限制在 TUI 中映射为焦点区域：文件列表 / 命令行 / 快速搜索框。
