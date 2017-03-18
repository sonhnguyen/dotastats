package main

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"dotastats"

	"github.com/gorilla/context"
	"github.com/julienschmidt/httprouter"
)

// Build all common params of an API endpoint.
func BuildAPIParams(req *http.Request) (dotastats.APIParams, error) {
	var apiParams dotastats.APIParams
	queryValues := req.URL.Query()
	apiParams.Limit = 200
	if value := queryValues.Get("limit"); value != "" {
		limitInt, err := strconv.Atoi(value)
		if err != nil {
			return dotastats.APIParams{}, newAPIError(300, "error when process limit info %s", err)
		}
		apiParams.Limit = limitInt
	}

	if value := queryValues.Get("skip"); value != "" {
		skipInt, err := strconv.Atoi(value)
		if err != nil {
			return dotastats.APIParams{}, newAPIError(300, "error when process skip info %s", err)
		}
		apiParams.Skip = skipInt
	}

	if value := queryValues.Get("fields"); value != "" {
		apiParams.Fields = strings.Split(value, ",")
	}

	dateFormat := "02012006"
	apiParams.TimeTo = time.Now()
	apiParams.TimeFrom = time.Date(1994, 11, 24, 0, 0, 0, 0, time.UTC)
	if value := queryValues.Get("time_from"); value != "" {
		t, err := time.Parse(dateFormat, value)
		if err != nil {
			return dotastats.APIParams{}, newAPIError(300, "error when process date info %s", err)
		}
		apiParams.TimeFrom = t
	}

	if value := queryValues.Get("time_to"); value != "" {
		t, err := time.Parse(dateFormat, value)
		if err != nil {
			return dotastats.APIParams{}, newAPIError(300, "error when process date info %s", err)
		}
		apiParams.TimeTo = t
	}

	if value := queryValues.Get("game"); value != "" {
		apiParams.Game = value
	}

	return apiParams, nil
}

// GetParamsObj returns a httprouter params object given the request.
func GetParamsObj(req *http.Request) httprouter.Params {
	ps := context.Get(req, Params).(httprouter.Params)
	return ps
}
