package main

import (
	"dotastats"
	"net/http"
)

func (a *App) RunCrawlerAndSave() error {
	var vpParams = dotastats.VPGameAPIParams{Page: "1", Status: "close"}
	closedMatches, err := dotastats.RunCrawlerVpgame(vpParams)
	if err != nil {
		return err
	}

	vpParams.Status = "open"
	openingMatches, err := dotastats.RunCrawlerVpgame(vpParams)
	if err != nil {
		return err
	}

	vpParams.Status = "start"
	liveMatches, err := dotastats.RunCrawlerVpgame(vpParams)
	if err != nil {
		return err
	}

	result := append(closedMatches, openingMatches...)
	result = append(result, liveMatches...)
	err = a.mongodb.SaveMatches(result)
	_, err = http.Get("http://dotabetstats.herokuapp.com")
	if err != nil {
		return err
	}
	return nil
}
