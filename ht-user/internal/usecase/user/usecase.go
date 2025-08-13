package user

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	error2 "github.com/qq754174349/ht/ht-frame/common/error"
	"github.com/qq754174349/ht/ht-frame/common/utils/asser"
	"github.com/qq754174349/ht/ht-frame/common/utils/jwt"
	grpcClient "github.com/qq754174349/ht/ht-frame/grpc/client"
	_ "github.com/qq754174349/ht/ht-frame/logger"
	log "github.com/qq754174349/ht/ht-frame/logger"
	"github.com/qq754174349/ht/ht-frame/orm/mysql"
	"github.com/qq754174349/ht/ht-frame/orm/tx"
	mallPd "github.com/qq754174349/ht/ht-notification/pkg/pd/mail"
	"github.com/qq754174349/ht/ht-user/internal/common/constant"
	"github.com/qq754174349/ht/ht-user/internal/domain/user/base"
	"github.com/qq754174349/ht/ht-user/internal/domain/user/wechat"
	wechatClient "github.com/qq754174349/ht/ht-user/internal/infrastructure/client/wechat"
	"github.com/qq754174349/ht/ht-user/internal/infrastructure/gorm/user"
	"github.com/qq754174349/ht/ht-user/internal/interface/dto/req"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
)

var mailBody = "<html lang=\"zh-CN\">\n<head>\n    <meta charset=\"UTF-8\">\n    <meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\">\n    <title>激活您的账户</title>\n    <style type=\"text/css\">\n        body {\n            font-family: Arial, Helvetica, sans-serif;\n            line-height: 1.6;\n            color: #333;\n            margin: 0;\n            padding: 0;\n            background-color: #f9f9f9;\n        }\n        .container {\n            max-width: 600px;\n            margin: 30px auto;\n            background-color: #ffffff;\n            border-radius: 8px;\n            padding: 30px;\n            box-shadow: 0 2px 8px rgba(0,0,0,0.05);\n        }\n        .header {\n            text-align: center;\n            margin-bottom: 25px;\n        }\n        .logo {\n            max-width: 140px;\n            height: auto;\n        }\n        h1 {\n            font-size: 20px;\n            color: #222;\n            margin-top: 15px;\n            font-weight: 600;\n        }\n        .button {\n            display: inline-block;\n            padding: 12px 28px;\n            background-color: #007BFF;\n            color: #ffffff !important;\n            text-decoration: none;\n            border-radius: 6px;\n            font-weight: bold;\n            margin: 20px 0;\n        }\n        .code-box {\n            background: #f3f4f6;\n            padding: 10px;\n            word-break: break-all;\n            border-radius: 4px;\n            font-family: monospace;\n            font-size: 13px;\n            color: #444;\n        }\n        .footer {\n            margin-top: 30px;\n            font-size: 12px;\n            color: #888;\n            text-align: center;\n            line-height: 1.5;\n        }\n        hr {\n            border: none;\n            border-top: 1px solid #eee;\n            margin: 25px 0;\n        }\n    </style>\n</head>\n<body>\n    <div class=\"container\">\n        <div class=\"header\">\n            <img src=\"https://yourwebsite.com/logo.png\" alt=\"Company Logo\" class=\"logo\">\n            <h1>请激活您的账户</h1>\n        </div>\n\n        <p>尊敬的<strong>#{userName}</strong>：</p>\n        <p>感谢您注册 <strong>ht百货</strong>！请点击下方按钮完成账户激活：</p>\n\n        <div style=\"text-align: center;\">\n            <a href=\"http://localhost:9000/api/user/activate?token=#{token}\" class=\"button\">立即激活账户</a>\n        </div>\n\n        <p>如果按钮无法点击，请复制以下链接到浏览器：</p>\n        <div class=\"code-box\">\n            http://localhost:9000/api/user/activate?token=#{token}\n        </div>\n\n        <p><strong>激活链接24小时内有效</strong>，请尽快完成操作。</p>\n\n        <hr>\n\n        <p>如有问题，请联系客服：<a href=\"mailto:support@yourdomain.com\">support@yourdomain.com</a></p>\n\n        <div class=\"footer\">\n            <p>© #{time} ht科技. 保留所有权利。</p>\n            <p>如果您并未注册此账户，请忽略本邮件。</p>\n        </div>\n    </div>\n</body>\n</html>\n"

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

	return jwt.Gen(map[string]interface{}{"userId": userWechatInfo.UserId}, time.Hour, ""), nil
}

