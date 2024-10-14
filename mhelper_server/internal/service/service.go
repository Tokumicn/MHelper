package service

import (
	"context"

	otgorm "github.com/eddycjy/opentracing-gorm"

	"mhelper_server/global"
	"mhelper_server/internal/dao"
)

type Service struct {
	ctx context.Context
	dao *dao.Dao
	// cache
	// mq
	// 这些需要再业务上使用到的组件可以通过该方式传入
}

func New(ctx context.Context) Service {
	svc := Service{ctx: ctx}
	svc.dao = dao.New(otgorm.WithContext(svc.ctx, global.DBEngine))
	return svc
}
