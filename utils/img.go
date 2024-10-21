package utils

import (
	"LYLChatBot/conf"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/wangluozhe/requests"
	"github.com/wangluozhe/requests/url"
	"io"
)

func GetImagesByGoogleSearch(start int, search string) (io.Reader, error) {
	req := url.NewRequest()
	req.Proxies = conf.ConfigureInstance.Proxy
	apiUrl := fmt.Sprintf("https://www.googleapis.com/customsearch/v1?key=%s&cx=%s&q=%s&searchType=image&num=1&start=%d",
		conf.ConfigureInstance.Google.ApiKey, conf.ConfigureInstance.Google.CX, search, start)
	r, err := requests.Get(apiUrl, req)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	jsonRes, err := r.SimpleJson()
	//log.Info(jsonRes.GetIndex())
	if jsonRes.Get("items") != nil {
		items := jsonRes.Get("items")
		imgUrl := items.GetIndex(0).Get("link").MustString()
		//log.Info(imgUrl)

		req.Headers = url.ParseHeaders(`
	Sec-Ch-Ua: "Not A(Brand";v="99", "Google Chrome";v="121", "Chromium";v="121"
	Web-Agent: web
	Sec-Ch-Ua-Mobile: ?0
	User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/121.0.0.0 Safari/537.36
	Sec-Ch-Ua-Platform: "macOS"
	Sec-Fetch-Site: same-origin
	Sec-Fetch-Mode: cors
	Sec-Fetch-Dest: empty
	Accept-Encoding: gzip, deflate
	Accept-Language: zh-CN,zh;q=0.9
    `)

		nr, err := requests.Get(imgUrl, req)
		if err != nil {
			log.Error(err)
			return nil, err
		}
		return io.Reader(nr.Body), nil
	}

	return nil, errors.New("goole search error")
}
