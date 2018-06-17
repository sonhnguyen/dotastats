package main

import (
	"dotastats"
	"fmt"
	"net/http"
	"strconv"
)

func (a *App) RunCrawler() ([]dotastats.Match, error) {
	var vpParams = dotastats.VPGameAPIParams{Page: "1", Status: "close", Limit: "400"}
	closedMatches, err := dotastats.RunCrawlerVpgame(vpParams)
	if err != nil {
		return []dotastats.Match{}, err
	}

	for i, closedMatch := range closedMatches {
		if closedMatch.Game != "dota" {
			continue
		}
		if closedMatch.Status == "Canceled" {
			continue
		}
		openDotaMatch, teamAIsRadiant, err := a.mongodb.GetOpenDotaMatch(closedMatch)
		if err != nil {
			fmt.Errorf("error getting matchid from vpgame crawler: %s", err)
			continue
		}
		closedMatches[i].DotaMatchID = openDotaMatch.MatchID
		if openDotaMatch.MatchID != 0 {
			closedMatches[i].OpenDotaURL = "https://www.opendota.com/matches/" + strconv.Itoa(openDotaMatch.MatchID)
			closedMatches[i].DotaBuffURL = "https://www.dotabuff.com/matches/" + strconv.Itoa(openDotaMatch.MatchID)
			closedMatches[i].PicksBans = openDotaMatch.PicksBans
			closedMatches[i].Duration = openDotaMatch.Duration
			closedMatches[i].TeamAIsRadiant = teamAIsRadiant
		}
	}

	vpParams.Status = "open"
	openingMatches, err := dotastats.RunCrawlerVpgame(vpParams)
	if err != nil {
		return []dotastats.Match{}, err
	}

	vpParams.Status = "start"
	liveMatches, err := dotastats.RunCrawlerVpgame(vpParams)
	if err != nil {
		return []dotastats.Match{}, err
	}

	result := append(closedMatches, openingMatches...)
	result = append(result, liveMatches...)
	return result, nil
}

func (a *App) RunCrawlerTeamInfo() ([]dotastats.TeamInfo, error) {
	dotaTeams, err := dotastats.RunCrawlerLiquidTeam("dota")
	if err != nil {
		return []dotastats.TeamInfo{}, err
	}
	csgoTeams, err := dotastats.RunCrawlerLiquidTeam("csgo")
	if err != nil {
		return []dotastats.TeamInfo{}, err
	}
	result := append(csgoTeams, dotaTeams...)
	return result, nil
}

func (a *App) RunCrawlerOpenDotaTeamAndSave() error {
	opendotaTeams, err := dotastats.RunCrawlerOpenDotaTeam()
	if err != nil {
		return err
	}
	err = a.mongodb.SaveOpenDotaTeam(opendotaTeams)
	if err != nil {
		return err
	}

	return nil
}

func (a *App) RunPingHeroku() error {
	_, err := http.Get("http://dotabetstats.herokuapp.com")
	if err != nil {
		return err
	}
	_, err = http.Get("http://f10k.herokuapp.com")
	if err != nil {
		return err
	}
	_, err = http.Get("http://dotastats-client.herokuapp.com")
	if err != nil {
		return err
	}
	_, err = http.Get("http://dotastats.me")
	if err != nil {
		return err
	}
	return nil
}

func (a *App) RunCrawlerTeamInfoAndSave() error {
	result, err := a.RunCrawlerTeamInfo()
	if err != nil {
		return err
	}

	err = a.mongodb.SaveTeamInfo(result)
	if err != nil {
		return err
	}

	err = dotastats.CreateTwitterList(result)
	if err != nil {
		return err
	}

	return nil
}

func (a *App) RunCrawlerOpenDotaProMatchesAndSave() error {
	var openDotaAPIParams = dotastats.OpenDotaAPIParams{}
	result, err := dotastats.RunCrawlerOpenDota(openDotaAPIParams)
	if err != nil {
		return err
	}

	err = a.mongodb.SaveOpenDotaProMatches(result)
	if err != nil {
		return err
	}

	return nil
}

func (a *App) RunCrawlerAndSave() error {
	err := a.RunPingHeroku()
	if err != nil {
		return err
	}
	result, err := a.RunCrawler()
	if err != nil {
		return err
	}
	err = a.mongodb.SaveMatches(result)
	if err != nil {
		return err
	}
	return nil
}
