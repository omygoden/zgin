package smsservice

import (
	"zgin/pkg/sflogger"
	"zgin/thirdparty/sms"
)

//异步批量发送短信
func batchSmsSend(mobile, content string) {
	smsServer := sms.SmsInstance(1)
	//发送短信
	err := smsServer.BatchSend("", mobile, content, "")
	if err != nil {
		sflogger.SmsSendErrorLog(1, 1, "", err.Error())
		return
	}

}
