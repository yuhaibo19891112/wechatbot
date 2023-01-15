package handlers

import (
	"github.com/869413421/wechatbot/config"
	"github.com/eatmoreapple/openwechat"
	"log"
	"os"
	"strings"
	"time"
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

	// 设置command（第一次设置command，用于第二次发送消息）
	if strings.EqualFold(message.Content, SendMsgCommand) {
		msgCommand = SendMsgCommand
		return nil
	}

	if strings.EqualFold(message.Content, "结束发送消息") {
		msgCommand = ""
		return nil
	}

	// 开始发送消息
	if strings.EqualFold(msgCommand, SendMsgCommand) {
		msgCommand = ""
		log.Printf("------------> 发送消息")
		self, _ := message.Bot().GetCurrentUser()
		groups, _ := self.Groups()
		if message.IsText() {
			for i := 0; i < len(groups); i++ {
				time.Sleep(5 * time.Second)
				groups[i].SendText(message.Content)
			}
		} else if message.IsPicture() {
			message.SaveFileToLocal("temp.png")
			for i := 0; i < len(groups); i++ {
				temp, _ := os.Open("temp.png")
				time.Sleep(5 * time.Second)
				groups[i].SendImage(temp)
			}
		}
	}
	return nil
}

// NewCommandSendMsgHandler 创建发送消息处理器
func NewCommandSendMsgHandler() CommandHandlerInterface {
	return &CommandSendMsgHandler{}
}
