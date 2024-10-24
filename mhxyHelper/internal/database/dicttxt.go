package database

import (
	"bufio"
	"fmt"
	"mhxyHelper/pkg/logger"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	cutStrFileName     = "./ignore.txt"
	dictFileName       = "./dict.txt"
	dictBackupFileName = "./dict_bak_%d.txt"
)

var (
	ignoreTxtCutSets = make([]string, 0)
	nameDictMap      = map[string]struct{}{}
)

// 构建、备份 字典
func DictBackup() error {
	// 读取dict文件
	dictF, err := ReadCurrentDirFile(dictFileName)
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

// BuildDict 分词按行输出字典数据
func BuildDict(textArr []string) ([]string, []StuffLog) {
	dictResults := []string{}
	pLogResults := []StuffLog{}

	tempLog := StuffLog{}
	// 按行处理识别的字符数据
	for i, _ := range textArr {
		if len(tempLog.Name) > 0 && tempLog.ValMH > 0 {
			pLogResults = append(pLogResults, tempLog)
			tempLog = StuffLog{} // 加入后就可以重置了
		}

		curText := textArr[i]
		// 清洗数据
		tempText := TextTrims(curText)
		if len(tempText) < 0 {
			return nil, nil
		}

		// 获取该行文本汇总的数字出现的位置
		numIndex := getNumIndex(curText)
		// 没有数字则本行直接保留
		if numIndex == -1 {
			dictResults = append(dictResults, tempText)
			tempLog.Name = tempText // 记录商品名
			continue
		} else {
			str1 := tempText[0:numIndex]
			str2 := tempText[numIndex:]

			if len(str1) > 0 {
				dictResults = append(dictResults, str1)
			}

			// 价格数字
			if len(str2) > 0 {
				index := getNumIndex(str2)
				if index == 0 {
					fmt.Println("这是价格：", str2)
					valueStr, err := replaceAllString(str2)
					if err != nil {
						fmt.Println("replaceAllString value string err:", err.Error())
						continue
					}
					valueFloat, err := ConvertStr2Float32(valueStr)
					if err != nil {
						fmt.Println("convertStr2Int value string to int err:", err.Error())
						continue
					}
					tempLog.ValMH = valueFloat
					continue
				}
			}
		}
	}

	logger.Log.Info("=============================================")
	for _, v := range dictResults {
		logger.Log.Info(v)
	}

	return dictResults, pLogResults
}

// 获取数字出现的第一个位置
func getNumIndex(str string) int {
	// 正则表达式匹配字符串中的第一个数字
	re := regexp.MustCompile(`\d+`)
	matches := re.FindStringIndex(str)

	if len(matches) == 2 {
		// matches[0] 是匹配的数字的起始索引，matches[1] 是结束索引
		return matches[0]
	}
	return -1
}

// 替换字符串中的非数字字符
func replaceAllString(str string) (string, error) {
	// 创建一个正则表达式对象，匹配所有非数字字符
	reg, err := regexp.Compile("[^0-9]+")
	if err != nil {
		logger.Log.Error("Error compiling regex:", err)
		return "", err
	}

	// 使用正则表达式的ReplaceAllString方法替换掉所有非数字字符
	cleanStr := reg.ReplaceAllString(str, "")
	return cleanStr, nil
}

// 将字符串转换为Int数字
func ConvertStr2Int(str string) (int, error) {

	if len(str) <= 0 {
		return 0, nil
	}

	num, err := strconv.Atoi(str)
	if err != nil {
		logger.Log.Error("convertStr2Int string to int conversion failed:", err)
		return -1, err
	}
	return num, nil
}

// 将字符串转换为Float数字
func ConvertStr2Float32(str string) (float32, error) {

	if len(str) <= 0 {
		return 0, nil
	}

	floatValue, err := strconv.ParseFloat(str, 32)
	if err != nil {
		logger.Log.Error("convertStr2Int string to int conversion failed:", err)
		return -1, err
	}
	return float32(floatValue), nil
}
