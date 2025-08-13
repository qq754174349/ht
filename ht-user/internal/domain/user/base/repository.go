package base

import (
	"context"
)

type UserBaseRepository interface {
	Save(ctx context.Context, entity *UserBaseInfo) (int64, error)
	UpdateById(ctx context.Context, entity *UserBaseInfo) error
	FindByEmail(ctx context.Context, email string) (*UserBaseInfo, error)
	FindByUsername(ctx context.Context, username string) (*UserBaseInfo, error)
	FindById(ctx context.Context, id int64) (*UserBaseInfo, error)
}
