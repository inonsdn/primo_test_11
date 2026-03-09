package handler

import (
	"encoding/json"
	"net/http"
)

type RouteHandlerFunc func(http.ResponseWriter, *http.Request)

type RoutePath struct {
	Method  string
	Path    string
	Handler RouteHandlerFunc
}

// Create handler function for serve http
// by wrapping function
// function must receive argument of route handler
func MakeHandler(rp RoutePath) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// verify method
		if r.Method != rp.Method {
			ResponseError(w, http.StatusMethodNotAllowed, "Invalid method")
		} else {
			// execute function
			rp.Handler(w, r)
		}
	}
}

// send response back in JSON format
func ResponseJSON(w http.ResponseWriter, status int, data any) {
	// set header of response to be type of json
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	// encoded data
	if err := json.NewEncoder(w).Encode(data); err != nil {
		ResponseError(w, http.StatusInternalServerError, err.Error())
	}
}

func ResponseError(w http.ResponseWriter, status int, message string) {
	http.Error(w, message, status)
}
