package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func respondWithErr(w http.ResponseWriter, code int, msg string) {
	if code > 499 {
		log.Printf("Responding to 5xx error: %v", msg)
	}
	type errResponse struct {
		Error string `json:"error"`
	}
	respondWithJson(w, code, errResponse{
		Error: msg,
	})
}

func respondWithJson(w http.ResponseWriter, code int, payload interface{}) {
	dat, err := json.Marshal(payload)
	if err != nil {
		log.Printf("failed to marshal json data %v", payload)
		w.WriteHeader(500)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(dat)
}
