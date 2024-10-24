package database

import (
	"log"
	"os"
)

// 读取工作目录下的文件
func ReadCurrentDirFile(fileName string) (*os.File, error) {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("[当前工作目录：%s]\n", dir)
	// 打开文件
	return os.Open(fileName)
}
