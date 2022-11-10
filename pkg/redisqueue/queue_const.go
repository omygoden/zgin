package redisqueue

const (
	QUEUE_MAX_LEN = 10000 // xdel只是做了删除标记，实际总长度并没有变，所以需要设置最大长度，避免消息堆积
	QUEUE_GROUP_NAME       = "jd-go-group"
	QUEUE_NAME_ORDER       = "jd-go-order-callback"
	QUEUE_NAME_SMS         = "jd-go-sms-callback"
	QUEUE_CONSUMER_NAME_01 = "jd-go-consumer-01"
	QUEUE_CONSUMER_NAME_02 = "jd-go-consumer-02"
	QUEUE_CONSUMER_NAME_03 = "jd-go-consumer-03"
)
