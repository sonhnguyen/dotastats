package main

import (
	"net/http"

	"github.com/gorilla/context"
	"github.com/julienschmidt/httprouter"
)

// GetParamsObj returns a httprouter params object given the request.
func GetParamsObj(req *http.Request) httprouter.Params {
	ps := context.Get(req, Params).(httprouter.Params)
	return ps
}
