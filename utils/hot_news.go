package utils

import (
	"LYLChatBot/conf"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"math/rand"
	"net/http"
	"time"
)

// https://open.tophub.today/hot
// https://hot.iczkj.com/#/
// https://hot.iczkj.com/api/bilibili/

var NewsTypes = []string{"bilibili", "baidu", "36kr", "zhihu", "douyin",
	"weibo", "newsqq", "juejin", "tieba", "toutiao",
	"thepaper", "sspai", "genshin", "github", "ngabbs",
	"douban_new", "douban_group", "weread", "netease", "kuaishou", "lol",
}

var MusicTypes = []string{"qq_music_toplist/?type=1", "netease_music_toplist/?type=1"}

type HotNewsResp struct {
	Code       int       `json:"code"`
	Message    string    `json:"message"`
	Name       string    `json:"name"`
	Title      string    `json:"title"`
	Subtitle   string    `json:"subtitle"`
	From       string    `json:"from"`
	Total      int       `json:"total"`
	UpdateTime time.Time `json:"updateTime"`
	Data       []struct {
		Title     string      `json:"title"`
		Desc      string      `json:"desc"`
		Pic       string      `json:"pic"`
		Hot       interface{} `json:"hot"`
		Url       string      `json:"url"`
		MobileUrl string      `json:"mobileUrl"`
	} `json:"data"`
}

func GetHostNews(newsType string) (hostNews *HotNewsResp, err error) {
	resp, err := http.Get(conf.ConfigureInstance.HowNewsApi + newsType)
	if err != nil {
		log.Error(err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
		return
	}

	err = json.Unmarshal(body, &hostNews)
	if err != nil {
		log.Error("解码错误: ", err)
		return
	}
	return
}

// https://hot.miykah.top/#/list?type=weibo&page=1
func FormatHostNews(hostNews *HotNewsResp, newsType string, limit int) (news string) {
	news += fmt.Sprintf("%s\n%s/#/list?type=%s&page=1\n", hostNews.Title, conf.ConfigureInstance.HowNewsSite, newsType)
	i := 0
	for _, v := range hostNews.Data {
		hot, ok := v.Hot.(float64)
		i = i + 1
		if ok {
			news += fmt.Sprintf("%d : %s  [hot:%d] \n", i, v.Title, int(hot))
		} else {
			news += fmt.Sprintf("%d :%s \n", i, v.Title)
		}
		if i > limit {
			break
		}
	}
	return
}

func RandomNews(roundTypesLen, newLimit int) string {

	newNewsTypes := make([]string, len(NewsTypes))
	copy(newNewsTypes, NewsTypes)
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(newNewsTypes), func(i, j int) { newNewsTypes[i], newNewsTypes[j] = newNewsTypes[j], newNewsTypes[i] })

	res := GetNewsByTypes(newNewsTypes[:roundTypesLen], newLimit)
	return res
}

func GetNewsByTypes(newNewsTypes []string, newLimit int) string {

	res := "今日热点：\n"
	for _, t := range newNewsTypes {
		hostNews, err := GetHostNews(t)
		if err != nil {
			log.Error(err)
			continue
		}
		res += FormatHostNews(hostNews, t, newLimit)
		res += "---------------\n"

	}

	return res

}
