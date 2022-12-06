package main

import (
	"os"
	"os/signal"

	_ "miraiGoDo/modules/master"
	_ "miraiGoDo/modules/pixiv"
	_ "miraiGoDo/modules/wzry"

	"github.com/Logiase/MiraiGo-Template/bot"
	"github.com/Logiase/MiraiGo-Template/config"
	_ "github.com/Logiase/MiraiGo-Template/modules/logging"
	"github.com/Logiase/MiraiGo-Template/utils"
)

func init() {
	utils.WriteLogToFS(utils.LogInfoLevel, utils.LogWithStack, utils.LogDebugLevel)
	config.Init()
}

func main() {
	// 快速初始化
	bot.Init()

	// 初始化 Modules
	bot.StartService()

	// 使用协议
	// 不同协议可能会有部分功能无法使用
	// 在登陆前切换协议
	bot.UseProtocol(bot.AndroidWatch)

	// 登录
	bot.QrcodeLogin()

	// 刷新好友列表，群列表
	bot.RefreshList()
	//监听事件
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	<-ch
	bot.Stop()
}
