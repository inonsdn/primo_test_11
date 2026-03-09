package router

import (
	"encoding/json"
	"net/http"
	"primo_test_11/internal/database"
)

// type of function to received route handler which is wrapper of connection
type RouteHandlerFunc func(*RouteHandler) error

type RoutePathHandler struct {
	Method  string
	Path    string
	Handler RouteHandlerFunc
}

type RouteHandler struct {
	w         http.ResponseWriter
	r         *http.Request
	dbHandler *database.DBHandler
}

// Create handler function for serve http
// by wrapping function
// function must receive argument of route handler
func makeHandler(h RoutePathHandler, dbHandler *database.DBHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// new handler for this response
		rh := &RouteHandler{
			w:         w,
			r:         r,
			dbHandler: dbHandler,
		}

		// verify method
		if r.Method != h.Method {
			rh.ResponseError(http.StatusMethodNotAllowed, "Invalid method")
		} else {
			// execute function
			if err := h.Handler(rh); err != nil {
				rh.Response(http.StatusInternalServerError, err.Error())
			}
		}
	}
}

// wrapper function to get query params from url
// eg. url is https://example.com/getUser?id=1
// and format is /getUser
// use this function to get id by GetQuery("id")
// and this function will return value of given name of parameters
func (rh *RouteHandler) GetQuery(name string) string {
	return rh.r.URL.Query().Get(name)
}

func (rh *RouteHandler) GetJSON(requestBody any) error {
	if err := json.NewDecoder(rh.r.Body).Decode(requestBody); err != nil {
		return err
	}
	return nil
}

func (rh *RouteHandler) ResponseError(status int, message string) {
	http.Error(rh.w, message, status)
}

// send response back in JSON format
func (rh *RouteHandler) ResponseJSON(status int, data any) {
	// set header of response to be type of json
	rh.w.Header().Set("Content-Type", "application/json")
	rh.w.WriteHeader(status)

	// encoded data
	if err := json.NewEncoder(rh.w).Encode(data); err != nil {
		rh.ResponseError(http.StatusInternalServerError, err.Error())
	}
}

// response function to just response plain text
// with http status
func (rh *RouteHandler) Response(status int, body string) {
	// set header of response
	rh.w.Header().Set("Content-Type", "text/plain")
	rh.w.WriteHeader(status)
	rh.w.Write([]byte(body))
}
