package base

import "context"

type UserBaseRepository interface {
	Save(ctx context.Context, entity *UserBaseInfo) (int64, error)
}
