package database

import (
	"context"
	"fmt"
	"gorm.io/gorm"
)

// 账单
type Account struct {
	Model
	UserId     int64   `gorm:"userId" json:"user_id"`
	StuffName  string  `gorm:"stuffName" json:"stuff_name"`
	BuyValMH   float32 `gorm:"column:buy_val_mh"  json:"bMHCoin"`    // MH W为单位
	BuyValRM   float32 `gorm:"column:buy_val_rm" json:"bRMCoin"`     // RM yuan为单位
	SellValMH  float32 `gorm:"column:sell_val_mh"  json:"sMHCoin"`   // MH W为单位
	SellValRM  float32 `gorm:"column:sell_val_rm" json:"sRMCoin"`    // RM yuan为单位
	RegionID   int     `gorm:"column:region_id" json:"-"`            // 所在区 通过 map 翻译即可
	RegionName string  `gorm:"column:region_name" json:"regionName"` // 所在区 通过 map 翻译即可
}

// 查询该用户单个物品信息  name + userId 字段为查询组合   一个用户可以买多个
func (ac Account) FindUserAccountInfo(ctx context.Context) (Account, error) {
	res := Account{}
	if err := LocalDB().
		WithContext(ctx).
		Model(Account{}).
		Where("user_id = ? AND stuff_name = ?", ac.UserId, ac.StuffName).
		First(&res).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return Account{}, nil
		}
		return Account{}, fmt.Errorf("find one stuff info by qname err: %v", err)
	}

	return res, nil
}

// 创建账单信息
func (ac Account) Create(ctx context.Context) (uint, error) {
	if err := LocalDB().
		WithContext(ctx).
		Create(&ac).Error; err != nil {
		return 0, fmt.Errorf("create account info err: %v", err)
	}
	return ac.ID, nil
}

// 更新账单信息
func (ac Account) Update(ctx context.Context) (uint, error) {

	updateMap := map[string]interface{}{}

	if ac.BuyValMH > 0 {
		updateMap["buy_val_mh"] = ac.BuyValMH
	}

	if ac.BuyValRM > 0 {
		updateMap["buy_val_rm"] = ac.BuyValRM
	}

	if ac.SellValMH > 0 {
		updateMap["sell_val_mh"] = ac.SellValMH
	}

	if ac.SellValRM > 0 {
		updateMap["sell_val_rm"] = ac.SellValRM
	}

	if err := LocalDB().WithContext(ctx).
		Model(Account{}).
		Where("id = ?", ac.ID).
		Error; err != nil {
		return 0, fmt.Errorf("update user account value by updates:[%s] err: %v", updateMap, err)
	}
	return ac.ID, nil
}

// 获取账单列表 目前仅提供通过用户和物品名称
func (ac Account) List(ctx context.Context, offset, limit int) (int64, []Account, error) {
	DB := LocalDB()
	vals := make([]Account, 0)

	var total int64
	DB = DB.WithContext(ctx).
		Model(Account{})

	// 物品名称查询
	if len(ac.StuffName) > 0 {
		DB = DB.Where("stuff_name = ?", ac.StuffName)
	}

	// 用户查询
	if ac.UserId > 0 {
		DB = DB.Where("user_id = ?", ac.UserId)
	}

	// id查询
	if ac.ID > 0 {
		DB = DB.Where("id = ?", ac.ID)
	}

	if err := DB.Count(&total).Error; err != nil {
		return -1, nil, fmt.Errorf("get list account info count err: %v", err)
	}

	if err := DB.
		WithContext(ctx).
		Offset(offset).
		Limit(limit).
		Find(&vals).Error; err != nil {
		return -1, nil, fmt.Errorf("get list account info err: %v", err)
	}
	return total, vals, nil
}
