package utils

import (
	"testing"
)

func TestRemoveAt(t *testing.T) {
	//fmt.Println(FormatMessage("@三胖子 qw:今天有什么新闻吗", []string{"qw:"}))
	//fmt.Println(time.Now().Unix())
	//fmt.Println(RemoveAt("@三胖子\\u2005表情包 你好"))
	GetKnowledgeFromResp("文章类别:渗透经验\n文章摘要：一团队在谷歌举办的“LLM bugSWAT”漏洞赏金活动中，发现并利用多个安全漏洞，获得总计 50,000 美元的奖金。这些漏洞包括图片未经授权访问、Google Cloud Console 服务拒绝、Google Workspace 数据泄露等。团队详细介绍了漏洞发现和利用的过程，并分享了他们与谷歌安全团队之间的互动和学习经验。")
}
