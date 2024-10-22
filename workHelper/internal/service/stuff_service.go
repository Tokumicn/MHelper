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

// 查询商品信息
func QueryStuff(inStr string) (int64, []database.Stuff, error) {
	var (
		ctx    = context.Background()
		qStuff database.Stuff
		err    error
		offset int
		limit  = 50
	)

	qStuff, err = buildQuery(inStr)
	if err != nil {
		return 0, nil, err
	}

	// TODO log
	fmt.Printf("QueryStuff-buildQuery qStuff: %+v", qStuff)

	return qStuff.List(ctx, offset, limit)
}

// 根据输入字符串构建更为精准的查询条件  数据量不大该过程可以通过map映射完成
func buildQuery(inStr string) (database.Stuff, error) {
	res := database.Stuff{}

	questionStr, ok := common.QueryNameMap[inStr]
	if ok {
		res.Name = questionStr
	}

	if len(questionStr) == 0 {
		fmt.Printf("[优化内容经常查看] buildQuery [inStr: %s] 出现无法映射的用户输入, 请查看后收录进查询映射表\n", inStr)
		res.Name = inStr // 都无法命中的尝试用name字段查询
	}

	return res, nil
}

// 传入商品信息字符串 构建数据和日志
func BuildStuffByStr(stuffArr []string) error {
	ctx := context.Background()

	stuffs := make([]database.Stuff, len(stuffArr))

	// 读取本地文件增加到本次处理商品信息中
	tempStuffs, err := readCSVFromStuffData()
	if err != nil {
		logger.Log.ErrorContext(ctx, "BuildStuffByStr-readCSVFromStuffData err: %v\n", err)
		return err
	}
	stuffArr = append(stuffArr, tempStuffs...)

	for _, stuf := range stuffArr {

		temp, err := str2Stuff(stuf)
		if err != nil {
			logger.Log.ErrorContext(ctx, "buildStuff [temp: %v] err: %v \n", temp, err)
			continue
		}

		stuffs = append(stuffs, temp)
	}

	// 存储数据
	err = saveStuffs(ctx, stuffs)
	if err != nil {
		logger.Log.ErrorContext(ctx, "saveStuffs err: %v", err.Error())
		return err
	}

	return nil
}

// 从stuff数据文件中读取数据
func readCSVFromStuffData() ([]string, error) {
	res := make([]string, 0)
	f, err := os.Open("/Users/zhangrui/Workspace/goSpace/src/Tokumicn/theBookofChangesEveryDay/mhxyHelper/config/stuff_data.csv")
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
func str2Stuff(stuf string) (database.Stuff, error) {

	var (
		err    error
		qName  string
		name   string
		vMH    float32
		vRM    float32
		order  int
		region int
		empty  database.Stuff
	)

	splits := strings.Split(stuf, ",")

	// 必填字段  没有则报错
	qName, err = utils.ArrGetWithCheck(splits, 0)
	if err != nil {
		return empty, err
	}

	// 必填字段  没有则报错
	name, err = utils.ArrGetWithCheck(splits, 1)
	if err != nil {
		return empty, err
	}

	// 非必填字段
	vMHStr, _ := utils.ArrGetWithCheck(splits, 2)
	vMH, err = utils.ConvertStr2Float32(vMHStr)
	if err != nil {
		return empty, err
	}

	// 非必填字段
	vRMStr, _ := utils.ArrGetWithCheck(splits, 3)
	vRM, err = utils.ConvertStr2Float32(vRMStr)
	if err != nil {
		return empty, err
	}

	orderStr, _ := utils.ArrGetWithCheck(splits, 4)
	order, err = utils.ConvertStr2Int(orderStr)
	if err != nil {
		return empty, err
	}

	regionStr, _ := utils.ArrGetWithCheck(splits, 5)
	region, err = utils.ConvertStr2Int(regionStr)
	if err != nil {
		return empty, err
	}

	temp := database.Stuff{
		QName:    qName,
		Name:     name,
		ValMH:    vMH,
		ValRM:    vRM,
		Order:    order,
		RegionID: region,
	}

	return temp, nil
}

// 存储Stuff信息，根据Name判断是否已经存放，该段为全库表唯一
func saveStuffs(ctx context.Context, list []database.Stuff) error {

	for _, s := range list {
		isExist, id, err := s.ExistByQName(ctx)
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
