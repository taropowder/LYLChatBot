package utils

import (
	"LYLChatBot/conf"
	"LYLChatBot/constant"
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/google/generative-ai-go/genai"
	"github.com/googleapis/gax-go/v2/apierror"
	log "github.com/sirupsen/logrus"
	"google.golang.org/api/option"
	"net/http"
	"net/url"
)

func NewGeminitGpt(prompt string, systemPrompt string, historyMsgs []GptMessage) (reply string, err error) {

	if len(prompt)/3 > constant.MaxPromptLen {
		prompt = string([]rune(prompt)[0:constant.MaxPromptLen])
	}

	log.Debugf("new Geminit gpt %s", prompt)

	ctx := context.Background()
	//Access your API key as an environment variable (see "Set up your API key" above)
	c := &http.Client{Transport: &APIKeyProxyTransport{
		APIKey:    conf.ConfigureInstance.Gpt.Gemini.ApiKey,
		Transport: nil,
		ProxyURL:  conf.ConfigureInstance.Proxy,
	}}

	client, err := genai.NewClient(ctx, option.WithHTTPClient(c))

	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	// For text-only input, use the gemini-pro model
	model := client.GenerativeModel("gemini-1.5-flash")

	model.SafetySettings = []*genai.SafetySetting{
		{
			Category:  genai.HarmCategoryHateSpeech,
			Threshold: genai.HarmBlockNone,
		},
		{
			Category:  genai.HarmCategorySexuallyExplicit,
			Threshold: genai.HarmBlockNone,
		},
		{
			Category:  genai.HarmCategoryDangerousContent,
			Threshold: genai.HarmBlockNone,
		},
		{
			Category:  genai.HarmCategoryHarassment,
			Threshold: genai.HarmBlockNone,
		},
	}

	// Initialize the chat
	cs := model.StartChat()
	cs.History = []*genai.Content{
		&genai.Content{
			Parts: []genai.Part{
				genai.Text(systemPrompt),
			},
			Role: "user",
		}, &genai.Content{
			Parts: []genai.Part{
				genai.Text("好的,我明白了"),
			},
			Role: "model",
		},
	}

	for _, message := range historyMsgs {
		role := "user"
		if message.Role == "assistant" {
			role = "model"
		}
		cs.History = append(cs.History, &genai.Content{
			Parts: []genai.Part{
				genai.Text(message.Content),
			},
			Role: role,
		})
	}

	defer func() {
		if err := recover(); err != nil {
			log.Error(err)
		}
	}()

	resp, err := cs.SendMessage(ctx, genai.Text(prompt))
	if err != nil {
		var e *apierror.APIError
		if errors.As(err, &e) {
			log.Errorf("gemini http error %v", e.Unwrap())
		}
		return "", err
	}

	if resp.Candidates != nil {
		for _, candidate := range resp.Candidates {
			//log.Info(candidate.Content)
			if candidate.Content != nil && candidate.Content.Parts != nil {
				for _, part := range candidate.Content.Parts {
					return fmt.Sprintf("%s", part), nil
				}
			} else {
				log.Errorf("error for candidate %v", candidate)
			}

		}
	}

	return
}

func NewGeminitGptWithImg(prompt string, imgData []byte) (reply string, err error) {

	log.Debugf("new Geminit  img  gpt %s", prompt)

	ctx := context.Background()

	//Access your API key as an environment variable (see "Set up your API key" above)
	c := &http.Client{Transport: &APIKeyProxyTransport{
		APIKey:    conf.ConfigureInstance.Gpt.Gemini.ApiKey,
		Transport: nil,
		ProxyURL:  conf.ConfigureInstance.Gpt.Gemini.Proxy,
	}}

	client, err := genai.NewClient(ctx, option.WithHTTPClient(c))

	if err != nil {
		log.Error(err)
	}
	defer client.Close()

	// For text-only input, use the gemini-pro model
	model := client.GenerativeModel("gemini-1.5-flash")

	model.SafetySettings = []*genai.SafetySetting{
		{
			Category:  genai.HarmCategoryHateSpeech,
			Threshold: genai.HarmBlockNone,
		},
		{
			Category:  genai.HarmCategorySexuallyExplicit,
			Threshold: genai.HarmBlockNone,
		},
		{
			Category:  genai.HarmCategoryDangerousContent,
			Threshold: genai.HarmBlockNone,
		},
		{
			Category:  genai.HarmCategoryHarassment,
			Threshold: genai.HarmBlockNone,
		},
	}

	resp, err := model.GenerateContent(ctx, genai.Text(prompt), genai.ImageData("jpeg", imgData))
	if err != nil {
		log.Error(err)
		return "", err
	}

	if resp.Candidates != nil {
		for _, candidate := range resp.Candidates {
			//log.Info(candidate.Content)
			if candidate.Content != nil && candidate.Content.Parts != nil {
				for _, part := range candidate.Content.Parts {
					return fmt.Sprintf("%s", part), nil
				}
			} else {
				log.Errorf("error for candidate %v", candidate)
			}

		}
	}

	return
}

// 整合了APIKey和代理的 Transport
type APIKeyProxyTransport struct {
	// APIKey is the API Key to set on requests.
	APIKey string

	// Transport is the underlying HTTP transport.
	// If nil, http.DefaultTransport is used.
	Transport http.RoundTripper

	// ProxyURL is the URL of the proxy server. If empty, no proxy is used.
	ProxyURL string
}

func (t *APIKeyProxyTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	rt := t.Transport
	if rt == nil {
		rt = http.DefaultTransport
	}

	// 如果提供了 ProxyURL，则对 Transport 设置代理
	if t.ProxyURL != "" {
		proxyURL, err := url.Parse(t.ProxyURL)
		if err != nil {
			return nil, err
		}
		if transport, ok := rt.(*http.Transport); ok {
			// 只有当 rt 为 *http.Transport 类型时，才设置代理
			transport.Proxy = http.ProxyURL(proxyURL)
			transport.TLSClientConfig = &tls.Config{
				InsecureSkipVerify: true,
			}
		} else {
			// 如果 rt 不是 *http.Transport 类型，则创建一个新的带代理的 http.Transport
			rt = &http.Transport{
				Proxy: http.ProxyURL(proxyURL),
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			}
		}
	}

	// 克隆请求以避免修改原始请求
	newReq := *req
	args := newReq.URL.Query()
	args.Set("key", t.APIKey)
	newReq.URL.RawQuery = args.Encode()

	// 执行 HTTP 请求，并处理可能的错误
	resp, err := rt.RoundTrip(&newReq)
	if err != nil {
		// 返回网络请求中的错误
		return nil, fmt.Errorf("error during round trip: %v", err)
	}

	return resp, nil
}
