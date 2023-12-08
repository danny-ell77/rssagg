package main

import "net/http"

func errorHandler(w http.ResponseWriter, r *http.Request) {
	errorResponder(w, 400, "Something went wrong")
}