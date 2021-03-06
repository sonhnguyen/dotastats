package main

import (
	"encoding/json"
	"net/http"

	"dotastats"
)

func (a *App) GetMatchesListHandler() HandlerWithError {
	return func(w http.ResponseWriter, req *http.Request) error {

		queryValues := req.URL.Query()
		status := queryValues.Get("status")
		apiParams, err := BuildAPIParams(req)
		if err != nil {
			a.logr.Log("error when  building params %s", err)
			return newAPIError(500, "error when building params", err)
		}
		result, err := dotastats.GetMatchesList(status, apiParams, a.mongodb)
		if err != nil {
			a.logr.Log("error when return json %s", err)
			return newAPIError(500, "error when return json", err)
		}
		seriesList := ConvertMatchesToSeries(result)
		err = json.NewEncoder(w).Encode(seriesList)
		if err != nil {
			a.logr.Log("error when return json %s", err)
			return newAPIError(500, "error when return json", err)
		}
		return nil
	}
}

func (a *App) GetMatchByIDHandler() HandlerWithError {
	return func(w http.ResponseWriter, req *http.Request) error {

		params := GetParamsObj(req)
		matchID := params.ByName("id")
		result, err := dotastats.GetMatchByID(matchID, a.mongodb)
		if err != nil {
			a.logr.Log("error when return json %s", err)
			return newAPIError(500, "error when return json", err)
		}

		err = json.NewEncoder(w).Encode(result)
		if err != nil {
			a.logr.Log("error when return json %s", err)
			return newAPIError(500, "error when return json", err)
		}
		return nil
	}
}
