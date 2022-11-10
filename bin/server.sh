#! /bin/bash
BIN_PATH=`pwd`

PROJECT_PATH=".."
LOGPATH=runlog
if [ ! -d $LOGPATH  ];then
    mkdir $LOGPATH
fi

PROJECT_NAME=`echo $BIN_PATH | awk 'BEGIN{FS="/"}END{print $(NF-1)}'`
case $1 in
"main")
    BIN_NAME="$PROJECT_NAME-main"
    cd $PROJECT_PATH
    go build -o $BIN_PATH/$BIN_NAME main.go
    echo $BIN_NAME "build success"
    ;;
"cron")
   BIN_NAME="$PROJECT_NAME-cron"
   cd $PROJECT_PATH
   go build -o $BIN_PATH/$BIN_NAME ./console/cron.go
   echo $BIN_NAME "build success"
    ;;
"status")
    ps aux | grep -v grep | grep "go-"
    ;;
"")
echo "格式有误"
echo "示例：./server.sh main/cron"
;;
esac
