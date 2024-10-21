package utils

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type ScorePlayer struct {
	Sort        string `json:"sort"`
	WardsKilled int    `json:"wardsKilled"`
	AssisNum    int    `json:"assis_num"`
	PlayerId    int    `json:"player_id"`
	HeroID      int    `json:"heroID"`
	Spells      []int  `json:"spells"`
	HeroName    string `json:"hero_name"`
	PlayerName  string `json:"player_name"`
	LasthitNum  int    `json:"lasthit_num"`
	WardsPlaced int    `json:"wardsPlaced"`
	KillNum     int    `json:"kill_num"`
	DeadNum     int    `json:"dead_num"`
	ChampLevel  int    `json:"champLevel"`
	HeroId      int    `json:"hero_id"`
	PositionID  int    `json:"positionID"`
	Dl          []struct {
		Image string `json:"image"`
	} `json:"dl"`
	HeroDamage int   `json:"hero_damage"`
	HeroWound  int   `json:"hero_wound"`
	Economics  int   `json:"economics"`
	Devices    []int `json:"devices"`
	Perks      struct {
	} `json:"perks"`
	EconomicsRate float64 `json:"economics_rate"`
	HeroImage     string  `json:"hero_image"`
	PlayerImage   string  `json:"player_image"`
}

type ScoreGGTeam struct {
	WardsKilled    int    `json:"wardsKilled"`
	First10Kill    int    `json:"first10Kill"`
	Deaths         int    `json:"deaths"`
	First5Kill     int    `json:"first5Kill"`
	TeamImageThumb string `json:"team_image_thumb"`
	FirstTowerKill int    `json:"firstTowerKill"`
	WardsPlaced    int    `json:"wardsPlaced"`
	ClanId         int    `json:"clan_id"`
	Damages        string `json:"damages"`
	FirstBloodKill int    `json:"firstBloodKill"`
	Kills          int    `json:"kills"`
	TeamImage      string `json:"team_image"`
	Towers         int    `json:"towers"`
	TeamShortName  string `json:"team_short_name"`
	Assists        int    `json:"assists"`
	Dragons        []struct {
		DragonImage string `json:"dragon_image"`
		GroupId     int    `json:"group_id"`
		GameTimeTxt string `json:"game_time_txt"`
		DragonType  int    `json:"dragon_type"`
		GameTime    int    `json:"game_time"`
		Type        string `json:"type"`
		Id          int    `json:"id"`
		GroupId1    int    `json:"groupId"`
	} `json:"dragons"`
	FirstHerald     int           `json:"firstHerald"`
	GroupId         int           `json:"groupId"`
	Economics       string        `json:"economics"`
	Players         []ScorePlayer `json:"players"`
	FirstDragonKill int           `json:"firstDragonKill"`
	GroupId1        int           `json:"group_id"`
	FirstBaronKill  int           `json:"firstBaronKill"`
	Position        string        `json:"position"`
}

type TeamBattleData struct {
	GroupId       int `json:"group_id"`
	WoundAll      int `json:"wound_all"`
	DamageAll     int `json:"damage_all"`
	HeroID        int `json:"heroID"`
	TeamBattleSeq int `json:"team_battle_seq"`
	GroupID       int `json:"groupId"`
	HeroId        int `json:"hero_id"`
}

