package dotastats

import (
	"regexp"
	"strconv"
	"strings"
)

func scoreProcess(score string) []float64 {
	var result []float64
	score = strings.TrimSpace(score)
	scoreArr := strings.Split(score, ":")
	for _, scoreString := range scoreArr {
		if scoreString != " " {
			scoreInt, err := strconv.Atoi(strings.TrimSpace(scoreString))
			if err != nil {
				return []float64{}
			}
			result = append(result, float64(scoreInt))
		}
	}
	return result
}

func matchIDProcess(matchID string) int {
	re := regexp.MustCompile("[^0-9]")
	result, err := strconv.Atoi(re.ReplaceAllString(matchID, ""))
	if err != nil {
		return 0
	}
	return result
}

func trimLine(line string) string {
	s := strings.TrimSpace(line)
	return s
}

func ratioProcess(ratio string) float64 {
	re := regexp.MustCompile("\\d+\\.\\d+")
	number := re.FindAllString(ratio, -1)
	if number != nil {
		result, err := strconv.ParseFloat(number[0], 64)
		if err != nil {
			return 0
		}
		return result
	}
	return 0
}
func winnerProcess(winner string) string {
	if len(winner) < 10 {
		return "TBD"
	}
	return winner[9:]
}
