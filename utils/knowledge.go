package utils

import (
	"LYLChatBot/model"
	"LYLChatBot/pkg/database"
	log "github.com/sirupsen/logrus"
	"regexp"
	"strings"
)

func GetKnowledgeFromResp(resp string) (t, a string) {
	//**文章类别:** /windows安全
	//
	//**文章摘要:**
	//这篇文章探讨了打破Active Directory中森林信任的安全边界的问题。攻击者可以通过利用默认的配置和一个名为“printer bug”的漏洞，通过一个受损的域控制器或拥有不受约束委托的服务器，破坏另一个受信任的森林中的资源。该漏洞允许攻击者在不安全的信任边界的情况下，访问和泄露敏感的凭据材料。微软将此修复问题划分为下一版本的 Windows，但提供了缓解措施，如选择性身份验证和禁用信任域之间的 Kerberos 完全委派。
	regex := regexp.MustCompile(`文章类别\**\s*[:：]\**\s*(\S+)\s*\n\s*文章摘要\**\s*[:：]\**\s*(.+)`)
	matches := regex.FindStringSubmatch(resp)
	if matches != nil {
		t = strings.Trim(matches[1], "/")
		a = matches[2]
	}
	return
}

func GetAllKnowledgeTypes() string {
	db := database.GetDB()
	ts := make([]model.KnowledgeTypeRecord, 0)
	err := db.Find(&ts).Error
	if err != nil {
		return ""
	} else {
		res := ""
		for _, t := range ts {
			res = t.TypeName + "/" + res
		}
		return res
	}
}

func GetHistoryByKnowledgeId(kid int) []GptMessage {
	db := database.GetDB()
	qas := make([]model.KnowledgeQARecord, 0)
	msgs := make([]GptMessage, 0)

	err := db.Where("knowledge_id = ? and answer!='' ", kid).Find(&qas).Error
	if err != nil {
		log.Error(err)
		return nil
	} else {
		for _, t := range qas {
			msgs = append(msgs, GptMessage{
				Role:    "user",
				Content: t.Answer,
			})
			msgs = append(msgs, GptMessage{
				Role:    "assistant",
				Content: t.Answer,
			})
		}
	}
	return msgs
}
