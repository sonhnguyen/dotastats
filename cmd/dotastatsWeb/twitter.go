package main

import (
	"dotastats"
	"fmt"
	"log"
	"net/http"

	"github.com/spf13/viper"
)

func (a *App) RemoveAllTwitterList() HandlerWithError {
	return func(w http.ResponseWriter, req *http.Request) error {
		twitterID := viper.GetString("twitter.twitterID")
		fmt.Println(twitterID)
		err := dotastats.RemoveAllListFromTwitter(twitterID)
		if err != nil {
			log.Println("error removing twitter list %s", err)
			return err
		}

		return nil
	}
}
func (a *App) CreateAllTwitterList() HandlerWithError {
	return func(w http.ResponseWriter, req *http.Request) error {
		teams, err := a.mongodb.GetAllTeamInfo()
		if err != nil {
			log.Println("error querying teams %s", err)
			return err
		}

		err = dotastats.CreateTwitterList(teams)
		if err != nil {
			log.Println("error creating twitter list %s", err)
			return err
		}

		return nil
	}
}