type ScoreGGGameDetail struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
	Data    struct {
		DragonEvent []struct {
			DragonType  int    `json:"dragon_type"`
			GameTime    int    `json:"game_time"`
			Id          int    `json:"id"`
			GameTimeTxt string `json:"game_time_txt"`
			GroupId     int    `json:"group_id"`
			Type        string `json:"type"`
			GroupID     int    `json:"groupId"`
		} `json:"dragon_event"`
		TeamBattle []struct {
			Seq          int              `json:"seq"`
			EndTimeTxt   string           `json:"end_time_txt"`
			BeginTimeTxt string           `json:"begin_time_txt"`
			EndTime      int              `json:"end_time"`
			BeginTime    int              `json:"begin_time"`
			Data         []TeamBattleData `json:"data"`
		} `json:"team_battle"`
		EcoList   []int  `json:"eco_list"`
		ExpInfo   []int  `json:"exp_info"`
		UpdatedAt string `json:"updated_at"`
		Game      struct {
			Status       int    `json:"status"`
			BlueClanId   int    `json:"blue_clan_id"`
			RedClanName  string `json:"red_clan_name"`
			BlueClanName string `json:"blue_clan_name"`
			MatchId      int    `json:"match_id"`
			TmBattleId   int    `json:"tm_battle_id"`
			RedClanId    int    `json:"red_clan_id"`
			NumberTxt    string `json:"number_txt"`
			BlueTeamID   int    `json:"blue_teamID"`
			StartTime    int    `json:"start_time"`
			GameTime     int    `json:"game_time"`
			WinClanId    int    `json:"win_clan_id"`
			Number       int    `json:"number"`
			WinTeamID    int    `json:"win_teamID"`
			Source       string `json:"source"`
			EndTime      int    `json:"end_time"`
			GameTimeTxt  string `json:"game_time_txt"`
			SourceId     int    `json:"source_id"`
			RedTeamID    int    `json:"red_teamID"`
			BattleId     int    `json:"battle_id"`
		} `json:"game"`
		KillPos []struct {
			AssistantCount  int `json:"assistant_count"`
			KillerHeroID    int `json:"killer_heroID"`
			AssistantIdList []struct {
				HeroId int `json:"hero_id"`
			} `json:"assistant_id_list"`
			GameTime            int `json:"game_time"`
			DeadHeroID          int `json:"dead_heroID"`
			KillerId            int `json:"killer_id"`
			DeadId              int `json:"dead_id"`
			AxisX               int `json:"axis_x"`
			AxisY               int `json:"axis_y"`
			GroupId             int `json:"groupId"`
			AssistantHeroIDList []struct {
				HeroID int `json:"heroID"`
			} `json:"assistant_heroID_list"`
			GroupId1    int    `json:"group_id"`
			Type        string `json:"type"`
			Id          int    `json:"id"`
			GameTimeTxt string `json:"game_time_txt"`
		} `json:"kill_pos"`
		TeamA  ScoreGGTeam `json:"team_a"`
		TeamB  ScoreGGTeam `json:"team_b"`
		EyePos []struct {
			DeadType  int `json:"dead_type"`
			EyeType   int `json:"eye_type"`
			StartTime int `json:"start_time"`
			KillerId  int `json:"killer_id"`
			EndTime   int `json:"end_time"`
			AxisX     int `json:"axis_x"`
			AxisY     int `json:"axis_y"`
			GroupId   int `json:"group_id"`
			GroupId1  int `json:"groupId,omitempty"`
			HeroId    int `json:"hero_id"`
		} `json:"eye_pos"`
	} `json:"data"`
}

func (d *ScoreGGGameDetail) getTeamByGroupId(groupId int) ScoreGGTeam {
	if groupId == 100 {
		d.Data.TeamA.Position = "蓝色方"
		return d.Data.TeamA
	} else {
		d.Data.TeamB.Position = "红色方"
		return d.Data.TeamB
	}
}
func (d *ScoreGGGameDetail) getTeamByTeamName(name string) ScoreGGTeam {
	if name == d.Data.TeamA.TeamShortName {
		d.Data.TeamA.Position = "蓝色方"
		return d.Data.TeamA
	} else {
		d.Data.TeamB.Position = "红色方"
		return d.Data.TeamB
	}
}

func (d *ScoreGGGameDetail) GetDragonReportsFromDetail() []string {
	res := make([]string, 0)
	events := d.Data.DragonEvent
	for _, event := range events {
		t := d.getTeamByGroupId(event.GroupID)
		dragonEventStr := fmt.Sprintf("[%s]  %s%s 击杀了 %s",
			event.GameTimeTxt, t.Position, t.TeamShortName, d.getDragonNameById(event.DragonType))
		res = append(res, dragonEventStr)
	}
	return res
}

func (d *ScoreGGGameDetail) getDragonNameById(dragonId int) string {

	switch dragonId {
	case 1:
		return "火龙"
	}

	return "呆呆龙"
}

