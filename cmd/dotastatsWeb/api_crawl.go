package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"dotastats"
)

func (a *App) GetCustomCrawlHandler() HandlerWithError {
	return func(w http.ResponseWriter, req *http.Request) error {
		queryValues := req.URL.Query()
		pageFrom, err := strconv.Atoi(queryValues.Get("page_from"))
		if err != nil {
			a.logr.Log("error when converting params %s", err)
			return newAPIError(500, "error when converting params %s", err)
		}
		var pageTo int
		if queryValues.Get("page_to") != "" {
			pageTo, err = strconv.Atoi(queryValues.Get("page_to"))
			if err != nil {
				a.logr.Log("error when converting params %s", err)
				return newAPIError(500, "error when converting params %s", err)
			}
		} else {
			pageTo = pageFrom
		}

		status := queryValues.Get("status")
		var matchesResults []dotastats.Match
		for i := pageFrom; i <= pageTo; i++ {
			pageNum := strconv.Itoa(i)
			var vpParams = dotastats.VPGameAPIParams{Page: pageNum, Status: status}
			matches, err := dotastats.RunCrawlerVpgame(vpParams)
			if err != nil {
				a.logr.Log("error when crawling manually %s", err)
				return newAPIError(500, "error when crawling manually %s", err)
			}
			matchesResults = append(matchesResults, matches...)
		}
		err = a.mongodb.SaveMatches(matchesResults)
		if err != nil {
			a.logr.Log("error when saving json crawled manually %s", err)
			return newAPIError(500, "error when saving json crawled manually %s", err)
		}

		err = json.NewEncoder(w).Encode(matchesResults)
		if err != nil {
			a.logr.Log("error when return json %s", err)
			return newAPIError(500, "error when return json %s", err)
		}

		return nil
	}
}
