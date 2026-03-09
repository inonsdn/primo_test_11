package handler

import (
	"encoding/json"
	"net/http"
)

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
