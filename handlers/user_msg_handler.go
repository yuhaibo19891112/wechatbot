package handlers

import (
	"github.com/eatmoreapple/openwechat"
	"log"
	"strings"
)

var msgCommand = ""

var _ MessageHandlerInterface = (*UserMessageHandler)(nil)

// UserMessageHandler 私聊消息处理
type UserMessageHandler struct {
}

// handle 处理消息
func (g *UserMessageHandler) handle(msg *openwechat.Message) error {
	if msg.IsText() || msg.IsPicture() {
		content := msg.Content
		// 配置群消息发送命令
		if strings.EqualFold(SendMsgCommand, content) || strings.EqualFold(msgCommand, SendMsgCommand) {
			log.Printf("send msg command")
			return cmmdHandles[SendMsgCommand].handle(msg)
		}
		// 配置新闻发送命令
		if strings.EqualFold(ConfigNews, content) || strings.EqualFold(msgCommand, ConfigNews) {
			log.Printf("config news command")
			return cmmdHandles[ConfigNews].handle(msg)
		}
	}
	return nil
}

// NewUserMessageHandler 创建私聊处理器
func NewUserMessageHandler() MessageHandlerInterface {
	return &UserMessageHandler{}
}

// ReplyText 发送文本消息到群
func (g *UserMessageHandler) ReplyText(msg *openwechat.Message) error {
	return nil
}
