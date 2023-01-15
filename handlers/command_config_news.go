package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/869413421/wechatbot/botcron"
	"github.com/869413421/wechatbot/config"
	"github.com/eatmoreapple/openwechat"
	"github.com/robfig/cron/v3"
	"log"
	"strings"
)

var _ CommandHandlerInterface = (*CommandConfigNewsHandler)(nil)

var sendType = ""
var sendName = ""

var sendNewsCron *cron.Cron

// CommandConfigNewsHandler 新闻配置处理器
type CommandConfigNewsHandler struct {
}

type ConfigNewsData struct {
	TimeCron string `json:"TimeCron"`
	// 发送类型，1-好友，0-群
	SendType string `json:"SendType"`
	// 发送对应还有或群名称，多个用;隔开；特殊：群如果不填写，则发所有群
	SendName string `json:"SendName"`
}

func (c CommandConfigNewsHandler) handle(message *openwechat.Message) error {
	sender, _ := message.Sender()
	// 非系统用户直接返回，没有权限
	if !strings.Contains(config.Config.SystemUser, sender.NickName) {
		log.Printf("config news no system user")
		return nil
	}

	content := message.Content

	// 设置command（第一次设置command，用于第二次发送消息）
	if strings.EqualFold(content, ConfigNews) {
		msgCommand = ConfigNews
		return nil
	}

	if !strings.EqualFold(msgCommand, ConfigNews) {
		return nil
	}

	// 重置command
	msgCommand = ""

	// 执行任务
	configNewsData := &ConfigNewsData{}
	json.Unmarshal([]byte(content), configNewsData)
	/*if err != nil {
		return nil
	}
	sendType = configNewsData.SendType
	sendName = configNewsData.SendName*/
	if sendNewsCron != nil {
		log.Printf("stop old cron")
		sendNewsCron.Stop()
	}
	log.Printf("create new cron")
	sendNewsCron = botcron.NewWeChatBotCron("*/5 * * * * ?", sendNewsTask)
	return nil
}

func sendNewsTask() {
	fmt.Printf("cron test")
}

// NewCommandConfigNewsHandler 创建新闻配置处理器
func NewCommandConfigNewsHandler() CommandHandlerInterface {
	return &CommandConfigNewsHandler{}
}
