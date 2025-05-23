package main

import (
	"bufio"
	"context"
	"fmt"
	"mhxyHelper/internal/data"
	"mhxyHelper/internal/service"
	"mhxyHelper/internal/service/query/mhjl_query"
	"mhxyHelper/pkg/database"
	"mhxyHelper/pkg/logger"
	"os"
	"strings"
)

const (
	processLOCAL = "local" // 查询本地数据库 默认值
	processMHJL  = "mhjl"  // 梦幻精灵
	processRAG   = "rag"   // 问知识库
	processLLM   = "llm"   // 问大模型
	processExit  = "exit"  // 退出
)

func main() {
	logger.NewLogger()

	// 初始化数据库连接
	_, err := database.InitDB()
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	// 初始化表结构
	// database.InitDBWithAutoMigrate(true) // 初始化协助构建表结构
	// DictBuildToolV1() // 构建字典信息

	// 写入测试 测试根据csv文件构建数据并存储
	//err := service.BuildStuffByStr([]string{})
	//if err != nil {
	//	fmt.Println("err: ", err.Error())
	//}

	// 写入测试  测试根据excel文件构建属性数据并存储
	//err = service.BuildAttributeByStr([]string{})
	//if err != nil {
	//	logger.Log.Error("BuildAttributeByStr err: %v", err)
	//	return
	//}

	// 默认走本地查询
	processFlag := processLOCAL
	for {
		fmt.Println("帮助信息 [<命令: 直接输入问题 | Q | quit >:<空格> <问题>] [eg: Q: 160武器]") // <命令: Q LLM RAG quit>
		fmt.Print("输入查询关键词: ")

		// 读取用户输入
		cmdStr, inputStr := buildInput()

		// 检查用户输入尝试路由
		processFlag = checkUserInput(cmdStr)

		// 退出
		if processFlag == processExit {
			break
		}

		total, typeStr, data, err := processAnswer(ctx, processFlag, inputStr)
		if err != nil {
			fmt.Println("err: ", err.Error())
			continue
		}

		printQueryResult(total, typeStr, data)
	}

}

func processAnswer(ctx context.Context, processFlag string, inputStr string) (int64, string, interface{}, error) {
	var (
		err     error
		total   int64
		typeStr string
		data    interface{}
	)

	switch processFlag {
	case processLOCAL:
		total, typeStr, data, err = service.QueryLocal(inputStr)
		if err != nil {
			fmt.Println("err: ", err.Error())
			return 0, "", "", err
		}
	case processMHJL:
		total = 0
		typeStr = service.TypeMHJL
		data, err = mhjl_query.QueryMHJL(ctx, inputStr)
		if err != nil {
			fmt.Println("err: ", err.Error())
			return 0, "", "", err
		}
	default:
		fmt.Println("processFlag: ", processFlag)
		fmt.Println("暂未实现.")
	}
	return total, typeStr, data, nil
}

// 构建输入
func buildInput() (string, string) {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	rawInputStr := scanner.Text()

	var (
		cmdStr   string // 命令
		inputStr string // 问题
	)
	rawInputStr = strings.TrimSpace(rawInputStr)
	splitInput := strings.Split(rawInputStr, " ")
	if len(splitInput) >= 2 {
		cmdStr = splitInput[0]
		inputStr = splitInput[1]
	}

	if len(splitInput) == 1 {
		cmdStr = rawInputStr
		inputStr = rawInputStr
	}

	return cmdStr, inputStr
}

func checkUserInput(cmdStr string) string {

	if len(cmdStr) <= 0 {
		return processLOCAL
	}

	cmdStr = strings.ToLower(cmdStr)

	// 退出
	if cmdStr == "quit" || cmdStr == "exit" {
		return processExit
	}

	// 问精灵
	if hasPrefix4Map(cmdStr, []string{"q", "q:", "q: "}) {
		return processMHJL
	}

	// 问知识库
	if hasPrefix4Map(cmdStr, []string{"rag", "rag:", "rag: "}) {
		return processRAG
	}

	// 问大模型
	if hasPrefix4Map(cmdStr, []string{"llm", "llm:", "llm: "}) {
		return processLLM
	}

	return processLOCAL
}

// 基于 Map 匹配前缀
func hasPrefix4Map(inputStr string, prefixArr []string) bool {
	for _, prefix := range prefixArr {
		// 有一个就可以
		if strings.HasPrefix(inputStr, prefix) {
			return true
		}
	}
	return false
}

func printQueryResult(total int64, typeStr string, result interface{}) {

	fmt.Println("\n查询结果: ")
	fmt.Println("=====================================")
	if total > 0 {
		fmt.Println("total: ", total)
	}

	switch typeStr {
	case service.TypeAttribute:
		attributes := result.([]data.Attribute)
		for _, at := range attributes {
			fmt.Println(at.ToString())
			fmt.Println("")
		}

	case service.TypeStuff:
		stuffs := result.([]data.Stuff)
		for _, st := range stuffs {
			fmt.Println(st.ToString())
			fmt.Println("")
		}

	case service.TypeMHJL:
		fmt.Println(result)
	default:
		fmt.Println("没查到.")
	}

	fmt.Println("=====================================")
}

func DictBuildToolV1() {
	// 初始化数据清理字典
	data.InitCutSets()

	// 备份dict.txt
	data.DictBackup()

	// 接收多行输入  回车结束
	inputArr := scanInputText()
	tempDict, tempProducts := data.BuildDict(inputArr)

	logger.Log.Info("============================")
	for _, v := range tempProducts {
		fmt.Println(v)
	}
	logger.Log.Info("============================")

	data.SaveDict2Txt(tempDict)
}

// 按行接收输入的多行数据 直到回车结束
func scanInputText() []string {
	// 创建一个bufio.Scanner，用于读取控制台输入
	scanner := bufio.NewScanner(os.Stdin)

	// 打印提示信息
	fmt.Println("请输入多行数据，输入空行结束：")
	var inputTextArr []string
	// 使用循环读取每一行输入
	for scanner.Scan() {
		// 读取的文本赋值给text变量
		text := scanner.Text()
		inputTextArr = append(inputTextArr, text)
		// 检查是否输入了空行
		if text == "" {
			break
		}
	}

	// 检查是否有可能发生的错误
	if err := scanner.Err(); err != nil {
		logger.Log.Error("读取输入时发生错误: ", err.Error())
		return nil
	}

	return inputTextArr
}
