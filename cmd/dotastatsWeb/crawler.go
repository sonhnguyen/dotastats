package main

import (
	"dotastats"
	"fmt"
	"time"
)

func (a *App) doEvery(d time.Duration) error {
	for x := range time.Tick(d) {
		result, err := dotastats.RunCrawlerDota2BestYolo()
		if err != nil {
			return err
		}
		err = a.mongodb.SaveMatches(result)
		if err != nil {
			return err
		}
		fmt.Print(x)
	}
	return nil
}

func (a *App) RunCrawlerAndSave() error {
	err := a.doEvery(20 * time.Second)
	if err != nil {
		return err
	}
	return nil
}
