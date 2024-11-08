package utils

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	log "github.com/sirupsen/logrus"
	"github.com/wangluozhe/requests"
	request_url "github.com/wangluozhe/requests/url"
	"io"
	"math/rand"
	"strings"
	"time"
)

func GetImgFromUrl(imgUrl string, num int) (io.Reader, error) {
	req := request_url.NewRequest()
	r, err := requests.Get(imgUrl, req)
	if err != nil {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(r.Text))
	if err != nil {
		log.Fatal(err)
	}

	resImg := io.Reader(r.Body)

	doc.Find("img.ui.image.lazy").Each(func(i int, s *goquery.Selection) {
		src, exists := s.Attr("data-original")
		if exists {
			//fmt.Println(src)
			if i == num {
				log.Info(src)
				reqImg := request_url.NewRequest()
				//proxyUrl := "http://127.0.0.1:8080" // 替换为你的代理地址
				//reqImg.Proxies = proxyUrl
				reqImg.Headers = request_url.ParseHeaders(`
Host: img.soutula.com
Sec-Ch-Ua: "Google Chrome";v="129", "Not=A?Brand";v="8", "Chromium";v="129"
Sec-Ch-Ua-Mobile: ?0
Sec-Ch-Ua-Platform: "macOS"
Upgrade-Insecure-Requests: 1
User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/129.0.0.0 Safari/537.36
Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7
Sec-Fetch-Site: cross-site
Sec-Fetch-Mode: navigate
Sec-Fetch-User: ?1
Sec-Fetch-Dest: document
Referer: https://fabiaoqing.com/
Accept-Encoding: gzip, deflate, br
Accept-Language: zh-CN,zh;q=0.9
Priority: u=0, i
				`)
				rImg, err := requests.Get(src, reqImg)
				if err == nil {
					resImg = io.Reader(rImg.Body)
				}
			}

		}
	})
	return resImg, nil
}

func GetFaBiaoQingImagesByParameters(start int, search string) (io.Reader, error) {

	rand.Seed(time.Now().UnixNano())
	if search == "" {
		start = rand.Intn(360) + 1
		apiUrl := fmt.Sprintf("https://fabiaoqing.com/biaoqing/lists/page/%d.html", start)
		start = rand.Intn(15) + 1
		return GetImgFromUrl(apiUrl, start)
	} else {
		apiUrl := fmt.Sprintf("https://fabiaoqing.com/search/bqb/keyword/%s/type/bq/page/1.html", search)
		if start == 0 {
			start = rand.Intn(15) + 1
		}
		return GetImgFromUrl(apiUrl, start)
	}

	//if search != "" {
	//	search = "&w=" + url.QueryEscape(search)
	//}
	//
	//req := request_url.NewRequest()
	//req.Headers = request_url.ParseHeaders(`
	//Host: www.dbbqb.com
	//Sec-Ch-Ua: "Not A(Brand";v="99", "Google Chrome";v="121", "Chromium";v="121"
	//Web-Agent: web
	//Sec-Ch-Ua-Mobile: ?0
	//User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/121.0.0.0 Safari/537.36
	//Content-Type: application/json
	//Accept: application/json
	//Client-Id:
	//Sec-Ch-Ua-Platform: "macOS"
	//Sec-Fetch-Site: same-origin
	//Sec-Fetch-Mode: cors
	//Sec-Fetch-Dest: empty
	//Referer: https://www.dbbqb.com/
	//Accept-Encoding: gzip, deflate
	//Accept-Language: zh-CN,zh;q=0.9
	//`)
	//
	//apiUrl := fmt.Sprintf("https://www.dbbqb.com/api/search/json?start=%d&size=1%s", start, search)
	//r, err := requests.Get(apiUrl, req)
	//if err != nil {
	//	return nil, err
	//}
	//jsonRes, err := r.SimpleJson()
	//
	//res := jsonRes.GetIndex(0)
	//
	//path, err := res.Get("path").String()
	//if err != nil {
	//	return nil, err
	//
	//}
	//
	//imgUrl := fmt.Sprintf("https://image.dbbqb.com/%s", path)
	//
	//r, err = requests.Get(imgUrl, nil)
	//if err != nil {
	//	return nil, err
	//}

	//return io.Reader(r.Body), nil
	//return fmt.Sprintf("https://image.dbbqb.com/%s", path)
	return nil, nil
}
