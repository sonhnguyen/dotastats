package dotastats

func GetTeamMatches(teamName string, apiParams APIParams, mongodb Mongodb) ([]Match, error) {
	result, err := mongodb.GetTeamMatches(teamName, apiParams)
	if err != nil {
		return []Match{}, err
	}

	return result, nil
}
