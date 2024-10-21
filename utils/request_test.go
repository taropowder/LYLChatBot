package utils

import (
	"fmt"
	"testing"
)

func TestChrome(t *testing.T) {
	//html, _ := RequestWithChrome("https://mp.weixin.qq.com/s/krQFGQs4OHxgPWk5uulvlA")
	html, err := RequestWechatWithChrome("http://mp.weixin.qq.com/s?__biz=MzkyNzUzOTQzOA==&amp;mid=2247489487&amp;idx=1&amp;sn=fac3492cf99ee0ddde7770e9ae5b80a3&amp;chksm=c227de4ef5505758b2e0b4849913606cd972a31eee10596a0e06510935ef216d76c868d9ffe0&amp;mpshare=1&amp;scene=1&amp;srcid=0730qTvn2ulEKxIiJ3b1BlTC&amp;sharer_shareinfo=80639ef3675881f84329f0950ee41a14&amp;sharer_shareinfo_first=80639ef3675881f84329f0950ee41a14#rd")
	fmt.Println("err", err)
	fmt.Println(html)
	r := GetHtmlText(html)
	fmt.Println(r)
}
