package dotastats

func GetTeamMatches(teamName, limit string, mongodb Mongodb) ([]Match, error) {
	result, err := mongodb.GetTeamMatches(teamName, limit)
	if err != nil {
		return []Match{}, err
	}

	return result, nil
}

func GetTeamF10kMatches(teamName, limit string, mongodb Mongodb) ([]Match, error) {
	result, err := mongodb.GetTeamF10kMatches(teamName, limit)
	if err != nil {
		return []Match{}, err
	}

	return result, nil
}

func GetF10kResult(teamName, limit string, mongodb Mongodb) (F10kResult, error) {
	data, err := mongodb.GetTeamF10kMatches(teamName, limit)
	if err != nil {
		return F10kResult{}, err
	}
	if len(data) == 0 {
		return F10kResult{}, nil
	}

	var result F10kResult
	var kill, death int
	var ratio float64
	var totalKill, totalDeath, win int
	var avgKill, avgDeath, winrate, avgOdds, ratioKill float64
	for _, match := range data {
		if match.ScoreA == 0 {
			match.ScoreA = 1
		}
		if match.ScoreB == 0 {
			match.ScoreB = 1
		}
		if teamName == match.TeamA {
			kill = match.ScoreA
			death = match.ScoreB
			ratio = match.RatioA
		} else {
			kill = match.ScoreB
			death = match.ScoreA
			ratio = match.RatioB
		}
		totalKill += kill
		totalDeath += death
		if teamName == match.Winner {
			win++
		}
		avgOdds += ratio
		ratioKill += float64(kill) / float64(death)
	}
	avgKill = float64(totalKill) / float64(len(data))
	avgDeath = float64(totalDeath) / float64(len(data))
	winrate = float64(win) / float64(len(data))
	avgOdds = float64(avgOdds) / float64(len(data))
	ratioKill = ratioKill / float64(len(data))
	return result, nil
}
