package sms

import (
	"zgin/thirdparty/sms/chuanglan"
)

const (
	SMS_INSTANCE_TYPE = 1 //创蓝
)

func SmsInstance(smsType int) SmsInterface {
	switch smsType {
	case SMS_INSTANCE_TYPE:
		return chuanglan.NewChuangLanSms()
	default:
		return nil
	}
}
