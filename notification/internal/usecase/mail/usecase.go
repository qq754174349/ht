package mail

import (
	"fmt"

	htError "github.com/qq754174349/ht/ht-frame/common/error"
	"github.com/qq754174349/ht/notification/internal/infrastructure/config/mail"
	"github.com/qq754174349/ht/notification/internal/interface/dto/mail/req"
	"gopkg.in/gomail.v2"
)

type UseCase struct {
}

func NewUseCase() *UseCase {
	return &UseCase{}
}

func (*UseCase) SendTextMail(req *req.TextMailReq) *htError.HtError {
	mailConfig := mail.GetConfig()

	// 创建邮件消息
	m := gomail.NewMessage()
	m.SetHeader("From", mailConfig.Username)
	m.SetHeader("To", req.To)
	m.SetHeader("Subject", req.Subject)
	m.SetBody("text/html", req.Body)

	// 发送邮件
	d := gomail.NewDialer(mailConfig.Host, mailConfig.Port, mailConfig.Username, mailConfig.Password)
	d.SSL = true
	if err := d.DialAndSend(m); err != nil {
		fmt.Errorf("邮件发送失败: %v", err)
	}
	return nil
}
