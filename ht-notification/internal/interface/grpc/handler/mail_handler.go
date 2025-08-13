package handler

import (
	"context"

	"github.com/qq754174349/ht/ht-notification/internal/interface/dto/mail/req"
	"github.com/qq754174349/ht/ht-notification/internal/usecase/mail"
	pb "github.com/qq754174349/ht/ht-notification/pkg/pd/mail"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

var mailUseCase = mail.NewUseCase()

type Service struct {
	pb.UnimplementedMailServiceServer
}

type Register struct {
}

func (c *Register) Register(server *grpc.Server) {
	pb.RegisterMailServiceServer(server, &Service{})
}

// SendTextMail 实现邮件发送逻辑
func (*Service) SendTextMail(ctx context.Context, q *pb.TextMailReq) (*emptypb.Empty, error) {
	mailUseCase.SendTextMail(&req.TextMailReq{Body: q.Body, Subject: q.Subject, To: q.To})
	return &emptypb.Empty{}, nil
}
