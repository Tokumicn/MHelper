package database

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"time"
)

type Model struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// 问题
type Helper struct {
	Model
	Question string `gorm:"column:question"  json:"question"` // 提问
	Answers  string `gorm:"column:answers" json:"answers"`    // 答案列表Map[string]string key: Hash(answer)
}

func (h Helper) ToString() string {
	return fmt.Sprintf("[question: %s, answer: %s]", h.Question, h.Answers)
}

func (h Helper) ExistByQuestion(ctx context.Context) (bool, uint, error) {

	helper, err := h.FindByQuestion(ctx)
	if err == gorm.ErrRecordNotFound {
		return false, 0, nil
	}
	if err != nil {
		return false, 0, err
	}

	if helper.ID > 0 {
		return true, helper.ID, nil
	}

	return false, 0, nil
}

// 查询单个帮助信息  question字段为全表唯一索引
func (h Helper) FindByQuestion(ctx context.Context) (Helper, error) {
	res := Helper{}
	if err := LocalDB().
		WithContext(ctx).
		Model(Helper{}).
		Where("question = ?", h.Question).
		First(&res).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return Helper{}, nil
		}
		return Helper{}, fmt.Errorf("find one helper info by question[%s] err: %v",
			h.Question, err)
	}

	return res, nil
}

// 创建帮助信息
func (h Helper) Create(ctx context.Context) (uint, error) {
	if err := LocalDB().
		WithContext(ctx).
		Create(&h).Error; err != nil {
		return 0, fmt.Errorf("create helper info err: %v", err)
	}
	return h.ID, nil
}

// 更新帮助信息 待更新
func (h Helper) Update(ctx context.Context) (uint, error) {

	updateMap := map[string]interface{}{}

	if len(h.Question) > 0 {
		updateMap["question"] = h.Question
	}

	if len(h.Answers) > 0 {
		updateMap["answers"] = h.Answers
	}

	if err := LocalDB().WithContext(ctx).
		Model(Helper{}).
		Where("id = ?", h.ID).
		Error; err != nil {
		return 0, fmt.Errorf("update helper info by updates:[%s] err: %v", updateMap, err)
	}
	return h.ID, nil
}

// 获取列表 目前仅提供通过问题查询
func (h Helper) List(ctx context.Context, offset, limit int) (int64, []Helper, error) {
	DB := LocalDB()
	vals := make([]Helper, 0)
	var total int64
	DB = DB.WithContext(ctx).
		Model(Helper{})

	// 组名称查询
	if len(h.Question) > 0 {
		DB = DB.Where("question = ?", h.Question)
	}

	// id查询
	if h.ID > 0 {
		DB = DB.Where("id = ?", h.ID)
	}

	if err := DB.Count(&total).Error; err != nil {
		return -1, nil, fmt.Errorf("get list helper info count err: %v", err)
	}

	if err := DB.
		WithContext(ctx).
		Offset(offset).
		Limit(limit).
		Find(&vals).Error; err != nil {
		return -1, nil, fmt.Errorf("get list helper info err: %v", err)
	}
	return total, vals, nil
}
