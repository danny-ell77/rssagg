package main

import (
	"fmt"
	"net/http"

	"github.com/danny-ell77/rssagg/internal/auth"
	"github.com/danny-ell77/rssagg/internal/database"
)

type authedHandler func(http.ResponseWriter, *http.Request, database.User)


func (apiCfg apiConfig) authMiddleware(handler authedHandler)  http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey, err := auth.GetAPIKey(r.Header)

		if err != nil {
			errorResponder(w, 403, fmt.Sprintf("Authentication failed: %v", err))
		}

		user, err := apiCfg.DB.GetUserByAPIKey(r.Context(), apiKey)

		if err != nil {
			errorResponder(w, 404, fmt.Sprintf("No user for APIKey: %v", err))
		}

		handler(w, r, user)
	}
}