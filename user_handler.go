package main

import (
	"encoding/json"
	"fmt"
	"net/http"
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
		ID: uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name: params.Name,
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