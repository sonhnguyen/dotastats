package main

import (
	"dotastats"
	"encoding/json"
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"

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

		userFromDB, err := dotastats.GetUserByEmail(user.Email, a.mongodb)
		if err != nil {
			return newAPIError(401, "Unauthorized", nil)
		}

		ss, err := a.store.Get(req, SessionName)
		if err != nil {
			return newAPIError(500, "error getting store", nil)
		}

		err = bcrypt.CompareHashAndPassword([]byte(userFromDB.Password), []byte(user.Password))
		if err != nil {
			ss.Save(req, w)
			return newAPIError(401, "Unauthorized", nil)
		}

		sess, err := dotastats.CreateSessionForUser(userFromDB.Email, a.mongodb)
		if err != nil {
			return newAPIError(500, "Error creating session for user", nil)
		}

		ss.Values[SessionKeyName] = sess.SessionKey
		ss.Save(req, w)
		user.Password = ""

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

		_, err = dotastats.GetUserByEmail(user.Email, a.mongodb)
		if err == nil {
			return newAPIError(400, "this email is already associated with an account", nil)
		}

		user, err = dotastats.CreateUser(user, a.mongodb)
		if err != nil {
			return newAPIError(500, "error when creating user", nil)
		}

		ss, err := a.store.Get(req, SessionName)
		if err != nil {
			return newAPIError(500, "error getting store", nil)
		}

		sess, err := dotastats.CreateSessionForUser(user.Email, a.mongodb)
		if err != nil {
			a.logr.Log("Error creating session for user: %s", err)
			return newAPIError(500, "Error creating session for user", nil)
		}

		ss.Values[SessionKeyName] = sess.SessionKey
		ss.Save(req, w)
		user.Password = ""

		err = json.NewEncoder(w).Encode(user)
		if err != nil {
			return newAPIError(500, "error when return json", nil)
		}

		return nil
	}
}
