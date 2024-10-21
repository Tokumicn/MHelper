package main

import (
	"bufio"
	"fmt"
	"github.com/Tokumicn/theBookofChangesEveryDay/mhxyHelper/internal/database"
	"github.com/Tokumicn/theBookofChangesEveryDay/mhxyHelper/internal/service"
	"github.com/Tokumicn/theBookofChangesEveryDay/mhxyHelper/pkg/logger"
	"github.com/Tokumicn/theBookofChangesEveryDay/mhxyHelper/pkg/utils"
	"os"
)

func main() {
	logger.NewLogger()

	// 初始化数据库连接
	_, err := database.InitDB()
	if err != nil {
		panic(err)
	}

	// 初始化表结构
	// db.InitDBWithAutoMigrate(true) // 初始化协助构建表结构
	// DictBuildToolV1() // 构建字典信息

	// 写入测试 测试根据csv文件构建数据并存储
	//err := service.BuildStuffByStr([]string{})
	//if err != nil {
	//	fmt.Println("err: ", err.Error())
	//}

	for {
		fmt.Print("输入查询关键词: ")
		var inputStr string
		_, err = fmt.Scan(&inputStr)
		if err != nil {
			panic(err)
		}

		if inputStr == "q" || inputStr == "quit" || inputStr == "exit" {
			break
		}

		// 查询测试
		total, stuffs, err := service.QueryStuff(inputStr)
		if err != nil {
			fmt.Println("err: ", err.Error())
			return
		}

		fmt.Println("\n查询结果: ")
		fmt.Println("=====================================")
		fmt.Println("total: ", total)

		for _, st := range stuffs {
			fmt.Println(st.ToString())
		}
		fmt.Println("=====================================")
	}

}

func DictBuildToolV1() {
	// 初始化数据清理字典
	utils.InitCutSets()

	// 备份dict.txt
	utils.DictBackup()

	// 接收多行输入  回车结束
	inputArr := scanInputText()
	tempDict, tempProducts := utils.BuildDict(inputArr)

	logger.Log.Info("============================")
	for _, v := range tempProducts {
		fmt.Println(v)
	}
	logger.Log.Info("============================")

	utils.SaveDict2Txt(tempDict)
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
