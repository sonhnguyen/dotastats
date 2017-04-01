package dotastats

func GetTeamInfo(teamSlug string, apiParams APIParams, mongodb Mongodb) ([]TeamInfo, error) {
	result, err := mongodb.GetTeamInfo(teamSlug, apiParams)
	if err != nil {
		return []TeamInfo{}, err
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
