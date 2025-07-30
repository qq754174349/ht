package tx

import (
	"context"
	"gorm.io/gorm"
)

var txKey = "ormTx"

type manager struct {
	db *gorm.DB
}

func NewTxManager(db *gorm.DB) *manager {
	return &manager{db: db}
}

func (m *manager) Do(ctx context.Context, fn func(ctx context.Context) error) error {
	return m.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txCtx := WithTx(ctx, tx)
		return fn(txCtx)
	})
}

func WithTx(ctx context.Context, tx *gorm.DB) context.Context {
	return context.WithValue(ctx, txKey, tx)
}

func GetTx(ctx context.Context, db *gorm.DB) *gorm.DB {
	if tx, ok := ctx.Value(txKey).(*gorm.DB); ok {
		return tx
	}
	return db
}
