package main

import (
	"dotastats"

	"encoding/json"
	"net/http"
)

func (a *App) GetTeamInfoHandler() HandlerWithError {
	return func(w http.ResponseWriter, req *http.Request) error {

		params := GetParamsObj(req)
		teamSlug := params.ByName("slug")
		apiParams, err := BuildAPIParams(req)
		if err != nil {
			a.logr.Log("error when  building params %s", err)
			return newAPIError(300, "error when building params %s", err)
		}
		result, err := dotastats.GetTeamInfo(teamSlug, apiParams, a.mongodb)
		if err != nil {
			a.logr.Log("error when return json %s", err)
			return newAPIError(500, "error when return json %s", err)
		}
		err = json.NewEncoder(w).Encode(result)
		if err != nil {
			a.logr.Log("error when return json %s", err)
			return newAPIError(500, "error when return json %s", err)
		}
		return nil
	}
}

func (a *App) GetTeamMatchesHandler() HandlerWithError {
	return func(w http.ResponseWriter, req *http.Request) error {

		params := GetParamsObj(req)
		teamName := params.ByName("name")
		apiParams, err := BuildAPIParams(req)
		if err != nil {
			a.logr.Log("error when  building params %s", err)
			return newAPIError(300, "error when building params %s", err)
		}
		result, err := dotastats.GetTeamMatches(teamName, apiParams, a.mongodb)
		if err != nil {
			a.logr.Log("error when return json %s", err)
			return newAPIError(500, "error when return json %s", err)
		}
		seriesList := ConvertMatchesToSeries(result)
		err = json.NewEncoder(w).Encode(seriesList)
		if err != nil {
			a.logr.Log("error when return json %s", err)
			return newAPIError(500, "error when return json %s", err)
		}
		return nil
	}
}
