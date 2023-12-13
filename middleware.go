package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/danny-ell77/rssagg/internal/auth"
	"github.com/danny-ell77/rssagg/internal/database"
)

type authedHandler func(http.ResponseWriter, *http.Request, database.User)


func (apiCfg apiConfig) authMiddleware(next authedHandler)  http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey, err := auth.GetAPIKey(r.Header)

		if err != nil {
			errorResponder(w, 403, fmt.Sprintf("Authentication failed: %v", err))
		}

		user, err := apiCfg.DB.GetUserByAPIKey(r.Context(), apiKey)

		if err != nil {
			errorResponder(w, 404, fmt.Sprintf("No user for APIKey: %v", err))
		}

		next(w, r, user)
	}
}


func loggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next(w, r)
		// Figure out how to get the response status code
		log.Printf("%s %s %v", r.Method, r.URL.Path, time.Since(start))
	}
}

func (apiCfg apiConfig) newAuthMiddleware(next http.HandlerFunc)  http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey, err := auth.GetAPIKey(r.Header)

		if err != nil {
			errorResponder(w, 403, fmt.Sprintf("Authentication failed: %v", err))
			return
		}

		user, err := apiCfg.DB.GetUserByAPIKey(r.Context(), apiKey)

		if err != nil {
			errorResponder(w, 404, fmt.Sprintf("No user for APIKey: %v", err))
			return
		}

    	ctxWithUser := context.WithValue(r.Context(), AuthUserKey, user)
    	rWithUser := r.WithContext(ctxWithUser)

		next(w, rWithUser)
	}
}