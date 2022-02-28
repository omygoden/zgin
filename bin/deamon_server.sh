#! /bin/bash
BIN_PATH=`pwd`

PROJECT_PATH=".."
LOGPATH=runlog
if [ ! -d $LOGPATH  ];then
    mkdir $LOGPATH
fi
case $1 in
"main")
    BIN_NAME="go-main"
    cd $PROJECT_PATH
    go build -o $BIN_PATH/$BIN_NAME main.go
    echo $BIN_NAME "build success"
    ;;
"sms")
    BIN_NAME="go-sms-queue"
    cd $PROJECT_PATH
    go build -o $BIN_PATH/$BIN_NAME ./queue/sms_callback_customer_queue.go
    echo $BIN_NAME "build success"
    ;;
"cron")
   BIN_NAME="go-cron"
   cd $PROJECT_PATH
   go build -o $BIN_PATH/$BIN_NAME ./crontab/cron.go
   echo $BIN_NAME "build success"
    ;;
"status")
    ps aux | grep -v grep | grep "go-"
    ;;
"")
echo "格式有误"
echo "示例：./queue_server.sh main/sms/cron"
;;
esac
