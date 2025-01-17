package handlers

import (
	"github.com/869413421/wechatbot/config"
	"github.com/eatmoreapple/openwechat"
	"os"
	"regexp"
	"strings"
)

var _ MessageHandlerInterface = (*GroupMessageHandler)(nil)

// GroupMessageHandler 群消息处理
type GroupMessageHandler struct {
}

// handle 处理消息
func (g *GroupMessageHandler) handle(msg *openwechat.Message) error {
	// 别人加入群聊
	sender, _ := msg.Sender()
	if joinGroup(msg) && config.Config.JiekeTip != "" && (strings.HasPrefix(sender.NickName, "V起来") || strings.HasPrefix(sender.NickName, "OKR之剑") || strings.HasPrefix(sender.NickName, "ytest")) {
		// 只查找通过二维码扫描进来的用户
		userName := findJoinUserName(msg.Content)
		jiekeTip := strings.ReplaceAll(config.Config.JiekeTip, "xxx", userName)

		img, err := os.Open(config.Config.QunUrl)

		msg.ReplyText(jiekeTip)
		if err == nil {
			msg.ReplyImage(img)
		}
		return nil
	}
	return nil
}

// NewGroupMessageHandler 创建群消息处理器
func NewGroupMessageHandler() MessageHandlerInterface {
	return &GroupMessageHandler{}
}

// ReplyText 发送文本消息到群
func (g *GroupMessageHandler) ReplyText(msg *openwechat.Message) error {
	return nil
}

func joinGroup(m *openwechat.Message) bool {
	return m.IsSystem() && (strings.Contains(m.Content, "加入了群聊") || strings.Contains(m.Content, "加入群聊")) && m.IsSendByGroup()
}

func findJoinUserName(content string) string {
	reg := regexp.MustCompile(`^"[^"]*"通过`)
	match := reg.FindStringSubmatch(content)
	if len(match) > 0 {
		temp := strings.ReplaceAll(match[0], "\"", "")
		temp = strings.TrimSpace(temp)
		temp = strings.ReplaceAll(temp, "通过", "")
		return temp
	}

	reg2 := regexp.MustCompile(`邀请"[^"]*"加入`)
	match2 := reg2.FindStringSubmatch(content)
	if len(match2) > 0 {
		temp := strings.ReplaceAll(match2[0], "邀请\"", "")
		temp = strings.ReplaceAll(temp, "\"加入", "")
		return temp
	}
	return ""
}
