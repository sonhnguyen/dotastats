package main

import (
	"encoding/json"
	"net/http"
	"strings"

	"dotastats"
)

func (a *App) GetTeamMatchesHandler() HandlerWithError {
	return func(w http.ResponseWriter, req *http.Request) error {

		params := GetParamsObj(req)
		queryValues := req.URL.Query()
		teamName := params.ByName("name")
		limit := queryValues.Get("limit")
		skip := queryValues.Get("skip")
		var fields []string
		if value := queryValues.Get("fields"); value != "" {
			fields = strings.Split(value, ",")
		}

		result, err := dotastats.GetTeamMatches(teamName, limit, skip, fields, a.mongodb)
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

func (a *App) GetTeamF10kMatchesHandler() HandlerWithError {
	return func(w http.ResponseWriter, req *http.Request) error {

		params := GetParamsObj(req)
		queryValues := req.URL.Query()
		teamName := params.ByName("name")
		limit := queryValues.Get("limit")
		skip := queryValues.Get("skip")
		var fields []string
		if value := queryValues.Get("fields"); value != "" {
			fields = strings.Split(value, ",")
		}

		result, err := dotastats.GetTeamF10kMatches(teamName, limit, skip, fields, a.mongodb)
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

func (a *App) GetF10kResultHandler() HandlerWithError {
	return func(w http.ResponseWriter, req *http.Request) error {

		params := GetParamsObj(req)
		queryValues := req.URL.Query()
		teamName := params.ByName("name")
		limit := queryValues.Get("limit")
		skip := queryValues.Get("skip")

		result, err := dotastats.GetF10kResult(teamName, limit, skip, a.mongodb)
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
