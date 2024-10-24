package service

import (
	"context"
	"fmt"
	"mhxyHelper/internal/database"
	"mhxyHelper/internal/utils"
	"mhxyHelper/pkg/common"
	"mhxyHelper/pkg/logger"
	"strings"
)

// 查询属性信息
func QueryAttribute(inStr string) (int64, []database.Attribute, error) {
	var (
		ctx    = context.Background()
		qAttr  database.Attribute
		err    error
		offset int
		limit  = 50
	)

	qAttr, err = buildQueryAttr(inStr)
	if err != nil {
		return 0, nil, err
	}

	logger.Log.InfoContext(ctx, "QueryAttribute-buildQueryAttr qAttr: %+v", qAttr)

	return qAttr.List(ctx, offset, limit)
}

// 根据输入字符串构建更为精准的查询条件  数据量不大该过程可以通过map映射完成
func buildQueryAttr(inStr string) (database.Attribute, error) {
	res := database.Attribute{}

	qNameStr, ok := common.QueryQNameMapAttribute[inStr]
	if ok {
		res.QName = qNameStr
	}

	nameStr, ok := common.QueryNameMapAttribute[inStr]
	if ok {
		res.Name = nameStr
	}

	if len(nameStr) == 0 && len(qNameStr) == 0 {
		fmt.Printf("[优化内容经常查看] buildQueryAttr [inStr: %s] 出现无法映射的用户输入, 请查看后收录进查询映射表\n", inStr)
		res.Name = inStr // 都无法命中的尝试用name字段查询
	}

	return res, nil
}

// 传入商品信息字符串 构建数据和日志
func BuildAttributeByStr(attributeArr []string) error {
	ctx := context.Background()

	attributes := make([]database.Attribute, len(attributeArr))

	// 读取本地文件增加到本次处理商品信息中
	var (
		tempAttrs [][]string
		err       error
	)
	tempAttrs, err = ReadAttributeFromExcel(ctx)
	if err != nil {
		logger.Log.ErrorContext(ctx, "BuildAttributeByStr-ReadAttributeFromExcel err: %v\n", err)
		return err
	}

	for _, attr := range attributeArr {
		tempSplit := strings.Split(attr, ",")
		tempAttrs = append(tempAttrs, tempSplit)
	}

	for _, attrArr := range tempAttrs {

		temp, err := str2Attribute(attrArr)
		if err != nil {
			logger.Log.ErrorContext(ctx, "buildStuff [temp: %v] err: %v \n", temp, err)
			continue
		}

		attributes = append(attributes, temp)
	}

	// 存储数据
	err = saveAttributes(ctx, attributes)
	if err != nil {
		logger.Log.ErrorContext(ctx, "saveStuffs err: %v", err.Error())
		return err
	}

	return nil
}

// 将字符串转换为对象
func str2Attribute(attrArr []string) (database.Attribute, error) {

	var (
		err    error
		qName  string
		name   string
		maxStr string
		desc   string
		order  int
		empty  database.Attribute
	)

	// 必填字段  没有则报错
	qName, err = utils.ArrGetWithCheck(attrArr, 1)
	if err != nil {
		return empty, err
	}

	// 必填字段  没有则报错
	name, err = utils.ArrGetWithCheck(attrArr, 2)
	if err != nil {
		return empty, err
	}

	// 非必填字段
	maxStr, err = utils.ArrGetWithCheck(attrArr, 3)
	if err != nil {
		return empty, err
	}

	// 非必填字段
	desc, err = utils.ArrGetWithCheck(attrArr, 4)
	if err != nil {
		return empty, err
	}

	orderStr, _ := utils.ArrGetWithCheck(attrArr, 5)
	order, err = database.ConvertStr2Int(orderStr)
	if err != nil {
		return empty, err
	}

	temp := database.Attribute{
		QName: qName,
		Name:  name,
		Max:   maxStr,
		Desc:  desc,
		Order: order,
	}

	return temp, nil
}

// 存储Stuff信息，根据Name判断是否已经存放，该段为全库表唯一
func saveAttributes(ctx context.Context, list []database.Attribute) error {

	for _, s := range list {
		isExist, id, err := s.ExistByQName(ctx)
		if err != nil {
			logger.Log.ErrorContext(ctx, "saveAttributes ExistByQName [Attribute: %+v] err: %v", s, err)
			return err
		}

		if isExist { // 更新
			s.ID = id
			_, err = s.Update(ctx)
			if err != nil {
				logger.Log.ErrorContext(ctx, "saveAttributes Update [Attribute: %+v] err: %v", s, err)
				return err
			}
		} else {
			_, err = s.Create(ctx)
			if err != nil {
				logger.Log.ErrorContext(ctx, "saveAttributes Create [Attribute: %+v] err: %v", s, err)
				return err
			}
		}
	}

	return nil
}
