package main

import (
	"encoding/json"
	"fmt"
	"github.com/omygoden/gotools/encrypts"
	"github.com/omygoden/gotools/sfconst"
	"github.com/streadway/amqp"
	"gorm.io/gorm"
	"net/url"
	"strconv"
	"strings"
	"time"
	"zgin/global"
	"zgin/internal/model"
	"zgin/internal/model/branchmodel"
	"zgin/internal/services/statistics"
	"zgin/pkg/constant"
	"zgin/pkg/rabbitmq/rbmqserver/customers"
	"zgin/pkg/redisclient"
	"zgin/pkg/setting"
	"zgin/pkg/sflogger"
	"zgin/thirdparty/sms/chuanglan"
)

const (
	SMS_CALLBACK_MAX_DB_CONN = 300  //数据库最大连接数
	SMS_CALLBACK_MIN_DB_CONN = 200  //数据库空闲连接数
	SMS_CALLBACK_MAX_RB_CHAN = 200  //rabbitmq最大channel
	SMS_CALLBACK_MAX_RB_QPS  = 1000 //rabbitmq消费qp，最多消息处理量，超过该值则会阻塞

)

func init() {
	setting.QueueInitSettting(SMS_CALLBACK_MAX_DB_CONN, SMS_CALLBACK_MIN_DB_CONN)
	global.RabbitmqClient.Config.ChannelMax = SMS_CALLBACK_MAX_RB_CHAN
}

//短信回调信息处理队列 -- 放supervisor里执行
func main() {
	//账户冻结余额处理
	go userAccountHandle()

	_ = customers.QueueCustomerFanout(constant.RB_KEY_EXCHANGE_SMS_CALLBACK, constant.RB_KEY_QUEUE_SMS_CALLBACK, SMS_CALLBACK_MAX_RB_QPS, func(msg *amqp.Delivery) {
		var uniqueKey = encrypts.Sha1(string(msg.Body))
		var response chuanglan.ClCallbackResponse
		var redisKey = fmt.Sprintf(constant.REDIS_KEY_CL_SMS_CALLBACK_MIDENG, uniqueKey)
		_ = json.Unmarshal(msg.Body, &response)

		//幂等性处理
		redisclient.SetNxWait(redisKey, 1, time.Minute)
		defer func() {
			redisclient.Del(redisKey)
			_ = msg.Ack(false)
		}()

		/**
		短信uid格式： userId:任务表id:任务类型:短信类型:短信表ID
		任务类型：1短信群发任务，2提醒短信任务，3营销短信任务
		短信类型 ：0测试 1营销 2催付 3发货提醒 4催评提醒 5付款关怀 6延迟发货提醒 7下单提醒 8签收提醒 9好评提醒 400物流发货提醒 401抵达同城提醒 402派送提醒 403取件提醒 404快递签收 405物流异常  601申请售后 602同意售后 603拒绝售后 604退款成功'
		短信表ID：针对任务类型为提醒短信的时候，需要传，其他类型默认传0
		*/
		for _, v := range response.ReportList {
			uid := strings.Split(v.Uid, ":")
			if len(uid) != 5 {
				sflogger.SmsCallbackHandleErrorLog(string(msg.Body), "短信回调格式有误")
				continue
			}
			//1:75:1:1:0
			userId, _ := strconv.ParseInt(uid[0], 10, 64)
			jobId, _ := strconv.ParseInt(uid[1], 10, 64)
			jobType := uid[2]
			smsType := uid[3]
			smsId, _ := strconv.ParseInt(uid[4], 10, 64)

			var smsMarketJobModel = model.SmsMarketMessageJobs{ID: jobId}
			var smsMarketOrderModel = model.SmsMarketMessageOrders{JobId: jobId, Phone: v.Mobile}
			var smsCustomJobModel = model.SmsCustomMessageJobs{ID: jobId}
			var smsCustomOrderModel = model.SmsCustomMessageOrders{JobId: jobId, Phone: v.Mobile}
			//var smsRemindJobModel = model.SmsRemindMessageJobs{ID: jobId}
			var smsRemindOrderModel = model.SmsRemindMessageOrders{ID: smsId, JobId: jobId, Phone: v.Mobile}
			var smsCustomerModel = branchmodel.SmsCustomers{UserId: userId, RealPhone: v.Mobile}

			var status int
			//todo 其他错误码信息：https://code.253.com/
			if v.Status != "DELIVRD" {
				status = 3
			} else {
				status = 2
			}

			//根据不同任务类型，修改对应的订单状态
			var tempId int64
			var charges int
			var res bool
			var sendTime time.Time

			statusDesc, _ := url.QueryUnescape(v.StatusDesc)
			switch jobType {
			case "1":
				orderModel := smsCustomOrderModel.GetOne([]string{"id", "status", "template_id", "real_charges", "need_charges"})
				if orderModel.ID == 0 {
					sflogger.SmsCallbackHandleErrorLog(string(msg.Body), "未发现短信发送记录"+smsCustomOrderModel.TableName())
					return
				}
				if orderModel.Status != 2 && orderModel.Status != 3 {
					smsCustomOrderModel.MarkStatus(status, statusDesc)
					tempId = orderModel.TemplateId
					charges = orderModel.NeedCharges
					res = true
					sendTime = orderModel.SendTime.Time

					//短信任务计数
					if status == 2 {
						smsCustomJobModel.IncrSuccess()
					} else {
						smsCustomJobModel.IncrFail()
					}
				}
			case "2":
				orderModel := smsRemindOrderModel.GetOne([]string{"id", "status", "template_id", "charges"})
				if orderModel.ID == 0 {
					sflogger.SmsCallbackHandleErrorLog(string(msg.Body), "未发现短信发送记录"+smsRemindOrderModel.TableName())
					return
				}
				if orderModel.Status != 2 && orderModel.Status != 3 {
					smsRemindOrderModel.MarkStatusById(status, statusDesc)
					tempId = orderModel.TemplateId
					charges = orderModel.Charges
					res = true
					sendTime = orderModel.SendTime.Time

				}
			case "3":
				orderModel := smsMarketOrderModel.GetOne([]string{"id", "status", "charges"})
				if orderModel.ID == 0 {
					sflogger.SmsCallbackHandleErrorLog(string(msg.Body), "未发现短信发送记录"+smsMarketOrderModel.TableName())
					return
				}
				if orderModel.Status != 2 && orderModel.Status != 3 {
					smsMarketOrderModel.MarkStatus(status, statusDesc)
					charges = orderModel.Charges
					res = true
					sendTime = orderModel.SendTime.Time

					//短信任务计数
					if status == 2 {
						smsMarketJobModel.IncrSuccess()
					} else {
						smsMarketJobModel.IncrFail()
					}
				}
			}

			//模版统计
			if tempId > 0 {
				smsTempModel := model.SmsTemplateOperate{Id: tempId}
				smsTempModel.MarkSuccess()
			}

			//会员营销短信计数
			if res {
				var updateData = map[string]interface{}{
					"send_msg_total": gorm.Expr("send_msg_total + 1"),
				}
				smsCustomerModel.GetOne()
				if smsCustomerModel.LatestMsgTime.Time.Unix() < sendTime.Unix() {
					updateData["latest_msg_time"] = sendTime.Format(sfconst.GO_TIME_FULL)
				}
				smsCustomerModel.UpdateColumns(updateData)
			}

			sflogger.SmsCallbackHandleLog(userId, jobId, v.Mobile)

			//短信发送成功和失败数量lpsuh推入到redis,再定时监控redis llen,读取长度对应的数量进行冻结余额的操作
			//这样的目的主要是为了频繁对用户账号加锁，影响正常功能使用
			//测试短信不计入余额计算
			if smsType != "0" && charges > 0 {
				//计费几条就推送几条
				for i := 0; i < charges; i++ {
					redisclient.Rpush(constant.REDIS_KEY_SMS_SEND_RESULT, fmt.Sprintf("%d:%d", userId, status))
				}
				go redisSmsStatistics(userId, status, charges)
			}
		}

	})
}

