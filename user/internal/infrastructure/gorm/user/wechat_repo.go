package user

import (
	"context"

	"github.com/qq754174349/ht/ht-frame/logger"
	"github.com/qq754174349/ht/ht-frame/orm/mysql"
	"github.com/qq754174349/ht/ht-frame/orm/tx"
	"github.com/qq754174349/ht/user/internal/domain/user/wechat"
	"gorm.io/gorm"
)

type WechatRepo struct {
	db *gorm.DB
}

func NewWechatRepo() *WechatRepo {
	db, err := mysql.Get()
	if err != nil {
		logger.Fatalf("db获取失败，msg:%s", err)
	}
	return &WechatRepo{db: db}
}

func (repo *WechatRepo) FindByOpenId(ctx context.Context, openId string) (*wechat.UserWechatInfo, error) {
	db := repo.db
	wechatInfo := wechat.UserWechatInfo{}
	tx := db.Where("open_id=?", openId).Take(&wechatInfo)
	return &wechatInfo, tx.Error
}

func (repo *WechatRepo) Save(ctx context.Context, info *wechat.UserWechatInfo) (int64, error) {
	db := repo.db
	db = tx.GetTx(ctx, db)
	db.Create(info)
	return info.ID, db.Error
}
