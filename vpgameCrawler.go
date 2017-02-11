package dotastats

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

const (
	MatchAPI  = "http://www.vpgame.com/gateway/v1/match/"
	SeriesAPI = "http://www.vpgame.com/gateway/v1/match/schedule"
)

type VPGameAPIParams struct {
	Page     string
	Status   string
	Lang     string
	Limit    string
	Category string
	TID      string
}
type VPGameTournament struct {
	Name     string `json:"name"`
	Category string `json:"category"`
}
type VPGameTeamDetail struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	NameShort string `json:"short_name"`
}
type VPGameOddDetail struct {
	ID      string `json:"id"`
	Item    string `json:"item"`
	Victory string `json:"victory"`
}

type VPGameTeam struct {
	Left  VPGameTeamDetail `json:"left"`
	Right VPGameTeamDetail `json:"right"`
}

type VPGameOdd struct {
	Left  VPGameOddDetail `json:"left"`
	Right VPGameOddDetail `json:"right"`
}
type VPGameSchedule struct {
	LeftTeamID    string `json:"left_team_id"`
	RightTeamID   string `json:"right_team_id"`
	LeftTeamName  string `json:"left_team_name"`
	RightTeamName string `json:"right_team_name"`
}

type VPgameMatch struct {
	Id             string           `json:"id"`
	Round          string           `json:"round"`
	Category       string           `json:"category"`
	ModeName       string           `json:"mode_name"`
	ModeDesc       string           `json:"name"`
	HandicapAmount string           `json:"handicap"`
	HandicapTeam   string           `json:"handicap_team"`
	SeriesID       string           `json:"tournament_schedule_id"`
	GameTime       string           `json:"game_time"`
	LeftTeam       string           `json:"left_team"`
	RightTeam      string           `json:"right_team"`
	LeftTeamScore  string           `json:"left_team_score"`
	RightTeamScore string           `json:"right_team_score"`
	Status         string           `json:"status_name"`
	Schedule       VPGameSchedule   `json:"schedule"`
	Odd            VPGameOdd        `json:"odd"`
	Team           VPGameTeam       `json:"team"`
	Tournament     VPGameTournament `json:"tournament"`
}

type VPgameAPIResult struct {
	Status  int           `json:"status"`
	Message string        `json:"message"`
	Body    []VPgameMatch `json:"body"`
}

func parseUnixTime(unixString string) (time.Time, error) {
	i, err := strconv.ParseInt(unixString, 10, 64)
	if err != nil {
		return time.Time{}, err
	}
	timeParsed := time.Unix(i, 0)
	return timeParsed, nil
}

func processStatus(status string) string {
	status = strings.ToLower(status)
	if status == "live" {
		return "Live"
	} else if status == "settled" {
		return "Settled"
	}
	return "Upcoming"
}
func RunCrawlerVpgame(vpParams VPGameAPIParams) ([]Match, error) {
	var result []Match
	var vpgameResult VPgameAPIResult
	//result, err = processMatchesDota2BY(listMatches)
	resp, err := VPGameGet(MatchAPI, vpParams)
	if err != nil {
		return []Match{}, fmt.Errorf("error in getting vpgame api: %s", err)
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&vpgameResult)
	if err != nil {
		return []Match{}, fmt.Errorf("error in parsing result from vpgame: %s", err)
	}
	for _, match := range vpgameResult.Body {
		var matchFinal Match
		var seriesResult VPgameAPIResult
		matchFinal.TeamAID = match.Team.Left.ID
		matchFinal.TeamBID = match.Team.Right.ID
		matchFinal.TeamA = match.Team.Left.Name
		matchFinal.TeamB = match.Team.Right.Name
		matchFinal.Tournament = match.Tournament.Name
		matchFinal.Game = match.Category
		matchFinal.BestOf = match.Round

		seriesParam := VPGameAPIParams{TID: match.SeriesID}
		respSeries, err := VPGameGet(SeriesAPI, seriesParam)
		if err != nil {
			return []Match{}, fmt.Errorf("error in getting vpgame api respSeries: %s", err)
		}
		defer respSeries.Body.Close()
		err = json.NewDecoder(respSeries.Body).Decode(&seriesResult)
		if err != nil {
			return []Match{}, fmt.Errorf("error in parsing result from vpgame respSeries: %s", err)
		}
		for _, match := range seriesResult.Body {
			subMatch := matchFinal
			subMatch.MatchID = match.Id
			subMatch.URL = "http://www.vpgame.com/match/" + subMatch.MatchID
			subMatch.Time, err = parseUnixTime(match.GameTime)
			if err != nil {
				fmt.Errorf("error in parsing time from vpgame: %s", err)
			}
			subMatch.MatchName = match.LeftTeam + " vs " + match.RightTeam + ", " + match.ModeName
			subMatch.ModeName = match.ModeName
			subMatch.ModeDesc = match.ModeDesc
			subMatch.HandicapAmount = match.HandicapAmount
			if match.HandicapTeam == "left" {
				subMatch.HandicapTeam = match.LeftTeam
			} else if match.HandicapTeam == "right" {
				subMatch.HandicapTeam = match.RightTeam
			}
			subMatch.RatioA = ratioProcess(match.Odd.Left.Item)
			subMatch.RatioB = ratioProcess(match.Odd.Right.Item)
			if match.Odd.Left.Victory == "win" && match.Odd.Right.Victory == "lose" {
				subMatch.Winner = match.LeftTeam
			} else if match.Odd.Left.Victory == "lose" && match.Odd.Right.Victory == "win" {
				subMatch.Winner = match.RightTeam
			}
			subMatch.Status = processStatus(match.Status)
			clientVPGame := &http.Client{
				Timeout: time.Second * 5,
			}

			reqVPGame, err := http.NewRequest("GET", subMatch.URL, nil)
			if err != nil {
				fmt.Errorf("error createrq crawl match vpgame: %s", err)
			}
			respVPGame, err := clientVPGame.Do(reqVPGame)
			if err != nil {
				fmt.Errorf("error crawl match vpgame: %s", err)
				continue
			}
			defer respVPGame.Body.Close()

			doc, err := goquery.NewDocumentFromResponse(respVPGame)
			if err != nil {
				fmt.Errorf("error in crawling from vpgame: %s", err)
			}
			score := doc.Find("div.pic-mid p:nth-child(1)").Text()
			fmt.Println("SCORE DOC", score)
			scoreArray := scoreProcess(score)
			if len(scoreArray) > 0 {
				subMatch.ScoreA = scoreArray[0]
				subMatch.ScoreB = scoreArray[1]
			}
			if match.Round == "" {
				subMatch.BestOf = "BO1"

			}
			subMatch.TeamAShort = match.Team.Left.NameShort
			subMatch.TeamBShort = match.Team.Right.NameShort
			result = append(result, subMatch)

		}
	}
	fmt.Println("%v", len(result))
	return result, err
}