func (d *ScoreGGGameDetail) GetTeamBattleReportsFromDetail() []string {
	res := make([]string, 0)
	events := d.Data.TeamBattle
	for _, event := range events {

		dragonEventStr := fmt.Sprintf("[%s-%s]  团战对比: \n",
			event.BeginTimeTxt, event.EndTimeTxt)

		APlayers := make(map[TeamBattleData]ScorePlayer, 0)
		BPlayers := make(map[TeamBattleData]ScorePlayer, 0)

		for _, hero := range event.Data {

			p := d.GetPlayerFromHeroId(hero.HeroId)
			t := d.getTeamByGroupId(hero.GroupID)

			if t.Position == "蓝色方" {
				APlayers[hero] = p
			}

			if t.Position == "红色方" {
				BPlayers[hero] = p
			}
		}

		for i := 1; i <= 5; i++ {
			aBattleData, aPlayer := d.getPlayerFromTeamBattleDataSort(APlayers, strconv.Itoa(i))
			bBattleData, bPlayer := d.getPlayerFromTeamBattleDataSort(BPlayers, strconv.Itoa(i))

			dragonEventStr = dragonEventStr + fmt.Sprintf("%s(%s) 伤害%d 承伤%d--%s(%s) 伤害: %d 承伤%d\n",
				aPlayer.HeroName, aPlayer.PlayerName, aBattleData.DamageAll, aBattleData.WoundAll,
				bPlayer.HeroName, bPlayer.PlayerName, bBattleData.DamageAll, bBattleData.WoundAll,
			)
		}

		//t := d.getTeamByGroupId(event.)

		res = append(res, dragonEventStr)
	}
	return res
}

func (d *ScoreGGGameDetail) GetPlayerFromHeroId(heroID int) ScorePlayer {

	p := ScorePlayer{}

	for _, player := range d.Data.TeamA.Players {
		if player.HeroId == heroID {
			return player
		}
	}

	for _, player := range d.Data.TeamB.Players {
		if player.HeroId == heroID {
			return player
		}
	}

	return p
}

func (d *ScoreGGGameDetail) getPlayerFromTeamBattleDataSort(players map[TeamBattleData]ScorePlayer, sort string) (TeamBattleData, ScorePlayer) {

	for b, player := range players {
		if player.Sort == sort {
			return b, player
		}
	}
	return TeamBattleData{}, ScorePlayer{}
}

func (d *ScoreGGGameDetail) GetResultFromDetail() string {
	g := d.Data.Game
	team := map[int]string{
		g.BlueClanId: g.BlueClanName,
		g.RedClanId:  g.RedClanName,
	}

	winTeamName := team[g.WinClanId]
	teamA := d.getTeamByTeamName(g.BlueClanName)
	teamB := d.getTeamByTeamName(g.RedClanName)
	res := fmt.Sprintf("%s VS %s Round %d 比赛时长 %s  获胜方 %s\n"+
		"经济: %v: %v\n"+
		"伤害: %v: %v\n"+
		"人头: %v:%v\n"+
		"推塔: %v:%v\n"+
		"小龙: %v:%v",
		g.BlueClanName, g.RedClanName, g.Number, g.GameTimeTxt, winTeamName,
		teamA.Economics, teamB.Economics,
		teamA.Damages, teamB.Damages,
		teamA.Kills, teamB.Kills,
		teamA.Towers, teamB.Towers,
		len(teamA.Dragons), len(teamB.Dragons),
	)

	return res
}

func (d *ScoreGGGameDetail) GetTeamPlayersReportFromTeamName(name string) string {
	res := fmt.Sprintf("%s战队数据:\n", name)

	team := d.getTeamByTeamName(name)

	for _, player := range team.Players {
		res = res + fmt.Sprintf("%s(%s) %d/%d/%d  经济:%v 伤害%d\n",
			player.HeroName, player.PlayerName, player.KillNum, player.DeadNum, player.AssisNum,
			player.Economics, player.HeroDamage,
		)
	}

	return res
}

