package user

import (
	"context"
	"github.com/qq754174349/ht/ht-frame/orm/mysql"
	"github.com/qq754174349/ht/ht-frame/orm/tx"
	"ht-user/internal/domain/user/wechat"
)

type WechatRepo struct {
}

func NewWechatRepo() *WechatRepo {
	//mysqlDb, err := mysql.Get()
	//if err != nil {
	//	logger.Fatalf("db获取失败，msg:%s", err)
	//}
	return &WechatRepo{}
}

func (repo *WechatRepo) FindByOpenId(ctx context.Context, openId string) (*wechat.UserWechatInfo, error) {
	db, _ := mysql.Get()
	wechatInfo := wechat.UserWechatInfo{}
	tx := db.Where("open_id=?", openId).Take(&wechatInfo)
	return &wechatInfo, tx.Error
}

func (repo *WechatRepo) Save(ctx context.Context, info *wechat.UserWechatInfo) (int64, error) {
	db, _ := mysql.Get()
	db = tx.GetTx(ctx, db)
	db.Create(info)
	return info.ID, db.Error
}
