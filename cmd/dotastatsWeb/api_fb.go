package main

import (
	"encoding/json"
	"net/http"

	"dotastats"
)

func (a *App) GetTeamFBMatchesHandler() HandlerWithError {
	return func(w http.ResponseWriter, req *http.Request) error {
		params := GetParamsObj(req)
		teamName := params.ByName("name")
		apiParams, err := BuildAPIParams(req)
		if err != nil {
			a.logr.Log("error when  building params %s", err)
			return newAPIError(500, "error when building params %s", err)
		}

		result, err := dotastats.GetTeamFBMatches(teamName, apiParams, a.mongodb)
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
