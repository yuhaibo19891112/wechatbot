package main

import (
	"github.com/869413421/wechatbot/bootstrap"
	"github.com/869413421/wechatbot/config"
)

func main() {
	go config.Timer()
	//task.Run()
	bootstrap.Run()
}
