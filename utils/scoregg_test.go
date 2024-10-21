package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func TestTeamBattle(t *testing.T) {
	file, err := os.Open("others/bisai.json") // 替换为你的文件名
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	detail := ScoreGGGameDetail{}
	err = json.Unmarshal(bytes, &detail)
	if err != nil {
		fmt.Println("Error unmarshalling:", err)
		return
	}

	//fmt.Println(detail.GetTeamBattleReportsFromDetail())
	//fmt.Println(detail.GetDragonReportsFromDetail())
	//fmt.Println(detail.GetResultFromDetail())
	fmt.Println(detail.GetTeamPlayersReportFromTeamName("AL"))
	//fmt.Println(detail)
}

func TestGames(t *testing.T) {
	//file, err := os.Open("others/games.json") // 替换为你的文件名
	//if err != nil {
	//	fmt.Println("Error opening file:", err)
	//	return
	//}
	//defer file.Close()
	//
	//bytes, err := ioutil.ReadAll(file)
	//if err != nil {
	//	fmt.Println("Error reading file:", err)
	//	return
	//}
	//
	//g := ScoreGGGames{}
	//err = json.Unmarshal(bytes, &g)
	//if err != nil {
	//	fmt.Println("Error unmarshalling:", err)
	//	return
	//}

	//g := NewGames()
	//fmt.Println(g.GetGamesByTournamentId("629"))
	//fmt.Println(detail.GetTeamBattleReportsFromDetail())
	//fmt.Println(detail.GetDragonReportsFromDetail())
	//fmt.Println(detail.GetResultFromDetail())
	//fmt.Println(detail.GetTeamPlayersReportFromTeamName("AL"))
	//fmt.Println(detail)

	//m := NewScoreGGMatches("39095")
	//fmt.Println(m.Data[0].Url)
	//fmt.Println(m.Data[1].Url)

	//ms, key := GetLastLiveTextMessageByMatchIdAndLastKey("39099", "1")
	//fmt.Println(key)
	//
	//for _, m := range ms {
	//	fmt.Println(m)
	//}

	//ms, key = GetLastLiveTextMessageByMatchIdAndLastKey("39098", "170609347372092964")
	//fmt.Println(key)
	//
	//for _, m := range ms {
	//	fmt.Println(m)
	//}

	//str := "lol文字直播:22222"
	//
	//// 创建正则表达式
	//regex := regexp.MustCompile("l.l文字直播.(\\d+)")
	//
	//// 使用正则表达式匹配字符串
	//matches := regex.FindStringSubmatch(str)
	//
	//// 获取匹配项
	////number := matches[0]
	//fmt.Println(matches)

	m := NewScoreGGMatches("39099")
	fmt.Println(m.Data)

}
