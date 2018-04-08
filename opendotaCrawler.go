package dotastats

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

const (
	PRO_MATCHES_API   = "https://api.opendota.com/api/proMatches/"
	MATCH_DETAILS_API = "https://api.opendota.com/api/matches/"
)

type OpenDotaAPIParams struct {
	LessThanMatchID string `json:"less_than_match_id"`
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

type TeamInfoOpenDota struct {
	TeamID  int    `json:"team_id,omitempty"`
	Name    string `json:"name,omitempty"`
	Tag     string `json:"tag,omitempty"`
	LogoURL string `json:"logo_url,omitempty"`
}

type MatchDetailsOD struct {
	PicksBans   []PicksBans      `json:"picks_bans"`
	RadiantTeam TeamInfoOpenDota `json:"radiant_team"`
	DireTeam    TeamInfoOpenDota `json:"dire_team"`
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
		var matchDetails MatchDetailsOD

		matchID := strconv.Itoa(match.MatchID)

		respMatchDetails, err := OpenDotaGet(MATCH_DETAILS_API+matchID, OpenDotaAPIParams{})
		if err != nil {
			fmt.Errorf("error in getting opendota api respMatchDetails: %s", err)
			continue
		}
		defer respMatchDetails.Body.Close()
		err = json.NewDecoder(respMatchDetails.Body).Decode(&matchDetails)
		if err != nil {
			return []OpenDotaMatch{}, fmt.Errorf("error in parsing result from opendota respMatchDetails: %s", err)
		}

		openDotaMatch, err := createOpenDotaMatch(match, matchDetails)
		if err != nil {
			return []OpenDotaMatch{}, fmt.Errorf("error in getting opendota api respMatchDetails: %s", err)
		}
		result = append(result, openDotaMatch)

	}

	fmt.Println("crawling %d matches from opendota, from ID %d to %d", len(result), result[len(result)-1].MatchID, result[0].MatchID)
	return result, nil
}

func createOpenDotaMatch(match ProMatchesOD, matchDetails MatchDetailsOD) (OpenDotaMatch, error) {

	var openDotaMatch OpenDotaMatch
	var err error

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

	openDotaMatch.RadiantTag = matchDetails.RadiantTeam.Tag
	openDotaMatch.RadiantLogoURL = matchDetails.RadiantTeam.LogoURL
	openDotaMatch.DireTag = matchDetails.DireTeam.Tag
	openDotaMatch.DireLogoURL = matchDetails.DireTeam.LogoURL
	openDotaMatch.PicksBans = matchDetails.PicksBans
	return openDotaMatch, nil
}
