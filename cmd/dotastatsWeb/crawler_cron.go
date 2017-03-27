package main

import (
	"dotastats"
	"fmt"
	"net/http"
	"os"

	"github.com/spf13/viper"
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

func (a *App) SaveTeamListToTwitter(teams []dotastats.TeamInfo) error {
	c, err := dotastats.CreateOAuth()
	if err != nil {
		return err
	}
	var twitterID string
	if viper.GetBool("isDevelopment") {
		twitterID = viper.GetString("twitter.twitterID")
	} else {
		twitterID = os.Getenv("twitterID")
	}

	fmt.Println(len(teams))
	for _, team := range teams {
		nameSlug := "t-" + team.NameSlug
		if len(nameSlug) > 10 {
			nameSlug = nameSlug[:10]
		}

		nameSlug = nameSlug[1 : len(nameSlug)-1]
		err = dotastats.RemoveListFromTwitter(c, dotastats.TwitterRemoveListRequest{
			OwnerScreenName: twitterID,
			Slug:            nameSlug + "1",
		})

		err = dotastats.RemoveListFromTwitter(c, dotastats.TwitterRemoveListRequest{
			OwnerScreenName: twitterID,
			Slug:            nameSlug + "2",
		})
		err = dotastats.RemoveListFromTwitter(c, dotastats.TwitterRemoveListRequest{
			OwnerScreenName: twitterID,
			Slug:            nameSlug + "3",
		})
		err = dotastats.RemoveListFromTwitter(c, dotastats.TwitterRemoveListRequest{
			OwnerScreenName: twitterID,
			Slug:            nameSlug + "4",
		})
		err = dotastats.RemoveListFromTwitter(c, dotastats.TwitterRemoveListRequest{
			OwnerScreenName: twitterID,
			Slug:            nameSlug + "5",
		})
		err = dotastats.RemoveListFromTwitter(c, dotastats.TwitterRemoveListRequest{
			OwnerScreenName: twitterID,
			Slug:            nameSlug + "6",
		})
		err = dotastats.RemoveListFromTwitter(c, dotastats.TwitterRemoveListRequest{
			OwnerScreenName: twitterID,
			Slug:            nameSlug + "7",
		})
		err = dotastats.RemoveListFromTwitter(c, dotastats.TwitterRemoveListRequest{
			OwnerScreenName: twitterID,
			Slug:            nameSlug + "8",
		})
		err = dotastats.RemoveListFromTwitter(c, dotastats.TwitterRemoveListRequest{
			OwnerScreenName: twitterID,
			Slug:            nameSlug + "9",
		})
		err = dotastats.RemoveListFromTwitter(c, dotastats.TwitterRemoveListRequest{
			OwnerScreenName: twitterID,
			Slug:            nameSlug + "10",
		})
		if err != nil {
			return err
		}

	}

	return nil
}

func (a *App) RunCrawlerTeamInfoAndSave() error {
	result, err := a.RunCrawlerTeamInfo()
	if err != nil {
		return err
	}

	err = a.SaveTeamListToTwitter(result)
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
