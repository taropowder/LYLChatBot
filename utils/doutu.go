package utils

import (
	"bytes"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/wangluozhe/requests"
	request_url "github.com/wangluozhe/requests/url"
	"io"
	"net/url"
)

type DouTuResponse struct {
	Path   string `json:"path"`
	Width  int    `json:"width"`
	HashId string `json:"hashId"`
	Height int    `json:"height"`
}

func GetDouTuImagesByParameters(start int, search string) (io.Reader, error) {

	if search != "" {
		search = "&w=" + url.QueryEscape(search)
	}

	req := request_url.NewRequest()
	req.Headers = request_url.ParseHeaders(`
	Host: www.dbbqb.com
	Sec-Ch-Ua: "Not A(Brand";v="99", "Google Chrome";v="121", "Chromium";v="121"
	Web-Agent: web
	Sec-Ch-Ua-Mobile: ?0
	User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/121.0.0.0 Safari/537.36
	Content-Type: application/json
	Accept: application/json
	Client-Id: 
	Sec-Ch-Ua-Platform: "macOS"
	Sec-Fetch-Site: same-origin
	Sec-Fetch-Mode: cors
	Sec-Fetch-Dest: empty
	Referer: https://www.dbbqb.com/
	Accept-Encoding: gzip, deflate
	Accept-Language: zh-CN,zh;q=0.9
    `)

	apiUrl := fmt.Sprintf("https://www.dbbqb.com/api/search/json?start=%d&size=1%s", start, search)
	r, err := requests.Get(apiUrl, req)
	if err != nil {
		return nil, err
	}
	jsonRes, err := r.SimpleJson()

	res := jsonRes.GetIndex(0)

	path, err := res.Get("path").String()
	if err != nil {
		return nil, err

	}

	imgUrl := fmt.Sprintf("https://image.dbbqb.com/%s", path)

	r, err = requests.Get(imgUrl, nil)
	if err != nil {
		return nil, err
	}

	return io.Reader(r.Body), nil
	//return fmt.Sprintf("https://image.dbbqb.com/%s", path)
}

func IsGif(src io.Reader) bool {

	dst := new(bytes.Buffer)
	_, err := io.Copy(dst, src)
	if err != nil {
		log.Error(err)
		return false
	}

	// 读取文件头
	data := make([]byte, 6)
	_, err = dst.Read(data)
	if err != nil {
		return false
	}

	// 判断魔数
	return string(data[:3]) == "GIF"
}
