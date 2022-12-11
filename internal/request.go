package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/Logiase/MiraiGo-Template/config"
	"github.com/Mrs4s/MiraiGo/client"
	"github.com/Mrs4s/MiraiGo/message"
)

var cli *http.Client

const (
	USER_AGENT = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/107.0.0.0 Safari/537.36 Edg/107.0.1418.24"
	POST       = "POST"
	GET        = "GET"
	PRIVATE    = iota
	GROUP
)

// 根据图片URL返回图片二进制数据
func RequestImg(url string, method string) io.ReadCloser {
	req, _ := http.NewRequest(method, url, nil)
	req.Header.Add("user-agent", USER_AGENT)
	req.Header.Add("token", "eyJhbGciOiJIUzI1NiJ9.eyJzdWIiOiI1NDk3MDA0NTkiLCJ1dWlkIjoiNjMxYjE4MTIyZmZiNGE1OWI1ZmJjNTY2MDgzNmFkNWIiLCJpYXQiOjE2NTAxODA3NzEsImFjY291bnQiOiJ7XCJlbWFpbFwiOlwiNTQ5NzAwNDU5QHFxLmNvbVwiLFwiZ2VuZGVyXCI6LTEsXCJoYXNQcm9uXCI6MCxcImlkXCI6MTc4OCxcInBhc3NXb3JkXCI6XCI2N2JjZDAxZTNlYzc2MWY5ZjU2YzQyZDZkOTdkNGI1OFwiLFwic3RhdHVzXCI6MCxcInVzZXJOYW1lXCI6XCI1NDk3MDA0NTlcIn0iLCJqdGkiOiIxNzg4In0.nuuHLJeVCIfOg_1EEPHiL-nL8O82rCbxyI_PA4-QPBw")
	resp, _ := cli.Do(req)
	return resp.Body
}

// 根据JSON URL返回Map集合
func RequestJson(url string, method string) map[string]interface{} {
	req, _ := http.NewRequest(method, url, nil)
	req.Header.Add("user-agent", USER_AGENT)
	req.Header.Add("token", "eyJhbGciOiJIUzI1NiJ9.eyJzdWIiOiI1NDk3MDA0NTkiLCJ1dWlkIjoiNjMxYjE4MTIyZmZiNGE1OWI1ZmJjNTY2MDgzNmFkNWIiLCJpYXQiOjE2NTAxODA3NzEsImFjY291bnQiOiJ7XCJlbWFpbFwiOlwiNTQ5NzAwNDU5QHFxLmNvbVwiLFwiZ2VuZGVyXCI6LTEsXCJoYXNQcm9uXCI6MCxcImlkXCI6MTc4OCxcInBhc3NXb3JkXCI6XCI2N2JjZDAxZTNlYzc2MWY5ZjU2YzQyZDZkOTdkNGI1OFwiLFwic3RhdHVzXCI6MCxcInVzZXJOYW1lXCI6XCI1NDk3MDA0NTlcIn0iLCJqdGkiOiIxNzg4In0.nuuHLJeVCIfOg_1EEPHiL-nL8O82rCbxyI_PA4-QPBw")
	resp, err := cli.Do(req)
	if err != nil {
		fmt.Printf("error1: %v\n", err)
	}
	jsonData, err2 := io.ReadAll(resp.Body)
	if err2 != nil {
		fmt.Printf("error2: %v\n", err2)
	}
	var j map[string]interface{}
	json.Unmarshal(jsonData, &j)
	return j
}

func init() {
	//设置代理初始化请求头
	if config.GlobalConfig.Get("proxy") == nil {
		cli = &http.Client{Timeout: 0}
	} else {
		proxyURL, error := url.Parse(config.GlobalConfig.GetString("proxy"))
		if error != nil {
			fmt.Println("代理转换异常!!")
		}
		cli = &http.Client{Timeout: 0,
			Transport: &http.Transport{
				Proxy: http.ProxyURL(proxyURL),
			},
		}
	}
}

// 将图片数据上传QQ服务器并生成Message类型
func MakeImage(read io.Reader, c *client.QQClient, groupCode int64, msgtype int) message.IMessageElement {
	imageData, e1 := io.ReadAll(read)
	if e1 != nil {
		fmt.Printf("e1: %v\n", e1)
	}
	dataBuffer := bytes.NewReader(imageData)
	var img message.IMessageElement
	var err error
	if msgtype == 0 {
		img, err = c.UploadImage(message.Source{SourceType: message.SourcePrivate, PrimaryID: groupCode}, dataBuffer)
	} else {
		img, err = c.UploadImage(message.Source{SourceType: message.SourceGroup, PrimaryID: groupCode}, dataBuffer)
	}
	if err != nil {
		fmt.Printf("上传群图片失败: %v\n", err)
	}
	return img
}
