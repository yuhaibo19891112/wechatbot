#!/bin/sh

#进程名称
process_name=chatgpt-server-0.0.1-SNAPSHOT.jar

while [ 0 -eq 0 ]
do
    ps -ef|grep $process_name |grep -v grep
    # $? -ne 0 不存在，$? -eq 0 存在
    if [ $? -ne 0 ]
    then
    echo ">>>process is stop,to start"
    #启动进程
    nohup java -jar $process_name &
    
    break
    else
    echo ">>>process is runing,to kill"
    #停止进程
    ps -ef | grep $process_name | grep -v grep | awk '{print $2}' | xargs kill
    #休眠一秒后判断
    sleep 1
    fi
done
