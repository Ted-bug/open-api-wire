package repo

import (
	"api-gin/infra/log"
	"context"
	"gorm.io/gorm"
)

type UserRepo struct {
	baseRepo *BaseRepo
	logger   *log.Logger
}

func NewUserRepo(db *gorm.DB, logger *log.Logger) *UserRepo {
	return &UserRepo{
		baseRepo: NewBaseRepo(db),
		logger:   logger.NewLogger("UserRepo"),
	}
}

func (u *UserRepo) Hello(ctx context.Context) string {
	// 假设查询了数据库
	u.logger.Info(ctx, "hello")
	return "hello "
}
