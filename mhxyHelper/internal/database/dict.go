package database

import "gorm.io/gorm"

type DictV1 struct {
	gorm.Model
	Name    string `gorm:"name"`     // 字典名称
	IsFirst bool   `gorm:"is_first"` // 是否第一次构建
}
