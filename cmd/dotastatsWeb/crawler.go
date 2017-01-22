package main

import (
	"dotastats"
	"fmt"
	"net/http"
	"time"
)

func (a *App) doEvery(d time.Duration) error {
	for x := range time.Tick(d) {
		result, err := dotastats.RunCrawlerDota2BestYolo()
		if err != nil {
			return err
		}
		err = a.mongodb.SaveMatches(result)
		fmt.Print("done run at %s", x)
		_, err = http.Get("http://dotabetstats.herokuapp.com")
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *App) RunCrawlerAndSave() error {
	err := a.doEvery(90 * time.Second)
	if err != nil {
		return err
	}
	return nil
}
