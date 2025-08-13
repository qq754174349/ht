package user

import (
	"context"

	log "github.com/qq754174349/ht/ht-frame/logger"
	"github.com/qq754174349/ht/ht-frame/orm/mysql"
	"github.com/qq754174349/ht/ht-frame/orm/tx"
	"github.com/qq754174349/ht/ht-user/internal/domain/user/base"
	"gorm.io/gorm"
)

type BaseRepo struct {
	db *gorm.DB
}

func NewBaseRepo() *BaseRepo {
	return &BaseRepo{}
}

func (b *BaseRepo) Save(ctx context.Context, userBaseInfo *base.UserBaseInfo) (int64, error) {
	db := tx.GetTx(ctx, b.getDb())
	err := db.Create(&userBaseInfo).Error
	return userBaseInfo.ID, err
}

func (b *BaseRepo) UpdateById(ctx context.Context, userBaseInfo *base.UserBaseInfo) error {
	db := tx.GetTx(ctx, b.getDb())
	return db.Where("id = ?", userBaseInfo.ID).Updates(&userBaseInfo).Error
}

func (b *BaseRepo) FindByEmail(ctx context.Context, email string) (*base.UserBaseInfo, error) {
	db := b.getDb()
	var userBaseInfo base.UserBaseInfo
	err := db.Where("email = ?", email).Take(&userBaseInfo).Error
	return &userBaseInfo, err
}

func (b *BaseRepo) FindByUsername(ctx context.Context, username string) (*base.UserBaseInfo, error) {
	db := b.getDb()
	var userBaseInfo base.UserBaseInfo
	err := db.Where("username", username).Take(&userBaseInfo).Error
	return &userBaseInfo, err
}

func (b *BaseRepo) FindById(ctx context.Context, id int64) (*base.UserBaseInfo, error) {
	db := b.getDb()
	var userBaseInfo base.UserBaseInfo
	err := db.Where("id = ?", id).Take(&userBaseInfo).Error
	return &userBaseInfo, err
}

func (b *BaseRepo) getDb() *gorm.DB {
	if b.db == nil {
		db, err := mysql.Get()
		if err != nil {
			log.Fatal(err)
		}
		b.db = db
	}
	return b.db
}
