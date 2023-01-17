package task

import (
	"encoding/json"
	"github.com/869413421/wechatbot/botcron"
	"github.com/869413421/wechatbot/chatbot"
	"github.com/869413421/wechatbot/config"
	"github.com/869413421/wechatbot/database/model"
	"github.com/869413421/wechatbot/services"
	"github.com/robfig/cron/v3"
	"io/ioutil"
	"log"
	"net/http"
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
		log.Printf("remove old news cron")
		c.Stop()
		delete(botTaskMap, newsId)
	}
	// 新闻定时任务
	log.Printf("create news cron, %s", newsRuleCfg.SendTime)
	botTaskMap[newsId] = botcron.NewWeChatBotCron(newsRuleCfg.SendTime, newsSendTask)
}

// 新闻早餐任务
func newsSendTask() {
	if newsRuleCfg == nil || chatbot.GlobalBot == nil{
		return
	}
	user, err := chatbot.GlobalBot.GetCurrentUser()
	if err != nil {
		log.Printf("bot current user error, error: %s", err.Error())
		return
	}
	newsStr := RemoteNews()
	if newsStr == "" {
		log.Printf("query news empty")
		return
	}

	// 发送好友
	if newsRuleCfg.SendUser != "" {
		friends, _ := user.Friends()
		if friends == nil {
			log.Printf("command config news get friends error")
			return
		}

		temps := strings.Split(newsRuleCfg.SendUser, ";")
		for i := 0; i < len(temps); i++ {
			log.Printf("--------->send user news, user: %s", temps[i])
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

// RemoteConfigResponseBody 远程config.json参数体
type NewsResp struct {
	Code      int `json:"code"`
	Msg    string `json:"msg"`
	Data string `json:"data"`
}

// 新闻
func RemoteNews() string {

	resp, err := http.Get(config.Config.NewsApi)

	if err != nil {
		return ""
	}

	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return ""
	}

	if resp.StatusCode != 200 {
		return ""
	}

	newResp := &NewsResp{}
	errParse := json.Unmarshal(respBody, newResp)
	if errParse != nil {
		log.Printf("error: " + errParse.Error())
		return ""
	}

	log.Printf("msg: %s, data: %s", newResp.Msg, newResp.Data)

	return newResp.Data
}
