package dotastats

import "regexp"

func GetTeamFBMatches(teamName string, apiParams APIParams, mongodb Mongodb) (FBResult, error) {
	data, err := mongodb.GetTeamFBMatches(teamName, apiParams)
	if err != nil {
		return FBResult{}, err
	}
	if len(data) == 0 {
		return FBResult{}, nil
	}

	var result FBResult
	var win, winrate float64
	for _, match := range data {
		var winnerShort string
		if match.ScoreA == 0 {
			match.ScoreA = 1
		}
		if match.ScoreB == 0 {
			match.ScoreB = 1
		}
		rp := regexp.MustCompile("(?i)" + "\\b" + teamName + "\\b")
		if match.Winner == match.TeamA {
			winnerShort = match.TeamAShort
		}
		if rp.MatchString(match.Winner) || rp.MatchString(winnerShort) {
			win++
		}
		if match.TeamAShort == teamName {
			teamName = match.TeamA
		}
		if match.TeamBShort == teamName {
			teamName = match.TeamB
		}
	}
	winrate = win / float64(len(data))

	result = FBResult{Matches: data, Name: teamName, Winrate: winrate}
	return result, nil
}
