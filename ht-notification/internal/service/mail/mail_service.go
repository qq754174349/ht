package mail

import (
	"context"
	"fmt"
	log "github.com/qq754174349/ht-frame/logger"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"gopkg.in/gomail.v2"
	"ht-notification/internal/config/mail"
	pb "ht-notification/pkg/grpc/mail"
)

type Service struct {
	pb.UnimplementedMailServiceServer
}

type Register struct {
}

func (c *Register) Register(server *grpc.Server) {
	pb.RegisterMailServiceServer(server, &Service{})
}

func Send() {
	mailConfig := mail.GetConfig()
	body := `
		<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>请激活您的账户</title>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { text-align: center; margin-bottom: 20px; }
        .logo { height: 50px; }
        .button { 
            display: inline-block; padding: 12px 24px; background-color: #4CAF50; 
            color: white !important; text-decoration: none; border-radius: 4px; 
            font-weight: bold; margin: 20px 0; 
        }
        .footer { margin-top: 30px; font-size: 12px; color: #777; text-align: center; }
    </style>
</head>
<body>	
    <div class="header">
        <img src="https://your-website.com/logo.png" alt="爱农保" class="logo">
        <h2>欢迎加入 爱农保！</h2>
    </div>

    <p>感谢您注册 爱农保 账户。请点击下方按钮完成邮箱验证：</p>

    <div style="text-align: center;">
        <a href="{{.www.baidu.com}}" class="button">立即激活账户</a>
    </div>

    <p>如果按钮无效，请复制以下链接到浏览器地址栏访问：</p>
    <p style="word-break: break-all;"><a href="{{.www.baidu.com}}">{{.www.baidu.com}}</a></p>

    <div class="footer">
        <p>此链接将在 <strong>24 小时</strong> 后失效。</p>
        <p>如果您未注册 爱农保，请忽略此邮件。</p>
        <p>需要帮助？请联系 <a href="mailto:support@your-app.com">support@your-app.com</a></p>
    </div>
</body>
</html>
	`
	// 创建邮件消息
	m := gomail.NewMessage()
	m.SetHeader("From", mailConfig.Username)
	m.SetHeader("To", "754174349@qq.com")
	m.SetHeader("Subject", "【爱农保】请激活您的账户")
	m.SetBody("text/html", body)
	m.AddAlternative("text/plain",
		fmt.Sprintf("欢迎加入 爱农保！\n\n请访问以下链接激活账户：\n%s\n\n此链接24小时内有效。", "www.baidu.com"),
	)

	// 发送邮件
	d := gomail.NewDialer(mailConfig.Host, mailConfig.Port, mailConfig.Username, mailConfig.Password)
	d.SSL = true
	if err := d.DialAndSend(m); err != nil {
		fmt.Errorf("邮件发送失败: %v", err)
	}
}

// Send 实现邮件发送逻辑
func (*Service) Send(ctx context.Context, req *emptypb.Empty) (*emptypb.Empty, error) {
	log.Info("wo shi 1 hao ,cheng gong diao yong ")
	//Send()
	return &emptypb.Empty{}, nil
}
