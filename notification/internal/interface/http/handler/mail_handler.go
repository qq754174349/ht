package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/qq754174349/ht/notification/internal/usecase/mail"
)

var (
	mailService = mail.NewUseCase()
)

func SendTextMail(ctx *gin.Context) {
	mailService.SendTextMail(nil)
}
