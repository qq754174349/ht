package asser

import (
	"context"
	"errors"
	"fmt"

	htError "github.com/qq754174349/ht/ht-frame/common/error"
	"github.com/qq754174349/ht/ht-frame/logger"
	"gorm.io/gorm"
)

// That 断言条件为 true，否则返回指定 error
func That(condition bool, err error) error {
	if !condition {
		return err
	}
	return nil
}

// Thatf 断言条件为 true，否则返回格式化的 error
func Thatf(condition bool, format string, args ...any) error {
	if !condition {
		return fmt.Errorf(format, args...)
	}
	return nil
}

// CtxThat 断言条件为 true，否则返回指定的 HtError
func CtxThat(ctx context.Context, condition bool, code int, msg string) error {
	if !condition {
		return htError.NewHtError(code, msg)
	}
	return nil
}

// CtxThatf 断言条件为 true，否则返回格式化的 HtError
func CtxThatf(ctx context.Context, condition bool, code int, format string, args ...any) error {
	if !condition {
		return htError.NewHtError(code, fmt.Sprintf(format, args...))
	}
	return nil
}

// Must 断言条件为 true，否则 panic
func Must(condition bool, errMsg string) {
	if !condition {
		panic(errMsg)
	}
}

// Nil 断言 err == nil，否则返回它
func Nil(err error) error {
	if err != nil {
		return err
	}
	return nil
}

// GormErr 处理 GORM 查询错误，自动忽略 RecordNotFound 错误（可选返回 nil）
func GormErr(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil
	}
	return err
}

// CtxGormErr 和 HtError 结合
func CtxGormErr(err error, allowNotFound bool) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) && allowNotFound {
		return nil
	}
	return htError.NewHtError(500, err.Error())
}

// MustNoError 如果 err 不为 nil，panic 抛出
func MustNoError(err error) {
	if err != nil {
		logger.Fatal("unexpected error: %v", err)
	}
}

// Code 条件断言失败时返回带错误码的 error.HtError
func Code(condition bool, code int, msg string) error {
	if !condition {
		return htError.NewHtError(code, msg)
	}
	return nil
}

// Codef 条件断言失败时返回格式化 error.HtError
func Codef(condition bool, code int, format string, args ...any) error {
	if !condition {
		return htError.NewHtError(code, fmt.Sprintf(format, args...))
	}
	return nil
}

// CtxCode 条件断言失败时返回带上下文的 HtError
func CtxCode(condition bool, code int, msg string) error {
	if !condition {
		return htError.NewHtError(code, msg)
	}
	return nil
}

// CtxCodef 条件断言失败时返回格式化带上下文的 HtError
func CtxCodef(condition bool, code int, format string, args ...any) error {
	if !condition {
		return htError.NewHtError(code, fmt.Sprintf(format, args...))
	}
	return nil
}
