package app

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"strings"
	"zgin/global"
	"zgin/pkg/errcode"
)

type ValidError struct {
	Key     string
	Message string
}

type ValidErrors []*ValidError

func (v *ValidError) Error() string {
	return v.Message
}

func (v ValidErrors) Errors() []string {
	var errs []string
	for _, err := range v {
		errs = append(errs, err.Error())
	}

	return errs
}

func (v ValidErrors) Error() string {
	return strings.Join(v.Errors(), ",")
}

// BindAndValid binds and validates data
func BindAndValid(ctx *gin.Context, params interface{}) {
	var errs ValidErrors
	err := ctx.ShouldBind(params)
	if err != nil {
		pjson, _ := json.Marshal(params)
		verrs, ok := err.(validator.ValidationErrors)
		if !ok {
			//传参字段类型校验
			if jsonValue, jsonOk := err.(*json.UnmarshalTypeError); jsonOk {
				panic(errcode.ApiPanic{400, fmt.Sprintf("%s类型有误", jsonValue.Field), string(pjson)})
			}

			panic(errcode.ApiPanic{400, err.Error(), string(pjson)})
		}

		for key, value := range verrs.Translate(global.Trans) {
			errs = append(errs, &ValidError{
				Key:     key,
				Message: value,
			})
		}

		if errs != nil {
			panic(errcode.ApiPanic{400, errs.Error(), string(pjson)})
		}
	}
}
