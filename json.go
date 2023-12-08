package main

import (
	"encoding/json"
	"log"
	"net/http"
) 

func errorResponder(w http.ResponseWriter, code int, msg string) {
	if code < 499 {
		log.Printf("Responding with %v error: %s", code, msg)
	}
	type errResponse struct {
		Error string `json:"error"`
	}

	JSONResponder(w, code, errResponse{
		Error: msg,
	})
}

func JSONResponder(w http.ResponseWriter, code int, payload interface{}) {
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Failed to marsal JSON response  %v", payload)
		w.WriteHeader(500)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data)
}