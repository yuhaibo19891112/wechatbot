package main

import (
	"fmt"
	"github.com/869413421/wechatbot/bootstrap"
	"github.com/869413421/wechatbot/botcron"
	"github.com/869413421/wechatbot/config"
	"github.com/869413421/wechatbot/services"
)

func main() {
	go config.Timer()

	// 定时任务测试demo
	task := func() {
		fmt.Printf("cron test")
		rulesConfig := services.NewRulesConfigService().GetNewsRulesConfig()
		fmt.Printf("{},{}", rulesConfig.SendGroup, rulesConfig.SendTime)
	}
	go botcron.NewWeChatBotCron("*/5 * * * * ?", task)
	// 定时任务测试结束

	bootstrap.Run()
}
