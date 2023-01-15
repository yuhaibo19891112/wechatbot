package handlers

import "github.com/eatmoreapple/openwechat"

type cmmdHandleType string

const (
	SendMsgCommand = "sendMsgCommand"
	ConfigNews     = "configNewsCommand"
)

// CommandHandlerInterface 消息处理接口
type CommandHandlerInterface interface {
	handle(*openwechat.Message) error
}

var cmmdHandles map[cmmdHandleType]CommandHandlerInterface

func init() {
	cmmdHandles = make(map[cmmdHandleType]CommandHandlerInterface)
	cmmdHandles[SendMsgCommand] = NewCommandSendMsgHandler()
	cmmdHandles[ConfigNews] = NewCommandConfigNewsHandler()
}
