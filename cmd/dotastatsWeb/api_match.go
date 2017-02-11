package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func (a *App) GetMatchesHandler() HandlerWithError {
	return func(w http.ResponseWriter, req *http.Request) error {

		queryValues := req.URL.Query()
		limit := queryValues.Get("limit")
		skip := queryValues.Get("skip")
		var fields []string
		status := queryValues.Get("status")
		if value := queryValues.Get("fields"); value != "" {
			fields = strings.Split(value, ",")
		}
		result, err := dotastats.GetMatches(limit, skip, status, fields, a.mongodb)
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