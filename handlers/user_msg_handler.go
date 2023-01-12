package handlers

import (
	"github.com/869413421/wechatbot/config"
	"github.com/eatmoreapple/openwechat"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

var systemUser = ""

var _ MessageHandlerInterface = (*UserMessageHandler)(nil)

// UserMessageHandler 私聊消息处理
type UserMessageHandler struct {
}

// handle 处理消息
func (g *UserMessageHandler) handle(msg *openwechat.Message) error {
	if msg.IsText() || msg.IsPicture() {
		return g.sendMsgCommand(msg)
	}
	return nil
}

// NewUserMessageHandler 创建私聊处理器
func NewUserMessageHandler() MessageHandlerInterface {
	return &UserMessageHandler{}
}

func (g *UserMessageHandler) sendMsgCommand(msg *openwechat.Message) error {
	log.Printf("-------sendMsg, user:%s", systemUser)
	// 接收私聊消息
	sender, _ := msg.Sender()
	if !strings.Contains(config.Config.SystemUser, sender.NickName) {
		return nil
	}

	if msg.IsText() && strings.EqualFold(msg.Content, "发送消息") {
		log.Printf("-----------sendMsg, record")
		systemUser = sender.NickName
		return nil
	}

	if msg.IsText() && strings.EqualFold(msg.Content, "结束发送消息") {
		systemUser = ""
		return nil
	}
	// 发送消息
	if systemUser != "" && strings.EqualFold(systemUser, sender.NickName) {
		log.Printf("------------> 发送消息")
		self, _ := msg.Bot.GetCurrentUser()
		groups, _ := self.Groups()
		if msg.IsText() {
			for i := 0; i < len(groups); i++ {
				time.Sleep(2 * time.Second)
				groups[i].SendText(msg.Content)
			}
		} else if msg.IsPicture() {
			msg.SaveFileToLocal("temp.png")
			for i := 0; i < len(groups); i++ {
				temp, _ := os.Open("temp.png")
				time.Sleep(2 * time.Second)
				groups[i].SendImage(temp)
			}
		}
		systemUser = ""
	}
	return nil
}

// ReplyText 发送文本消息到群
func (g *UserMessageHandler) ReplyText(msg *openwechat.Message) error {
	return nil
}

func loadRemoteImg(finalUrl string, savePath string) (*os.File, error) {
	//读取url的信息，存入到文件
	resp, err := http.Get(finalUrl)
	if err != nil {
		log.Printf("get remote img error, %v", err)
		return nil, err
	}
	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)
	if err == nil {
		ioutil.WriteFile(savePath, content, 0666)
	}
	return os.Open(savePath)
}
