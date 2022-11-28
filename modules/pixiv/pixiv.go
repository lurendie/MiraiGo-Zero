package pixiv

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/Logiase/MiraiGo-Template/config"
	miraiGoCli "github.com/Mrs4s/MiraiGo/client"
	"github.com/Mrs4s/MiraiGo/message"
	"github.com/spf13/viper"

	"github.com/Logiase/MiraiGo-Template/bot"
	"github.com/Logiase/MiraiGo-Template/utils"
)

func init() {
	instance := &pixiv{}
	bot.RegisterModule(instance)
}

type pixiv struct {
}

func (m *pixiv) MiraiGoModule() bot.ModuleInfo {
	return bot.ModuleInfo{
		ID:       "lurendie.pixiv",
		Instance: instance,
	}
}

func (m *pixiv) Init() {
	// 初始化过程
	// 在此处可以进行 Module 的初始化配置
	// 如配置读取
	viper.SetConfigFile("application.yaml")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.ReadInConfig()
}

func (m *pixiv) PostInit() {
	// 第二次初始化
	// 再次过程中可以进行跨Module的动作
	// 如通用数据库等等
}

func (m *pixiv) Serve(b *bot.Bot) {
	// 注册服务函数部分
	register(b)
}

func (m *pixiv) Start(b *bot.Bot) {
	// 此函数会新开携程进行调用
	// ```go
	// 		go exampleModule.Start()
	// ```

	// 可以利用此部分进行后台操作
	// 如http服务器等等
}

func (m *pixiv) Stop(b *bot.Bot, wg *sync.WaitGroup) {
	// 别忘了解锁
	defer wg.Done()
	// 结束部分
	// 一般调用此函数时，程序接收到 os.Interrupt 信号
	// 即将退出
	// 在此处应该释放相应的资源或者对状态进行保存
}

var instance *pixiv

var pixivLogger = utils.GetModuleLogger("pixiv")

func register(b *bot.Bot) {
	b.GroupMessageEvent.Subscribe(func(client *miraiGoCli.QQClient, event *message.GroupMessage) {
		//gruopCode := strconv.FormatInt(event.GroupCode, 10)
		for _, v := range config.GlobalConfig.GetIntSlice("gruops") {
			if int(event.GroupCode) == v {
				if event.ToString() == "功能" {
					pixivLogger.Info("'功能'关键词触发")
					menu := "====菜单====\n1.排行\n2.查看图片<PID>\n3.查看画师<UID>\n4.以图搜图<图片>"
					m := message.NewSendingMessage().Append(message.NewAt(event.Sender.Uin)).Append(message.NewText("\n")).Append(message.NewText(menu))
					client.SendGroupMessage(event.GroupCode, m)
					return
				} else if event.ToString() == "排行" {
					pixivLogger.Info("'排行'关键词触发")
					m := Top50(client, event)
					client.SendGroupMessage(event.GroupCode, m)
					return
				} else if strings.HasPrefix(event.ToString(), "查看图片") {
					m := ShowIllust(client, event, event.ToString()[12:])
					client.SendGroupMessage(event.GroupCode, m)
					return
				} else if event.ToString() == "涩图" {
					pixivLogger.Info("'涩图'关键词触发")
					m := setu(client, event)
					client.SendGroupTempMessage(event.GroupCode, event.Sender.Uin, m)
					client.SendGroupMessage(event.GroupCode, message.NewSendingMessage().Append(message.NewAt(event.Sender.Uin)).Append(message.NewText("\n不可以涩涩哦!")))
					return
				} else if strings.HasPrefix(event.ToString(), "查看画师") {
					pixivLogger.Info("'查看画师'关键词触发")
					m := ShowUser(client, event, event.ToString()[12:])
					client.SendGroupMessage(event.GroupCode, m)
					return
				} else if strings.HasPrefix(event.ToString(), "以图搜图") {
					pixivLogger.Info("'以图搜图'关键词触发")
					imgURL := event.Elements[1].(*message.GroupImageElement).Url
					m := searchImg(client, event, imgURL)
					client.SendGroupMessage(event.GroupCode, m)
					return
				}
			}
		}
	})

}

// 每日排行 不是当日数据,延迟一天
func Top50(client *miraiGoCli.QQClient, event *message.GroupMessage) *message.SendingMessage {
	sendMsg := message.NewSendingMessage()
	sendMsg.Append(message.NewAt(event.Sender.Uin)).Append(message.NewText("\n"))
	t := time.Now()
	date := fmt.Sprintf("%d-%02d-%02d", t.Year(), t.Month(), t.Day()-2)
	rank := RequestJson(fmt.Sprintf(RankURL, date), GET)
	if len(rank["illusts"].([]interface{})) > 0 {
		for i, v := range rank["illusts"].([]interface{}) {
			if i == 5 {
				break
			}
			t := fmt.Sprintf("标题:%s", v.(map[string]interface{})["title"])
			p := fmt.Sprintf("PID:%1.0f", v.(map[string]interface{})["id"].(float64))
			makes := fmt.Sprintf(t + "\n" + p + "\n")
			ImgURL := v.(map[string]interface{})["image_urls"].(map[string]interface{})["large"].(string)
			imgData := RequestImg(fmt.Sprintf(ImgDataURL, ImgURL), GET)
			img := makeImage(imgData, client, event.GroupCode, Group)
			sendMsg.Append(message.NewText(makes)).Append(img)
		}
	} else {
		pixivLogger.Error("一个数据都没有咋回事小老弟?(API)")
		sendMsg.Append(message.NewText("获取数据失败!,等稍后再试...\n"))
	}
	sendMsg.Append(message.NewText("PS:排名不是实时更新!"))
	return sendMsg
}

