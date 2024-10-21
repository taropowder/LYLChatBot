package utils

import (
	"LYLChatBot/conf"
	"bytes"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

type QwenHandler struct {
}

type QwenMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type QwenResponse struct {
	Output struct {
		FinishReason string `json:"finish_reason"`
		Text         string `json:"text"`
	} `json:"output"`
	Usage struct {
		TotalTokens  int `json:"total_tokens"`
		OutputTokens int `json:"output_tokens"`
		InputTokens  int `json:"input_tokens"`
	} `json:"usage"`
	RequestId string `json:"request_id"`
}

type QwenErrorResponse struct {
	Code      string `json:"code"`
	Message   string `json:"message"`
	RequestId string `json:"request_id"`
}

type QwenRequest struct {
	Input struct {
		Messages []GptMessage `json:"messages"`
	} `json:"input"`
	Model      string `json:"model"`
	Parameters struct {
		EnableSearch bool `json:"enable_search"`
	} `json:"parameters"`
}

func NewQwenGpt(prompt string, systemPrompt string, historyMsgs []GptMessage) (reply string, err error) {

	logrus.Debugf("new qwen gpt %s", prompt)
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

	request := QwenRequest{
		Input: struct {
			Messages []GptMessage `json:"messages"`
		}{messages},
		Model: conf.ConfigureInstance.Gpt.Qwen.Model,
		Parameters: struct {
			EnableSearch bool `json:"enable_search"`
		}{EnableSearch: true},
	}

	body, err := json.Marshal(&request)
	if err != nil {
		return
	}

	client := &http.Client{}

	// 创建请求
	req, err := http.NewRequest("POST", "https://dashscope.aliyuncs.com/api/v1/services/aigc/text-generation/generation", bytes.NewReader(body))
	if err != nil {
		return
	}

	// 设置请求头
	req.Header.Set("Authorization", "Bearer "+conf.ConfigureInstance.Gpt.Qwen.ApiKey)
	req.Header.Set("Content-Type", "application/json")

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		return
	}

	// 读取响应
	response, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	if resp.StatusCode == http.StatusOK {
		// 解析JSON
		var result QwenResponse
		err = json.Unmarshal(response, &result)
		if err != nil {
			//fmt.Println(err)
			return
		}
		// 读取delta.content的值
		reply = result.Output.Text
	} else {
		// 解析JSON
		var result QwenErrorResponse
		err = json.Unmarshal(response, &result)
		if err != nil {
			//fmt.Println(err)
			return
		}
		// 读取delta.content的值
		reply = result.Message
	}

	return reply, nil

}
