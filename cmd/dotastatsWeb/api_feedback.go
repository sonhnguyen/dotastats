package main

import (
	"encoding/json"
	"net/http"
	"time"

	"dotastats"
)

func (a *App) PostFeedback() HandlerWithError {
	return func(w http.ResponseWriter, req *http.Request) error {
		var feedBack dotastats.Feedback
		err := json.NewDecoder(req.Body).Decode(&feedBack)
		if err != nil {
			return newAPIError(500, "error when decoding request body", err)
		}

		if feedBack.Name == "" || feedBack.Feedback == "" {
			return newAPIError(400, "name and feedback should not be empty", nil)
		}

		feedBack.Time = time.Now()

		err = a.mongodb.SaveFeedback(&feedBack)
		if err != nil {
			a.logr.Log("error when saving feedback %s", err)
			return newAPIError(500, "error when saving feedback %s", err)
		}

		return nil
	}
}

func (a *App) GetFeedback() HandlerWithError {
	return func(w http.ResponseWriter, req *http.Request) error {
		feedBackArr, err := a.mongodb.GetFeedback()
		if err != nil {
			a.logr.Log("error when getting feedback %s", err)
			return newAPIError(500, "error when getting feedback", nil)
		}

		response := struct {
			Data []dotastats.Feedback `json:"data"`
		}{
			feedBackArr,
		}
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			return newAPIError(500, "error when return json", nil)
		}

		return nil
	}
}
