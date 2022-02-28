package chuanglan

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"zgin/global"
)

type ChuangLanSms struct {
	Sign      string
	Content   string
	Mobile    string
	Uid       string
	smsConfig map[string]interface{}
}

const (
	batchSend  = "/msg/send/json"
	getBalance = "/msg/balance/json"
)

//var _ sms.SmsInterface = &ChuangLanSms{}

func NewChuangLanSms() *ChuangLanSms {
	return &ChuangLanSms{
		smsConfig: map[string]interface{}{
			"account":  global.Config.Sms.Account,
			"password": global.Config.Sms.Password,
		},
	}
}

//单发短信
func (this *ChuangLanSms) SignalSend(sign, mobile, content, extraParams string) error {
	return nil
}

//群发短信
func (this *ChuangLanSms) BatchSend(sign, mobile, content, extraParams string) error {
	var urlStr = global.Config.Sms.Domain + batchSend
	params := map[string]interface{}{
		"account":  this.smsConfig["account"],
		"password": this.smsConfig["password"],
		"msg":      content,
		"phone":    mobile,
		"report":   "true",
		"uid":      extraParams,
	}
	res, err := httprequest(urlStr, params, "POST")
	if err != nil {
		return err
	}

	var response ClBatchSendResponse
	_ = json.Unmarshal([]byte(res), &response)
	if response.Code != "0" {
		return errors.New(response.ErrorMsg)
	}

	return nil
}

//获取短信余额
func (this *ChuangLanSms) GetBalance() (int, error) {
	urlStr := global.Config.Sms.Domain + getBalance

	params := map[string]interface{}{
		"account":  this.smsConfig["account"],
		"password": this.smsConfig["password"],
	}
	res, err := httprequest(urlStr, params, "POST")
	if err != nil {
		log.Println(fmt.Sprintf("获取短信余额失败，错误信息：%s", err.Error()))
		return 0, err
	}

	var response ClGetBalanceResponse
	_ = json.Unmarshal([]byte(res), &response)
	if response.Code != "0" {
		log.Println(fmt.Sprintf("获取短信余额失败，错误信息：%s", response.ErrorMsg))
		return 0, errors.New(response.ErrorMsg)
	}

	return strconv.Atoi(response.Balance)
}

//修改配置
func (this *ChuangLanSms) SetConfig(config map[string]interface{}) {
	for k, _ := range this.smsConfig {
		if v, ok := config[k]; ok {
			this.smsConfig[k] = v
		}
	}
}
