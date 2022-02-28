package sflogger

import (
	"zgin/pkg/constant"
)

//服务异常处理
func ServerErrorLog(errMsg string) {
	Error(constant.LOG_PANIC, errMsg, nil)
}

//服务异常处理
func GoroutineLog(Msg string) {
	Error(constant.LOG_GOROUTINE, Msg, nil)
}

//redis服务器异常处理
func RedisErrorLog(methodName, errMsg string) {
	Error(constant.LOG_REDIS_ERR, "redis错误日志", map[string]interface{}{
		"【操作】":   methodName,
		"【错误信息】": errMsg,
	})
}

//mysql服务器异常处理
func MysqlErrorLog(tableName, errMsg, method, params string) {
	msg := map[string]interface{}{
		"【表名】":   tableName,
		"【错误信息】": errMsg,
		"【操作】":   method,
		"【参数】":   params,
	}
	Error(constant.LOG_MYSQL_ERR, "mysql错误日志", msg)
}

//推入队列失败日志
func QueueAddErrorLog(methodName, msg, errMsg string) {
	m := map[string]interface{}{
		"【操作】":   methodName,
		"【入队信息】": msg,
		"【错误信息】": errMsg,
	}
	Error(constant.LOG_QUEUE_ERR, "rabbitmq错误日志", m)
}

//短信发送日志
func SmsSendLog(userId, jobId int64, mobileStr string) {
	Info(constant.LOG_SMS_SEND, "短信发送日志", map[string]interface{}{
		"【用户ID】": userId,
		"【任务ID】": jobId,
		"【手机号】":  mobileStr,
	})
}

//短信发送日志
func SmsSendErrorLog(userId, jobId int64, mobileStr, errMsg string) {
	Error(constant.LOG_SMS_SEND_ERR, "短信错误日志", map[string]interface{}{
		"【用户ID】": userId,
		"【任务ID】": jobId,
		"【错误信息】": errMsg,
		"【手机号】":  mobileStr,
	})
}

//短信回调原始数据日志
func SmsCallbackLog(errMsg string) {
	Info(constant.LOG_SMS_CALLBACK, "短信回调日志", map[string]interface{}{
		"【相关数据】": errMsg,
	})
}

//短信回调处理成功
func SmsCallbackHandleLog(userId, jobId, mobile interface{}) {
	Info(constant.LOG_SMS_CALLBACK_HANDLE, "短信回调处理成功", map[string]interface{}{
		"【userID】": userId,
		"【jobID】":  jobId,
		"【手机号】":    mobile,
	})
}

//短信回调处理错误
func SmsCallbackHandleErrorLog(responseStr, errMsg string) {
	Error(constant.LOG_SMS_CALLBACK_HANDLE_ERR, "短信错误日志", map[string]interface{}{
		"【错误信息】": errMsg,
		"【相关数据】": responseStr,
	})
}