func NewGameDetailFromUrl(url string) ScoreGGGameDetail {
	d := ScoreGGGameDetail{}
	// 发起GET请求
	response, err := http.Get(url)
	if err != nil {
		logrus.Error("Error making GET request:", err)
	}
	defer response.Body.Close()

	// 读取响应的内容
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		logrus.Error(err)
	}

	err = json.Unmarshal(body, &d)
	if err != nil {
		logrus.Error(err)
	}

	return d
}

// games

type ScoreGGTournamentinfo struct {
	TournamentID        string `json:"tournamentID"`
	Name                string `json:"name"`
	EnName              string `json:"en_name"`
	ShortName           string `json:"short_name"`
	Image               string `json:"image"`
	ImageThumb          string `json:"image_thumb"`
	SsdbV               string `json:"ssdb_v"`
	NameEn              string `json:"name_en"`
	NameTw              string `json:"name_tw"`
	ShortNameEn         string `json:"short_name_en"`
	ShortNameTw         string `json:"short_name_tw"`
	NavigateNewsName    string `json:"navigate_news_name"`
	NavigateGeneralName string `json:"navigate_general_name"`
	NavigateMatchName   string `json:"navigate_match_name"`
	NavigateTeamName    string `json:"navigate_team_name"`
	NavigateDataName    string `json:"navigate_data_name"`
	TeamText            string `json:"team_text"`
}

type ScoreGGDescription struct {
	MatchId        string `json:"match_id"`
	TournamentID   string `json:"tournamentID"`
	TournamentName string `json:"tournament_name"`
	RoundName      string `json:"round_name"`
	TeamIDA        string `json:"teamID_a"`
	TeamIDB        string `json:"teamID_b"`
	TeamAWin       string `json:"team_a_win"`
	TeamBWin       string `json:"team_b_win"`
	StartDate      string `json:"start_date"`
	StartTime      string `json:"start_time"`
	Status         string `json:"status"`
	TeamAImage     string `json:"team_a_image"`
	TeamBImage     string `json:"team_b_image"`
	TeamAShortName string `json:"team_a_short_name"`
	TeamBShortName string `json:"team_b_short_name"`
	HomesiteA      string `json:"homesite_a"`
	HomesiteB      string `json:"homesite_b"`
	Homesite       string `json:"homesite"`
	Time           string `json:"time"`
	GameCount      string `json:"game_count"`
	LiveVideoUrl1  string `json:"live_video_url1"`
	LiveVideoName1 string `json:"live_video_name1"`
	LiveVideoUrl2  string `json:"live_video_url2"`
	LiveVideoName2 string `json:"live_video_name2"`
	LiveVideoUrl3  string `json:"live_video_url3"`
	LiveVideoName3 string `json:"live_video_name3"`
	SiteId         string `json:"site_id"`
	DynamicId      string `json:"dynamic_id"`
	//Video          []struct {
	//Title     string `json:"title"`
	//VideoLink string `json:"video_link"`
	//ImageUrl     string `json:"image_url"`
	//ThumbnailPic string `json:"thumbnail_pic"`
	//TitleEn      string `json:"title_en"`
	//TitleTw      string `json:"title_tw"`
	//} `json:"video"`
	HasRealTime int `json:"has_real_time"`
	Result      []struct {
		ResultID      string `json:"resultID"`
		WinTeamID     string `json:"win_teamID"`
		TeamName      string `json:"team_name"`
		TeamShortName string `json:"team_short_name"`
		TeamImage     string `json:"team_image"`
		Bo            string `json:"bo"`
		LivedataKey   string `json:"livedata_key"`
	} `json:"result"`
	NearTen struct {
		TeamAWinCount int `json:"team_a_win_count"`
		TeamBWinCount int `json:"team_b_win_count"`
		List          []struct {
			WinTeamId   string `json:"win_team_id"`
			WinTeamName string `json:"win_team_name"`
			StartTime   string `json:"start_time"`
			MatchID     string `json:"matchID"`
		} `json:"list"`
	} `json:"near_ten"`
	//Remind string `json:"remind"`
}

