package chat

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"

	"github.com/Logiase/MiraiGo-Template/config"

	"github.com/Logiase/MiraiGo-Template/bot"
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
		ID:       "ChatGPT",
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

var chatLogger = utils.GetModuleLogger("chatGPT-3.5")

func register(b *bot.Bot) {
	if config.GlobalConfig.GetBool("modules.ChatGPT.start") {
		b.PrivateMessageEvent.Subscribe(privateChatHandle)
	}
	//判断功能是否开启
	defer chatLogger.Debug("ChatGPT的功能状态是:", config.GlobalConfig.GetBool("modules.ChatGPT.start"))
}

// 私聊消息业务处理
func privateChatHandle(client *client.QQClient, event *message.PrivateMessage) {
	defer chatLogger.Info("Chat内容已发送")
	m := message.NewSendingMessage().Append(message.NewText(ChatGPT(event.ToString())))
	client.SendPrivateMessage(event.Sender.Uin, m)
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

// gpt-3.5-turbo 模式
type GPTturbo struct {
	Model    string     `json:"model"`
	Messages []Messages `json:"messages"`
}

type Messages struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Response struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int    `json:"created"`
	Choices []struct {
		Index   int `json:"index"`
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

// 通过消息
func ChatGPT(userMsg string) string {

	ur := "https://api.openai.com/v1/chat/completions"
	apiKey := config.GlobalConfig.GetString("modules.ChatGPT.APIKey")
	if apiKey == "" {
		chatLogger.Error("GPT的APIKay是空,无法执行")
		panic("GPT的APIKay是空,无法执行")
	}
	ms := append(make([]Messages, 0),
		Messages{Role: "user",
			Content: userMsg})
	request := GPTturbo{
		Model:    "gpt-3.5-turbo",
		Messages: ms,
	}
	jsonStr, _ := json.Marshal(request)
	//fmt.Printf("jsonStr: %v\n", string(jsonStr))
	req, err := http.NewRequest("POST", ur, bytes.NewBuffer(jsonStr))
	if err != nil {
		panic(err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)
	//转换代理
	client := &http.Client{Timeout: 0}
	if config.GlobalConfig.GetString("proxy") != "" {
		proxyURL, err := url.Parse(config.GlobalConfig.GetString("proxy"))
		if err != nil {
			fmt.Printf("err: %v\n", err)
		}
		client = &http.Client{Timeout: 0,
			Transport: &http.Transport{
				Proxy: http.ProxyURL(proxyURL),
			},
		}
	}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var response Response
	json.Unmarshal(body, &response)
	var str string
	for _, v := range response.Choices {
		str += v.Message.Content
	}
	return str
}
