package wzry

import (
	"fmt"
	"miraiGoDo/internal"
	"net/url"
	"strings"
	"sync"

	"github.com/Mrs4s/MiraiGo/client"
	"github.com/Mrs4s/MiraiGo/message"
	"github.com/spf13/viper"

	"github.com/Logiase/MiraiGo-Template/bot"
	"github.com/Logiase/MiraiGo-Template/utils"
)

func init() {
	instance := &wzry{}
	bot.RegisterModule(instance)
}

type wzry struct {
}

func (m *wzry) MiraiGoModule() bot.ModuleInfo {
	return bot.ModuleInfo{
		ID:       "lurendie.wzry",
		Instance: instance,
	}
}

func (m *wzry) Init() {
	// 初始化过程
	// 在此处可以进行 Module 的初始化配置
	// 如配置读取
}

func (m *wzry) PostInit() {
	// 第二次初始化
	// 再次过程中可以进行跨Module的动作
	// 如通用数据库等等
}

func (m *wzry) Serve(b *bot.Bot) {
	// 注册服务函数部分
	register(b)
}

func (m *wzry) Start(b *bot.Bot) {
	// 此函数会新开携程进行调用
	// ```go
	// 		go exampleModule.Start()
	// ```

	// 可以利用此部分进行后台操作
	// 如http服务器等等
}

func (m *wzry) Stop(b *bot.Bot, wg *sync.WaitGroup) {
	// 别忘了解锁
	defer wg.Done()
	// 结束部分
	// 一般调用此函数时，程序接收到 os.Interrupt 信号
	// 即将退出
	// 在此处应该释放相应的资源或者对状态进行保存
}

var instance *wzry

var wzryLogger = utils.GetModuleLogger("wzry")

func register(b *bot.Bot) {
	b.GroupMessageEvent.Subscribe(func(client *client.QQClient, event *message.GroupMessage) {
		for _, v := range viper.GetIntSlice("gruops") {
			if event.GroupCode == int64(v) {
				if strings.HasPrefix(event.ToString(), "查询战力") {
					wzryLogger.Info("查询战力触发")
					m := searchAtk(client, event)
					client.SendGroupMessage(event.GroupCode, m)
					return
				}
			}
		}
	})

}

func searchAtk(client *client.QQClient, event *message.GroupMessage) *message.SendingMessage {
	arr := strings.Split(event.ToString(), " ")
	heroName := arr[1]
	zone := getType(arr[2])
	m := message.NewSendingMessage()
	URL := fmt.Sprintf(API_URL, token, url.QueryEscape(heroName), zone)
	fmt.Printf("url: %v\n", URL)
	response := internal.RequestJson(URL, internal.GET)
	fmt.Printf("response: %v\n", response)
	//区标
	area := response["area"].(map[string]interface{})
	//市区
	city := response["city"].(map[string]interface{})
	//省区
	province := response["province"].(map[string]interface{})
	server := response["server"].(string)
	name := response["name"].(string)
	updatetime := response["updatetime"]
	make := fmt.Sprintf("服务器:%s\n省标:\n%s-%s\n市标:\n%s-%s\n区标:\n%s-%s\n最后更新时间:%v", server, province["name"], province["power"], city["name"], city["power"], area["name"], area["power"], updatetime)
	m.Append(message.NewAt(event.Sender.Uin)).Append(message.NewText("\n英雄名称:" + name + "\n")).Append(message.NewText(make))
	return m

}

func getType(t string) string {
	if t == "安卓QQ" {
		return Zone.qq
	} else if t == "安卓WX" {
		return Zone.wx
	} else if t == "苹果QQ" {
		return Zone.ios_qq
	} else if t == "苹果WX" {
		return Zone.ios_wx
	}
	return ""
}
