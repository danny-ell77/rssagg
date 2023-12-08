package main

import "net/http"

func readinessHandler(w http.ResponseWriter, r *http.Request) {
	JSONResponder(w, 200, struct{}{})
}