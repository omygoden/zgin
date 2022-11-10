package constant

// redis key
const (
	REDIS_KEY_PREFIX     = "go-ks-"
	REDIS_KEY_ORDER_LOCK = REDIS_KEY_PREFIX + "order-%d-%d-lock"
)

// 日志名称
const (
	LOG_REQUEST     = "request.log"       // 接口请求日志
	LOG_REQUEST_ERR = "request_error.log" // 接口请求错误日志
	LOG_PANIC       = "panic.log"         // 异常日志--主要记录异步处理时候的异常
	LOG_REDIS_ERR   = "redis_error.log"   // redis错误日志
	LOG_MYSQL_ERR   = "mysql_error.log"   // mysql错误日志
	LOG_MYSQL       = "mysql.log"         // mysql日志
	LOG_MYSQL_SLOW  = "mysql_slow.log"    // mysql慢日志
	LOG_QUEUE_ERR   = "queue_error.log"   // 入队列错误日志
	LOG_GOROUTINE   = "goroutine.log"     // 协程变化日志

	// 短信相关日志
	LOG_SMS_SEND                = "sms_send.log"                // 短信发送日志
	LOG_SMS_SEND_ERR            = "sms_send_error.log"          // 短信发送日志
	LOG_SMS_CALLBACK            = "sms_callback.log"            // 短信回调日志
	LOG_SMS_CALLBACK_HANDLE     = "sms_callback_handle.log"     // 短信回调处理成功
	LOG_SMS_CALLBACK_HANDLE_ERR = "sms_callback_handle_err.log" // 短信回调处理错误日志
)

// rabbitmq--exchange/queue/route对应名称
const (
	RB_KEY_EXCHANGE_SMS_CALLBACK = "go.sms.exchange.sms.callback" // 短信回调通知 -- 交换器名称
	RB_KEY_QUEUE_SMS_CALLBACK    = "go.sms.queue.sms.callback"    // 短信回调通知 -- 队列名称
)
