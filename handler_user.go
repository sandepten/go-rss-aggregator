package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/sandepten/go-rss-aggregator/internal/database"
)

// we can use this apiConfig struct to pass the database connection to our handlers as a method
func (apiCfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	decoder := json.NewDecoder(r.Body)

	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	user, err := apiCfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		Email:     params.Email,
		Password:  params.Password,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error creating user: %v", err))
		return
	}

	respondWithJSON(w, http.StatusCreated, databaseUserToUser(user))
}

func (apiCfg *apiConfig) handlerGetUserByEmail(w http.ResponseWriter, r *http.Request) {
	userEmail, contextErr := r.Context().Value("email").(string)
	if contextErr {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error getting user email from context"))
		return
	}
	user, err := apiCfg.DB.GetUserByEmail(r.Context(), userEmail)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error getting user: %v", err))
		return
	}

	respondWithJSON(w, http.StatusOK, databaseUserToUser(user))
}

func (apiCfg *apiConfig) handlerGetUserByAPIKey(w http.ResponseWriter, r *http.Request, user database.User) {
	respondWithJSON(w, http.StatusOK, databaseUserToUser(user))
}

func (apiCfg *apiConfig) handlerGetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := apiCfg.DB.GetAllUsers(r.Context())
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error getting user: %v", err))
		return
	}
	respondWithJSON(w, http.StatusOK, users)
}