//redis统计短信发送成功/失败/计费数
func redisSmsStatistics(userId int64, status, charges int) {
	smsModel := model.SmsUsers{ID: userId}
	shopId := smsModel.GetShopId()
	smsStatistics := statistics.NewSmsStatistics(shopId)
	if status == 2 {
		smsStatistics.SmsSendSuccess(1)
		smsStatistics.SmsSendCharges(int64(charges))
	} else {
		smsStatistics.SmsSendFail(1)
	}
}

//根据短信成功失败数量，处理用户账户余额
func userAccountHandle() {
	var sleepTime = 10
	var maxLen int64 = 1000
	var allUser = make(map[int64]int)
	var successUser = make(map[int64]int)
	var failUser = make(map[int64]int)

	for {
		if llen := redisclient.Llen(constant.REDIS_KEY_SMS_SEND_RESULT); llen > 0 {
			for i := 0; i < int(llen); i++ {
				result := strings.Split(redisclient.LPop(constant.REDIS_KEY_SMS_SEND_RESULT), ":")
				userId, _ := strconv.ParseInt(result[0], 10, 64)
				allUser[userId] = 1
				if result[1] == "2" {
					if _, ok := successUser[userId]; ok {
						successUser[userId]++
					} else {
						successUser[userId] = 1
					}
				} else {
					if _, ok := failUser[userId]; ok {
						failUser[userId]++
					} else {
						failUser[userId] = 1
					}
				}
			}

			//统计完每个账户的成功和失败数量，开始解冻余额
			for k := range allUser {
				_ = global.Mysql.Transaction(func(tx *gorm.DB) error {
					smsAccountBalanceModel := model.SmsAccountBalances{UserId: k, Tx: tx}
					smsAccountBalanceModel.GetOne()
					freezingSmsNum := smsAccountBalanceModel.FreezingSmsNum - successUser[k] - failUser[k]
					if freezingSmsNum < 0 {
						freezingSmsNum = 0
					}
					err := smsAccountBalanceModel.UpdateColumnsBySql(map[string]interface{}{
						"remaining_sms_num": gorm.Expr("remaining_sms_num + ?", failUser[k]),
						"freezing_sms_num":  freezingSmsNum,
					}, "id = ?", smsAccountBalanceModel.Id)
					msg := "结算成功"
					if err != nil {
						msg = fmt.Sprintf("结算失败，失败原因:%s", err.Error())
					}
					sflogger.SmsSettleLog(k, successUser[k], failUser[k], msg)
					return nil
				})

			}

			if llen >= 10*maxLen {
				sleepTime = 1
			} else if llen >= 5*maxLen {
				sleepTime = 2
			} else if llen >= 3*maxLen {
				sleepTime = 3
			} else if llen >= maxLen {
				sleepTime = 5
			} else {
				sleepTime = 10
			}
			successUser = map[int64]int{}
			failUser = map[int64]int{}
		}
		time.Sleep(time.Second * time.Duration(sleepTime))
	}
}
