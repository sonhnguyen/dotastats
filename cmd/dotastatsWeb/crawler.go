package main

import (
	"dotastats"
	"net/http"
)

func (a *App) RunCrawlerAndSave() error {
	result, err := dotastats.RunCrawlerDota2BestYolo()
	if err != nil {
		return err
	}
	err = a.mongodb.SaveMatches(result)
	_, err = http.Get("http://dotabetstats.herokuapp.com")
	if err != nil {
		return err
	}
	return nil
}
