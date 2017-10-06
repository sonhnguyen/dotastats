package main

import (
	"dotastats"
	"encoding/json"
	"net/http"

	"github.com/asaskevich/govalidator"
)

func (a *App) LoginPostHandler() HandlerWithError {
	return func(w http.ResponseWriter, req *http.Request) error {
		email, pass := req.FormValue("email"), req.FormValue("password")

		if !govalidator.IsEmail(email) {
			return newAPIError(400, "email provided is not valid", nil)
		}

		user, err := dotastats.GetUser(email, pass, a.mongodb)
		if err != nil {
			return newAPIError(400, "can not find user with provided email and password", err)
		}

		err = json.NewEncoder(w).Encode(user)
		if err != nil {
			return newAPIError(500, "error when return json", err)
		}

		return nil
	}
}

func (a *App) RegisterPostHandler() HandlerWithError {
	return func(w http.ResponseWriter, req *http.Request) error {
		name := req.FormValue("name")
		email, pass := req.FormValue("email"), req.FormValue("password")
		key := req.FormValue("register_key")

		if a.config.RegisterKey != key {
			return newAPIError(400, "register key is not correct", nil)
		}

		user, err := dotastats.CreateUser(name, email, pass, a.mongodb)
		if err != nil {
			return newAPIError(500, "error when creating user", err)
		}

		err = json.NewEncoder(w).Encode(user)
		if err != nil {
			return newAPIError(500, "error when return json", err)
		}

		return nil
	}
}
