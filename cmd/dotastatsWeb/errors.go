package main

import (
	"fmt"
)

// Error represents a handler error. It provides methods for a HTTP status
// code and embeds the built-in error interface.
type Error interface {
	error
	Status() int
}

type APIError struct {
	Code    int    `json:"-"`
	Err     error  `json:"-"`
	Message string `json:"message"`
}

// StatusError represent an error with an associated HTTP status code
type StatusError struct {
	Code int
	Err  error
}

// Allows StatusError to satisfy the error interface
func (se StatusError) Error() string {
	return se.Err.Error()
}

// Returns our HTTP status code.
func (se StatusError) Status() int {
	return se.Code
}

// Returns our HTTP status code on API.
func (ae APIError) Status() int {
	return ae.Code
}

func (ae APIError) Error() string {
	return ae.Err.Error()
}

func newAPIError(code int, msg string, err error) *APIError {
	if err != nil {
		return &APIError{Code: code, Err: fmt.Errorf(msg+": %s", err), Message: err.Error()}
	} else {
		return &APIError{Code: code, Err: fmt.Errorf(msg), Message: msg}
	}
}

func newError(code int, msg string, err error) *StatusError {
	if err != nil {
		return &StatusError{Code: code, Err: fmt.Errorf(msg+": %s", err)}
	} else {
		return &StatusError{Code: code, Err: fmt.Errorf(msg)}
	}
}

func newSessionSaveError(err error) *StatusError {
	return &StatusError{Code: 500, Err: fmt.Errorf("problem saving to cookie store: %s", err)}
}

func newRenderErrMsg(err error) string {
	return fmt.Sprintf("error rendering HTML: %s", err)
}
