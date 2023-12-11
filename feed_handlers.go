package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/danny-ell77/rssagg/internal/database"
	"github.com/google/uuid"
)


func (apiCfg *apiConfig) createFeed(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		Name string `json:"name"`
		Url string `json:"url"`
	}
	decoder := json.NewDecoder(r.Body)

	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		errorResponder(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	feed, err := apiCfg.DB.CreateFeed(r.Context(), database.CreateFeedParams{
		ID: uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name: params.Name,
		Url: params.Url,
		UserID: user.ID,
	})
	if err != nil {
		errorResponder(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}
	JSONResponder(w, 201, feedSerializer(feed))
}