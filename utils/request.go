package utils

import (
	"LYLChatBot/conf"
	"LYLChatBot/pkg/redis_conn"
	"context"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"github.com/gomodule/redigo/redis"
	log "github.com/sirupsen/logrus"
	"time"
)

func RequestWithChrome(url string) (string, error) {

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),                        // 不开启图像界面
		chromedp.Flag("mute-audio", false),                     // 关闭声音
		chromedp.Flag("blink-settings", "imagesEnabled=false"), //开启图像界面,重点是开启这个

	)

	if conf.ConfigureInstance.Proxy != "" {
		opts = append(opts, chromedp.ProxyServer(conf.ConfigureInstance.Proxy))
	}

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel := chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))
	defer cancel()

	timeoutCtx, cancel := context.WithTimeout(ctx, 40*time.Second)
	defer cancel()

	var body string
	err := chromedp.Run(timeoutCtx,
		chromedp.Navigate(url),
		// 等待页面加载
		chromedp.Sleep(5*time.Second),
		chromedp.OuterHTML("html", &body),
	)
	if err != nil {
		log.Error(err)
		return "", err
	}

	return body, nil

}

func SetCookie(name, value, domain, path string, httpOnly, secure bool) chromedp.Action {
	return chromedp.ActionFunc(func(ctx context.Context) error {
		expr := cdp.TimeSinceEpoch(time.Now().Add(180 * 24 * time.Hour))
		err := network.SetCookie(name, value).
			WithExpires(&expr).
			WithDomain(domain).
			WithPath(path).
			WithHTTPOnly(httpOnly).
			WithSecure(secure).
			Do(ctx)
		if err != nil {
			return err
		}
		return nil
	})
}
func RequestWechatWithChrome(url string) (string, error) {

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),                        // 不开启图像界面
		chromedp.Flag("mute-audio", false),                     // 关闭声音
		chromedp.Flag("blink-settings", "imagesEnabled=false"), //开启图像界面,重点是开启这个
	)

	if conf.ConfigureInstance.Proxy != "" {
		opts = append(opts, chromedp.ProxyServer(conf.ConfigureInstance.Proxy))
	}

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel := chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))
	defer cancel()

	timeoutCtx, cancel := context.WithTimeout(ctx, 40*time.Second)
	defer cancel()

	var body string
	//var result string

	redisPool := redis_conn.RedisConnPool.Get()
	defer redisPool.Close()

	WechatPocSid, err := redis.String(redisPool.Do("GET", "wechat_poc_sid"))
	if err != nil {
		log.Debug("获取 value 失败:", err)
	}

	err = chromedp.Run(timeoutCtx,
		SetCookie("poc_sid", WechatPocSid, "mp.weixin.qq.com", "/", false, false),
		chromedp.Navigate(url),
		// 等待页面加载
		//chromedp.Sleep(5*time.Second),

		//chromedp
		//chromedp.EvaluateAsDevTools("document.getElementById('js_verify').click()", &result),
		//chromedp.Click(`js_verify`, chromedp.ByID),
		//chromedp.Sleep(10*time.Second),
		//chromedp.Click(`//*[@id="js_verify"]`, chromedp.BySearch),
		chromedp.WaitVisible(`img-content`, chromedp.ByID),
		chromedp.OuterHTML("html", &body),
	)

	if err != nil {
		log.Error(err)
		return "", err
	}

	return body, nil

}
