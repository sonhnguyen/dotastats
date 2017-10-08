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

func FindSeriesIDExists(seriesID string, series []dotastats.Series) int {
	for index, value := range series {
		if seriesID == value.SeriesID {
			return index
		}
	}
	return -1
}

func ConvertMatchesToSeries(matches []dotastats.Match) []dotastats.Series {
	var result []dotastats.Series

	for _, value := range matches {
		if indexSeries := FindSeriesIDExists(value.SeriesID, result); indexSeries == -1 {
			newSeries := dotastats.Series{Matches: []dotastats.Match{value}, SeriesID: value.SeriesID}
			result = append(result, newSeries)
		} else {
			result[indexSeries].Matches = append(result[indexSeries].Matches, value)
		}
	}
	return result
}

// Build all common params of an API endpoint.
func BuildAPIParams(req *http.Request) (dotastats.APIParams, error) {
	var apiParams dotastats.APIParams
	queryValues := req.URL.Query()
	if value := queryValues.Get("limit"); value != "" {
		limitInt, err := strconv.Atoi(value)
		if err != nil {
			return dotastats.APIParams{}, newAPIError(500, "error when process limit info %s", err)
		}
		apiParams.Limit = limitInt
	}

	if value := queryValues.Get("skip"); value != "" {
		skipInt, err := strconv.Atoi(value)
		if err != nil {
			return dotastats.APIParams{}, newAPIError(500, "error when process skip info %s", err)
		}
		apiParams.Skip = skipInt
	}

	if value := queryValues.Get("fields"); value != "" {
		apiParams.Fields = strings.Split(value, ",")
	}

	dateFormat := "02012006"
	apiParams.TimeTo = time.Now().AddDate(0, 0, 2)
	apiParams.TimeFrom = time.Date(1994, 11, 24, 0, 0, 0, 0, time.UTC)
	if value := queryValues.Get("time_from"); value != "" {
		t, err := time.Parse(dateFormat, value)
		if err != nil {
			return dotastats.APIParams{}, newAPIError(500, "error when process date info %s", err)
		}
		apiParams.TimeFrom = t
	}

	if value := queryValues.Get("time_to"); value != "" {
		t, err := time.Parse(dateFormat, value)
		if err != nil {
			return dotastats.APIParams{}, newAPIError(500, "error when process date info %s", err)
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

// getUser return the user in request
func getUser(req *http.Request) *dotastats.User {
	if rv := context.Get(req, UserKeyName); rv != nil {
		res := rv.(*dotastats.User)
		return res
	}

	return nil
}
