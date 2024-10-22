package utils

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
	"workHelper/pkg/logger"
)

// 构建、备份 字典

func DictBackup() error {
	// 读取dict文件
	dictF, err := readCurrentDirFile(dictFileName)
	defer dictF.Close() // 确保在函数结束时关闭文件
	if err != nil {
		logger.Log.Error("[ERROR] DictBackup readCurrentDirFile [%s] err: %v \n", dictFileName, err)
		return err
	}

	backupFileName := fmt.Sprintf(dictBackupFileName, time.Now().Unix())
	logger.Log.Info("[INFO] DictBackup 开始备份文件 文件名: ", backupFileName)

	backupFile, err := os.Create(backupFileName)
	defer backupFile.Close()
	if err != nil {
		logger.Log.Error("[ERROR] DictBackup 备份字典文件错误，err: %s\n", err)
		return err
	}

	// 创建 Scanner 来按行读取
	scanner := bufio.NewScanner(dictF)
	// 使用 Scan() 方法按行迭代文件
	for scanner.Scan() {
		line := scanner.Text() // 获取当前行的文本
		// 加载到内存中的Map
		nameDictMap[line] = struct{}{}

		// 写入备份文件
		_, err = backupFile.WriteString(line + "\n")
		if err != nil {
			logger.Log.Error("写入备份文件时发生错误: %v", err)
			continue
		}
	}
	// 检查是否有可能的错误
	if err := scanner.Err(); err != nil {
		logger.Log.Error("scanner", err)
	}

	logger.Log.Info("[备份完成!!!] 文件名: ", backupFileName)
	return nil
}

// SaveDict2Txt 保存dict.txt
func SaveDict2Txt(tempDict []string) {
	for _, v := range tempDict {
		// 将新的字典数据覆盖之前的数据
		nameDictMap[v] = struct{}{}
	}

	// 创建文件，如果文件已存在，它将被截断（覆盖）
	file, err := os.Create(dictFileName)
	if err != nil {
		logger.Log.Error("create backup dict file err: %v", err)
	}
	defer file.Close() // 确保在函数结束时关闭文件

	var fileLines []string
	for v, _ := range nameDictMap {
		fileLines = append(fileLines, v)
	}
	fileText := strings.Join(fileLines, "\n")

	// 写入数据到文件
	_, err = file.WriteString(fileText)
	if err != nil {
		logger.Log.Error("覆盖写入 dict.txt 错误 , err: ", err.Error())
		return
	}
}
