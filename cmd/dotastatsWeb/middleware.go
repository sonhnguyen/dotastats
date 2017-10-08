package main

import (
	"dotastats"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/context"
)

const UserKeyName = "user-dotastats-4829123"
const SessionKeyName = "session_key-9638182"
const SessionName = "session-dotastats-8429623"

// appLogger is an interface for logging.
// Used to introduce a seam into the app, for testing
type appLogger interface {
	Log(str string, v ...interface{})
}

// dotastatsLogger is a wrapper for long.Logger
type dotastatsLogger struct {
	*log.Logger
}

// Log produces a log entry with the current time prepended
func (ml *dotastatsLogger) Log(str string, v ...interface{}) {
	// Prepend current time to the slice of arguments
	v = append(v, 0)
	copy(v[1:], v[0:])
	v[0] = dotastats.TimeNow().Format(time.RFC3339)
	ml.Printf("[%s] "+str, v...)
}

// newMiddlewareLogger returns a new middlewareLogger.
func newLogger() *dotastatsLogger {
	return &dotastatsLogger{log.New(os.Stdout, "[dotastats] ", 0)}
}

// loggerHanderGenerator prduces a loggingHandler middleware
// loggingHandler middleware logs all request
func (a *App) loggingHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		t1 := dotastats.TimeNow()
		a.logr.Log("Started %s %s", req.Method, req.URL.Path)

		next.ServeHTTP(w, req)

		rw := w.(ResponseWriter)
		a.logr.Log("Completed %v %s in %v", rw.Status(), http.StatusText(rw.Status()), time.Since(t1))
	}
	return http.HandlerFunc(fn)
}

// recoverHandlerGenerator products a recoverHandler middleware
// recoverHander is an middleware that captures and recovers from panics
func (a *App) recoverHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				a.logr.Log("Panic: %+v", err)
				//
				http.Error(w, http.StatusText(500), 500)
			}
		}()
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

// Auth middleware
func (a *App) authMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		user := getUser(req)

		if user == nil {
			a.logr.Log("unauthorized access")
			e := newAPIError(401, "you have to login to access this", nil)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(e.Code)
			err := json.NewEncoder(w).Encode(e)
			if err != nil {
				a.logr.Log("error when return json %s", e)
			}
			return
		}
		next.ServeHTTP(w, req)
	}
	return http.HandlerFunc(fn)
}

// userMiddleware is the middleware wrapper that detects and provides the user
func (a *App) UserMiddlewareGenerator(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		fmt.Println("inside user middleware")
		session, err := a.store.Get(req, SessionName)
		if err != nil {
			a.logr.Log("error retrieving session from store", err)
			next.ServeHTTP(w, req)
			return
		}

		sessionKey, ok := session.Values[SessionKeyName]
		if ok {
			ssk := sessionKey.(string)
			s, err := dotastats.GetSessionBySessionKey(ssk, a.mongodb)
			if err != nil {
				a.logr.Log("error getting session with session key %s: %s", ssk, err)
				delete(session.Values, sessionKey)
				session.Save(req, w)
			} else {
				user, err := dotastats.GetUserByEmail(s.Email, a.mongodb)
				if err != nil {
					a.logr.Log("error getting user with session key %s: %s", ssk, err)
					delete(session.Values, sessionKey)
					session.Save(req, w)
				}
				context.Set(req, UserKeyName, &user)
			}
		}
		next.ServeHTTP(w, req)
	}
	return http.HandlerFunc(fn)
}
