# API 文档

完整的 Win32Utils 库 API 参考。

## 目录

1. [TrayApp - 高级 API](#trayapp---高级-api)
2. [TrayIcon & PopupMenu - 中级 API](#trayicon--popupmenu---中级-api)
3. [窗口和 UI 工具函数](#窗口和-ui-工具函数)
4. [通知系统](#通知系统)
5. [对话框](#对话框)
6. [常量和类型](#常量和类型)

---

## TrayApp - 高级 API

`TrayApp` 是使用库的推荐方式。它提供了一个高级抽象，处理了大部分复杂的 Windows API 细节。

### 创建应用

#### `NewTrayAppBuilder(appID string) *TrayAppBuilder`

创建一个新的 TrayApp 构建器。

**参数:**
- `appID`: 应用程序的唯一标识符（通常是反向域名，如 `com.example.myapp`）

**返回:**
- `*TrayAppBuilder`: 可用于配置应用的构建器

**示例:**
```go
builder := win32utils.NewTrayAppBuilder("com.example.myapp")
```

### 构建器方法

构建器使用流式 API 进行配置。所有配置方法都返回 `*TrayAppBuilder` 以支持方法链接。

#### `Name(name string) *TrayAppBuilder`

设置应用程序名称。

**参数:**
- `name`: 应用名称

**返回:**
- `*TrayAppBuilder`: 用于方法链接

```go
builder.Name("My Application")
```

#### `IconID(id uint16) *TrayAppBuilder`

设置要使用的系统图标 ID。

**参数:**
- `id`: 系统图标 ID（如 `32516` 表示 `IDI_INFORMATION`）

**常用图标 ID:**
- `32515`: `IDI_HAND` - 错误/停止
- `32516`: `IDI_QUESTION` - 问题/信息
- `32517`: `IDI_EXCLAMATION` - 警告
- `32516`: `IDI_INFORMATION` - 信息

**返回:**
- `*TrayAppBuilder`: 用于方法链接

```go
builder.IconID(32516)
```

#### `IconTip(tip string) *TrayAppBuilder`

设置托盘图标的提示文本（悬停时显示）。

**参数:**
- `tip`: 提示文本（最多 128 个字符）

**返回:**
- `*TrayAppBuilder`: 用于方法链接

```go
builder.IconTip("My Application")
```

#### `OnLeftClick(callback func()) *TrayAppBuilder`

设置左键单击回调。

**参数:**
- `callback`: 单击托盘图标左键时调用的函数

**返回:**
- `*TrayAppBuilder`: 用于方法链接

```go
builder.OnLeftClick(func() {
    fmt.Println("Left clicked!")
})
```

#### `OnDoubleClick(callback func()) *TrayAppBuilder`

设置双击回调。

**参数:**
- `callback`: 双击托盘图标时调用的函数

**返回:**
- `*TrayAppBuilder`: 用于方法链接

```go
builder.OnDoubleClick(func() {
    fmt.Println("Double clicked!")
})
```

#### `Build() (*TrayApp, error)`

构建 TrayApp 实例。

**返回:**
- `*TrayApp`: 构建的应用实例
- `error`: 如果构建失败则返回错误

**示例:**
```go
app, err := builder.Build()
if err != nil {
    panic(err)
}
defer app.Close()
```

### TrayApp 方法

#### `AddMenuItem(label string, onClick func()) error`

添加一个没有图标的菜单项。

**参数:**
- `label`: 菜单项的显示文本
- `onClick`: 单击菜单项时调用的回调函数

**返回:**
- `error`: 如果添加失败则返回错误

**示例:**
```go
err := app.AddMenuItem("点击我", func() {
    fmt.Println("Menu item clicked!")
})
```

#### `AddMenuItemWithEmoji(emoji string, label string, onClick func()) error`

添加一个带 Emoji 图标的菜单项。

**参数:**
- `emoji`: 要显示在菜单项前的 Emoji 字符
- `label`: 菜单项的显示文本
- `onClick`: 单击菜单项时调用的回调函数

**返回:**
- `error`: 如果添加失败则返回错误

**示例:**
```go
err := app.AddMenuItemWithEmoji("⚙️", "设置", func() {
    fmt.Println("Settings clicked!")
})
```

**常用 Emoji:**
- `✅` - 成功/完成
- `⚠️` - 警告
- `❌` - 错误/取消
- `ℹ️` - 信息
- `⚙️` - 设置
- `📋` - 显示/列表
- `👋` - 退出/再见
- `💾` - 保存
- `🔄` - 刷新
- `📁` - 文件/文件夹

#### `AddMenuSeparator() error`

添加一个菜单分隔符（水平线）。

**返回:**
- `error`: 如果添加失败则返回错误

**示例:**
```go
err := app.AddMenuSeparator()
```

### 通知方法

所有通知方法都是线程安全的，可以从任何线程调用。

通知会在一定时间后自动关闭：
- **短时长** (`DurationShort`) - 约 5 秒后自动关闭
- **长时长** (`DurationLong`) - 约 10 秒后自动关闭（默认）

#### `ShowNotificationSuccess(title string, message string) error`

显示成功通知（✅ 图标），使用默认时长（~10秒）。

**参数:**
- `title`: 通知标题
- `message`: 通知消息

**返回:**
- `error`: 如果显示失败则返回错误

```go
err := app.ShowNotificationSuccess("完成", "操作成功！")
```

#### `ShowNotificationSuccessEx(title string, message string, duration NotificationDuration) error`

显示成功通知（✅ 图标），支持自定义自动关闭时长。

**参数:**
- `title`: 通知标题
- `message`: 通知消息
- `duration`: 通知持续时长 - `DurationShort` (~5秒) 或 `DurationLong` (~10秒)

**返回:**
- `error`: 如果显示失败则返回错误

```go
// 快速关闭（5秒）
err := app.ShowNotificationSuccessEx("完成", "操作成功！", win32utils.DurationShort)

// 长时长（10秒）
err := app.ShowNotificationSuccessEx("完成", "操作成功！", win32utils.DurationLong)
```

#### `ShowNotificationWarning(title string, message string) error`

显示警告通知（⚠️ 图标），使用默认时长。

**参数:**
- `title`: 通知标题
- `message`: 通知消息

**返回:**
- `error`: 如果显示失败则返回错误

```go
err := app.ShowNotificationWarning("警告", "请检查您的输入")
```

#### `ShowNotificationWarningEx(title string, message string, duration NotificationDuration) error`

显示警告通知（⚠️ 图标），支持自定义自动关闭时长。

**参数:**
- `title`: 通知标题
- `message`: 通知消息
- `duration`: 通知持续时长 - `DurationShort` 或 `DurationLong`

```go
// 快速关闭的警告
err := app.ShowNotificationWarningEx("警告", "即将超时", win32utils.DurationShort)
```

#### `ShowNotificationError(title string, message string) error`

显示错误通知（❌ 图标），使用默认时长。

**参数:**
- `title`: 通知标题
- `message`: 通知消息

**返回:**
- `error`: 如果显示失败则返回错误

```go
err := app.ShowNotificationError("错误", "发生了错误，请重试")
```

#### `ShowNotificationErrorEx(title string, message string, duration NotificationDuration) error`

显示错误通知（❌ 图标），支持自定义自动关闭时长。

**参数:**
- `title`: 通知标题
- `message`: 通知消息
- `duration`: 通知持续时长 - `DurationShort` 或 `DurationLong`

```go
// 长时长错误提示
err := app.ShowNotificationErrorEx("发生错误", "请联系管理员", win32utils.DurationLong)
```

#### `ShowNotificationInfo(title string, message string) error`

显示信息通知（ℹ️ 图标），使用默认时长。

**参数:**
- `title`: 通知标题
- `message`: 通知消息

**返回:**
- `error`: 如果显示失败则返回错误

```go
err := app.ShowNotificationInfo("信息", "这是一条信息消息")
```

#### `ShowNotificationInfoEx(title string, message string, duration NotificationDuration) error`

显示信息通知（ℹ️ 图标），支持自定义自动关闭时长。

**参数:**
- `title`: 通知标题
- `message`: 通知消息
- `duration`: 通知持续时长 - `DurationShort` 或 `DurationLong`

```go
// 快速关闭的提示
err := app.ShowNotificationInfoEx("状态", "已就绪", win32utils.DurationShort)
```

### 对话框方法

#### `ShowDialog(title string, label1 string, label2 string, default1 string, default2 string) (string, string, bool, error)`

显示文本输入对话框。

**参数:**
- `title`: 对话框标题
- `label1`: 第一个输入框的标签
- `label2`: 第二个输入框的标签
- `default1`: 第一个输入框的默认值
- `default2`: 第二个输入框的默认值

**返回:**
- `string`: 第一个输入框的值
- `string`: 第二个输入框的值
- `bool`: 是否被取消（true = 取消，false = 确定）
- `error`: 如果显示失败则返回错误

**示例:**
```go
text1, text2, cancelled, err := app.ShowDialog(
    "输入信息",
    "用户名:",
    "密码:",
    "默认用户",
    "",
)

if err != nil {
    fmt.Printf("对话框错误: %v\n", err)
} else if cancelled {
    fmt.Println("用户取消了对话框")
} else {
    fmt.Printf("用户名: %s, 密码: %s\n", text1, text2)
}
```

### 生命周期方法

#### `Run() (int32, error)`

启动消息循环。此方法会阻塞，直到应用退出。

**返回:**
- `int32`: 退出代码
- `error`: 如果发生错误则返回错误

**示例:**
```go
exitCode, err := app.Run()
if err != nil {
    fmt.Printf("消息循环错误: %v\n", err)
}
fmt.Printf("应用已退出，代码: %d\n", exitCode)
```

#### `Close() error`

关闭应用并清理资源。应该在 `defer` 中调用。

**返回:**
- `error`: 如果清理失败则返回错误

**示例:**
```go
app, _ := builder.Build()
defer app.Close()
```

#### `Exit()`

从消息循环中退出应用。可以在回调中调用。

**示例:**
```go
err := app.AddMenuItemWithEmoji("👋", "退出", func() {
    app.Exit()
})
```

---

## TrayIcon & PopupMenu - 中级 API

如果需要更低级的控制，可以直接使用 `TrayIcon` 和 `PopupMenu`。

### TrayIcon

#### `NewTrayIcon() *TrayIcon`

创建新的托盘图标。

```go
tray := win32utils.NewTrayIcon()
```

#### `(ti *TrayIcon) Add(icon windows.Handle, tooltip string) error`

添加托盘图标到系统托盘。

**参数:**
- `icon`: 图标句柄
- `tooltip`: 图标提示文本

```go
hIcon := // ... 加载图标
err := tray.Add(hIcon, "我的应用")
```

#### `(ti *TrayIcon) Update(icon windows.Handle, tooltip string) error`

更新现有的托盘图标。

```go
err := tray.Update(hIcon, "新的提示")
```

#### `(ti *TrayIcon) Remove() error`

从系统托盘移除图标。

```go
err := tray.Remove()
```

#### `(ti *TrayIcon) ShowMenu(x, y int32, menu *PopupMenu) error`

在指定位置显示弹出菜单。

```go
err := tray.ShowMenu(100, 100, menu)
```

### PopupMenu

#### `NewPopupMenu() *PopupMenu`

创建新的弹出菜单。

```go
menu := win32utils.NewPopupMenu()
```

#### `(pm *PopupMenu) Append(id uint32, text string, callback func()) error`

添加菜单项。

**参数:**
- `id`: 菜单项 ID
- `text`: 菜单项文本
- `callback`: 单击时的回调

```go
err := menu.Append(1, "选项 1", func() {
    fmt.Println("选项 1 被点击")
})
```

#### `(pm *PopupMenu) AppendSeparator() error`

添加分隔符。

```go
err := menu.AppendSeparator()
```

#### `(pm *PopupMenu) Clear() error`

清除菜单中的所有项目。

```go
err := menu.Clear()
```

#### `(pm *PopupMenu) Destroy() error`

销毁菜单并释放资源。

```go
err := menu.Destroy()
```

---

## 窗口和 UI 工具函数

### 初始化

#### `ToHighDPI()`

启用应用程序的高 DPI 支持。应该在 `main()` 的开始处调用。

```go
func main() {
    win32utils.ToHighDPI()
    // ... 其余代码
}
```

### 窗口创建（高级用户）

#### `CreateMessageOnlyWindow(className string, wndProc func(hwnd windows.Handle, msg uint32, w uintptr, l uintptr) uintptr) (windows.Handle, error)`

创建仅用于消息的窗口。

**参数:**
- `className`: 窗口类名
- `wndProc`: 窗口过程回调

**返回:**
- `windows.Handle`: 窗口句柄
- `error`: 如果创建失败则返回错误

```go
hwnd, err := win32utils.CreateMessageOnlyWindow("MyWindowClass", func(hwnd windows.Handle, msg uint32, w uintptr, l uintptr) uintptr {
    // 处理消息
    return 0
})
```

### DPI 相关函数

#### `GetDPIScaleFactor() float32`

获取系统DPI缩放因子。

```go
scale := win32utils.GetDPIScaleFactor()
scaledWidth := int32(float32(width) * scale)
```

### 消息循环

#### `MessageLoop() (int32, error)`

启动标准消息循环。阻塞直到 `PostQuitMessage` 被调用。

**返回:**
- `int32`: 退出代码
- `error`: 如果发生错误则返回错误

```go
exitCode, err := win32utils.MessageLoop()
```

---

## 通知系统

通知通过 PowerShell 和 Windows Runtime (WinRT) 实现。

### SimpleNotification

#### `ShowSimpleNotification(appID string, title string, message string, icon string) error`

显示简单通知。

**参数:**
- `appID`: 应用程序 ID
- `title`: 通知标题
- `message`: 通知消息
- `icon`: 图标（支持 Emoji、文件路径或系统图标）

```go
err := win32utils.ShowSimpleNotification(
    "com.example.app",
    "标题",
    "这是消息",
    "✅",
)
```

### AdvancedNotification

如需自定义 Toast 通知（按钮、图像、特殊布局等），请参考 `notification.go` 中的 `AdvancedNotificationBuilder`。

---

## 对话框

### TwoTextInputDialog

#### `TwoTextInputDialog(title string, label1 string, label2 string, default1 string, default2 string) (string, string, bool, error)`

显示两个文本输入字段的对话框。

**参数:**
- `title`: 对话框标题
- `label1`: 第一个输入框的标签
- `label2`: 第二个输入框的标签
- `default1`: 第一个输入框的默认值
- `default2`: 第二个输入框的默认值

**返回:**
- `string`: 第一个文本框的输入
- `string`: 第二个文本框的输入
- `bool`: 是否被取消
- `error`: 如果出错则返回错误

```go
text1, text2, cancelled, err := win32utils.TwoTextInputDialog(
    "输入数据",
    "名称:",
    "年龄:",
    "John",
    "25",
)

if !cancelled && err == nil {
    fmt.Printf("姓名: %s, 年龄: %s\n", text1, text2)
}
```

---

## 常量和类型

### 通知持续时长

```go
type NotificationDuration string

const (
    DurationShort NotificationDuration = "short"  // ~5秒后自动关闭
    DurationLong  NotificationDuration = "long"   // ~10秒后自动关闭（默认）
)
```

**使用示例:**
```go
// 快速关闭（5秒）
err := app.ShowNotificationSuccessEx("完成", "操作成功！", win32utils.DurationShort)

// 长时长（10秒）
err := app.ShowNotificationErrorEx("错误", "请检查输入", win32utils.DurationLong)
```

### 通知状态 Emoji

```
✅ - 成功
⚠️ - 警告
❌ - 错误
ℹ️ - 信息
```

### 常用系统图标 ID

```go
const (
    IDI_HAND        = 32515 // 错误/停止
    IDI_QUESTION    = 32514 // 问题
    IDI_EXCLAMATION = 32517 // 警告
    IDI_INFORMATION = 32516 // 信息
)
```

### NOTIFYICONDATAW 结构

用于与 `Shell_NotifyIconW` 交互的结构体（通常不需要直接使用）。

---

## 错误处理

所有方法都返回 `error` 类型的错误。始终检查和处理错误：

```go
if err := app.AddMenuItem("选项", func() {}); err != nil {
    fmt.Printf("添加菜单失败: %v\n", err)
}
```

## 线程安全性

`TrayApp` 的大多数方法都是线程安全的：
- ✅ `AddMenuItem*` - 线程安全
- ✅ `ShowNotification*` - 线程安全
- ✅ `ShowDialog` - 线程安全
- ✅ `Exit` - 线程安全
- ⚠️ `Run` - 应该从主线程调用

可以从后台线程安全地调用所有 UI 更新方法。

---

[返回到 README](../README.md)
