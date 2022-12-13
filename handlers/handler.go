package handlers

import (
	"github.com/869413421/wechatbot/config"
	"github.com/869413421/wechatbot/service"
	"github.com/eatmoreapple/openwechat"
	"log"
)

// MessageHandlerInterface 消息处理接口
type MessageHandlerInterface interface {
	handle(*openwechat.Message) error
	ReplyText(*openwechat.Message) error
}

type HandlerType string

const (
	GroupHandler = "group"
	UserHandler  = "user"
)

// handlers 所有消息类型类型的处理器
var handlers map[HandlerType]MessageHandlerInterface
var UserService service.UserServiceInterface

func init() {
	handlers = make(map[HandlerType]MessageHandlerInterface)
	handlers[GroupHandler] = NewGroupMessageHandler()
	handlers[UserHandler] = NewUserMessageHandler()

	UserService = service.NewUserService()
}

// Handler 全局处理入口
func Handler(msg *openwechat.Message) {
	log.Printf("hadler Received msg : %v", msg.Content)
	// 处理群消息
	if msg.IsSendByGroup() {
		handlers[GroupHandler].handle(msg)
		return
	}

	// 好友申请
	if msg.IsFriendAdd() {
		if config.Config.AutoPass {
			friend, err := msg.Agree("你好我是【V起来】微信群聊版 ChatGPT机器人。\n 请进群体验，私聊小窗不再回复：https://mp.weixin.qq.com/s/n-zjrRsa8lNrzhZV9iFMww")
			if err != nil {
				log.Printf("add friend agree error : %v", err)
				return
			}
			// 加好友默认提示语
			img, err := loadRemoteImg(config.Config.QunUrl, "qun.png")
			reply := "自动回复：由于线路限流，暂时关闭私聊功能，非常抱歉！\n \n不过我们仍然支持群聊，建议您邀请朋友一起关注【V起来】，然后拉机器人进群，进行群聊。或者直接加入官方群体>验，关注进群方式请点击：\n \nhttps://mp.weixin.qq.com/s/n-zjrRsa8lNrzhZV9iFMww"
			if img != nil {
				reply = "自动回复：由于线路限流，暂时关闭私聊功能，非常抱歉！\n \n不过我们仍然支持群聊，建议您邀请朋友一起关注【V起来】，然后拉机器人进群，进行群聊。或者直接加入官方群体验。进群方式请扫下方二维码"
			}
			_, err = friend.SendText(reply)
			if img != nil {
				_, err = friend.SendImage(img)
			}
			return
		}
	}

	// 私聊
	handlers[UserHandler].handle(msg)
}
