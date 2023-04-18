package chat

import (
	"sync"

	"github.com/Logiase/MiraiGo-Template/bot"
	"github.com/Logiase/MiraiGo-Template/config"
	"github.com/Logiase/MiraiGo-Template/utils"
	"github.com/Mrs4s/MiraiGo/client"
	"github.com/Mrs4s/MiraiGo/message"
)

func init() {
	instance := &chat{}
	bot.RegisterModule(instance)
}

type chat struct {
}

func (m *chat) MiraiGoModule() bot.ModuleInfo {
	return bot.ModuleInfo{
		ID:       "lurendie.chat",
		Instance: instance,
	}
}

func (m *chat) Init() {
	// 初始化过程
	// 在此处可以进行 Module 的初始化配置
	// 如配置读取
}

func (m *chat) PostInit() {
	// 第二次初始化
	// 再次过程中可以进行跨Module的动作
	// 如通用数据库等等
}

func (m *chat) Serve(b *bot.Bot) {
	// 注册服务函数部分
	register(b)
}

func (m *chat) Start(b *bot.Bot) {
	// 此函数会新开携程进行调用
	// ```go
	// 		go exampleModule.Start()
	// ```

	// 可以利用此部分进行后台操作
	// 如http服务器等等
}

func (m *chat) Stop(b *bot.Bot, wg *sync.WaitGroup) {
	// 别忘了解锁
	defer wg.Done()
	// 结束部分
	// 一般调用此函数时，程序接收到 os.Interrupt 信号
	// 即将退出
	// 在此处应该释放相应的资源或者对状态进行保存
}

var instance *chat

var chatLogger = utils.GetModuleLogger("chat")

func register(b *bot.Bot) {
	b.GroupMessageEvent.Subscribe(groupHandle)
}

// 群消息业务处理
func groupHandle(client *client.QQClient, event *message.GroupMessage) {
	//循环群号
	for _, v := range config.GlobalConfig.GetIntSlice("gruops") {
		//如果群是功能群
		if int(event.GroupCode) == v {
			//业务代码
			if event.ToString() == "chat" {
				chatLogger.Info("ChatGPT启动!!!")
				config.GlobalConfig.Viper.Set("chat", 1)
				m := message.NewSendingMessage().Append(message.NewText("好的,已启动,接下来将有我为你解答问题!!"))
				client.SendGroupMessage(event.GroupCode, m)
			} else if config.GlobalConfig.Viper.GetInt("chat") == 1 {
				m := message.NewSendingMessage().Append(message.NewText(ChatGPT(event.ToString())))
				client.SendGroupMessage(event.GroupCode, m)
			}
		}
	}

}

// 删除切片元素
func DeleteSlice(a []int64, elem int64) []int64 {
	j := 0
	for _, v := range a {
		if v != elem {
			a[j] = v
			j++
		}
	}
	return a[:j]
}
