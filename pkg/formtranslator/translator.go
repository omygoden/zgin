package formtranslator

import (
	"fmt"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zh_trans "github.com/go-playground/validator/v10/translations/zh"
	"github.com/omygoden/gotools/sfslice"
	"log"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
	"zgin/global"
)

func InitTrans() {
	translator, _ := ut.New(zh.New()).GetTranslator("zh")

	validatorEngine, _ := binding.Validator.Engine().(*validator.Validate)

	_ = zh_trans.RegisterDefaultTranslations(validatorEngine, translator)

	//自定义定义标题名称
	validatorEngine.RegisterTagNameFunc(func(field reflect.StructField) string {
		label := field.Tag.Get("label")
		if label == "" {
			return field.Name
		}
		return label
	})

	addMobileRule(validatorEngine, &translator)
	addIntLenRule(validatorEngine, &translator)
	addDatetimeRule(validatorEngine, &translator)
	addIntInRule(validatorEngine, &translator)

	global.Trans = translator
}

//添加手机号规则
func addMobileRule(validatorEngine *validator.Validate, translator *ut.Translator) {
	registerCustomRule(validatorEngine, translator, "mobile_rule", "格式错误", func(fl validator.FieldLevel) bool {
		var res string
		reg := regexp.MustCompile("^1[\\d]{10}$")

		switch fl.Field().Kind().String() {
		case "int64":
			res = reg.FindString(strconv.Itoa(int(fl.Field().Int())))
		case "string":
			res = reg.FindString(fl.Field().String())
		}
		if res != "" {
			return true
		}
		return false
	})
}

//添加int类型长度限制规则
func addIntLenRule(validatorEngine *validator.Validate, translator *ut.Translator) {
	registerCustomRule(validatorEngine, translator, "int_len_rule", "长度有误", func(fl validator.FieldLevel) bool {
		ruleValue := strconv.Itoa(int(fl.Field().Int()))
		ruleParam, _ := strconv.Atoi(fl.Param())
		if len(ruleValue) == ruleParam {
			return true
		}
		return false
	})
}

//添加datetime校验规格 -- 有值的情况下才校验，否不校验，示例：date_time_rule=2006-01-02 15:04:05
func addDatetimeRule(validatorEngine *validator.Validate, translator *ut.Translator) {
	registerCustomRule(validatorEngine, translator, "datetime_rule", "时间格式有误", func(fl validator.FieldLevel) bool {
		if fl.Field().String() != "" {
			_, err := time.Parse(fl.Param(), fl.Field().String())
			if err != nil {
				return false
			}
		}
		return true
	})
}

//添加int包含规则 -- 规则参数用-连接，例如：int_in_rule=1-2-3,判断该参数是否在123中
func addIntInRule(validatorEngine *validator.Validate, translator *ut.Translator) {
	registerCustomRule(validatorEngine, translator, "int_in_rule", "参数有误", func(fl validator.FieldLevel) bool {
		ruleValue := strconv.Itoa(int(fl.Field().Int()))
		ruleParams := strings.Split(fl.Param(), "-")
		if len(ruleParams) == 0 {
			return true
		}

		if sfslice.SliceContainString(ruleParams, ruleValue) {
			return true
		}
		return false
	})
}

func registerCustomRule(validatorEngine *validator.Validate, translator *ut.Translator, ruleName string, errHint string, fn validator.Func) {
	//定义规则
	err := validatorEngine.RegisterValidation(ruleName, fn)
	//翻译规则
	_ = validatorEngine.RegisterTranslation(ruleName, *translator, func(ut ut.Translator) error {
		return ut.Add(ruleName, "{0}"+errHint, false)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T(ruleName, fe.Field())
		return t
	})

	if err != nil {
		log.Println(fmt.Sprintf("自定义【%s】规则注册失败)：%s", ruleName, err.Error()))
		os.Exit(1)
	}
}
