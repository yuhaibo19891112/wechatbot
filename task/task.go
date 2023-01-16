package task

import (
	"github.com/869413421/wechatbot/bootstrap"
	"github.com/869413421/wechatbot/botcron"
	"github.com/869413421/wechatbot/database/model"
	"github.com/869413421/wechatbot/handlers"
	"github.com/869413421/wechatbot/services"
	"github.com/robfig/cron/v3"
	"log"
	"strconv"
	"strings"
	"time"
)

var botTaskMap map[string]*cron.Cron

var newsRuleCfg *model.RulesConfig

func Run() {
	botTaskMap = make(map[string]*cron.Cron)
	initTask := botcron.NewWeChatBotCron("0 0 4 * * ?", CreateNewsTask)
	botTaskMap["initTask"] = initTask
	// 启动的时候，立即执行
	CreateNewsTask()
}

func CreateNewsTask() {
	newsRuleCfg = services.NewRulesConfigService().GetNewsRulesConfig()
	if newsRuleCfg == nil {
		return
	}
	newsId := strconv.Itoa(int(newsRuleCfg.ID))
	c, ok := botTaskMap[newsId]
	// remove old cron
	if ok {
		c.Stop()
		delete(botTaskMap, newsId)
	}
	// 新闻定时任务
	botTaskMap[newsId] = botcron.NewWeChatBotCron(newsRuleCfg.SendTime, newsSendTask)
}

// 新闻早餐任务
func newsSendTask() {
	if newsRuleCfg == nil {
		return
	}
	bot := bootstrap.GetChatBot()
	user, err := bot.GetCurrentUser()
	if err != nil {
		log.Printf("bot current user error, error: %s", err.Error())
		return
	}
	newsStr := handlers.RemoteNews()
	// 发送好友
	if newsRuleCfg.SendUser != "" {
		friends, _ := user.Friends()
		if friends == nil {
			log.Printf("command config news get friends error")
			return
		}

		temps := strings.Split(newsRuleCfg.SendUser, ";")
		for i := 0; i < len(temps); i++ {
			temp := friends.GetByNickName(temps[i])
			if temp != nil {
				temp.SendText(newsStr)
				time.Sleep(5 * time.Second)
			}
		}

	}

	// 发送群
	if newsRuleCfg.SendGroup != "" {
		groups, _ := user.Groups()
		if groups == nil {
			log.Printf("command config news get groups error")
			return
		}

		temps := strings.Split(newsRuleCfg.SendGroup, ";")
		for i := 0; i < len(temps); i++ {
			log.Printf("----------->send group news, %s", temps[i])
			temp := groups.GetByNickName(temps[i])
			if temp != nil {
				temp.SendText(newsStr)
				time.Sleep(5 * time.Second)
			}
		}
	}
}
