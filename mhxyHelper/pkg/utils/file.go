package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// 定义支持的图片扩展名（小写格式）
var imageExtensions = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".gif":  true,
	".bmp":  true,
	".webp": true,
}

func IsImageExt(entry os.DirEntry) bool {
	// 获取文件扩展名并转换为小写
	ext := strings.ToLower(filepath.Ext(entry.Name()))

	// 检查是否为图片文件
	if !imageExtensions[ext] {
		return false
	}

	// 验证是否为常规文件（可选）
	if info, err := entry.Info(); err == nil && !info.Mode().IsRegular() {
		fmt.Printf("跳过非常规文件: %s\n", entry.Name())
		return false
	}

	return true
}

// 将相对路径转换为绝对路径
func ConvRelative2FullPath(relativePath, fileName string) string {

	// 直接获取绝对路径（自动处理当前工作目录）
	fullDirPath, err := filepath.Abs(relativePath)
	if err != nil {
		fmt.Printf("路径转换失败: %v\n", err)
		return ""
	}

	// 获取完整文件路径
	return filepath.Join(fullDirPath, fileName)
}

func WriteLinesToFile(lines []string, filename string) error {
	// 创建或打开文件
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("创建文件失败: %v", err)
	}
	defer file.Close() // 确保关闭文件

	// 遍历字符串数组，逐行写入
	for _, line := range lines {
		// 使用 `fmt.Fprintln` 写入行并自动添加换行符
		_, err := fmt.Fprintln(file, line)
		if err != nil {
			return fmt.Errorf("写入文件失败: %v", err)
		}
	}
	return nil
}