type ScoreGGDayGames struct {
	Info map[string]struct {
		Tournamentinfo ScoreGGTournamentinfo `json:"tournamentinfo"`
		List           []ScoreGGDescription  `json:"list"`
	} `json:"info"`
	//Today bool `json:"today"`
}

type ScoreGGGames struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Data    struct {
		List         map[string]ScoreGGDayGames `json:"list"`
		ErrorMessage string                     `json:"error_message"`
	} `json:"data"`
	TaskData struct {
	} `json:"task_data"`
	Badge []interface{} `json:"badge"`
	Event []interface{} `json:"event"`
}

func (g *ScoreGGGames) GetGamesByTournamentId(id string) string {

	res := ""
	for _, tournaments := range g.Data.List {
		for tournamentId, tournament := range tournaments.Info {
			if tournamentId == id {
				for _, game := range tournament.List {
					gameDescription := fmt.Sprintf("%s %s  %s VS %s (%s)\n", game.StartDate, game.StartTime, game.TeamAShortName, game.TeamBShortName, game.MatchId)
					res = res + gameDescription
				}
			}

		}

	}
	return res
}

func NewGames() ScoreGGGames {
	g := ScoreGGGames{}

	apiUrl := "https://www.scoregg.com/services/api_url.php"
	payload := url.Values{
		"api_path":      {"services/match/web_math_list.php"},
		"gameID":        {"1"},
		"date":          {""},
		"tournament_id": {""},
		"api_version":   {"9.9.9"},
		"platform":      {"web"},
	}

	reqBody := strings.NewReader(payload.Encode())

	req, err := http.NewRequest("POST", apiUrl, reqBody)
	if err != nil {
		logrus.Error("Error creating request:", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("Referer", "https://www.scoregg.com/schedule")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")

	client := &http.Client{Transport: &http.Transport{
		TLSHandshakeTimeout: 10 * time.Second, // 设置超时时间为10秒
	}}
	resp, err := client.Do(req)
	if err != nil {
		logrus.Error("Error sending request:", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Error("Error reading response body:", err)
	}

	err = json.Unmarshal(body, &g)
	if err != nil {
		logrus.Error(err)
	}

	return g
}

// match

type ScoreGGMatches struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
	Data    []struct {
		Status       int    `json:"status"`
		RedClanName  string `json:"red_clan_name"`
		RedClanId    int    `json:"red_clan_id"`
		WinClanName  string `json:"win_clan_name"`
		WebUrl       string `json:"web_url"`
		BlueClanName string `json:"blue_clan_name"`
		BlueClanId   int    `json:"blue_clan_id"`
		WinTeamID    int    `json:"win_teamID"`
		Url          string `json:"url"`
		BlueTeamID   int    `json:"blue_teamID"`
		Source       string `json:"source"`
		RedTeamID    int    `json:"red_teamID"`
		WinClanColor string `json:"win_clan_color"`
		WinClanId    int    `json:"win_clan_id"`
	} `json:"data"`
}

func NewScoreGGMatches(matchId string) (m ScoreGGMatches) {

	apiUrl := fmt.Sprintf("https://img.scoregg.com/lol/livedata/%s.json", matchId)

	req, err := http.NewRequest("GET", apiUrl, nil)
	if err != nil {
		logrus.Error("Error creating request:", err)
		return
	}

	req.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Origin", "https://www.scoregg.com")
	req.Header.Set("Referer", "https://www.scoregg.com/")
	req.Header.Set("Sec-Ch-Ua", `"Not_A Brand";v="8", "Chromium";v="120", "Google Chrome";v="120"`)
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua-Platform", `"macOS"`)
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-site")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logrus.Error("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Error("Error reading response body:", err)
		return
	}

	m = ScoreGGMatches{}

	err = json.Unmarshal(body, &m)
	if err != nil {
		logrus.Error(err)
		return
	}

	return m

}

// live txt

type ScoreGGLiveTextMatch struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Data    struct {
		SsdbV                 string `json:"ssdb_v"`
		MatchID               string `json:"matchID"`
		TournamentID          string `json:"tournamentID"`
		TeamIDA               string `json:"teamID_a"`
		TeamIDB               string `json:"teamID_b"`
		StartTime             string `json:"start_time"`
		RoundName             string `json:"round_name"`
		TournamentName        string `json:"tournament_name"`
		TournamentShortName   string `json:"tournament_short_name"`
		RoundID               string `json:"roundID"`
		Status                string `json:"status"`
		Prospect              string `json:"prospect"`
		Situation             string `json:"situation"`
		TeamAName             string `json:"team_a_name"`
		TeamBName             string `json:"team_b_name"`
		TeamAImage            string `json:"team_a_image"`
		TeamBImage            string `json:"team_b_image"`
		TeamACountryImage     string `json:"team_a_country_image"`
		TeamBCountryImage     string `json:"team_b_country_image"`
		TeamAImageThumb       string `json:"team_a_image_thumb"`
		TeamBImageThumb       string `json:"team_b_image_thumb"`
		TeamAShortName        string `json:"team_a_short_name"`
		TeamBShortName        string `json:"team_b_short_name"`
		TeamAEnName           string `json:"team_a_en_name"`
		TeamBEnName           string `json:"team_b_en_name"`
		TeamAWin              string `json:"team_a_win"`
		TeamBWin              string `json:"team_b_win"`
		TeamAStarId           string `json:"team_a_star_id"`
		TeamBStarId           string `json:"team_b_star_id"`
		GameCount             string `json:"game_count"`
		StartTimeString       string `json:"start_time_string"`
		DynamicId             string `json:"dynamic_id"`
		DynamicMemberId       string `json:"dynamic_member_id"`
		LiveStatus            string `json:"live_status"`
		LiveUrl               string `json:"live_url"`
		LiveName              string `json:"live_name"`
		LiveUrl2              string `json:"live_url2"`
		LiveName2             string `json:"live_name2"`
		LiveUrl3              string `json:"live_url3"`
		LiveName3             string `json:"live_name3"`
		LiveSourceUrl         string `json:"live_source_url"`
		IsHaveVideoLink       string `json:"is_have_video_link"`
		ResultCount           string `json:"result_count"`
		TournamentImageThumb  string `json:"tournament_image_thumb"`
		DefaultBetId          string `json:"default_bet_id"`
		TournamentDescription string `json:"tournament_description"`
		ExternalUrl           string `json:"external_url"`
		RoundSonId            string `json:"round_son_id"`
		RoundSonName          string `json:"round_son_name"`
		RoundList             []struct {
			RoundID string `json:"roundID"`
			Name    string `json:"name"`
			NameEn  string `json:"name_en"`
			NameTw  string `json:"name_tw"`
		} `json:"round_list"`
		LiveVideoUrl1  string `json:"live_video_url1"`
		LiveVideoName1 string `json:"live_video_name1"`
		LiveVideoUrl2  string `json:"live_video_url2"`
		LiveVideoName2 string `json:"live_video_name2"`
		LiveVideoUrl3  string `json:"live_video_url3"`
		LiveVideoName3 string `json:"live_video_name3"`
		GameID         string `json:"gameID"`
		SiteId         string `json:"site_id"`
		Homesite       string `json:"homesite"`
		HomesiteA      string `json:"homesite_a"`
		HomesiteB      string `json:"homesite_b"`
		CircleId       string `json:"circle_id"`
		Picture        struct {
		} `json:"picture"`
		BuyTicketList      []interface{} `json:"buy_ticket_list"`
		VideoList          []interface{} `json:"video_list"`
		TournamentTourTag  string        `json:"tournament_tour_tag"`
		FirstRoundID       string        `json:"first_roundID"`
		SplendidVideoList  []interface{} `json:"splendid_video_list"`
		EndTime            string        `json:"end_time"`
		TeamABindID3       string        `json:"team_a_bindID_3"`
		TeamBBindID3       string        `json:"team_b_bindID_3"`
		ResultList         []interface{} `json:"result_list"`
		Title              string        `json:"title"`
		MoreTeamWinList    []interface{} `json:"more_team_win_list"`
		MoreTeamList       []interface{} `json:"more_team_list"`
		WinTeamID          string        `json:"win_teamID"`
		IsReallyEnd        string        `json:"is_really_end"`
		TopicCount         int           `json:"topic_count"`
		ChatCommentCount   string        `json:"chat_comment_count"`
		MatchForecastCount string        `json:"match_forecast_count"`
		MatchCommentCount  string        `json:"match_comment_count"`
		RankingA           string        `json:"ranking_a"`
		RankingB           string        `json:"ranking_b"`
		MaxData            struct {
			MINUTEECONOMIC       float64 `json:"MINUTE_ECONOMIC"`
			MINUTEWARDSPLACED    string  `json:"MINUTE_WARDSPLACED"`
			AVERAGEMinionsKilled string  `json:"AVERAGE_MinionsKilled"`
			AVERAGELife          string  `json:"AVERAGE_Life"`
			AVERAGEKILLS         string  `json:"AVERAGE_KILLS"`
			AVERAGEASSISTS       string  `json:"AVERAGE_ASSISTS"`
		} `json:"max_data"`
		ForecastA          string `json:"forecast_a"`
		ForecastB          string `json:"forecast_b"`
		TotalForecastCount string `json:"total_forecast_count"`
		ForecastAP         string `json:"forecast_a_p"`
		ForecastBP         string `json:"forecast_b_p"`
		BetList            []struct {
			BetId        string `json:"bet_id"`
			Title        string `json:"title"`
			CategoryName string `json:"category_name"`
			Items        []struct {
				BetItemId    string `json:"bet_item_id"`
				InitPrice    string `json:"init_price"`
				Price        string `json:"price"`
				MemberMaxBet string `json:"member_max_bet"`
				ItemName     string `json:"item_name"`
				ItemNameEn   string `json:"item_name_en"`
				ItemNameTw   string `json:"item_name_tw"`
				WinRate      string `json:"win_rate"`
				Odds         string `json:"odds"`
				IsBetting    string `json:"is_betting"`
				IsBet        string `json:"is_bet"`
			} `json:"items"`
			BetEndTime    string `json:"bet_end_time"`
			BetEndTimeTxt string `json:"bet_end_time_txt"`
			DateTxt       string `json:"date_txt"`
			//Status          int    `json:"status"`
			TotalPrice      string `json:"total_price"`
			PeopleNum       string `json:"people_num"`
			ResultItemId    string `json:"result_item_id"`
			ViewType        int    `json:"view_type"`
			Image           string `json:"image"`
			Name            string `json:"name"`
			MatchTeamA      string `json:"match_team_a"`
			MatchTeamB      string `json:"match_team_b"`
			TeamImageThumbA string `json:"team_image_thumb_a"`
			TeamImageThumbB string `json:"team_image_thumb_b"`
			MatchDate       string `json:"match_date"`
			MatchStartTime  string `json:"match_start_time"`
			MatchId         string `json:"match_id"`
			TeamAWin        string `json:"team_a_win"`
			TeamBWin        string `json:"team_b_win"`
			MatchStatus     string `json:"match_status"`
			MatchBetCount   string `json:"match_bet_count"`
			GameCount       string `json:"game_count"`
		} `json:"bet_list"`
		ScoreCount   string `json:"score_count"`
		ScoreCount1  string `json:"score_count_1"`
		ScoreCount2  string `json:"score_count_2"`
		ScoreCount3  string `json:"score_count_3"`
		ScoreCount4  string `json:"score_count_4"`
		ScoreCount5  string `json:"score_count_5"`
		ScoreAverage string `json:"score_average"`
		IsRemind     string `json:"is_remind"`
		Forecast     struct {
		} `json:"forecast"`
		Comment struct {
		} `json:"comment"`
		Member []interface{} `json:"member"`
	} `json:"data"`
	TaskData struct {
	} `json:"task_data"`
	Badge []interface{} `json:"badge"`
	Event []interface{} `json:"event"`
}

