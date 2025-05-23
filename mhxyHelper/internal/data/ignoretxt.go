package data

import (
	"io"
	"log"
	"mhxyHelper/pkg/logger"
	"strings"
)

func InitCutSets() {
	// 读取需要清洗的字符字典
	cutSets := readCutSetsFromFile()
	logger.Log.Debug("ignore-text 读取需要清洗的字典数据 cutSets:[%v] \n", cutSets)
	// 写入到全局变量备用
	ignoreTxtCutSets = append(ignoreTxtCutSets, cutSets...)
}

// TextTrims 文本清洗
func TextTrims(pendingTxt string) string {
	logger.Log.Debug("本文清理前: %s", pendingTxt)
	res := strTrims(pendingTxt)
	logger.Log.Debug("本文清理后: %s", res)
	return res
}

// 对预处理文本进行清理
func strTrims(str string) string {
	// 去两边空格
	result := strings.TrimSpace(str)

	// 清理额外字符
	cutSets := []string{"(", ")", " ", "!", "！", "单价", "单馆", "。"}
	if len(ignoreTxtCutSets) > 0 {
		cutSets = append(cutSets, ignoreTxtCutSets...)
	}

	for _, cut := range cutSets {
		result = strings.Trim(result, cut)
	}

	return result
}

// 从文件中读取需要清理的额外字符
func readCutSetsFromFile() []string {
	res := make([]string, 0)

	// 读取 ignore.txt 文件
	cutStrF, err := ReadCurrentDirFile(cutStrFileName)
	defer cutStrF.Close() // 确保在函数结束时关闭文件
	if err != nil {
		log.Printf("[ERROR] readCurrentDirFile [%s] err: %v", dictFileName)
		return nil
	}
	allCutStrs, err := io.ReadAll(cutStrF)
	if err != nil {
		log.Printf("[ERROR] io.ReadAll [%s] err: %v", dictFileName)
		return nil
	}

	cutStrs := string(allCutStrs)
	log.Printf("cutStrs : %s", cutStrs)
	res = strings.Split(cutStrs, "\n")
	log.Printf("cutArr : %s", res)
	return res
}
