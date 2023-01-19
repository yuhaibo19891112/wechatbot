package task

import (
	"encoding/json"
	"github.com/869413421/wechatbot/botcron"
	"github.com/869413421/wechatbot/config"
	"github.com/869413421/wechatbot/services"
	"github.com/eatmoreapple/openwechat"
	"github.com/robfig/cron/v3"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var botTaskMap map[string]*cron.Cron

func Run() {
	botTaskMap = make(map[string]*cron.Cron)
	initTask := botcron.NewWeChatBotCron("0 0 4 * * ?", createTask)
	botTaskMap["initTask"] = initTask
	// 启动的时候，立即执行
	createTask()
}

func createTask() {
	// 新闻早操任务
	CreateNewsTask()

	// 消息发送任务
	CreateMsgTask()
}

func CreateNewsTask() {
	newsRuleCfgList := services.NewRulesConfigService().GetRulesConfigList("1")
	if newsRuleCfgList == nil {
		return
	}
	for i := 0; i < len(newsRuleCfgList); i++ {
		newsId := strconv.Itoa(int(newsRuleCfgList[i].ID))
		c, ok := botTaskMap[newsId]
		// remove old cron
		if ok {
			log.Printf("remove old news cron, newsId: %s", newsId)
			c.Stop()
			delete(botTaskMap, newsId)
		}
		// 新闻定时任务
		log.Printf("create news cron, %s, key: %s", newsRuleCfgList[i].SendTime, newsId)
		newsSendJob := SendNewsCronJob{SendUser: newsRuleCfgList[i].SendUser, SendGroup: newsRuleCfgList[i].SendGroup}
		botTaskMap[newsId] = botcron.NewWeChatBotCronJob(newsRuleCfgList[i].SendTime, newsSendJob)
	}
}

func CreateMsgTask() {
	msgRuleCfgList := services.NewRulesConfigService().GetRulesConfigList("2")
	if msgRuleCfgList == nil {
		return
	}
	for i := 0; i < len(msgRuleCfgList); i++ {
		msgCronId := strconv.Itoa(int(msgRuleCfgList[i].ID))
		c, ok := botTaskMap[msgCronId]
		// remove old cron
		if ok {
			log.Printf("remove old msg cron, cron key: %s", msgCronId)
			c.Stop()
			delete(botTaskMap, msgCronId)
		}
		// 新闻定时任务
		log.Printf("create msg cron, %s, key: %s", msgRuleCfgList[i].SendTime, msgCronId)
		sendMsgJob := SendMsgCronJob{SendUser: msgRuleCfgList[i].SendUser, SendGroup: msgRuleCfgList[i].SendGroup, Content: msgRuleCfgList[i].Content}
		botTaskMap[msgCronId] = botcron.NewWeChatBotCronJob(msgRuleCfgList[i].SendTime, sendMsgJob)
	}
}

func sendMsg(bot *openwechat.Bot, msg string, img io.Reader, sendUser string, sendGroup string) {
	user, err := bot.GetCurrentUser()
	if err != nil {
		log.Printf("bot current user error, error: %s", err.Error())
		return
	}
	// 发送好友
	if sendUser != "" {
		friends, _ := user.Friends()
		if friends == nil {
			log.Printf("command config msg get friends error")
			return
		}

		temps := strings.Split(sendUser, ";")
		for i := 0; i < len(temps); i++ {
			log.Printf("--------->send user msg, user: %s", temps[i])
			temp := friends.GetByNickName(temps[i])
			if temp != nil {
				if msg != "" {
					temp.SendText(msg)
				}
				if img != nil {
					temp.SendImage(img)
				}
				time.Sleep(5 * time.Second)
			}
		}

	}

	// 发送群
	if sendGroup != "" {
		groups, _ := user.Groups()
		if groups == nil {
			log.Printf("command config msg get groups error")
			return
		}

		temps := strings.Split(sendGroup, ";")
		for i := 0; i < len(temps); i++ {
			log.Printf("----------->send group msg, %s", temps[i])
			temp := groups.GetByNickName(temps[i])
			if temp != nil {
				if msg != "" {
					temp.SendText(msg)
				}
				if img != nil {
					temp.SendImage(img)
				}
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