type ScoreGGLiveText struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Data    struct {
		List []struct {
			Content string `json:"content"`
			//ReplyId    int    `json:"reply_id"`
			MemberId   int `json:"member_id"`
			Timestamp  int `json:"timestamp"`
			Attachment []struct {
				Type   string `json:"type"`
				TeamId string `json:"team_id"`
			} `json:"attachment"`
			Picture []struct {
				ImageUrl  string `json:"image_url"`
				Width     int    `json:"width"`
				Height    int    `json:"height"`
				ImageUrl5 string `json:"image_url5"`
			} `json:"picture"`
			OtherCacheKey string `json:"other_cache_key"`
			Avatar        string `json:"avatar"`
			Nickname      string `json:"nickname"`
			MemberType    int    `json:"member_type"`
			Rank          string `json:"rank"`
			BadgeImage    string `json:"badge_image"`
			BadgeTitle    string `json:"badge_title"`
			IsGirl        string `json:"is_girl"`
			MemberStatus  string `json:"member_status"`
			PageKey       string `json:"page_key"`
			PraiseCount   string `json:"praise_count"`
			IsPraise      int    `json:"is_praise"`
		} `json:"list"`
		LastRefreshTime string        `json:"last_refresh_time"`
		GiftData        []interface{} `json:"gift_data"`
		GiftLastPageKey string        `json:"gift_last_page_key"`
		Events          []interface{} `json:"events"`
		Scoreboard      struct {
		} `json:"scoreboard"`
	} `json:"data"`
	TaskData struct {
	} `json:"task_data"`
	Badge []interface{} `json:"badge"`
	Event []interface{} `json:"event"`
}

