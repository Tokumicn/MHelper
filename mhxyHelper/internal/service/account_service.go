package service

import (
	"context"
	"mhxyHelper/internal/database"
	"mhxyHelper/pkg/logger"
	"mhxyHelper/pkg/utils"
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
	userId, err := utils.ConvertStr2Int(userIdStr)
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
	buyValMH, err = utils.ConvertStr2Float32(buyMHStr)
	if err != nil {
		return empty, err
	}

	// 非必填字段
	buyRMStr, _ := utils.ArrGetWithCheck(splits, 3)
	buyValRM, err = utils.ConvertStr2Float32(buyRMStr)
	if err != nil {
		return empty, err
	}

	// 非必填字段
	sellMHStr, _ := utils.ArrGetWithCheck(splits, 4)
	sellValMH, err = utils.ConvertStr2Float32(sellMHStr)
	if err != nil {
		return empty, err
	}

	// 非必填字段
	sellRMStr, _ := utils.ArrGetWithCheck(splits, 5)
	sellValRM, err = utils.ConvertStr2Float32(sellRMStr)
	if err != nil {
		return empty, err
	}

	regionStr, _ := utils.ArrGetWithCheck(splits, 5)
	region, err = utils.ConvertStr2Int(regionStr)
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
