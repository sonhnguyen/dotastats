package dotastats

func GetMatchesList(status string, apiParams APIParams, mongodb Mongodb) ([]Match, error) {
	result, err := mongodb.GetMatchesList(status, apiParams)
	if err != nil {
		return []Match{}, err
	}

	return result, nil
}

func GetMatchByID(matchID string, mongodb Mongodb) (Match, error) {
	result, err := mongodb.GetMatchByID(matchID)
	if err != nil {
		return Match{}, err
	}

	return result, nil
}
