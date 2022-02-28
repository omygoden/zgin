package chuanglan

import (
	"zgin/common"
)

func httprequest(urlStr string, params map[string]interface{}, method string) (string, error) {
	header := map[string]string{
		"Content-Type": "application/json;charset=UTF-8",
	}

	res, err := common.HttpRequest(urlStr, method, header, params, "", 5)

	return res, err
}
