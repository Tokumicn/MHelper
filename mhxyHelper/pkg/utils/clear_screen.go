package utils

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

func ClearScreen() {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls") // Windows 使用 cls
	} else {
		cmd = exec.Command("clear") // Linux/macOS 使用 clear
	}

	cmd.Stdout = os.Stdout
	err := cmd.Run() // 执行命令
	if err != nil {
		fmt.Println("Error clearing screen:", err)
	}
}

func ClearScreenV2() {
	switch runtime.GOOS {
	case "windows":
		// Windows 下发送 ANSI 转义码（需终端支持，如 Windows 10+）
		fmt.Print("\033[2J\033[H") // 清除屏幕并移动光标到左上角
	default:
		// Linux/macOS 或其他类 Unix 系统
		fmt.Print("\033[2J\033[H")
	}
}
