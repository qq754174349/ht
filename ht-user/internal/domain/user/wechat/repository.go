package wechat

import (
	"context"
)

type UserWechatRepository interface {
	FindByOpenId(ctx context.Context, openId string) (*UserWechatInfo, error)
	Save(ctx context.Context, entity *UserWechatInfo) (int64, error)
}
