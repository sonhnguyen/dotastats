package dotastats

func GetTeamInfo(teamSlug string, apiParams APIParams, mongodb Mongodb) (TeamInfo, error) {
	result, err := mongodb.GetTeamInfo(teamSlug, apiParams)
	if err != nil {
		return TeamInfo{}, err
	}

	return result, nil
}

func GetTeamMatches(teamName string, apiParams APIParams, mongodb Mongodb) ([]Match, error) {
	result, err := mongodb.GetTeamMatches(teamName, apiParams)
	if err != nil {
		return []Match{}, err
	}

	return result, nil
}
func GetTeamHistory(teamA, teamB string, apiParams APIParams, mongodb Mongodb) ([]Match, error) {
	result, err := mongodb.GetTeamHistoryMatches(teamA, teamB, apiParams)
	if err != nil {
	}
	return result, nil
}
