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
			return nil
		}

		if feedBack.Name == "" || feedBack.Feedback == "" {
			return nil
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
