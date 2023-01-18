package handlers

import (
	"encoding/json"
	"github.com/869413421/wechatbot/config"
	"github.com/869413421/wechatbot/services"
	"github.com/869413421/wechatbot/task"
	"github.com/eatmoreapple/openwechat"
	"log"
	"strings"
)

var _ CommandHandlerInterface = (*CommandConfigNewsHandler)(nil)

// CommandConfigNewsHandler 新闻配置处理器
type CommandConfigNewsHandler struct {
}

type ConfigNewsData struct {
	TimeCron string
	// 发送对应还有或群名称，多个用;隔开
	SendGroup string
	// 发送对应用户
	SendUser string
	// 内容
	Content  string
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
	err := json.Unmarshal([]byte(content), configNewsData)
	if err != nil {
		return nil
	}
	services.NewRulesConfigService().UpdateRulesConfig("1", configNewsData.TimeCron, configNewsData.SendUser, configNewsData.SendGroup, configNewsData.Content)
	log.Printf("bot news config cmd set")
	task.CreateNewsTask()
	return nil
}

// NewCommandConfigNewsHandler 创建新闻配置处理器
func NewCommandConfigNewsHandler() CommandHandlerInterface {
	return &CommandConfigNewsHandler{}
}
