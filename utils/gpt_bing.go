package utils

import (
	"LYLChatBot/conf"
	"bytes"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

type NewBingMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type NewBingResponse struct {
	Choices []struct {
		Delta struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"delta"`
		Message NewBingMessage `json:"message"`
	} `json:"choices"`
}

type NewBingRequest struct {
	Messages []GptMessage `json:"messages"`
	Stream   bool         `json:"stream"`
	Model    string       `json:"model"`
}

func NewBingGpt(prompt string, systemPrompt string, historyMsgs []GptMessage) (reply string, err error) {

	logrus.Debugf("new bing gpt %s", prompt)
	messages := make([]GptMessage, 0)
	messages = append(messages, GptMessage{
		Role:    "system",
		Content: systemPrompt,
	})

	messages = append(messages, historyMsgs...)

	messages = append(messages, GptMessage{
		Role:    "user",
		Content: prompt,
	})

	request := NewBingRequest{
		Messages: messages,
		Stream:   false,
		// Creative、Balanced、Precise、gpt-4
		Model: "Creative",
	}
	// 构建请求体
	body, err := json.Marshal(&request)
	if err != nil {
		logrus.Error(err)
		return
	}

	client := &http.Client{
		Timeout: time.Duration(conf.ConfigureInstance.Gpt.NewBing.Timeout) * time.Second,
	}

	if conf.ConfigureInstance.Gpt.NewBing.Proxy != "" {
		// 创建一个代理URL
		proxyURL, err := url.Parse(conf.ConfigureInstance.Gpt.NewBing.Proxy)
		if err != nil {
			logrus.Warnf("new bing proxy error %s", err)
			return "", err
		}

		// 创建一个自定义的Transport，并设置代理
		transport := &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		}

		client.Transport = transport
	}

	tryTimes := 0

	for tryTimes <= conf.ConfigureInstance.Gpt.NewBing.RetryTimes {
		tryTimes = tryTimes + 1

		logrus.Infof("try %v times for %v", tryTimes, prompt)

		// 创建请求
		req, err := http.NewRequest("POST", conf.ConfigureInstance.Gpt.NewBing.URL, bytes.NewReader(body))
		if err != nil {
			logrus.Error(err)
			continue
		}

		// 设置请求头
		req.Header.Set("Content-Type", "application/json")

		if conf.ConfigureInstance.Gpt.NewBing.Headers != nil {

			for key, value := range conf.ConfigureInstance.Gpt.NewBing.Headers {
				req.Header.Set(key, value)
			}

		}

		// 发送请求
		resp, err := client.Do(req)
		if err != nil {
			continue
		}

		// 读取响应
		response, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			logrus.Error(string(response))
			continue
		}

		// 解析JSON
		var result NewBingResponse
		err = json.Unmarshal(response, &result)
		if err != nil {
			logrus.Error(string(response))
			//fmt.Println(err)
			continue
		}

		// 读取delta.content的值
		content := result.Choices[0].Delta.Content
		return content, nil

	}
	return "", nil
}