func GetLastLiveTextMessageByMatchIdAndLastKey(matchId string, lastKey string) ([]string, string) {

	messages := make([]string, 0)
	last_key := ""
	payload := url.Values{
		"api_path":    {"/services/match/match_info_new.php"},
		"method":      {"post"},
		"platform":    {"web"},
		"api_version": {"9.9.9"},
		"language_id": {"1"},
		"matchID":     {matchId},
	}

	if respBody, err := postScoreGGApiUrl(payload); err == nil {
		m := ScoreGGLiveTextMatch{}

		err = json.Unmarshal(respBody, &m)
		if err != nil {
			logrus.Error(err)
			return messages, "error"
		}

		dynamicId := m.Data.DynamicId

		orderSwitch := "0"
		//if lastKey == "" {
		//	orderSwitch = "1"
		//}

		payload := url.Values{
			"api_path":      {"/services/dynamic/dynamic_text_live.php"},
			"method":        {"post"},
			"platform":      {"web"},
			"api_version":   {"9.9.9"},
			"language_id":   {"1"},
			"dynamic_id":    {dynamicId},
			"limit":         {"3"},
			"landlord_only": {"0"},
			"order_switch":  {orderSwitch},
			"last_page_key": {lastKey},
		}
		if respBody, err = postScoreGGApiUrl(payload); err == nil {
			texts := ScoreGGLiveText{}
			err = json.Unmarshal(respBody, &texts)
			if err != nil {
				logrus.Error(err)
				return messages, "error"
			}

			for _, s := range texts.Data.List {
				messages = append(messages, fmt.Sprintf("[%s]: %s",
					time.Unix(int64(s.Timestamp), 0).Format("2006-01-02 15:04:05"), s.Content))
				for _, picture := range s.Picture {
					messages = append(messages, picture.ImageUrl)
				}
				last_key = s.PageKey
			}

		}
	}

	return messages, last_key

}

func postScoreGGApiUrl(payload url.Values) ([]byte, error) {
	apiUrl := "https://www.scoregg.com/services/api_url.php"

	reqBody := strings.NewReader(payload.Encode())

	req, err := http.NewRequest("POST", apiUrl, reqBody)
	if err != nil {
		logrus.Error("Error creating request:", err)
		return nil, err

	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("Referer", "https://www.scoregg.com/schedule")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")

	client := &http.Client{Transport: &http.Transport{
		TLSHandshakeTimeout: 10 * time.Second, // 设置超时时间为10秒
	}}
	resp, err := client.Do(req)
	if err != nil {
		logrus.Error("Error sending request:", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Error("Error reading response body:", err)
		return nil, err

	}
	return body, nil

}
