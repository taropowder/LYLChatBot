package utils

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	log "github.com/sirupsen/logrus"
	"regexp"
	"strings"
)

// WechatArticle struct
type WechatArticle struct {
	Title       string
	Author      string
	AppName     string
	AppID       string
	Cover       string
	Intro       string
	HtmlContent string // html内容
	MdContent   string // md内容
	TextContent string // text内容
	PubAt       string
	URL         string // 链接地址
	SourceURL   string // 原文地址（当公众号的文章时转载时）
	RoundHead   string
	OriHead     string
	WxID        string
	WxIntro     string
	Copyright   string
	Video       string
	Audio       string
	Images      []string
}

// 抓取公众号文章
func CollectWechatArticle(url string) (article WechatArticle, err error) {
	res, err := RequestWechatWithChrome(url)
	log.Debug(res)
	if err != nil {
		res, err = RequestWithChrome(url)
		log.Debug(res)
		if err != nil {
			return
		}
	}

	reader := strings.NewReader(res)

	document, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return
	}
	article, err = collectWechatArticleCommon(document)
	if err != nil {
		return
	}

	article.SourceURL = url

	title, htmlContent, textContent := collectWechatArticleContent(document)
	if len(title) == 0 || len(htmlContent) == 0 {
		err = errors.New("标题或内容为空")
		return
	}

	article.Title = title
	article.HtmlContent = htmlContent
	article.TextContent = textContent
	return
}

func findString(pattern, value, key string) string {
	Map := selectString(pattern, value)
	if val, ok := Map[key]; ok {
		return val
	}
	return ""
}

// selectString 数组匹配
func selectString(pattern, value string) map[string]string {
	Exp := myRegexp{regexp.MustCompile(pattern)}
	Map := Exp.findStringSubmatchMap(value)
	return Map
}

type myRegexp struct {
	*regexp.Regexp
}

func (r *myRegexp) findStringSubmatchMap(s string) map[string]string {
	captures := make(map[string]string)

	match := r.FindStringSubmatch(s)
	if match == nil {
		return captures
	}

	for i, name := range r.SubexpNames() {
		if i == 0 {
			continue
		}
		captures[name] = match[i]
	}
	return captures
}

// 抓取通用信息
func collectWechatArticleCommon(document *goquery.Document) (article WechatArticle, err error) {
	html, _ := document.Html()

	// 视频地址
	article.Video, _ = document.Find("iframe").Eq(0).Attr("data-src")
	if article.Video == `` {
		article.Video, _ = document.Find("video").Eq(0).Attr("src")
	}

	// 音频地址
	audio, _ := document.Find("mpvoice").Eq(0).Attr("voice_encode_fileid")
	if audio != `` {
		article.Audio = fmt.Sprintf("https://res.wx.qq.com/voice/getvoice?mediaid=%v", audio)
	} else {
		article.Audio, _ = document.Find("audio").Eq(0).Attr("src")
	}

	article.AppID = strings.TrimSpace(findString(`var user_name = "(?P<user_name>[^"]+)";`, html, "user_name"))

	article.AppName = strings.TrimSpace(findString(`var nickname = "(?P<nickname>[^"]+)";`, html, "nickname"))

	//if article.AppName == "" {
	//	err = errors.New("无法获取文章信息")
	//	return
	//}

	article.Title = strings.TrimSpace(findString(`var msg_title = "(?P<title>[^"]+)";`, html, "title"))

	article.Intro = strings.TrimSpace(findString(`var msg_desc = "(?P<intro>[^"]+)";`, html, "intro"))

	article.WxID = strings.TrimSpace(findString(`<label class="profile_meta_label">微信号</label>(?P<intro>[\s]+)<span class="profile_meta_value">(?P<wxid>[^"]+)</span>`, html, "wxid"))

	article.WxIntro = strings.TrimSpace(findString(`<label class="profile_meta_label">功能介绍</label>(?P<intro>[\s]+)<span class="profile_meta_value">(?P<wxintro>[^"]+)</span>`, html, "wxintro"))

	article.Cover = strings.TrimSpace(findString(`var msg_cdn_url = "(?P<cover>[^"]+)";`, html, "cover"))

	article.RoundHead = strings.TrimSpace(findString(`var round_head_img = "(?P<round_head>[^"]+)";`, html, "round_head"))

	article.OriHead = strings.TrimSpace(findString(`var ori_head_img_url = "(?P<ori_head>[^"]+)";`, html, "ori_head"))

	article.PubAt = strings.TrimSpace(findString(`var ct = "(?P<date>\d+)";`, html, "date"))

	article.Copyright = strings.TrimSpace(findString(`var _copyright_stat = "(?P<copyright>\d+)";`, html, "copyright"))

	article.Author = strings.TrimSpace(findString(`<span class="rich_media_meta rich_media_meta_text">(?P<author>[^<]+)</span>`, html, "author"))

	link := strings.TrimSpace(findString(`var msg_link = "(?P<url>[^"]+)";`, html, "url"))

	link2 := strings.TrimSpace(findString(`var msg_source_url = '(?P<url>[^']+)';`, html, "url"))

	article.SourceURL = strings.Replace(link2, `\x26amp;`, "&", -1)

	// 处理特殊字符
	article.URL = strings.Replace(link, `\x26amp;`, "&", -1)
	article.URL = strings.Replace(article.URL, `http://`, `https://`, 1)
	article.URL = strings.Replace(article.URL, `&amp;`, `&`, -1)
	article.URL = strings.Replace(article.URL, `#rd`, "&scene=27#wechat_redirect", 1)

	article.Title = strings.Replace(article.Title, `\x26quot;`, `"`, -1)
	article.Title = strings.Replace(article.Title, `\x0a`, "\n", -1)
	article.Title = strings.Replace(article.Title, `\x26gt;`, `>`, -1)
	article.Title = strings.Replace(article.Title, `\x26lt;`, `<`, -1)
	article.Title = strings.Replace(article.Title, `\x26amp;`, `&`, -1)
	article.Title = strings.Replace(article.Title, `\x26#39;`, `'`, -1)

	article.Intro = strings.Replace(article.Intro, `\x26quot;`, `"`, -1)
	article.Intro = strings.Replace(article.Intro, `\x0a`, "\n", -1)
	article.Intro = strings.Replace(article.Intro, `\x26gt;`, `>`, -1)
	article.Intro = strings.Replace(article.Intro, `\x26lt;`, `<`, -1)
	article.Intro = strings.Replace(article.Intro, `\x26amp;`, `&`, -1)
	article.Intro = strings.Replace(article.Intro, `\x26#39;`, `'`, -1)

	return article, nil
}

