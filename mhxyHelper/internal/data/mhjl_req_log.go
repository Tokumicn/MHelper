package data

import (
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"log/slog"
	"mhxyHelper/pkg/database"
)

type MHJLResponseLog struct {
	Model
	UserId       uint32 `gorm:"column:userId" json:"user_id"`
	QueryMd5     string `gorm:"column:query_md5" json:"query_md5"`
	Query        string `gorm:"column:query" json:"query"`
	RawAnswerMd5 string `gorm:"column:raw_answer_md5" json:"raw_answer_md5"`
	RawAnswer    string `gorm:"column:raw_answer" json:"raw_answer"`
	FormatAnswer string `gorm:"column:format_answer"  json:"format_answer"` // MH W为单位
}

func (at MHJLResponseLog) Save(ctx context.Context) (uint, error) {

	// 回答存在过 就不再存了
	exist, err := at.existRawAnswer(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "[MHJLResponseLog] existRawAnswer err: ", err.Error())
	}
	if exist != nil {
		return exist.ID, nil
	}

	if err = database.LocalDB().
		WithContext(ctx).
		Save(&at).Error; err != nil {
		return 0, fmt.Errorf("[MHJLResponseLog]  save 梦幻精灵 request log err: %v", err)
	}
	return at.ID, nil
}

// 重复问题
func (at MHJLResponseLog) existQuery(ctx context.Context) (*MHJLResponseLog, error) {
	res := &MHJLResponseLog{}
	if err := database.LocalDB().
		WithContext(ctx).
		Model(MHJLResponseLog{}).
		Find("query_md5 = ?", at.QueryMd5).
		First(res).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return res, nil
}

// 重复回答
func (at MHJLResponseLog) existRawAnswer(ctx context.Context) (*MHJLResponseLog, error) {
	res := &MHJLResponseLog{}
	if err := database.LocalDB().
		WithContext(ctx).
		Model(MHJLResponseLog{}).
		Where("raw_answer_md5 = ?", at.RawAnswerMd5).
		Find(res).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return res, nil
}
