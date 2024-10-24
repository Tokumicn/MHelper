package database

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"mhxyHelper/internal/utils"
)

// 物品属性信息 范围属性的物品 全服一致的属性 如：装备、灵饰等
type Attribute struct {
	Model
	QName string `gorm:"column:q_name" json:"qName"` // 搜索名 一类商品总名称 如：月亮石
	Name  string `gorm:"column:name"  json:"name"`   // 实际商品名 TODO: name字段添加全表唯一索引
	Max   string `gorm:"column:max"  json:"max"`     // 最大值  比如：总伤害、范围属性最大值等
	Desc  string `gorm:"column:desc" json:"desc"`    // 详细描述
	Order int    `gorm:"column:order"  json:"order"` // 顺序
}

func (at Attribute) ToString() string {

	res := ""
	if utils.IsMultiline(at.Desc) {
		res = fmt.Sprintf("[qName: %s, name: %s, max: %s, order: %d]\n",
			at.QName, at.Name, at.Max, at.Order)

		// TODO 针对复杂的desc特殊输出格式待开发
		res += at.Desc
	} else {
		res = fmt.Sprintf("[qName: %s, name: %s, max: %s, order: %d, desc: %s]",
			at.QName, at.Name, at.Max, at.Order, at.Desc)
	}

	return res
}

func (at Attribute) ExistByQName(ctx context.Context) (bool, uint, error) {

	attr, err := at.FindByName(ctx)
	if err == gorm.ErrRecordNotFound {
		return false, 0, nil
	}
	if err != nil {
		return false, 0, err
	}

	if attr.ID > 0 {
		return true, attr.ID, nil
	}

	return false, 0, nil
}

// 查询单个物品属性信息  name字段为全表唯一索引
func (at Attribute) FindByName(ctx context.Context) (Attribute, error) {
	res := Attribute{}
	if err := LocalDB().
		WithContext(ctx).
		Model(Attribute{}).
		Where("name = ?", at.Name).
		First(&res).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return Attribute{}, nil
		}
		return Attribute{}, fmt.Errorf("find one attribute info by name err: %v", err)
	}

	return res, nil
}

// 创建物品属性信息
func (at Attribute) Create(ctx context.Context) (uint, error) {
	if err := LocalDB().
		WithContext(ctx).
		Create(&at).Error; err != nil {
		return 0, fmt.Errorf("create attribute info err: %v", err)
	}
	return at.ID, nil
}

// 更新物品属性信息
func (at Attribute) Update(ctx context.Context) (uint, error) {

	updateMap := map[string]interface{}{}

	if len(at.Max) > 0 {
		updateMap["max"] = at.Max
	}

	if len(at.Desc) > 0 {
		updateMap["desc"] = at.Desc
	}

	if err := LocalDB().WithContext(ctx).
		Model(Attribute{}).
		Where("id = ?", at.ID).
		Error; err != nil {
		return 0, fmt.Errorf("update attribute info by updates:[%s] err: %v", updateMap, err)
	}
	return at.ID, nil
}

// 获取列表 目前仅提供通过名称查询
func (at Attribute) List(ctx context.Context, offset, limit int) (int64, []Attribute, error) {
	DB := LocalDB()
	vals := make([]Attribute, 0)
	var total int64
	DB = DB.WithContext(ctx).
		Model(Attribute{})

	// 组名称查询
	if len(at.QName) > 0 {
		DB = DB.Where("q_name = ?", at.QName)
	}

	// 唯一名称查询
	if len(at.Name) > 0 {
		DB = DB.Where("name = ?", at.Name)
	}

	// id查询
	if at.ID > 0 {
		DB = DB.Where("id = ?", at.ID)
	}

	if err := DB.Count(&total).Error; err != nil {
		return -1, nil, fmt.Errorf("get list attribute value count err: %v", err)
	}

	if err := DB.
		WithContext(ctx).
		Offset(offset).
		Limit(limit).
		Find(&vals).Error; err != nil {
		return -1, nil, fmt.Errorf("get list attribute value  err: %v", err)
	}
	return total, vals, nil
}
