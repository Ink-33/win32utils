# Win32Utils

一个 Go 语言库，用于简化 Windows 系统托盘应用程序开发。提供了对 Windows API 的高级封装，使得创建具有系统托盘图标、菜单、对话框和 Toast 通知的应用程序变得简单直观。

## 特性

- **🎯 系统托盘图标** - 快速创建和管理系统托盘图标
- **📋 菜单支持** - 支持右键菜单，带有 Emoji 图标和自定义回调
- **🎨 DPI 缩放** - 自动处理高 DPI 显示器的缩放
- **📝 文本输入对话框** - 现代的 DPI 感知文本输入对话框
- **🔔 Toast 通知** - Windows Toast 通知，支持 Emoji 图标和自定义消息
- **🚀 高级 API** - 流式构建器 API，简化应用程序创建
- **⚙️ 无 CGO** - 使用 `golang.org/x/sys/windows` 纯 Go 实现，无 CGO 依赖

## 安装

```bash
go get github.com/Ink-33/win32utils
```

或使用 `go.mod` 中的模块：

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
	// 启用高 DPI 支持
	win32utils.ToHighDPI()

	// 创建托盘应用
	app, err := win32utils.NewTrayAppBuilder("com.example.myapp").
		Name("My Application").
		IconID(32516). // IDI_INFORMATION
		IconTip("My App Tray Icon").
		Build()
	if err != nil {
		panic(err)
	}
	defer app.Close()

	// 添加菜单项
	_ = app.AddMenuItemWithEmoji("⚙️", "Settings", func() {
		app.ShowNotificationInfo("Settings", "Opening settings...")
	})

	_ = app.AddMenuItemWithEmoji("👋", "Exit", func() {
		app.Exit()
	})

	// 运行消息循环
	exitCode, err := app.Run()
	if err != nil {
		panic(err)
	}

	fmt.Printf("应用已退出，代码: %d\n", exitCode)
}
```

## 核心概念

### TrayApp - 高级应用抽象

`TrayApp` 是一个高级封装，用于管理托盘应用的完整生命周期。它包括：
- 托盘图标管理
- 菜单管理
- 通知
- 对话框
- 消息循环

### 构建器模式

使用流式 `TrayAppBuilder` API 配置您的应用：

```go
builder := win32utils.NewTrayAppBuilder("appID").
	Name("App Name").
	IconID(32516).
	IconTip("Tooltip").
	OnLeftClick(func() { /* ... */ }).
	OnDoubleClick(func() { /* ... */ })

app, err := builder.Build()
```

### DPI 感知

所有 UI 元素都是 DPI 感知的：

```go
// 启用系统范围的高 DPI 支持
win32utils.ToHighDPI()
```

### 通知系统

支持四种预定义的通知类型，使用 Emoji 图标：

- `ShowNotificationSuccess()` - ✅ 成功通知
- `ShowNotificationWarning()` - ⚠️ 警告通知
- `ShowNotificationError()` - ❌ 错误通知
- `ShowNotificationInfo()` - ℹ️ 信息通知

```go
app.ShowNotificationSuccess("标题", "操作成功！")
app.ShowNotificationError("标题", "发生错误！")
```

### 对话框

创建模态文本输入对话框：

```go
text1, text2, cancelled, err := app.ShowDialog(
	"对话框标题",
	"第一个输入框标签:",
	"第二个输入框标签:",
	"默认值1",
	"默认值2",
)

if err != nil {
	// 处理错误
} else if !cancelled {
	fmt.Printf("输入: %s, %s\n", text1, text2)
}
```

## 项目结构

```
.
├── README.md                 # 项目文档
├── go.mod                    # Go 模块定义
├── LICENSE                   # 许可证
│
├── trayapp.go               # 高级 TrayApp 抽象层
├── trayicon.go              # 托盘图标和菜单管理
├── notification.go          # Toast 通知系统
├── dialog.go                # 文本输入对话框
├── window.go                # 窗口创建和管理
├── dll.go                   # Windows DLL 句柄
├── winbase.go               # Windows 结构和常量
│
├── cmd/
│   └── main.go              # 示例应用程序
│
└── *_test.go                # 单元测试
```

## 主要 API

### TrayApp

- `NewTrayAppBuilder(appID string)` - 创建构建器
- `Build()` - 构建 TrayApp 实例
- `AddMenuItem(label, callback)` - 添加菜单项
- `AddMenuItemWithEmoji(emoji, label, callback)` - 添加带 Emoji 的菜单项
- `AddMenuSeparator()` - 添加分隔符
- `ShowNotificationSuccess/Warning/Error/Info()` - 显示通知
- `ShowDialog()` - 显示对话框
- `Run()` - 启动消息循环（阻塞）
- `Close()` - 关闭应用并清理资源
- `Exit()` - 退出应用

### 构建器选项

- `Name(string)` - 设置应用名称
- `IconID(uint16)` - 设置系统图标 ID
- `IconTip(string)` - 设置托盘图标提示
- `OnLeftClick(callback)` - 左键单击回调
- `OnDoubleClick(callback)` - 双击回调

## 系统要求

- Windows 7 或更高版本
- Go 1.21 或更高版本
- 用于 Toast 通知的 PowerShell

## 示例

### 基本托盘应用

查看 [cmd/main.go](cmd/main.go) 获取完整的示例应用程序。

### 特性演示

示例应用程序展示了以下特性：
- 创建系统托盘图标
- 右键菜单
- 带 Emoji 的菜单项
- Toast 通知（成功、警告、错误、信息）
- 文本输入对话框
- 事件处理回调

## 线程安全

`TrayApp` 是线程安全的。从任何线程调用 `AddMenuItem`, `ShowNotification*` 和 `ShowDialog` 方法是安全的。

## 常见问题

**Q: 可以从其他线程显示通知吗？**

A: 是的，所有通知方法都是线程安全的。

**Q: 如何自定义菜单项的图标？**

A: 使用 `AddMenuItemWithEmoji()` 方法并传递您选择的 Emoji。

**Q: Toast 通知在哪里显示？**

A: Toast 通知显示在 Windows 10/11 的通知中心。需要 PowerShell 支持。

**Q: 可以在应用运行时添加菜单项吗？**

A: 是的，可以在任何时刻调用 `AddMenuItem`。菜单将在下次右键单击时更新。

## 许可证

查看 [LICENSE](LICENSE) 文件了解详情。

## 贡献

欢迎提交 Issue 和 Pull Request。

---

**下一步**: 阅读 [API 文档](docs/API.md) 了解详细的 API 参考。