// 抓取内容
func collectWechatArticleContent(document *goquery.Document) (title, htmlContent, textContent string) {
	articleSelection := document.Find("#js_article #page-content #img-content")
	titleSelection := articleSelection.Find(".rich_media_title#activity-name")
	contentSelection := articleSelection.Find("div.rich_media_content#js_content")

	// 处理内容
	handleContent(contentSelection)

	title = strings.TrimSpace(titleSelection.Text())
	htmlContentTemp, _ := contentSelection.Html()
	htmlContent = strings.TrimSpace(htmlContentTemp)
	textContent = strings.TrimSpace(contentSelection.Text())
	return
}

// 处理内容
func handleContent(selection *goquery.Selection) {
	selection.Find("img").Each(func(i int, selection *goquery.Selection) {
		src, exists := selection.Attr("data-src")
		if !exists || src == "" {
			return
		}

		//output, err := CopyImage(src)
		//if err == nil {
		//	selection.SetAttr("src", output)
		//} else {
		//	log.Error(err)
		//}
	})

	attrs := []string{"class", "id", "onclick", "onmouseover", "data-mpa-powered-by", "data-mpa-template-id",
		"data-mpa-color", "data-mpa-category", "data-tools", "data-id", "data-ratio", "data-s", "data-src", "data-type",
		"data-h", "data-w", "data-backh", "data-backw", "data-before-oversubscription-url", "data-linktype",
		"data-itemshowtype", "data-miniprogram-appid", "data-miniprogram-path", "data-miniprogram-nickname",
		"data-miniprogram-title", "data-miniprogram-imageurl", "data-miniprogram-type", "data-miniprogram-servicetype",
		"data-width", "data-height", "data-bdless", "data-bdlessp", "data-miniprogram-avatar", "data-bdopacity"}
	cleanAttrs(selection, attrs...)
}

// 删除html中不需要的attr
func cleanAttrs(selection *goquery.Selection, attrs ...string) {
	selection.Each(func(i int, selection *goquery.Selection) {
		for _, attr := range attrs {
			selection.RemoveAttr(attr)
		}
	})
	children := selection.Children()
	if children.Size() > 0 {
		cleanAttrs(children, attrs...)
	}
}
