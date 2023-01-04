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
	if systemUser != "" && strings.EqualFold(systemUser, sender.NickName){
		log.Printf("------------> 发送消息")
		self, _ := msg.Bot.GetCurrentUser()
		groups, _ := self.Groups()
		if msg.IsText() {
			for i := 0; i < len(groups); i++ {
				time.Sleep(2 * time.Second)
				groups[i].SendText(msg.Content)
			}
		}else if msg.IsPicture() {
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
	// 接收私聊消息
	sender, err := msg.Sender()
	log.Printf("Received User %v Text Msg : %v", sender.NickName, msg.Content)


        // 设置上下文，回复用户
	img, err := loadRemoteImg(config.Config.QunUrl, "qun.png")
	reply := "自动回复：由于线路限流，暂时关闭私聊功能，非常抱歉！\n \n不过我们仍然支持群聊，建议您邀请朋友一起关注【V起来】，然后拉机器人进群，进行群聊。或者直接加入官方群体>验，关注进群方式请点击：\n \nhttps://mp.weixin.qq.com/s/n-zjrRsa8lNrzhZV9iFMww"
	if img != nil {
		reply = "自动回复：由于线路限流，暂时关闭私聊功能，非常抱歉！\n \n不过我们仍然支持群聊，建议您邀请朋友一起关注【V起来】，然后拉机器人进群，进行群聊。或者直接加入官方群体验。进群方式请扫下方二维码"
	}
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
