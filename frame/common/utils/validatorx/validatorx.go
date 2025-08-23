package validatorx

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
)

// 全局翻译器
var trans ut.Translator

// 初始化翻译器
func init() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// 中文翻译器
		zhCn := zh.New()
		uni := ut.New(zhCn, zhCn)
		trans, _ = uni.GetTranslator("zh")

		// 注册默认中文翻译
		if err := zh_translations.RegisterDefaultTranslations(v, trans); err != nil {
			panic(err)
		}

		// 统一替换所有字段为中文别名
		v.RegisterTagNameFunc(func(field reflect.StructField) string {
			if alias := field.Tag.Get("label"); alias != "" {
				return alias
			}
			return field.Name
		})
	}
}

// TranslateError 翻译错误信息
func TranslateError(err error) string {
	var errs validator.ValidationErrors
	if errors.As(err, &errs) {
		return errs[0].Translate(trans)
	}
	return fmt.Sprintf("参数错误: %v", err)
}
