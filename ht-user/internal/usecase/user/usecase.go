package user

import (
	"context"
	"errors"
	error2 "github.com/qq754174349/ht/ht-frame/common/error"
	grpcClient "github.com/qq754174349/ht/ht-frame/grpc/client"
	_ "github.com/qq754174349/ht/ht-frame/logger"
	log "github.com/qq754174349/ht/ht-frame/logger"
	"github.com/qq754174349/ht/ht-frame/orm/mysql"
	"github.com/qq754174349/ht/ht-frame/orm/tx"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
	mallPd "ht-notification/pkg/grpc/mail"
	"ht-user/internal/common/constant"
	"ht-user/internal/common/utils"
	"ht-user/internal/domain/user/base"
	"ht-user/internal/domain/user/wechat"
	wechatClient "ht-user/internal/infrastructure/client/wechat"
	"ht-user/internal/infrastructure/gorm/user"
	"ht-user/internal/interface/dto/req"
)

type UseCase struct {
	wechatRepo   wechat.UserWechatRepository
	baseRepo     base.UserBaseRepository
	wechatClient *wechatClient.Client
}

func NewUserUseCase() *UseCase {
	return &UseCase{
		wechatRepo:   user.NewWechatRepo(),
		baseRepo:     user.NewBaseRepo(),
		wechatClient: wechatClient.NewClient(),
	}
}

var (
	mallService = grpcClient.New("hn-notification-grpc", mallPd.NewMailServiceClient)
)

// WechatLogin 微信登录
func (user *UseCase) WechatLogin(ctx context.Context, code string) (string, error) {
	session, _ := user.wechatClient.Code2Session(code)

	userWechatInfo, err := user.wechatRepo.FindByOpenId(ctx, session.OpenId)
	if err != nil {
		return "", error2.NewHtErrorFromTemplate(ctx, constant.NoReg)
	}

	return utils.JwtGen(userWechatInfo.UserId), nil
}

// WechatReg 微信注册
func (user *UseCase) WechatReg(ctx context.Context, req req.WechatRegReq) error {
	session, _ := user.wechatClient.Code2Session(req.Code)
	_, err := user.wechatRepo.FindByOpenId(ctx, session.OpenId)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		mysqlDb, _ := mysql.Get()
		tx.NewTxManager(mysqlDb).Do(ctx, func(ctx context.Context) error {
			baseInfo := base.UserBaseInfo{AvatarUrl: req.AvatarUrl, NickName: req.Nickname}
			userId, err := user.baseRepo.Save(ctx, &baseInfo)
			if err != nil {
				return err
			}
			wechatInfo := &wechat.UserWechatInfo{UserId: userId, AvatarUrl: baseInfo.AvatarUrl, NickName: baseInfo.NickName, OpenId: session.OpenId}
			_, err = user.wechatRepo.Save(ctx, wechatInfo)
			return err
		})
	} else if err != nil {
		return err
	} else {
		return error2.NewHtErrorFromTemplate(ctx, constant.RepeatReg)
	}

	return nil
}

// EMailReg 邮箱注册
func EMailReg(ctx context.Context, regReq req.EmailRegReq) {
	resp, err := grpcClient.Invoke(mallService, ctx, func(client mallPd.MailServiceClient, ctx context.Context) (*emptypb.Empty, error) {
		return client.Send(ctx, nil)
	})

	if err != nil {
		log.Info(err)
	}

	log.Info(resp)
}
