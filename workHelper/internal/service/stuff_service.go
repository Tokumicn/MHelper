package service

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"workHelper/internal/database"
	"workHelper/pkg/common"
	"workHelper/pkg/logger"
	"workHelper/pkg/utils"
)

// 查询问题帮助
func QueryQuestion(inStr string) (int64, []database.Helper, error) {
	var (
		ctx     = context.Background()
		qHelper database.Helper
		err     error
		offset  int
		limit   = 50
	)

	qHelper, err = buildQuery(inStr)
	if err != nil {
		return 0, nil, err
	}

	// TODO log
	fmt.Printf("QueryHelper-buildQuery qHelper: %+v", qHelper)

	return qHelper.List(ctx, offset, limit)
}

// 根据输入字符串构建更为精准的查询条件  数据量不大该过程可以通过map映射完成
func buildQuery(inStr string) (database.Helper, error) {
	res := database.Helper{}

	questionStr, ok := common.QueryNameMap[inStr]
	if ok {
		res.Question = questionStr
	}

	if len(questionStr) == 0 {
		fmt.Printf("[优化内容经常查看] buildQuery [inStr: %s] 出现无法映射的用户输入, 请查看后收录进查询映射表\n", inStr)
		res.Question = inStr // 都无法命中的尝试用name字段查询
	}

	return res, nil
}

// 传入商品信息字符串 构建数据和日志
func BuildHelperByStr(helperArr []string) error {
	ctx := context.Background()

	helpers := make([]database.Helper, len(helperArr))

	// 读取本地文件增加到本次处理帮助信息中
	tempHelpers, err := readCSVFromStuffData()
	if err != nil {
		logger.Log.ErrorContext(ctx, "BuildStuffByStr-readCSVFromStuffData err: %v\n", err)
		return err
	}
	helperArr = append(helperArr, tempHelpers...)

	for _, stuf := range helperArr {

		temp, err := str2Helper(stuf)
		if err != nil {
			logger.Log.ErrorContext(ctx, "str2Helper [temp: %v] err: %v \n", temp, err)
			continue
		}

		helpers = append(helpers, temp)
	}

	// 存储数据
	err = saveHelpers(ctx, helpers)
	if err != nil {
		logger.Log.ErrorContext(ctx, "saveHelpers err: %v", err.Error())
		return err
	}

	return nil
}

// 从stuff数据文件中读取数据
func readCSVFromStuffData() ([]string, error) {
	res := make([]string, 0)
	f, err := os.Open("/Users/zhangrui/Workspace/goSpace/src/Tokumicn/theBookofChangesEveryDay/workHelper/config/helper_data.csv")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	r := bufio.NewReader(f)
	r.ReadLine() // 丢弃第一行
	for {
		line, _, err := r.ReadLine()
		if err != nil {
			break
		}
		res = append(res, string(line))
	}
	return res, nil
}

// 将字符串转换为对象
func str2Helper(stuf string) (database.Helper, error) {

	var (
		err     error
		quetion string
		answers string
		empty   database.Helper
	)

	splits := strings.Split(stuf, ",")

	// 必填字段  没有则报错
	quetion, err = utils.ArrGetWithCheck(splits, 0)
	if err != nil {
		return empty, err
	}

	// 必填字段  没有则报错
	answers, err = utils.ArrGetWithCheck(splits, 1)
	if err != nil {
		return empty, err
	}

	temp := database.Helper{
		Question: quetion,
		Answers:  answers,
	}

	return temp, nil
}

// 存储Stuff信息，根据Name判断是否已经存放，该段为全库表唯一
func saveHelpers(ctx context.Context, list []database.Helper) error {

	for _, s := range list {
		isExist, id, err := s.ExistByQuestion(ctx)
		if err != nil {
			// TODO log
			return err
		}

		if isExist { // 更新
			s.ID = id
			_, err = s.Update(ctx)
			if err != nil {
				// TODO log
				return err
			}
		} else {
			_, err = s.Create(ctx)
			if err != nil {
				// TODO log
				return err
			}
		}
	}

	return nil
}
