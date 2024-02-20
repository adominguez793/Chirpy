package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

func (cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request) {
	authHead := r.Header.Get("Authorization")
	if authHead == "" {
		respondWithError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}
	bearerToken := strings.Split(authHead, " ")
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
		respondWithError(w, http.StatusUnauthorized, "Failed to validate token")
		return
	}
	userID, err := token.Claims.GetSubject()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to retrieve subject from token")
		return
	}
	userIDInt, err := strconv.Atoi(userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to convert userID from string to int")
		return
	}

	type parameters struct {
		Body string `json:"body"`
	}
	type returnVals struct {
		AuthorID int    `json:"author_id"`
		Body     string `json:"body"`
		ID       int    `json:"id"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	CleanedBody := naughtyChecker(params.Body)

	chirp, err := cfg.DB.CreateChirp(CleanedBody, userIDInt)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create chirp")
		return
	}

	respondWithJSON(w, 201, returnVals{
		AuthorID: userIDInt,
		Body:     chirp.Body,
		ID:       chirp.ID,
	})
}
