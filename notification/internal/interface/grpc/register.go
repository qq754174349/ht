package grpc

import (
	"github.com/qq754174349/ht/ht-frame/grpc/service"
	"github.com/qq754174349/ht/notification/internal/interface/grpc/handler"
)

func Register() {
	service.RegisterRegistrant(&handler.Register{})
	service.Bootstrap()
}
