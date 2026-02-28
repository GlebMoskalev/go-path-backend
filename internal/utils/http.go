package utils

import (
	"encoding/json"
	"net/http"
)

func ResponseWithJSON(w http.ResponseWriter, code int, payload any) {
	w.Header().Set("Content-Type", "application/json")

	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
		return
	}
	w.WriteHeader(code)
	w.Write(response)
}

func ResponseWithError(w http.ResponseWriter, code int, message string) {
	ResponseWithJSON(w, code, map[string]string{"error": message})
}
