package main

import "dotastats"

func (a *App) RunCrawlerAndSave() error {
	_, err := dotastats.RunCrawlerDota2BestYolo()
	if err != nil {
		return err
	}
	return nil
}
