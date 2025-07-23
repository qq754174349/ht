package user

import (
	"context"
	error2 "github.com/qq754174349/ht-frame/common/error"
	grpcClient "github.com/qq754174349/ht-frame/grpc/client"
	log "github.com/qq754174349/ht-frame/logger"
	"github.com/qq754174349/ht-frame/mysql"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
	notificationPd "ht-notification/pkg/grpc/mail"
	"ht-user/internal/common/constant"
	"ht-user/internal/common/utils"
	"ht-user/internal/dao/user"
	req2 "ht-user/internal/dto/req"
	model2 "ht-user/internal/model"
	"ht-user/internal/service"
	"time"
)

// WechatLogin 微信登录
func WechatLogin(ctx context.Context, code string) (string, error) {
	session, _ := service.Code2Session(code)

	userWechatInfo := user.QueryUserWechatInfo(session.OpenId)
	if userWechatInfo == nil {
		return "", error2.NewHtErrorFromTemplate(ctx, constant.NoReg)
	}

	return utils.JwtGen(userWechatInfo.UserId), nil
}

// WechatReg 微信注册
func WechatReg(ctx context.Context, req req2.WechatRegReq) error {
	session, _ := service.Code2Session(req.Code)
	userWechatInfo := user.QueryUserWechatInfo(session.OpenId)
	if userWechatInfo != nil {
		return error2.NewHtErrorFromTemplate(ctx, constant.RepeatReg)
	}
	baseInfo := model2.BaseInfo{AvatarUrl: req.AvatarUrl, Nickname: req.Nickname}
	mysqlDb, _ := mysql.Get()
	err := mysqlDb.Transaction(func(tx *gorm.DB) error {
		tx.Create(&baseInfo)
		wechatInfo := model2.UserWechatInfo{UserId: baseInfo.ID, AvatarUrl: baseInfo.AvatarUrl, NickName: baseInfo.Nickname, OpenId: session.OpenId}
		tx.Create(&wechatInfo)
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

// EMailReg 邮箱注册
func EMailReg(ctx context.Context, regReq req2.EmailRegReq) {
	conn := grpcClient.GetConn("hn-notification-grpc")
	defer conn.Close()
	c := notificationPd.NewMailServiceClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := c.Send(ctx, &emptypb.Empty{})
	if err != nil {
		log.Error(err)
	}
}
