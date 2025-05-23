package data

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"mhxyHelper/pkg/database"
)

// 物品信息 无特殊属性的通用物品 如：宝石、兽决、灵饰指南、书铁、珍珠 等
type Stuff struct {
	Model
	QName    string  `gorm:"column:q_name" json:"qName"`   // 搜索名 一类商品总名称 如：月亮石
	Name     string  `gorm:"column:name"  json:"name"`     // 实际商品名 TODO: name字段添加全表唯一索引
	Order    int     `gorm:"column:order"  json:"order"`   // 顺序
	ValMH    float32 `gorm:"column:val_mh"  json:"mhCoin"` // MH W为单位
	ValRM    float32 `gorm:"column:val_rm" json:"rmCoin"`  // RM yuan为单位
	RegionID int     `gorm:"column:region_id" json:"-"`    // 所在区 通过 map 翻译即可
}

func (st Stuff) ToString() string {
	return fmt.Sprintf("[qName: %s, name: %s, mhCoin: %.2fW, rmCoin: %.2f元, order: %d]",
		st.QName, st.Name, st.ValMH, st.ValRM, st.Order)
}

type StuffLog struct {
	Stuff
}

// 搜索名：月亮石/月亮/月
// 月亮石  1级    8W 4元
// 元宵   攻击   64W ?元
// 低兽决 连击  100W  ?元
// 高兽决 高连击 500W ?元
// 灵饰指南 100级佩饰 55W ?元

func (pVal Stuff) ExistByQName(ctx context.Context) (bool, uint, error) {

	stuff, err := pVal.FindByName(ctx)
	if err == gorm.ErrRecordNotFound {
		return false, 0, nil
	}
	if err != nil {
		return false, 0, err
	}

	if stuff.ID > 0 {
		return true, stuff.ID, nil
	}

	return false, 0, nil
}

// 查询单个物品信息  name字段为全表唯一索引
func (pVal Stuff) FindByName(ctx context.Context) (Stuff, error) {
	res := Stuff{}
	if err := database.LocalDB().
		WithContext(ctx).
		Model(Stuff{}).
		Where("name = ?", pVal.Name).
		First(&res).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return Stuff{}, nil
		}
		return Stuff{}, fmt.Errorf("find one stuff info by qname err: %v", err)
	}

	return res, nil
}

// 创建商品价格信息
func (pVal Stuff) Create(ctx context.Context) (uint, error) {
	if err := database.LocalDB().
		WithContext(ctx).
		Create(&pVal).Error; err != nil {
		return 0, fmt.Errorf("create stuff info err: %v", err)
	}
	return pVal.ID, nil
}

// 更新商品信息
func (pVal Stuff) Update(ctx context.Context) (uint, error) {

	updateMap := map[string]interface{}{}

	if pVal.ValMH > 0 {
		updateMap["val_mh"] = pVal.ValMH
	}

	if pVal.ValRM > 0 {
		updateMap["val_rm"] = pVal.ValRM
	}

	if err := database.LocalDB().WithContext(ctx).
		Model(Stuff{}).
		Where("id = ?", pVal.ID).
		Error; err != nil {
		return 0, fmt.Errorf("update product value by updates:[%s] err: %v", updateMap, err)
	}
	return pVal.ID, nil
}

// 获取列表 目前仅提供通过名称查询
func (pVal Stuff) List(ctx context.Context, offset, limit int) (int64, []Stuff, error) {
	DB := database.LocalDB()
	vals := make([]Stuff, 0)
	var total int64
	DB = DB.WithContext(ctx).
		Model(Stuff{})

	// 组名称查询
	if len(pVal.QName) > 0 {
		DB = DB.Where("q_name = ?", pVal.QName)
	}

	// 唯一名称查询
	if len(pVal.Name) > 0 {
		DB = DB.Where("name = ?", pVal.Name)
	}

	// id查询
	if pVal.ID > 0 {
		DB = DB.Where("id = ?", pVal.ID)
	}

	if err := DB.Count(&total).Error; err != nil {
		return -1, nil, fmt.Errorf("get list product value count err: %v", err)
	}

	if err := DB.
		WithContext(ctx).
		Offset(offset).
		Limit(limit).
		Find(&vals).Error; err != nil {
		return -1, nil, fmt.Errorf("get list product value  err: %v", err)
	}
	return total, vals, nil
}

func (pVal Stuff) CreateStuffLog(ctx context.Context) (uint, error) {
	if err := database.LocalDB().WithContext(ctx).
		Table("stuff_log").
		Create(&pVal).
		Error; err != nil {
		return 0, fmt.Errorf("create stuff log err: %v", err)
	}
	return pVal.ID, nil
}
