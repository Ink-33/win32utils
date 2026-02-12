# 使用示例

本文档包含了 Win32Utils 库的各种使用示例。

## 目录

1. [最小示例](#最小示例)
2. [基本托盘应用](#基本托盘应用)
3. [完整功能示例](#完整功能示例)
4. [菜单和通知](#菜单和通知)
5. [对话框操作](#对话框操作)
6. [后台操作](#后台操作)
7. [错误处理](#错误处理)

---

## 最小示例

最简单的 Win32Utils 应用程序：

```go
package main

import (
	"repo.smlk.org/win32utils"
)

func main() {
	win32utils.ToHighDPI()

	app, _ := win32utils.NewTrayAppBuilder("com.example.app").
		Name("My App").
		Build()
	defer app.Close()

	app.Run()
}
```

这将创建一个系统托盘图标，当用户右键单击时显示一个空菜单。

---

## 基本托盘应用

一个完整的基本托盘应用程序：

```go
package main

import (
	"fmt"
	"repo.smlk.org/win32utils"
)

func main() {
	// 启用高 DPI 支持
	win32utils.ToHighDPI()

	// 创建应用
	app, err := win32utils.NewTrayAppBuilder("com.example.basic").
		Name("Basic Tray App").
		IconID(32516). // 使用信息图标
		IconTip("Basic App").
		OnLeftClick(func() {
			fmt.Println("Left clicked on tray icon!")
		}).
		OnDoubleClick(func() {
			fmt.Println("Double clicked on tray icon!")
		}).
		Build()

	if err != nil {
		panic(fmt.Errorf("创建应用失败: %v", err))
	}
	defer app.Close()

	// 添加菜单项
	app.AddMenuItemWithEmoji("ℹ️", "About", func() {
		app.ShowNotificationInfo("About", "Basic Tray App v1.0")
	})

	app.AddMenuSeparator()

	app.AddMenuItemWithEmoji("👋", "Exit", func() {
		app.Exit()
	})

	fmt.Println("应用已启动...")
	exitCode, _ := app.Run()
	fmt.Printf("应用已退出，代码: %d\n", exitCode)
}
```

---

## 完整功能示例

展示库的所有主要功能的完整示例：

```go
package main

import (
	"fmt"
	"repo.smlk.org/win32utils"
)

func main() {
	win32utils.ToHighDPI()

	// 构建应用
	builder := win32utils.NewTrayAppBuilder("com.example.complete").
		Name("完整功能演示").
		IconID(32516).
		IconTip("点击右键查看菜单")

	app, err := builder.Build()
	if err != nil {
		panic(err)
	}
	defer app.Close()

	// 通知演示部分
	app.AddMenuItemWithEmoji("🔔", "显示通知", func() {
		fmt.Println("显示各种通知类型...")
		app.ShowNotificationInfo("信息", "这是一条信息通知")
	})

	// 成功通知
	app.AddMenuItemWithEmoji("✅", "成功通知", func() {
		app.ShowNotificationSuccess("成功", "操作已成功完成！")
	})

	// 警告通知
	app.AddMenuItemWithEmoji("⚠️", "警告通知", func() {
		app.ShowNotificationWarning("警告", "请小心处理！")
	})

	// 错误通知
	app.AddMenuItemWithEmoji("❌", "错误通知", func() {
		app.ShowNotificationError("错误", "发生了什么问题！")
	})

	app.AddMenuSeparator()

	// 对话框演示
	app.AddMenuItemWithEmoji("📝", "用户输入对话框", func() {
		text1, text2, cancelled, err := app.ShowDialog(
			"用户信息",
			"输入您的名字:",
			"输入您的职位:",
			"张三",
			"工程师",
		)

		if err != nil {
			app.ShowNotificationError("对话框错误", err.Error())
			return
		}

		if cancelled {
			app.ShowNotificationWarning("已取消", "对话框被取消")
		} else {
			message := fmt.Sprintf("感谢提交！\n名字: %s\n职位: %s", text1, text2)
			app.ShowNotificationSuccess("提交成功", message)
		}
	})

	app.AddMenuSeparator()

	// 设置演示
	app.AddMenuItemWithEmoji("⚙️", "设置", func() {
		width, height, cancelled, err := app.ShowDialog(
			"应用设置",
			"窗口宽度:",
			"窗口高度:",
			"800",
			"600",
		)

		if !cancelled && err == nil {
			app.ShowNotificationSuccess(
				"设置已保存",
				fmt.Sprintf("分辨率: %sx%s", width, height),
			)
		}
	})

	app.AddMenuSeparator()

	// 退出
	app.AddMenuItemWithEmoji("👋", "退出", func() {
		app.ShowNotificationInfo("再见", "应用正在关闭...")
		app.Exit()
	})

	fmt.Println("应用已启动，右键单击托盘图标查看菜单")
	exitCode, err := app.Run()
	if err != nil {
		fmt.Printf("错误: %v\n", err)
	}
	fmt.Printf("应用已退出，代码: %d\n", exitCode)
}
```

---

## 菜单和通知

### 动态菜单管理

```go
package main

import (
	"fmt"
	"repo.smlk.org/win32utils"
)

func main() {
	win32utils.ToHighDPI()

	app, _ := win32utils.NewTrayAppBuilder("com.example.menu").
		Name("菜单演示").
		Build()
	defer app.Close()

	// 计数器变量
	var count int = 0

	// 添加计数器菜单项
	app.AddMenuItemWithEmoji("🔢", "增加计数", func() {
		count++
		msg := fmt.Sprintf("计数: %d", count)
		app.ShowNotificationInfo("计数更新", msg)
	})

	// 仅当计数 > 0 时重置
	app.AddMenuItemWithEmoji("🔄", "重置计数", func() {
		if count > 0 {
			count = 0
			app.ShowNotificationSuccess("已重置", "计数已被重置为 0")
		} else {
			app.ShowNotificationInfo("信息", "计数已经是 0")
		}
	})

	app.AddMenuSeparator()

	// 显示当前计数
	app.AddMenuItemWithEmoji("📊", "显示计数", func() {
		msg := fmt.Sprintf("当前计数: %d", count)
		app.ShowNotificationInfo("当前值", msg)
	})

	app.AddMenuSeparator()

	app.AddMenuItemWithEmoji("👋", "Exit", func() {
		app.Exit()
	})

	app.Run()
}
```

### 不同的菜单项类型

```go
package main

import (
	"repo.smlk.org/win32utils"
)

func main() {
	win32utils.ToHighDPI()

	app, _ := win32utils.NewTrayAppBuilder("com.example.menutype").
		Name("菜单类型演示").
		Build()
	defer app.Close()

	// 不同功能类别的菜单项
	
	// 文件操作
	app.AddMenuItemWithEmoji("📁", "打开文件夹", func() {
		app.ShowNotificationInfo("文件操作", "打开文件夹...")
	})

	app.AddMenuItemWithEmoji("💾", "保存", func() {
		app.ShowNotificationSuccess("保存完成", "文件已保存")
	})

	app.AddMenuSeparator()

	// 编辑操作
	app.AddMenuItemWithEmoji("✏️", "编辑", func() {
		app.ShowNotificationInfo("编辑", "打开编辑对话框...")
	})

	app.AddMenuItemWithEmoji("🗑️", "删除", func() {
		app.ShowNotificationWarning("删除", "确认删除?")
	})

	app.AddMenuSeparator()

	// 视图选项
	app.AddMenuItemWithEmoji("🔍", "放大", func() {
		app.ShowNotificationInfo("放大", "放大 50%")
	})

	app.AddMenuItemWithEmoji("🔍", "缩小", func() {
		app.ShowNotificationInfo("缩小", "缩小 50%")
	})

	app.AddMenuSeparator()

	// 工具
	app.AddMenuItemWithEmoji("🔧", "工具", func() {
		app.ShowNotificationInfo("工具", "打开工具面板...")
	})

	app.AddMenuSeparator()

	// 帮助和退出
	app.AddMenuItemWithEmoji("❓", "帮助", func() {
		app.ShowNotificationInfo("帮助", "访问文档...")
	})

	app.AddMenuItemWithEmoji("👋", "退出", func() {
		app.Exit()
	})

	app.Run()
}
```

---

## 对话框操作

### 简单输入对话框

```go
package main

import (
	"fmt"
	"repo.smlk.org/win32utils"
)

func main() {
	win32utils.ToHighDPI()

	app, _ := win32utils.NewTrayAppBuilder("com.example.dialog").
		Name("对话框演示").
		Build()
	defer app.Close()

	app.AddMenuItemWithEmoji("📋", "输入数据", func() {
		// 显示对话框
		field1, field2, cancelled, err := app.ShowDialog(
			"数据输入",
			"第一个字段:",
			"第二个字段:",
			"",
			"",
		)

		if err != nil {
			app.ShowNotificationError("错误", fmt.Sprintf("对话框错误: %v", err))
			return
		}

		if cancelled {
			app.ShowNotificationWarning("已取消", "用户取消了操作")
		} else {
			// 处理输入
			result := fmt.Sprintf("字段1: %s\n字段2: %s", field1, field2)
			app.ShowNotificationSuccess("输入已接收", result)
		}
	})

	app.AddMenuSeparator()

	app.AddMenuItemWithEmoji("👋", "Exit", func() {
		app.Exit()
	})

	app.Run()
}
```

### 带默认值的对话框

```go
package main

import (
	"fmt"
	"repo.smlk.org/win32utils"
)

func main() {
	win32utils.ToHighDPI()

	app, _ := win32utils.NewTrayAppBuilder("com.example.defaults").
		Name("对话框默认值演示").
		Build()
	defer app.Close()

	// 模拟用户设置
	var currentSettings = map[string]string{
		"username": "user@example.com",
		"timeout":  "30",
	}

	app.AddMenuItemWithEmoji("⚙️", "编辑设置", func() {
		username, timeout, cancelled, _ := app.ShowDialog(
			"应用配置",
			"用户名/邮箱:",
			"超时时间(秒):",
			currentSettings["username"],
			currentSettings["timeout"],
		)

		if !cancelled {
			currentSettings["username"] = username
			currentSettings["timeout"] = timeout
			msg := fmt.Sprintf("用户名: %s\n超时: %s 秒", username, timeout)
			app.ShowNotificationSuccess("设置已保存", msg)
		}
	})

	app.AddMenuSeparator()

	app.AddMenuItemWithEmoji("📊", "查看设置", func() {
		msg := fmt.Sprintf("用户名: %s\n超时: %s 秒",
			currentSettings["username"],
			currentSettings["timeout"])
		app.ShowNotificationInfo("当前设置", msg)
	})

	app.AddMenuSeparator()

	app.AddMenuItemWithEmoji("👋", "Exit", func() {
		app.Exit()
	})

	app.Run()
}
```

---

## 后台操作

### 从菜单回调启动后台任务

```go
package main

import (
	"fmt"
	"sync"
	"time"
	"repo.smlk.org/win32utils"
)

func main() {
	win32utils.ToHighDPI()

	app, _ := win32utils.NewTrayAppBuilder("com.example.background").
		Name("后台任务演示").
		Build()
	defer app.Close()

	var (
		taskRunning = false
		mu          sync.Mutex
	)

	// 启动后台任务
	app.AddMenuItemWithEmoji("▶️", "启动任务", func() {
		mu.Lock()
		if taskRunning {
			mu.Unlock()
			app.ShowNotificationWarning("警告", "任务已在运行")
			return
		}
		taskRunning = true
		mu.Unlock()

		go func() {
			app.ShowNotificationInfo("任务已启动", "后台任务正在处理...")

			// 模拟长时间运行的任务
			time.Sleep(2 * time.Second)

			app.ShowNotificationSuccess("任务完成", "后台任务已完成!")

			mu.Lock()
			taskRunning = false
			mu.Unlock()
		}()
	})

	// 停止任务
	app.AddMenuItemWithEmoji("⏹️", "停止任务", func() {
		mu.Lock()
		defer mu.Unlock()

		if !taskRunning {
			app.ShowNotificationInfo("信息", "当前没有运行的任务")
		} else {
			taskRunning = false
			app.ShowNotificationWarning("已停止", "任务已停止")
		}
	})

	// 任务状态
	app.AddMenuItemWithEmoji("📊", "任务状态", func() {
		mu.Lock()
		status := "未运行"
		if taskRunning {
			status = "运行中..."
		}
		mu.Unlock()

		app.ShowNotificationInfo("任务状态", status)
	})

	app.AddMenuSeparator()

	app.AddMenuItemWithEmoji("👋", "Exit", func() {
		app.Exit()
	})

	app.Run()
}
```

---

## 错误处理

### 完整的错误处理示例

```go
package main

import (
	"fmt"
	"repo.smlk.org/win32utils"
)

func main() {
	win32utils.ToHighDPI()

	// 错误处理：应用创建
	builder := win32utils.NewTrayAppBuilder("com.example.errors")
	app, err := builder.
		Name("错误处理演示").
		Build()

	if err != nil {
		fmt.Printf("创建应用失败: %v\n", err)
		return
	}
	defer app.Close()

	// 错误处理：菜单操作
	if err := app.AddMenuItemWithEmoji("⚙️", "配置", func() {
		name, email, cancelled, err := app.ShowDialog(
			"用户信息",
			"姓名:",
			"邮箱:",
			"",
			"",
		)

		// 检查对话框错误
		if err != nil {
			fmt.Printf("对话框错误: %v\n", err)
			_ = app.ShowNotificationError("对话框错误", fmt.Sprintf("错误: %v", err))
			return
		}

		// 检查取消
		if cancelled {
			fmt.Println("对话框被取消")
			return
		}

		// 验证输入
		if name == "" || email == "" {
			_ = app.ShowNotificationWarning(
				"无效输入",
				"请填写所有字段",
			)
			return
		}

		// 处理输入
		fmt.Printf("名字: %s, 邮箱: %s\n", name, email)
		_ = app.ShowNotificationSuccess(
			"已保存",
			fmt.Sprintf("用户: %s <%s>", name, email),
		)
	}); err != nil {
		fmt.Printf("添加菜单项失败: %v\n", err)
		return
	}

	// 错误处理：通知
	if err := app.ShowNotificationInfo("应用已启动", "等待用户交互..."); err != nil {
		fmt.Printf("显示通知失败: %v\n", err)
	}

	// 错误处理：消息循环
	exitCode, err := app.Run()
	if err != nil {
		fmt.Printf("消息循环错误: %v\n", err)
	}

	fmt.Printf("应用已退出，代码: %d\n", exitCode)
}
```

---

## 常见模式

### 切换状态

```go
var isEnabled = true

app.AddMenuItemWithEmoji("🔘", "切换功能", func() {
	isEnabled = !isEnabled
	status := "禁用"
	if isEnabled {
		status = "启用"
	}
	app.ShowNotificationInfo("功能状态", status)
})
```

### 计数器

```go
var counter = 0

app.AddMenuItemWithEmoji("➕", "增加", func() {
	counter++
	app.ShowNotificationInfo("计数", fmt.Sprintf("值: %d", counter))
})

app.AddMenuItemWithEmoji("➖", "减少", func() {
	counter--
	app.ShowNotificationInfo("计数", fmt.Sprintf("值: %d", counter))
})
```

### 条件菜单

```go
app.AddMenuItemWithEmoji("🔓", "操作", func() {
	if !isAuthorized() {
		app.ShowNotificationWarning("未授权", "您没有权限执行此操作")
		return
	}
	
	performAction()
	app.ShowNotificationSuccess("完成", "操作已执行")
})
```

---

[返回到 README](../README.md) | [查看 API 文档](API.md)
