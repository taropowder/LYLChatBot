package utils

import (
	"LYLChatBot/constant"
	"github.com/PuerkitoBio/goquery"
	log "github.com/sirupsen/logrus"
	"strings"
	"unicode"
)

var BadChineseText = []string{"系统监测到您的网络环境存在异常"}

func GetChineseText(text string) string {

	runes := []rune(text)
	var chinese []rune
	isChinese := false
	for _, r := range runes {
		if unicode.Is(unicode.Han, r) {
			isChinese = true
			chinese = append(chinese, r)
		}
		if r == '。' || r == '，' || r == '？' || r == '！' || r == ',' || r == '.' {
			chinese = append(chinese, r)
		}
	}

	if !isChinese {
		return ""
	}

	for _, s := range BadChineseText {
		if strings.Contains(string(chinese), s) {
			return ""
		}
	}

	return string(chinese)
}

func GetHtmlText(text string) string {

	res := ""

	reader := strings.NewReader(text)

	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		log.Fatal(err)
	}

	// Remove script and style tags
	doc.Find("script, style").Each(func(i int, s *goquery.Selection) {
		s.Remove()
	})

	doc.Find("*").Each(func(i int, s *goquery.Selection) {
		//fmt.Println(s.Text())
		res = res + s.Text()
	})

	for _, s := range BadChineseText {
		if strings.Contains(string(res), s) {
			return ""
		}
	}

	res = strings.ReplaceAll(strings.ReplaceAll(res, "  ", ""), "\n", "")

	if len(res)/3 > constant.MaxPromptLen {
		return string([]rune(res)[0:constant.MaxPromptLen])
	}
	return res
}
