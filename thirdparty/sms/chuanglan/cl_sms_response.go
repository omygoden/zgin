package chuanglan

//群发短信返回结果
type ClBatchSendResponse struct {
	Code     string `json:"code"` //提交响应状态码，返回“0”表示提交成功（详细参考提交响应状态码）
	MsgId    string `json:"msgId"`
	Time     string `json:"time"` //响应时间,例如："time":"20180519161329"
	ErrorMsg string `json:"errorMsg"`
}

//获取余额返回结果
type ClGetBalanceResponse struct {
	Code     string `json:"code"` //提交响应状态码，返回“0”表示提交成功（详细参考提交响应状态码）
	Balance  string `json:"balance"`
	Time     string `json:"time"` //响应时间,例如："time":"20180519161329"
	ErrorMsg string `json:"errorMsg"`
}

type ClCallbackResponse struct {
	Receiver string `json:"receiver"` //接收验证的用户名，配置时不填写则为空
	Pswd     string `json:"pswd"`     //接收验证的密码，配置时不填则为空

	ReportList []reportList `json:"reportList"`
}

type reportList struct {
	MsgId      string `json:"msgId"`      //消息id
	ReportTime string `json:"reportTime"` //运营商返回的状态更新时间，格式YYMMddHHmm，其中YY=年份的最后两位（00-99）
	Mobile     string `json:"mobile"`     //接收短信的手机号码
	Status     string `json:"status"`     //运营商返回的状态,00000表示成功（详情请前往 code.253.com 查看）
	NotifyTime string `json:"notifyTime"` //253平台收到运营商回复状态报告的时间，格式yyMMddHHmmss
	StatusDesc string `json:"statusDesc"` //状态说明，内容经过URLEncode编码(UTF-8)
	Uid        string `json:"uid"`        //该条短信在您业务系统内的ID，如订单号或者短信发送记录流水号
	Length     string `json:"length"`     //下发短信计费条数
}
