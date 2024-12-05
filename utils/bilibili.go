package utils

import (
	"LYLChatBot/conf"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"net/url"
)

type BiliBiliPageListResp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Ttl     int    `json:"ttl"`
	Data    []struct {
		Cid       int    `json:"cid"`
		Page      int    `json:"page"`
		From      string `json:"from"`
		Part      string `json:"part"`
		Duration  int    `json:"duration"`
		Vid       string `json:"vid"`
		Weblink   string `json:"weblink"`
		Dimension struct {
			Width  int `json:"width"`
			Height int `json:"height"`
			Rotate int `json:"rotate"`
		} `json:"dimension"`
		FirstFrame string `json:"first_frame"`
	} `json:"data"`
}

type BliBliPlayerResp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Ttl     int    `json:"ttl"`
	Data    struct {
		Aid  int    `json:"aid"`
		Bvid string `json:"bvid"`

		Subtitle struct {
			AllowSubmit bool   `json:"allow_submit"`
			Lan         string `json:"lan"`
			LanDoc      string `json:"lan_doc"`
			Subtitles   []struct {
				Id          int64  `json:"id"`
				Lan         string `json:"lan"`
				LanDoc      string `json:"lan_doc"`
				IsLock      bool   `json:"is_lock"`
				SubtitleUrl string `json:"subtitle_url"`
				Type        int    `json:"type"`
				IdStr       string `json:"id_str"`
				AiType      int    `json:"ai_type"`
				AiStatus    int    `json:"ai_status"`
			} `json:"subtitles"`
		} `json:"subtitle"`
	} `json:"data"`
}

type BiliBliSubTitleResp struct {
	FontSize        float64 `json:"font_size"`
	FontColor       string  `json:"font_color"`
	BackgroundAlpha float64 `json:"background_alpha"`
	BackgroundColor string  `json:"background_color"`
	Stroke          string  `json:"Stroke"`
	Type            string  `json:"type"`
	Lang            string  `json:"lang"`
	Version         string  `json:"version"`
	Body            []struct {
		From     float64 `json:"from"`
		To       float64 `json:"to"`
		Sid      int     `json:"sid"`
		Location int     `json:"location"`
		Content  string  `json:"content"`
	} `json:"body"`
}

func DealWithBliBli(bid string) string {
	res := ""
	reqUrl := fmt.Sprintf("https://api.bilibili.com/x/player/pagelist?bvid=%s", bid)
	respContent := requestWithUa(reqUrl)
	plresp := BiliBiliPageListResp{}
	err := json.Unmarshal([]byte(respContent), &plresp)
	if err != nil {
		logrus.Error(err)
		return ""
	}
	if len(plresp.Data) > 0 {
		for _, data := range plresp.Data {
			reqUrl = fmt.Sprintf("https://api.bilibili.com/x/player/wbi/v2?bvid=%s&cid=%d", bid, data.Cid)
			respContent := requestWithUa(reqUrl)
			//logrus.Info(respContent)
			playerResp := BliBliPlayerResp{}
			err := json.Unmarshal([]byte(respContent), &playerResp)
			if err != nil {
				logrus.Error(err)
				return ""
			}

			for _, Subtitle := range playerResp.Data.Subtitle.Subtitles {
				logrus.Info(Subtitle.SubtitleUrl)
				reqUrl = fmt.Sprintf("https:%s", Subtitle.SubtitleUrl)
				respContent := requestWithUa(reqUrl)
				s := BiliBliSubTitleResp{}
				err := json.Unmarshal([]byte(respContent), &s)
				if err != nil {
					logrus.Error(err)
					return ""
				}
				for _, body := range s.Body {
					subtitle := fmt.Sprintf("%v: %s\n", formatSeconds(body.From), body.Content)
					res = res + subtitle
				}

			}
		}
	}

	return res
}

func formatSeconds(seconds float64) string {
	minutes := int(seconds / 60)
	remainingSeconds := int(seconds) % 60

	return fmt.Sprintf("%d分%d秒", minutes, remainingSeconds)
}

func requestWithUa(reqUrl string) string {
	ua := "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36"

	req, err := http.NewRequest("GET", reqUrl, nil)
	if err != nil {
		logrus.Error("创建请求失败:", err)
		return ""
	}

	cookieValue := conf.ConfigureInstance.Handlers.Abstract.Cookie

	cookie := &http.Cookie{
		Name:  "SESSDATA",
		Value: cookieValue,
	}
	req.AddCookie(cookie)

	req.Header.Set("User-Agent", ua)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		logrus.Error("发送请求失败:", err)
		return ""
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logrus.Error("Error reading response body:", err)
		return ""
	}
	return string(body)
}

type WechatAppMsg struct {
	XMLName xml.Name `xml:"msg"`
	URL     string   `xml:"appmsg>url"`
}

func GetBliBliBvInMessageContent(content string) string {
	var msg WechatAppMsg
	err := xml.Unmarshal([]byte(content), &msg)
	if err != nil {
		logrus.Error("解析 XML 失败:", err)
		return ""
	}

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	resp, err := client.Get(msg.URL)
	if err != nil {
		logrus.Error("发送请求失败:", err)
		return ""
	}
	defer resp.Body.Close()

	// 获取重定向后的地址
	location := resp.Header.Get("Location")

	u, err := url.Parse(location)
	if err != nil {
		logrus.Error("解析 URL 失败:", err)
		return ""
	}

	if len(u.Path) > 8 {
		return u.Path[7:]

	} else {
		return ""
	}
}
