package master

import (
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/Logiase/MiraiGo-Template/bot"
	"github.com/Logiase/MiraiGo-Template/config"
	"github.com/Logiase/MiraiGo-Template/utils"
	"github.com/Mrs4s/MiraiGo/client"
	"github.com/Mrs4s/MiraiGo/message"
)

func init() {
	instance := &master{}
	bot.RegisterModule(instance)
}

type master struct {
}

func (m *master) MiraiGoModule() bot.ModuleInfo {
	return bot.ModuleInfo{
		ID:       "lurendie.master",
		Instance: instance,
	}
}

func (m *master) Init() {
	// 初始化过程
	// 在此处可以进行 Module 的初始化配置
	// 如配置读取
}

func (m *master) PostInit() {
	// 第二次初始化
	// 再次过程中可以进行跨Module的动作
	// 如通用数据库等等
}

func (m *master) Serve(b *bot.Bot) {
	// 注册服务函数部分
	register(b)
}

func (m *master) Start(b *bot.Bot) {
	// 此函数会新开携程进行调用
	// ```go
	// 		go exampleModule.Start()
	// ```

	// 可以利用此部分进行后台操作
	// 如http服务器等等
}

func (m *master) Stop(b *bot.Bot, wg *sync.WaitGroup) {
	// 别忘了解锁
	defer wg.Done()
	// 结束部分
	// 一般调用此函数时，程序接收到 os.Interrupt 信号
	// 即将退出
	// 在此处应该释放相应的资源或者对状态进行保存
}

var instance *master

var masterLogger = utils.GetModuleLogger("master")

func register(b *bot.Bot) {
	b.PrivateMessageEvent.Subscribe(func(client *client.QQClient, event *message.PrivateMessage) {
		if event.Sender.Uin == config.GlobalConfig.GetInt64("master") {
			if strings.HasPrefix(event.ToString(), "添加功能群") {
				masterLogger.Infoln("添加群聊触发")
				array := make([]int64, len(config.GlobalConfig.GetIntSlice("gruops")))
				for _, v := range config.GlobalConfig.GetIntSlice("gruops") {
					array = append(array, int64(v))
				}
				str := event.ToString()[15:]

				group, _ := strconv.ParseInt(str, 10, 64)
				array = append(array, group)
				array = DeleteSlice(array, 0)
				config.GlobalConfig.Set("gruops", array)
				if err := config.GlobalConfig.WriteConfig(); err != nil {
					fmt.Printf("err: %v\n", err)
				}
				m := message.NewSendingMessage().Append(message.NewText(fmt.Sprintf("功能群:%v添加成功", group)))
				client.SendPrivateMessage(event.Sender.Uin, m)
				config.GlobalConfig.WatchConfig()
			} else if strings.HasPrefix(event.ToString(), "删除功能群") {
				masterLogger.Infoln("删除群聊触发")
				array := make([]int64, len(config.GlobalConfig.GetIntSlice("gruops"))-1)
				str := event.ToString()[15:]
				group, _ := strconv.ParseInt(str, 10, 64)
				for _, v := range config.GlobalConfig.GetIntSlice("gruops") {
					array = append(array, int64(v))
				}
				array = DeleteSlice(array, group)
				array = DeleteSlice(array, 0)
				config.GlobalConfig.Set("gruops", array)
				if err := config.GlobalConfig.WriteConfig(); err != nil {
					fmt.Printf("err: %v\n", err)
				}
				m := message.NewSendingMessage().Append(message.NewText(fmt.Sprintf("功能群:%v删除成功", group)))
				client.SendPrivateMessage(event.Sender.Uin, m)
				config.GlobalConfig.WatchConfig()
			} else if event.ToString() == "查看功能群" {
				m := message.NewSendingMessage().Append(message.NewText("当前功能群:\n"))
				for _, v := range config.GlobalConfig.GetIntSlice("gruops") {
					m.Append(message.NewText(strconv.FormatInt(int64(v), 10) + "\n"))
				}
				client.SendPrivateMessage(event.Sender.Uin, m)
			}
		}
	})
}
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
