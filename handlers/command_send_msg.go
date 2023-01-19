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

var _ CommandHandlerInterface = (*CommandSendMsgHandler)(nil)

var senderCmmder = ""

// CommandSendMsgHandler 群消息处理
type CommandSendMsgHandler struct {
}

func (c CommandSendMsgHandler) handle(message *openwechat.Message) error {
	sender, _ := message.Sender()

	// 非系统用户直接返回，没有权限
	if !strings.Contains(config.Config.SystemUser, sender.NickName) {
		return nil
	}

	content := message.Content
	// 设置command（第一次设置command，用于第二次发送消息）
	if strings.EqualFold(content, SendMsgCommand) {
		msgCommand = SendMsgCommand
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
	services.NewRulesConfigService().UpdateRulesConfig("2", configNewsData.TimeCron, configNewsData.SendUser, configNewsData.SendGroup, configNewsData.Content, 0)
	log.Printf("bot msg config cmd set")
	task.CreateMsgTask()
	return nil
}

// NewCommandSendMsgHandler 创建发送消息处理器
func NewCommandSendMsgHandler() CommandHandlerInterface {
	return &CommandSendMsgHandler{}
}
