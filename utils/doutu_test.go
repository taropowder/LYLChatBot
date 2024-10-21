package utils

import (
	"LYLChatBot/constant"
	log "github.com/sirupsen/logrus"
	"regexp"
	"testing"
)

func TestDoutu(t *testing.T) {
	//GetDouTuImagesByParameters(1, "")
	//fmt.Println(url)

	//content := "表情包 哈哈哈dd"
	//regex := regexp.MustCompile(constant.ImgKey)
	//
	//matches := regex.FindStringSubmatch(content)
	//
	//fmt.Println(matches)

	//re := regexp.MustCompile(constant.PlayerRoleMenu)
	//
	//match := re.FindStringSubmatch("扮演角色:(搅屎棍)")
	//log.Info(match)

	//ts := "「三胖子：KnowledgeID:1\n文章类别: windows安全\n\n文章摘要: 微软长期以来一直声称在Active Directory中，森林是安全边界。然而，研究人员发现，森林不再是安全边界。利用MS-RPRN滥用问题和各种信任场景，一个域中的管理员实际上可以破坏与其共享双向林间信任的域中的资源。这种攻击之所以可能，是因为默认的Active Directory林配置中存在四个主要“特性”。此问题已报告给微软安全响应中心，但微软认为这是一个最好通过v.Next解决的问题，这意味着它可能会在未来版本的Windows中得到修复。缓解这种攻击的建议包括选择性身份验证和禁用跨信任的Kerberos完全委派。」\n- - - - - - - - - - - - - - -\n什么是森林？"
	//ts := "「三胖子：KnowledgeID:1文章类别\n: windows安全。」\n- - - - - - - - - - - - - - -\n什么是森林？"
	//fmt.Println(ParseReplyMsgText(ts))
	ts := "KnowledgeID:1\n文章类别: windows安全\n\n文章摘要: 微软长期以来一直声称在Active Directory中，森林是安全边界。然而，研究人员发现，森林不再是安全边界。利用MS-RPRN滥用问题和各种信任场景，一个域中的管理员实际上可以破坏与其共享双向林间信任的域中的资源。这种攻击之所以可能，是因为默认的Active Directory林配置中存在四个主要“特性”。此问题已报告给微软安全响应中心，但微软认为这是一个最好通过v.Next解决的问题，这意味着它可能会在未来版本的Windows中得到修复。缓解这种攻击的建议包括选择性身份验证和禁用跨信任的Kerberos完全委派。"
	regex := regexp.MustCompile(constant.KnowledgeKey)

	matches := regex.FindStringSubmatch(ts)

	if matches != nil {
		log.Info("123")
	}

	//fmt.Println(matches[1])
	//fmt.Println(matches[2])
	//fmt.Println(matches[2] == "")

	//initConfig()
	//GetImagesByGoogleSearch(12, "cat")

	//req := url.NewRequest()
	//req.Proxies = "http://127.0.0.1:7890"
	//imgUrl := "https://i.natgeofe.com/n/548467d8-c5f1-4551-9f58-6817a8d2c45e/NationalGeographic_2572187_16x9.jpg"
	//r, err := requests.Get(imgUrl, req)
	//if err != nil {
	//	log.Error(err)
	//}
	//fmt.Println(r.Text)
}
