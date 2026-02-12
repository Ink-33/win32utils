package main

import (
	"repo.smlk.org/win32utils"
)

func main() {
	win32utils.ToHighDPI()

	app, _ := win32utils.NewTrayAppBuilder("com.example.autoclose").
		Name("自动关闭通知演示").
		IconID(32516).
		Build()
	defer app.Close()

	// 快速关闭（5秒）
	_ = app.AddMenuItemWithEmoji("⚡", "快速成功", func() {
		_ = app.ShowNotificationSuccessEx("完成", "操作成功！", win32utils.DurationShort)
	})

	// 长时关闭（10秒）
	_ = app.AddMenuItemWithEmoji("📌", "长时错误", func() {
		_ = app.ShowNotificationErrorEx("错误", "请立即处理！", win32utils.DurationLong)
	})

	_ = app.AddMenuItemWithEmoji("⚠️", "快速警告", func() {
		_ = app.ShowNotificationWarningEx("警告", "即将超时！", win32utils.DurationShort)
	})

	_ = app.AddMenuItemWithEmoji("ℹ️", "长时信息", func() {
		_ = app.ShowNotificationInfoEx("提示", "重要信息，请注意", win32utils.DurationLong)
	})

	_ = app.AddMenuSeparator()

	_ = app.AddMenuItemWithEmoji("👋", "Exit", func() {
		app.Exit()
	})

	_, _ = app.Run()
}
