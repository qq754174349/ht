package mail

import (
	"github.com/gin-gonic/gin"
	"ht-notification/internal/service/mail"
)

func Send(ctx *gin.Context) {
	mail.Send()
}