// 查看图片
func ShowIllust(client *miraiGoCli.QQClient, event *message.GroupMessage, illustId string) *message.SendingMessage {
	sendMsg := message.NewSendingMessage()
	sendMsg.Append(message.NewAt(event.Sender.Uin)).Append(message.NewText("\n"))
	illust := RequestJson(fmt.Sprintf(IllustURL, illustId), GET)
	send := "作品ID:" + illust["data"].(map[string]interface{})["illust"].(string) + "\n标题:" + illust["data"].(map[string]interface{})["title"].(string) + "\n画师UID:" + illust["data"].(map[string]interface{})["user"].(map[string]interface{})["id"].(string)
	imgURL := illust["data"].(map[string]interface{})["originals"].([]interface{})[0].(map[string]interface{})["url"].(string)
	imgData := RequestImg(fmt.Sprintf(ImgDataURL, imgURL), GET)
	img := makeImage(imgData, client, event.GroupCode, 1)
	sendMsg.Append(message.NewText(send)).Append(img)
	return sendMsg
}

// 查看作者
func ShowUser(client *miraiGoCli.QQClient, event *message.GroupMessage, userId string) *message.SendingMessage {
	sendMsg := message.NewSendingMessage()
	sendMsg.Append(message.NewAt(event.Sender.Uin)).Append(message.NewText("\n"))
	userInfo := RequestJson(fmt.Sprintf(UserURL, userId), GET)
	imgData := RequestImg(fmt.Sprintf(ImgDataURL, userInfo["user"].(map[string]interface{})["profile_image_urls"].(map[string]interface{})["medium"].(string)), GET)
	img := makeImage(imgData, client, event.GroupCode, 1)

	sendMsg.Append(message.NewText(fmt.Sprintf("画师ID:%s\n画师名称:%s", userId, userInfo["user"].(map[string]interface{})["name"].(string))))
	sendMsg.Append(img)
	for _, v := range userInfo["illusts"].([]interface{}) {
		sendMsg.Append(message.NewText(fmt.Sprintf("%1.0f\n", v.(map[string]interface{})["id"])))
	}
	return sendMsg
}

// 涩图
func setu(client *miraiGoCli.QQClient, event *message.GroupMessage) *message.SendingMessage {
	sendMsg := message.NewSendingMessage()
	illust := RequestJson(SetuURL, GET)
	illustID := illust["data"].([]interface{})[0].(map[string]interface{})["pid"]
	original := illust["data"].([]interface{})[0].(map[string]interface{})["urls"].(map[string]interface{})["original"].(string)
	sendMsg.Append(message.NewText(fmt.Sprintf("PID:%1.0f\n", illustID)))
	imgData := RequestImg(original, GET)
	img := makeImage(imgData, client, event.Sender.Uin, Private)
	sendMsg.Append(img)
	return sendMsg
}

// 以图搜图
func searchImg(client *miraiGoCli.QQClient, event *message.GroupMessage, url string) *message.SendingMessage {
	sendMsg := message.NewSendingMessage()
	fmt.Printf("URL: %v\n", fmt.Sprintf(searchImgURL, API_KEY, url))
	responseJson := RequestJson(fmt.Sprintf(searchImgURL, API_KEY, url), GET)
	//匹配度
	similarity := responseJson["results"].([]interface{})[0].(map[string]interface{})["header"].(map[string]interface{})["similarity"].(string)
	thumbnail := responseJson["results"].([]interface{})[0].(map[string]interface{})["header"].(map[string]interface{})["thumbnail"].(string)
	//标题
	title := responseJson["results"].([]interface{})[0].(map[string]interface{})["data"].(map[string]interface{})["title"].(string)
	//原图链接
	ext_urls := responseJson["results"].([]interface{})[0].(map[string]interface{})["data"].(map[string]interface{})["ext_urls"].([]interface{})[0].(string)
	img := makeImage(RequestImg(thumbnail, GET), client, event.GroupCode, Group)
	sendMsg.Append(message.NewAt(event.Sender.Uin)).Append(message.NewText("\n标题:" + title + "\n匹配度:" + similarity)).Append(img).Append(message.NewText("\n原图链接:" + ext_urls))
	return sendMsg
}
