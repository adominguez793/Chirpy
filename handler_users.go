package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}
	type returnVals struct {
		ID    int    `json:"id"`
		Email string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	user, err := cfg.DB.CreateUser(params.Email, params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create user")
		return
	}

	respondWithJSON(w, http.StatusCreated, returnVals{
		ID:    user.ID,
		Email: user.Email,
	})
}

func (cfg *apiConfig) handlerUpdateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	type returnVals struct {
		Email string `json:"email"`
		ID    int    `json:"id"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		respondWithError(w, http.StatusUnauthorized, "Authorization Header is required")
		return
	}
	bearerToken := strings.Split(authHeader, " ")
	if len(bearerToken) != 2 || bearerToken[0] != "Bearer" {
		respondWithError(w, http.StatusUnauthorized, "Invalid token format")
		return
	}
	strToken := bearerToken[1]
	claimsStruct := jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(
		strToken,
		&claimsStruct,
		func(token *jwt.Token) (interface{}, error) { return []byte(cfg.Secret), nil },
	)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Token failed to be validated")
		return
	}
	userIDString, err := token.Claims.GetSubject()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to retrieve subject from token")
		return
	}

	err = cfg.DB.UpdateUser(userIDString, params.Email, params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to update user")
		return
	}

	ID, err := strconv.Atoi(userIDString)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to convert ID from string to int")
		return
	}

	respondWithJSON(w, http.StatusOK, returnVals{
		Email: params.Email,
		ID:    ID,
	})
}
