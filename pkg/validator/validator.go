package validator

//https://github.com/Harsh-Katiyar/JetBrains-products-Activation-code-until-13-Dec-2024
//https://github.com/Naimul007A/jetbrains
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
	return regexpx.RegPassword(fl.Field().String())
}

var CustomTranslator ut.Translator
var once sync.Once

func registerTagNameJSON(validate *validator.Validate) {
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("label"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
}

func registerTranslator(locale string, validate *validator.Validate) {
	zhT := zh.New()
	enT := en.New()
	uni := ut.New(enT, zhT, enT)
	found := false
	CustomTranslator, found = uni.GetTranslator(locale)
	if !found {
		fmt.Println("translator not found")
		panic(found)
	}
	// 注册翻译器
	var err error
	switch locale {
	case "en":
		err = enTranslations.RegisterDefaultTranslations(validate, CustomTranslator)
	case "zh":
		err = zhTranslations.RegisterDefaultTranslations(validate, CustomTranslator)
	default:
		err = enTranslations.RegisterDefaultTranslations(validate, CustomTranslator)
	}
	if err != nil {
		panic((err).Error())
	}
}

func registerCustomsValidation(validate *validator.Validate) error {
	var data map[string][2]interface{} = map[string][2]interface{}{
		"isPassword": [2]interface{}{constant.IS_PASSWORD_ERR, isPassword},
	}

	// 在校验器注册自定义的校验方法
	for key, val := range data {
		err := validate.RegisterValidation(key, val[1].(validator.Func))
		if err != nil {
			return err
		}
		if err = validate.RegisterTranslation(key, CustomTranslator, func(trans ut.Translator) error { // registerTranslator 为自定义字段添加翻译功能
			if err := trans.Add(key, val[0].(string), false); err != nil {
				return err
			}
			return nil
		}, func(trans ut.Translator, fe validator.FieldError) string { // translate 自定义字段的翻译方法
			msg, err := trans.T(fe.Tag(), fe.Field())
			if err != nil {
				panic(fe.(error).Error())
			}
			return msg

		}); err != nil {
			panic((err).Error())
			return err
		}
	}
	return nil
}

func InitTrans(locale string) error {

	once.Do(func() {
		validate := binding.Validator.Engine().(*validator.Validate)
		registerTagNameJSON(validate)
		registerTranslator(locale, validate)
		_ = registerCustomsValidation(validate)
	})

	return nil
}
