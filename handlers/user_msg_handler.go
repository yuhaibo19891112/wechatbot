package handlers

import (
	"github.com/869413421/wechatbot/config"
	"github.com/eatmoreapple/openwechat"
	"io/ioutil"
	"log"
	"net/http"
	"os"
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


    // 设置上下文，回复用户
	downloadImg(config.Config.QunUrl, "qun.jpg")
	img, err := os.Open("qun.jpg")
	requestText := msg.Content
	reply := "自动回复：由于线路限流，暂时关闭私聊功能，非常抱歉！\n \n不过我们仍然支持群聊，建议您邀请朋友一起关注【V起来】，然后拉机器人进群，进行群聊。或者直接加入官方群体>验，关注进群方式请点击：\n \nhttps://mp.weixin.qq.com/s/n-zjrRsa8lNrzhZV9iFMww"
	if img != nil {
		reply = "自动回复：由于线路限流，暂时关闭私聊功能，非常抱歉！\n \n不过我们仍然支持群聊，建议您邀请朋友一起关注【V起来】，然后拉机器人进群，进行群聊。或者直接加入官方群体验。进群方式请扫下方二维码"
	}
	UserService.SetUserSessionContext(sender.ID(), requestText, reply)
	_, err = msg.ReplyText(reply)
	if img != nil {
		_, err = msg.ReplyImage(img)
	}
	if err != nil {
		log.Printf("response user error: %v \n", err)
	}
	return err

/*
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
*/
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

func downloadImg(finalUrl string, savePath string)  {
	//读取url的信息，存入到文件
	resp, err := http.Get(finalUrl)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)
	if err == nil {
		ioutil.WriteFile(savePath, content, 0666)
	}
}