// WechatReg 微信注册
func (user *UseCase) WechatReg(ctx context.Context, req req.WechatRegReq) error {
	session, _ := user.wechatClient.Code2Session(req.Code)
	_, err := user.wechatRepo.FindByOpenId(ctx, session.OpenId)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		mysqlDb, _ := mysql.Get()
		tx.NewTxManager(mysqlDb).Do(ctx, func(ctx context.Context) error {
			baseInfo := base.UserBaseInfo{AvatarUrl: &req.AvatarUrl, NickName: &req.Nickname}
			userId, err := user.baseRepo.Save(ctx, &baseInfo)
			if err != nil {
				return err
			}
			wechatInfo := &wechat.UserWechatInfo{UserId: userId, AvatarUrl: *baseInfo.AvatarUrl, NickName: *baseInfo.NickName, OpenId: &session.OpenId}
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
func (user *UseCase) EMailReg(ctx context.Context, regReq req.EmailRegReq) *error2.HtError {
	_, err := user.baseRepo.FindByEmail(ctx, regReq.Email)
	if err := asser.GormErr(err); err != nil {
		log.Errorf("查询用户失败, err:%s", err)
		return error2.NewBaseHtError(ctx)
	}
	if err := asser.CtxCode(ctx, err != nil && errors.Is(err, gorm.ErrRecordNotFound), 500, "用户邮箱已存在"); err != nil {
		return err
	}

	_, err = user.baseRepo.FindByUsername(ctx, regReq.Username)
	if err := asser.GormErr(err); err != nil {
		log.Errorf("查询用户失败, err:%s", err)
		return error2.NewBaseHtError(ctx)
	}
	if err := asser.CtxCode(ctx, err != nil && errors.Is(err, gorm.ErrRecordNotFound), 500, "用户邮箱已存在"); err != nil {
		return err
	}

	userId, err := user.baseRepo.Save(ctx, &base.UserBaseInfo{
		Email:    &regReq.Email,
		Username: &regReq.Username,
		Password: &regReq.Password,
	})

	resp, err := grpcClient.Invoke(mallService, ctx, func(client mallPd.MailServiceClient, ctx context.Context) (*emptypb.Empty, error) {
		token := jwt.Gen(map[string]interface{}{"userId": userId}, 24*time.Hour, "")
		mail := strings.ReplaceAll(mailBody, "#{userName}", regReq.Username)
		mail = strings.ReplaceAll(mail, "#{token}", token)
		mail = strings.ReplaceAll(mail, "#{time}", strconv.Itoa(time.Now().Year()))
		return client.SendTextMail(ctx, &mallPd.TextMailReq{To: regReq.Email, Subject: "请激活您的账户 - 立即完成注册！", Body: mail})
	})

	if err != nil {
		log.Info(err)
	}

	log.Info(resp)
	return nil
}

// UserActivate 激活用户
func (user *UseCase) UserActivate(ctx *gin.Context, token string) error {
	data, fla := jwt.Parse(token, "")
	if !fla {
		return error2.NewHtErrorFromTemplate(ctx, constant.ActivateExpire)
	}

	userId := data["userId"].(float64)
	userBaseInfo, err := user.baseRepo.FindById(ctx, int64(userId))
	if err := asser.CtxCode(ctx, !errors.Is(err, gorm.ErrRecordNotFound), constant.NoUser.Code, constant.NoUser.Msg); err != nil {
		return err
	}

	if userBaseInfo.Status == 1 {
		return error2.NewHtErrorFromTemplate(ctx, constant.RepeatActivate)
	}

	now := time.Now()
	userBaseInfo.ActivateTime = &now
	userBaseInfo.Status = 1
	err = user.baseRepo.UpdateById(ctx, userBaseInfo)
	if err != nil {
		return error2.NewHtErrorFromMsg(ctx, err.Error())
	}
	return nil
}
