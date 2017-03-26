package main

import (
	"dotastats"
	"net/http"
)

func (a *App) RunCrawler() ([]dotastats.Match, error) {
	var vpParams = dotastats.VPGameAPIParams{Page: "1", Status: "close", Limit: "200"}
	closedMatches, err := dotastats.RunCrawlerVpgame(vpParams)
	if err != nil {
		return []dotastats.Match{}, err
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

func (a *App) RunPingHeroku() error {
	_, err := http.Get("http://dotabetstats.herokuapp.com")
	if err != nil {
		return err
	}
	_, err = http.Get("http://f10k.herokuapp.com")
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
