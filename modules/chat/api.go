package chat

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

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

func ChatGPT(msg string) string {
	ur := "https://api.openai.com/v1/chat/completions"
	apiKey := "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"
	ms := append(make([]Messages, 0),
		Messages{Role: "user",
			Content: msg})
	request := GPTturbo{
		Model:    "gpt-3.5-turbo",
		Messages: ms,
	}
	jsonStr, _ := json.Marshal(request)
	fmt.Printf("jsonStr: %v\n", string(jsonStr))
	req, err := http.NewRequest("POST", ur, bytes.NewBuffer(jsonStr))
	if err != nil {
		panic(err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)
	proxyURL, err := url.Parse("http://127.0.0.1:10810")
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
	client := &http.Client{Timeout: 0,
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		},
	}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	var response Response
	json.Unmarshal(body, &response)
	var str string
	for _, v := range response.Choices {
		str += v.Message.Content
	}
	return str
}
