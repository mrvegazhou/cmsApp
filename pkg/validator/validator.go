package validator

import (
	"cmsApp/internal/constant"
	"cmsApp/pkg/utils/regexpx"
	"fmt"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	zhTranslations "github.com/go-playground/validator/v10/translations/zh"
	"reflect"
	"strings"
	"sync"
)

var isPassword validator.Func = func(fl validator.FieldLevel) bool {
	data, ok := fl.Field().Interface().(string)
	if ok {
		return regexpx.RegPassword(data)
	}
	return true
}

func InitCustomValidator(locale string) {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("isPassword", isPassword)
	}

}

var (
	trans    ut.Translator
	once     sync.Once
	validate *validator.Validate
)

func customZhTrans(v *validator.Validate, trans ut.Translator) {
	var data map[string]string = map[string]string{
		"isPassword": constant.IS_PASSWORD_ERR,
	}
	for key, val := range data {
		//自定义translate与我们的自定义validator配合使用（其实这里也需要把validator与translator进行绑定，v与global.Trans）
		v.RegisterTranslation(key, trans, func(ut ut.Translator) error {
			//自定义失败信息
			return ut.Add(key, val, true) // see universal-translator for details
		}, func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T(key, fe.Field())
			return t
		})
	}
}

func InitTrans(locale string) (ut.Translator, error) {
	var trans ut.Translator
	var err error
	// 修改gin框架中的Validator引擎属性，实现自定制
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("label"), ",", 2)[0]
			fmt.Println(name, "----name-----")
			if name == "-" {
				return ""
			}
			return name
		})

		zhT := zh.New()
		enT := en.New()
		uni := ut.New(enT, zhT, enT)

		var ok bool
		trans, ok = uni.GetTranslator(locale)
		if !ok {
			err = fmt.Errorf("uni.GetTranslator(%s) failed", locale)
			return nil, err
		} else {
			customZhTrans(v, trans)
		}

		// 注册翻译器
		switch locale {
		case "en":
			err = enTranslations.RegisterDefaultTranslations(v, trans)
		case "zh":
			err = zhTranslations.RegisterDefaultTranslations(v, trans)
		default:
			err = enTranslations.RegisterDefaultTranslations(v, trans)
		}
		return trans, err
	}
	return trans, nil
}
