package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/danny-ell77/rssagg/internal/database"
	"github.com/google/uuid"
)

func (apiCfg *apiConfig) createUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name string `json:"name"`
	}
	decoder := json.NewDecoder(r.Body)

	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		errorResponder(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	user, err := apiCfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID: 		uuid.New(),
		CreatedAt: 	time.Now().UTC(),
		UpdatedAt: 	time.Now().UTC(),
		Name: 		params.Name,
	})
	if err != nil {
		errorResponder(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}
	JSONResponder(w, 201, userSerializer(user))
}

func (apiCfg *apiConfig) getUser(w http.ResponseWriter, r *http.Request, user database.User) {
	JSONResponder(w, 201, userSerializer(user))

}


func (apiCfg *apiConfig) getPostsForUser(w http.ResponseWriter, r *http.Request) {
	var pageSizeInt int
	user := r.Context().Value(AuthUserKey).(database.User)


	pageSize := r.URL.Query().Get("page_size")
	pageSizeInt, err := strconv.Atoi(pageSize)
	if err != nil {
		if pageSize == "" {
			pageSizeInt = 100
		} else {
			errorResponder(w, 400, fmt.Sprintf("Error parsing Query param: %v", err))
			return
		}
	}

	posts, err := apiCfg.DB.GetPostsForUser(r.Context(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit: int32(pageSizeInt),	
	})
	if err != nil {
		errorResponder(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	JSONResponder(w, 200, manyPostsSerializer(posts))

}