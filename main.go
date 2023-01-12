package main

import (
	"fmt"
	"github.com/869413421/wechatbot/bootstrap"
	"github.com/869413421/wechatbot/botcron"
	"github.com/869413421/wechatbot/config"
)

func main() {
	go config.Timer()

	// 定时任务测试demo
	task := func() {
		fmt.Printf("cron test")
	}
	go botcron.NewWeChatBotCron("*/5 * * * * ?", task)
	// 定时任务测试结束

	bootstrap.Run()
}
