package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/danny-ell77/rssagg/internal/database"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (apiCfg apiConfig) followFeed(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		FeedID uuid.UUID `json:"feed_id"`
	}
	decoder := json.NewDecoder(r.Body)

	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		errorResponder(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	feed_follow, err := apiCfg.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		ID: uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		FeedID: params.FeedID,
		UserID: user.ID,
	})
	if err != nil {
		errorResponder(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}
	JSONResponder(w, 201, feedFollowSerializer(feed_follow))
}

func (apiCfg apiConfig) getFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {
	feed_follows, err := apiCfg.DB.GetFeedFollows(r.Context(), user.ID)
	if err != nil {
		errorResponder(w, 400, fmt.Sprintf("Could not get Feed Follows: %v", err))
		return
	}
	JSONResponder(w, 201, manyFeedFollowSerializer(feed_follows))
}

func (apiCfg apiConfig) unfollowFeed(w http.ResponseWriter, r *http.Request, user database.User) {
	feedFollowIdStr := chi.URLParam(r, "feedFollowId")

	feedFollowId, err := uuid.Parse(feedFollowIdStr)
	if err != nil {
		errorResponder(w, 400, fmt.Sprintf("Could not parse URL param: %v", err))
		return
	}

	err = apiCfg.DB.DeleteFeedFollows(r.Context(), database.DeleteFeedFollowsParams{
		ID: feedFollowId,
		UserID: user.ID,
	})
	if err != nil {
		errorResponder(w, 400, fmt.Sprintf("Could not delete : %v", err))
		return
	}
	JSONResponder(w, 204, struct{}{})

}