package sms

type SmsInterface interface {
	//单发短信
	SignalSend(sign,mobile, content, extraParams string) error
	//群发短信
	BatchSend(sign,mobile, content, extraParams string) error
	//获取短信余额
	GetBalance() (int,error)
	//修改配置
	SetConfig(config map[string]interface{})
}
