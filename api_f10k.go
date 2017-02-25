package dotastats

import (
	"fmt"
	"regexp"
	"strings"
)

func GetTeamMatches(teamName, limit, skip string, fields []string, mongodb Mongodb) ([]Match, error) {
	result, err := mongodb.GetTeamMatches(teamName, limit, skip, fields)
	if err != nil {
		return []Match{}, err
	}

	return result, nil
}

func GetTeamF10kMatches(teamName, limit, skip string, fields []string, mongodb Mongodb) ([]Match, error) {
	result, err := mongodb.GetTeamF10kMatches(teamName, limit, skip, fields)
	if err != nil {
		return []Match{}, err
	}

	return result, nil
}

func GetF10kResult(teamName, limit, skip string, mongodb Mongodb) (F10kResult, error) {
	data, err := mongodb.GetTeamF10kMatches(teamName, limit, skip, []string{})
	if err != nil {
		return F10kResult{}, err
	}
	if len(data) == 0 {
		return F10kResult{}, nil
	}

	var result F10kResult
	var f10kHistory []F10kHistory
	var kill, death, ratio, totalKill, totalDeath, win, avgKill, avgDeath, winrate, avgOdds, ratioKill float64
	for _, match := range data {
		var summary F10kHistory
		var winnerShort string
		if match.ScoreA == 0 {
			match.ScoreA = 1
		}
		if match.ScoreB == 0 {
			match.ScoreB = 1
		}
		rp := regexp.MustCompile("(?i)" + "\\b" + teamName + "\\b")
		fmt.Println("teamname", teamName, strings.ToLower(match.TeamA), strings.ToLower(match.TeamB))
		if rp.MatchString(match.TeamA) || rp.MatchString(match.TeamAShort) {
			kill = float64(match.ScoreA)
			death = float64(match.ScoreB)
			ratio = float64(match.RatioA)
			summary.Name = strings.ToLower(match.TeamB)
		} else if rp.MatchString(match.TeamB) || rp.MatchString(match.TeamBShort) {
			kill = float64(match.ScoreB)
			death = float64(match.ScoreA)
			ratio = float64(match.RatioB)
			summary.Name = strings.ToLower(match.TeamA)
		}
		totalKill += kill
		totalDeath += death
		summary.Kill = kill
		summary.Death = death
		summary.Winner = match.Winner
		summary.Time = match.Time
		if match.Winner == match.TeamA {
			winnerShort = match.TeamAShort
		}
		if rp.MatchString(match.Winner) || rp.MatchString(winnerShort) {
			win++
		}
		avgOdds += ratio
		ratioKill += kill / death
		f10kHistory = append(f10kHistory, summary)
	}
	avgKill = totalKill / float64(len(data))
	avgDeath = totalDeath / float64(len(data))
	winrate = win / float64(len(data))
	avgOdds = avgOdds / float64(len(data))
	ratioKill = ratioKill / float64(len(data))

	result = F10kResult{F10kHistory: f10kHistory, AverageDeath: avgDeath, AverageKill: avgKill, Name: teamName, RatioKill: ratioKill, TotalKill: totalKill, TotalDeath: totalDeath, Winrate: winrate, AverageOdds: avgOdds}
	return result, nil
}
