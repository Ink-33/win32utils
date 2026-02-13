# Win32Utils

一个面向 Windows 的纯 Go Win32 工具库，聚焦系统托盘应用开发与常见桌面能力封装（通知、对话框、消息框、剪贴板、控制台、窗口消息循环）。

## 特性

- 托盘应用：`TrayApp` + 构建器 API
- 托盘菜单：支持普通菜单项、分隔符、Emoji 前缀
- 通知：Toast 通知（含 `short/long` 自动关闭时长）
- 对话框：双文本输入、用户名密码输入（密码掩码）
- 消息框：`MessageBoxW` 与常用 `MB_*` / `ID*` 常量
- 控制台管理：显示/隐藏、标题读写、可见性检测
- 剪贴板：文本写入与读取
- DPI 与窗口工具：高 DPI、消息循环及 Win32 辅助 API
- 纯 Go 实现：基于 `golang.org/x/sys/windows`，无 CGO

## 安装

项目当前模块路径：

```bash
go get repo.smlk.org/win32utils
```

导入方式：

```go
import "repo.smlk.org/win32utils"
```

## 快速开始

```go
package main

import (
	"fmt"

	"repo.smlk.org/win32utils"
)

func main() {
	win32utils.ToHighDPIEx()

	builder := win32utils.NewTrayAppBuilder("com.example.myapp").
		Name("My Application").
		IconID(32516).
		IconTip("My App Tray Icon")

	app, err := builder.Build()
	if err != nil {
		panic(err)
	}
	defer app.Close()

	_ = app.AddMenuItemWithEmoji("⚙️", "Settings", func() {
		_ = app.ShowNotificationInfo("Settings", "Opening settings...")
	})
	_ = app.AddMenuSeparator()
	_ = app.AddMenuItemWithEmoji("👋", "Exit", func() {
		app.Exit()
	})

	exitCode, err := app.Run()
	if err != nil {
		panic(err)
	}

	fmt.Printf("应用退出码: %d\n", exitCode)
}
```

## 常用 API 概览

### TrayApp

- `NewTrayAppBuilder(appID string)`
- `(*TrayAppBuilder).Name/IconID/IconTip/OnLeftClick/OnDoubleClick`
- `(*TrayAppBuilder).AddMenuItem/AddMenuItemWithEmoji/AddMenuSeparator/Build`
- `(*TrayApp).AddMenuItem/AddMenuItemWithEmoji/AddMenuSeparator`
- `(*TrayApp).ShowNotification*`（含 `Ex` 版本）
- `(*TrayApp).ShowDialog/ShowUsernamePasswordDialog`
- `(*TrayApp).MessageBoxW/Run/Exit/Close`

### 通知与对话框

- `SimpleToast`, `NewToastBuilder`, `NewAdvancedToastBuilder`
- `NotifySuccess`, `NotifyWarning`, `NotifyError`, `NotifyProgress`
- `TwoTextInputDialog`, `UsernamePasswordDialog`

### 控制台与消息框

- `GetConsoleWindow`, `ShowConsole`, `HideConsole`, `ToggleConsole`
- `IsConsoleVisible`, `GetConsoleTitle`, `SetConsoleTitle`, `ShowWindow`
- `MessageBoxW`, `RunningByDoubleClick`

### 剪贴板

- `SetText`, `SetClipboardText`, `GetClipboardDataText`
- `OpenClipboard`, `CloseClipboard`, `EmptyClipboard`

## 文档与示例

- API 文档：[`docs/API.md`](docs/API.md)
- 控制台指南：[`docs/console_guide.md`](docs/console_guide.md)
- 示例集合：[`docs/EXAMPLES.md`](docs/EXAMPLES.md)
- 演示入口：[`cmd/main.go`](cmd/main.go)
- 示例代码：[`examples/console/console_demo.go`](examples/console/console_demo.go)、[`examples/autoclose/autoclose_demo.go`](examples/autoclose/autoclose_demo.go)

## 项目结构

```text
.
├── cmd/                    # 演示入口
├── docs/                   # 使用文档与示例说明
├── examples/               # 示例代码
├── trayapp.go              # 高级托盘应用封装
├── trayicon.go             # 托盘图标与菜单处理
├── notification.go         # Toast 通知相关
├── dialog.go               # 输入对话框
├── messagebox.go           # Win32 MessageBox + DPI 辅助
├── clipboard.go            # 剪贴板封装
├── console.go              # 控制台窗口管理
├── window.go               # 窗口与消息循环基础能力
└── *_test.go               # 测试
```

## 系统要求

- Windows（`//go:build windows`）
- Go 1.21+
- PowerShell（用于 Toast 通知）

## 测试

```bash
go test ./...
```

## 许可证

查看 [`LICENSE`](LICENSE)。

## 贡献

欢迎提交 Issue 和 Pull Request。
