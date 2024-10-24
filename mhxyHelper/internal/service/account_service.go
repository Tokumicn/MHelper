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

// 账单信息构建并存储
func BuildAccountInfo(accountArr []string) error {
	ctx := context.Background()

	accounts := make([]database.Account, len(accountArr))

	for _, acStr := range accountArr {

		temp, err := str2Account(acStr)
		if err != nil {
			logger.Log.ErrorContext(ctx, "buildAccount [temp: %v] err: %v \n", temp, err)
			continue
		}

		// 部分字段特殊处理一下
		temp, err = buildAccountFields(temp)
		if err != nil {
			logger.Log.ErrorContext(ctx, "buildVal [temp: %v] err: %v \n", temp, err)
			continue
		}

		if temp.UserId <= 0 || len(temp.StuffName) <= 0 {
			continue // 如果没有必要信息则不处理
		}

		accounts = append(accounts, temp)
	}

	// 存储数据
	err := saveAccounts(ctx, accounts)
	if err != nil {
		logger.Log.ErrorContext(ctx, "saveAccounts err: %v", err)
		return err
	}

	return nil
}

// 将字符串转换为对象
// [userId, stuffName, buyMH, buyRM, sellMH, sellRM, region, regionName]
// 1,100灵饰指南,70,0,0,0,,
func str2Account(acStr string) (database.Account, error) {

	var (
		err        error
		userIdStr  string
		stuffName  string
		buyValMH   float32
		buyValRM   float32
		sellValMH  float32
		sellValRM  float32
		region     int
		regionName string
		empty      database.Account
	)

	splits := strings.Split(acStr, ",")

	// 必填字段  没有则报错
	userIdStr, err = utils.ArrGetWithCheck(splits, 0)
	userId, err := database.ConvertStr2Int(userIdStr)
	if err != nil {
		return empty, err
	}

	// 必填字段  没有则报错
	stuffName, err = utils.ArrGetWithCheck(splits, 1)
	if err != nil {
		return empty, err
	}

	// 非必填字段
	buyMHStr, _ := utils.ArrGetWithCheck(splits, 2)
	buyValMH, err = database.ConvertStr2Float32(buyMHStr)
	if err != nil {
		return empty, err
	}

	// 非必填字段
	buyRMStr, _ := utils.ArrGetWithCheck(splits, 3)
	buyValRM, err = database.ConvertStr2Float32(buyRMStr)
	if err != nil {
		return empty, err
	}

	// 非必填字段
	sellMHStr, _ := utils.ArrGetWithCheck(splits, 4)
	sellValMH, err = database.ConvertStr2Float32(sellMHStr)
	if err != nil {
		return empty, err
	}

	// 非必填字段
	sellRMStr, _ := utils.ArrGetWithCheck(splits, 5)
	sellValRM, err = database.ConvertStr2Float32(sellRMStr)
	if err != nil {
		return empty, err
	}

	regionStr, _ := utils.ArrGetWithCheck(splits, 5)
	region, err = database.ConvertStr2Int(regionStr)
	if err != nil {
		return empty, err
	}

	// 必填字段
	regionName, _ = utils.ArrGetWithCheck(splits, 6)
	if err != nil {
		return empty, err
	}

	temp := database.Account{
		UserId:     int64(userId),
		StuffName:  stuffName,
		BuyValMH:   buyValMH,
		BuyValRM:   buyValRM,
		SellValMH:  sellValMH,
		SellValRM:  sellValRM,
		RegionID:   region,
		RegionName: regionName,
	}

	return temp, nil
}

// 填充对象内容
func buildAccountFields(s database.Account) (database.Account, error) {

	// 处理购入和卖出价格信息
	bMH, bRM := buildVal(s.BuyValMH, s.BuyValRM)
	sMH, sRM := buildVal(s.SellValMH, s.SellValRM)

	s.BuyValMH = bMH
	s.BuyValRM = bRM
	s.SellValMH = sMH
	s.SellValRM = sRM

	// TODO 构建服务器名称

	return s, nil
}

// 存储Stuff信息，根据Name判断是否已经存放，该段为全库表唯一
func saveAccounts(ctx context.Context, list []database.Account) error {

	for _, ac := range list {
		_, err := ac.Create(ctx)
		if err != nil {
			logger.Log.Error("[saveAccounts] [account: %+v] err: %v", ac, err)
			continue
		}
	}

	return nil
}

// 查询账单信息
func QueryUserAccountInfo(queryStr string, userId int64) (int64, []database.Account, error) {
	var (
		ctx      = context.Background()
		qAccount database.Account
		offset   int
		limit    = 50
	)

	if userId > 0 {
		qAccount.UserId = userId
	} else {
		// 通过关键词查询userId
		userid, ok := common.QueryAccountMap[queryStr]
		if ok {
			qAccount.UserId = userid
		} else {
			return 0, nil, fmt.Errorf("查询不到账户信息呀~")
		}
	}

	return qAccount.List(ctx, offset, limit)
}
