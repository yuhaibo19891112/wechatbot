package task

import (
	"github.com/869413421/wechatbot/chatbot"
	"github.com/robfig/cron/v3"
	"log"
)

var _ cron.Job = (*SendNewsCronJob)(nil)

// SendNewsCronJob 新闻早餐任务
type SendNewsCronJob struct {
	SendUser    string
	SendGroup   string
}

func (s SendNewsCronJob) Run() {
	if chatbot.GlobalBot == nil{
		return
	}

	newsStr := RemoteNews()
	if newsStr == "" {
		log.Printf("query news empty")
		return
	}

	sendMsg(chatbot.GlobalBot, newsStr, nil, s.SendUser, s.SendGroup)
}

