package user

import (
	"context"
	"github.com/qq754174349/ht/ht-frame/orm/mysql"
	"github.com/qq754174349/ht/ht-frame/orm/tx"
	"ht-user/internal/domain/user/base"
)

type BaseRepo struct {
}

func NewBaseRepo() *BaseRepo {
	return &BaseRepo{}
}

func (b *BaseRepo) Save(ctx context.Context, userBaseInfo *base.UserBaseInfo) (int64, error) {
	db, _ := mysql.Get()
	db = tx.GetTx(ctx, db)
	db.Create(&userBaseInfo)
	return userBaseInfo.ID, db.Error
}
