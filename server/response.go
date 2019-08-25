package server

import (
	"encoding/json"
	"net/http"
)

type SuccessMessage struct {
	Message string `json:"message"`
}

func respondWithError(w http.ResponseWriter, code int, pretty bool, message string) {
	respondWithJSON(w, code, pretty, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, pretty bool, payload interface{}) {
	response := []byte{}
	if pretty {
		response, _ = json.MarshalIndent(payload, "", "    ")
	} else {
		response, _ = json.Marshal(payload)
	}

	response = append(response, '\n')

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
