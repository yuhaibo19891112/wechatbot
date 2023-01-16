package handlers

import (
	"encoding/json"
	"github.com/869413421/wechatbot/botcron"
	"github.com/869413421/wechatbot/config"
	"github.com/eatmoreapple/openwechat"
	"github.com/robfig/cron/v3"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

var _ CommandHandlerInterface = (*CommandConfigNewsHandler)(nil)

var sendType = ""
var sendName = ""

var sendNewsCron *cron.Cron

var self *openwechat.Bot

// CommandConfigNewsHandler 新闻配置处理器
type CommandConfigNewsHandler struct {
}

type ConfigNewsData struct {
	TimeCron string `json:"TimeCron"`
	// 发送类型，1-好友，0-群
	SendType string `json:"SendType"`
	// 发送对应还有或群名称，多个用;隔开
	SendName string `json:"SendName"`
}

func (c CommandConfigNewsHandler) handle(message *openwechat.Message) error {
	sender, _ := message.Sender()
	self = message.Bot()
	// 非系统用户直接返回，没有权限
	if !strings.Contains(config.Config.SystemUser, sender.NickName) {
		log.Printf("config news no system user")
		return nil
	}

	content := message.Content

	// 设置command（第一次设置command，用于第二次发送消息）
	if strings.EqualFold(content, ConfigNews) {
		msgCommand = ConfigNews
		return nil
	}

	if !strings.EqualFold(msgCommand, ConfigNews) {
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
	sendType = configNewsData.SendType
	sendName = configNewsData.SendName
	if sendNewsCron != nil {
		log.Printf("stop old cron")
		sendNewsCron.Stop()
	}
	log.Printf("create new cron")
	sendNewsCron = botcron.NewWeChatBotCron(configNewsData.TimeCron, sendNewsTask)
	return nil
}

func sendNewsTask() {
	user, err := self.GetCurrentUser()
	if err != nil {
		log.Printf("get current user error, %s", err.Error())
	}

	if sendName == "" {
		log.Printf("command config news sendName is blank")
		return
	}

	newsStr := RemoteNews()
	if newsStr == "" {
		log.Printf("query news empty")
		return
	}
	log.Printf("send news task, sendType:%s, sendName:%s", sendType, sendName)

	// 发送好友
	if "1" == sendType {
		friends, _ := user.Friends()
		if friends == nil {
			log.Printf("command config news get friends error")
			return
		}

		users := strings.Split(sendName, ";")
		for i := 0; i < len(users); i++ {
			temp := friends.GetByNickName(users[i])
			if temp != nil {
				temp.SendText(newsStr)
				time.Sleep(5 * time.Second)
			}
		}

	}

	// 发送群
	if "0" == sendType {
		groups, _ := user.Groups()
		if groups == nil {
			log.Printf("command config news get groups error")
			return
		}

		temps := strings.Split(sendName, ";")
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

// NewCommandConfigNewsHandler 创建新闻配置处理器
func NewCommandConfigNewsHandler() CommandHandlerInterface {
	return &CommandConfigNewsHandler{}
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
