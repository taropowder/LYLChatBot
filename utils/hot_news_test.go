package utils

import (
	"fmt"
	"testing"
)

func TestGetHostNews(t *testing.T) {
	//news, _ := GetHostNews("bilibili")
	//news, _ := GetHostNews("weibo")
	//newsStr := FormatHostNews(news, 10)
	//newsStr := RandomNews(3, 2)
	initConfig()
	newsStr := GetNewsByTypes([]string{"bilibili"}, 3)
	fmt.Println(newsStr)
}
