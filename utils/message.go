package utils

import (
	"github.com/eatmoreapple/openwechat"
	"regexp"
	"strings"
)

func RemoveAt(content string) string {

	// 截取字符串
	return regexp.MustCompile("@.*?[\\s\u2005]").ReplaceAllString(content, "")
}

func FormatMessage(content string, keywords []string) (res string) {

	res = regexp.MustCompile(`@.*?\s`).ReplaceAllString(content, "")

	if len(keywords) > 0 {
		for _, keyword := range keywords {
			res = strings.Replace(res, keyword, "", -1)
		}
	}

	return res
}

func MessageMatchInstruct(s string, message *openwechat.Message) []string {

	if message.IsText() {
		re := regexp.MustCompile(s)

		match := re.FindStringSubmatch(RemoveAt(message.Content))
		if match != nil {
			return match
		} else {
			return nil
		}
	}

	return nil

}
