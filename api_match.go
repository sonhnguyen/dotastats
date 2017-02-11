package dotastats

func GetMatches(limit, skip, status string, fields []string, mongodb Mongodb) ([]Match, error) {
	result, err := mongodb.GetMatches(limit, skip, status, fields)
	if err != nil {
		return []Match{}, err
	}

	return result, nil
}
