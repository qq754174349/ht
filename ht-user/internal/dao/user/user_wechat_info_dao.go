package user

import (
	"github.com/qq754174349/ht-frame/mysql"
	"ht-user/internal/model"
)

func QueryUserWechatInfo(openId string) *model.UserWechatInfo {
	wechatInfo := model.UserWechatInfo{}
	mysqlDb, _ := mysql.Get()
	tx := mysqlDb.Where("open_id=?", openId).Take(&wechatInfo)
	if tx.Error != nil {
		return nil
	}
	return &wechatInfo
}
