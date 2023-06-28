package main

import (
	"fmt"
	"net/http"

	"github.com/sandepten/go-rss-aggregator/internal/auth"
	"github.com/sandepten/go-rss-aggregator/internal/database"
)

type authHandler func(http.ResponseWriter, *http.Request, database.User)

func (apiCfg *apiConfig) middlewareAuth(handler authHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) { // this will create a closure and we can access the inner function by using this middleware
		apiKey, err := auth.GetAPIKey(r.Header)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, fmt.Sprintf("Auth Error: %v", err))
			return
		}
		user, err := apiCfg.DB.GetUserByAPIKey(r.Context(), apiKey)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error getting user: %v", err))
			return
		}

		handler(w, r, user)
	}
}
