package dotastats

import (
	"encoding/json"
	"fmt"
	"time"
)

const (
	PRO_MATCHES_API = "https://api.opendota.com/api/proMatches"
)

type OpenDotaAPIParams struct {
	LessThanMatchID string
}

type ProMatchesOD struct {
	MatchID       int    `json:"match_id"`
	Duration      int    `json:"duration"`
	StartTime     int    `json:"start_time"`
	RadiantTeamID int    `json:"radiant_team_id"`
	RadiantName   string `json:"radiant_name"`
	DireTeamID    int    `json:"dire_team_id"`
	DireName      string `json:"dire_name"`
	LeagueID      int    `json:"leagueid"`
	LeagueName    string `json:"league_name"`
	SeriesID      int    `json:"series_id"`
	SeriesType    int    `json:"series_type"`
	RadiantScore  int    `json:"radiant_score"`
	DireScore     int    `json:"dire_score"`
	RadiantWin    bool   `json:"radiant_win"`
}

type ProMatchesAPIResult struct {
	Body []ProMatchesOD
}

func parseUnixTimeInt(unixInt int64) (*time.Time, error) {
	timeParsed := time.Unix(unixInt, 0)
	return &timeParsed, nil
}

func RunCrawlerOpenDota(openDotaAPIParams OpenDotaAPIParams) ([]OpenDotaMatch, error) {
	var result []OpenDotaMatch
	var proMatchesResult ProMatchesAPIResult
	resp, err := OpenDotaGet(PRO_MATCHES_API, openDotaAPIParams)
	if err != nil {
		return []OpenDotaMatch{}, fmt.Errorf("error in getting opendota api: %s", err)
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&proMatchesResult.Body)
	if err != nil {
		return []OpenDotaMatch{}, fmt.Errorf("error in parsing result from opendota: %s", err)
	}
	for _, match := range proMatchesResult.Body {
		var openDotaMatch OpenDotaMatch
		openDotaMatch.MatchID = match.MatchID
		openDotaMatch.Duration = match.Duration
		openDotaMatch.RadiantTeamID = match.RadiantTeamID
		openDotaMatch.RadiantName = match.RadiantName
		openDotaMatch.DireTeamID = match.DireTeamID
		openDotaMatch.DireName = match.DireName
		openDotaMatch.LeagueID = match.LeagueID
		openDotaMatch.LeagueName = match.LeagueName
		openDotaMatch.SeriesID = match.SeriesID
		openDotaMatch.SeriesType = match.SeriesType
		openDotaMatch.RadiantScore = match.RadiantScore
		openDotaMatch.DireScore = match.DireScore
		openDotaMatch.RadiantWin = match.RadiantWin

		openDotaMatch.StartTime, err = parseUnixTimeInt(int64(match.StartTime))
		if err != nil {
			fmt.Errorf("error in parsing time from opendota: %s", err)
		}
		result = append(result, openDotaMatch)

	}

	fmt.Println("%v", len(result))
	return result, nil
}
