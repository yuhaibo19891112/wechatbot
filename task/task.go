package task

import (
	"encoding/json"
	"github.com/869413421/wechatbot/botcron"
	"github.com/869413421/wechatbot/chatbot"
	"github.com/869413421/wechatbot/common"
	"github.com/869413421/wechatbot/config"
	"github.com/869413421/wechatbot/database/model"
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

var newsRuleCfg *model.RulesConfig

var msgRuleCfg *model.RulesConfig

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
	newsRuleCfg = services.NewRulesConfigService().GetRulesConfig("1")
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

func CreateMsgTask() {
	msgRuleCfg = services.NewRulesConfigService().GetRulesConfig("2")
	if msgRuleCfg == nil {
		return
	}
	newsId := strconv.Itoa(int(msgRuleCfg.ID))
	c, ok := botTaskMap[newsId]
	// remove old cron
	if ok {
		log.Printf("remove old msg cron")
		c.Stop()
		delete(botTaskMap, newsId)
	}
	// 新闻定时任务
	log.Printf("create msg cron, %s", msgRuleCfg.SendTime)
	botTaskMap[newsId] = botcron.NewWeChatBotCron(msgRuleCfg.SendTime, msgSendTask)
}

// 发送消息
func msgSendTask() {
	if msgRuleCfg == nil || chatbot.GlobalBot == nil {
		return
	}
	if msgRuleCfg.Content == "" || msgRuleCfg.Content == "nil" {
		return
	}
	contentData := &ContentData{}
	err := json.Unmarshal([]byte(msgRuleCfg.Content), contentData)
	if err != nil {
		log.Printf("msg content parse error, %s", err.Error())
		return
	}
	img, _:= common.LoadRemoteImg(contentData.ImgMsg, "img.png")
	sendMsg(chatbot.GlobalBot, contentData.TextMsg, img, msgRuleCfg.SendUser, msgRuleCfg.SendGroup)
}

// 新闻早餐任务
func newsSendTask() {
	if newsRuleCfg == nil || chatbot.GlobalBot == nil{
		return
	}

	newsStr := RemoteNews()
	if newsStr == "" {
		log.Printf("query news empty")
		return
	}

	sendMsg(chatbot.GlobalBot, newsStr, nil, newsRuleCfg.SendUser, newsRuleCfg.SendGroup)
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

type ContentData struct {
	TextMsg    string
	ImgMsg     string
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
