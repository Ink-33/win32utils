package main

import (
	"fmt"
	"log"
	"time"

	"repo.smlk.org/win32utils"
)

func main() {
	fmt.Println("=== Win32Utils Console Management Demo ===")
	
	// 获取当前控制台状态
	consoleHwnd, err := win32utils.GetConsoleWindow()
	if err != nil {
		log.Fatalf("Failed to get console window: %v", err)
	}
	fmt.Printf("Console window handle: %v\n", consoleHwnd)
	
	visible, err := win32utils.IsConsoleVisible()
	if err != nil {
		log.Fatalf("Failed to check console visibility: %v", err)
	}
	fmt.Printf("Console initially %s\n", map[bool]string{true: "visible", false: "hidden"}[visible])
	
	// 获取当前标题
	currentTitle, err := win32utils.GetConsoleTitle()
	if err != nil {
		log.Printf("Warning: Failed to get console title: %v", err)
		currentTitle = "Unknown"
	}
	fmt.Printf("Current console title: '%s'\n", currentTitle)
	
	// 演示标题更改
	fmt.Println("\n--- Demonstrating console title management ---")
	newTitle := "Win32Utils Console Demo"
	fmt.Printf("Setting console title to: '%s'\n", newTitle)
	err = win32utils.SetConsoleTitle(newTitle)
	if err != nil {
		log.Printf("Failed to set console title: %v", err)
	} else {
		// 验证标题更改
		updatedTitle, err := win32utils.GetConsoleTitle()
		if err != nil {
			log.Printf("Failed to verify title change: %v", err)
		} else {
			fmt.Printf("Title successfully updated to: '%s'\n", updatedTitle)
		}
	}
	
	// 演示显示/隐藏功能
	fmt.Println("\n--- Demonstrating show/hide functionality ---")
	fmt.Println("The console will be hidden for 3 seconds...")
	
	time.Sleep(2 * time.Second)
	
	// 隐藏控制台
	fmt.Println("Hiding console...")
	err = win32utils.HideConsole()
	if err != nil {
		log.Printf("Failed to hide console: %v", err)
	} else {
		fmt.Println("Console hidden! (You won't see this message)")
	}
	
	// 等待3秒
	time.Sleep(3 * time.Second)
	
	// 显示控制台
	fmt.Println("Showing console...")
	err = win32utils.ShowConsole()
	if err != nil {
		log.Printf("Failed to show console: %v", err)
	} else {
		fmt.Println("Console is now visible again!")
	}
	
	// 演示切换功能
	fmt.Println("\n--- Demonstrating toggle functionality ---")
	for i := 0; i < 3; i++ {
		fmt.Printf("Toggle #%d: ", i+1)
		newState, err := win32utils.ToggleConsole()
		if err != nil {
			log.Printf("Failed to toggle console: %v", err)
			break
		}
		
		stateStr := map[bool]string{true: "visible", false: "hidden"}[newState]
		fmt.Printf("Console is now %s\n", stateStr)
		
		time.Sleep(1 * time.Second)
	}
	
	// 最后确保控制台是可见的
	fmt.Println("\n--- Final cleanup ---")
	visible, err = win32utils.IsConsoleVisible()
	if err != nil {
		log.Printf("Failed to check final visibility: %v", err)
	} else if !visible {
		fmt.Println("Making sure console is visible...")
		err = win32utils.ShowConsole()
		if err != nil {
			log.Printf("Failed to show console during cleanup: %v", err)
		}
	}
	
	// 恢复原始标题
	fmt.Printf("Restoring original console title: '%s'\n", currentTitle)
	err = win32utils.SetConsoleTitle(currentTitle)
	if err != nil {
		log.Printf("Failed to restore original console title: %v", err)
	}
	
	fmt.Println("\n=== Demo completed successfully! ===")
	fmt.Println("Press Enter to exit...")
	fmt.Scanln()
}