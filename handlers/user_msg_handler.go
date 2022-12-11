package handlers

import (
	"github.com/869413421/wechatbot/config"
	"github.com/869413421/wechatbot/gtp"
	"github.com/eatmoreapple/openwechat"
	"log"
	"strings"
)

var _ MessageHandlerInterface = (*UserMessageHandler)(nil)
var warnUserFlg = true

// UserMessageHandler 私聊消息处理
type UserMessageHandler struct {
}

// handle 处理消息
func (g *UserMessageHandler) handle(msg *openwechat.Message) error {
	if msg.IsText() {
		go g.ReplyText(msg)
		return nil
	}
	return nil
}

// NewUserMessageHandler 创建私聊处理器
func NewUserMessageHandler() MessageHandlerInterface {
	return &UserMessageHandler{}
}

// ReplyText 发送文本消息到群
func (g *UserMessageHandler) ReplyText(msg *openwechat.Message) error {
	// 接收私聊消息
	sender, err := msg.Sender()
	log.Printf("Received User %v Text Msg : %v", sender.NickName, msg.Content)

	// 获取上下文，向GPT发起请求
	requestText := strings.TrimSpace(msg.Content)
	requestText = strings.Trim(msg.Content, "\n")

	requestText = UserService.GetUserSessionContext(sender.ID()) + requestText
	reply, err := gtp.Completions(requestText)
	if err != nil {
		log.Printf("gtp request error: %v \n", err)
		errorTip := "机器人累了要歇会儿，我很快就能V起来……"
		if reply == "429" {
		    warnFriend(msg)
		    errorTip = errorTip + "!!!!!!!"
		}
		msg.ReplyText(errorTip)
		return err
	}
	if reply == "" {
		return nil
	}

	// 设置上下文，回复用户
	reply = strings.TrimSpace(reply)
	reply = strings.Trim(reply, "\n")
	UserService.SetUserSessionContext(sender.ID(), requestText, reply)
	reply = "本消息由 V起来Bot 回复：\n " + reply
	_, err = msg.ReplyText(reply)
	if err != nil {
		log.Printf("response user error: %v \n", err)
	}
	return err
}

func warnFriend(msg *openwechat.Message) error{
	self, err := msg.Bot.GetCurrentUser()
	friends, err := self.Friends()
	alarmUser := friends.GetByRemarkName(config.Config.AlarmUserName)
	if alarmUser != nil && warnUserFlg{
		alarmUser.SendText("keys已过期，尽快重置")
		warnUserFlg = false
	}
	return err
}
