package task

import (
	"encoding/json"
	"github.com/869413421/wechatbot/chatbot"
	"github.com/869413421/wechatbot/common"
	"github.com/robfig/cron/v3"
	"log"
)

var _ cron.Job = (*SendMsgCronJob)(nil)

// SendMsgCronJob 消息发送任务
type SendMsgCronJob struct {
	SendUser    string
	SendGroup   string
	Content     string
}

func (s SendMsgCronJob) Run() {
	if chatbot.GlobalBot == nil {
		return
	}
	if s.Content == "" || s.Content == "nil" {
		return
	}
	contentData := &ContentData{}
	err := json.Unmarshal([]byte(s.Content), contentData)
	if err != nil {
		log.Printf("msg content parse error, %s", err.Error())
		return
	}
	img, _:= common.LoadRemoteImg(contentData.ImgMsg, "img.png")
	sendMsg(chatbot.GlobalBot, contentData.TextMsg, img, s.SendUser, s.SendGroup)
}

type ContentData struct {
	TextMsg    string
	ImgMsg     string
}

