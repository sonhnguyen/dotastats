package main

import (
	"dotastats"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/asaskevich/govalidator"
)

func (a *App) LoginPostHandler() HandlerWithError {
	return func(w http.ResponseWriter, req *http.Request) error {
		var user dotastats.User
		err := json.NewDecoder(req.Body).Decode(&user)
		if err != nil {
			return newAPIError(500, "error when decoding request body", nil)
		}

		if user.Email == "" || user.Password == "" {
			return newAPIError(400, "name, email and password should not be empty", nil)
		}

		user, err = dotastats.GetUserAndAuthenticate(user.Email, user.Password, a.mongodb)
		if err != nil {
			return newAPIError(400, "can not find user with provided email and password", nil)
		}

		err = json.NewEncoder(w).Encode(user)
		if err != nil {
			return newAPIError(500, "error when return json", nil)
		}

		return nil
	}
}

func (a *App) RegisterPostHandler() HandlerWithError {
	return func(w http.ResponseWriter, req *http.Request) error {
		var user dotastats.User
		err := json.NewDecoder(req.Body).Decode(&user)
		if err != nil {
			return newAPIError(500, "error when decoding request body", nil)
		}

		if user.Name == "" || user.Email == "" || user.Password == "" {
			return newAPIError(400, "name, email and password should not be empty", nil)
		}

		if !govalidator.IsEmail(user.Email) {
			return newAPIError(400, "email provided is not valid", nil)
		}

		if strings.TrimSpace(user.Password) == "" {
			return newAPIError(400, "password should not be empty", nil)
		}

		if a.config.RegisterKey != user.RegisterKey {
			return newAPIError(400, "register key is not correct", nil)
		}

		user, err = dotastats.CreateUser(user, a.mongodb)
		if err != nil {
			return newAPIError(500, "error when creating user", nil)
		}

		err = json.NewEncoder(w).Encode(user)
		if err != nil {
			return newAPIError(500, "error when return json", nil)
		}

		return nil
	}
}
