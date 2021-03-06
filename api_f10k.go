package dotastats

import "regexp"

func GetTeamF10kMatches(teamName string, apiParams APIParams, mongodb Mongodb) (F10kResult, error) {
	data, err := mongodb.GetTeamF10kMatches(teamName, apiParams)
	if err != nil {
		return F10kResult{}, err
	}
	if len(data) == 0 {
		return F10kResult{}, nil
	}

	var result F10kResult
	var kill, death, ratio, totalKill, totalDeath, win, avgKill, avgDeath, winrate, avgOdds, ratioKill float64
	for _, match := range data {
		var winnerShort string
		if match.ScoreA == 0 {
			match.ScoreA = 1
		}
		if match.ScoreB == 0 {
			match.ScoreB = 1
		}
		rp := regexp.MustCompile("(?i)" + "\\b" + teamName + "\\b")
		if rp.MatchString(match.TeamA) || rp.MatchString(match.TeamAShort) {
			kill = float64(match.ScoreA)
			death = float64(match.ScoreB)
			ratio = float64(match.RatioA)
		} else if rp.MatchString(match.TeamB) || rp.MatchString(match.TeamBShort) {
			kill = float64(match.ScoreB)
			death = float64(match.ScoreA)
			ratio = float64(match.RatioB)
		}
		totalKill += kill
		totalDeath += death
		if match.Winner == match.TeamA {
			winnerShort = match.TeamAShort
		}
		if rp.MatchString(match.Winner) || rp.MatchString(winnerShort) {
			win++
		}
		avgOdds += ratio
		ratioKill += kill / death
		if match.TeamAShort == teamName {
			teamName = match.TeamA
		}
		if match.TeamBShort == teamName {
			teamName = match.TeamB
		}
	}
	avgKill = totalKill / float64(len(data))
	avgDeath = totalDeath / float64(len(data))
	winrate = win / float64(len(data))
	avgOdds = avgOdds / float64(len(data))
	ratioKill = ratioKill / float64(len(data))

	result = F10kResult{Matches: data, AverageDeath: avgDeath, AverageKill: avgKill,
		Name: teamName, RatioKill: ratioKill, TotalKill: totalKill, TotalDeath: totalDeath,
		Winrate: winrate, AverageOdds: avgOdds}
	return result, nil
}
