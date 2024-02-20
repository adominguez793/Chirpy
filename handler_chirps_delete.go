package main

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
)

func (cfg *apiConfig) handlerDeleteChirp(w http.ResponseWriter, r *http.Request) {
	chirpIDString := chi.URLParam(r, "chirpID")
	chirpID, err := strconv.Atoi(chirpIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Failed to convert chirpID")
		return
	}

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

	userIDString, err := token.Claims.GetSubject()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to get subject from token")
		return
	}
	userIDInt, err := strconv.Atoi(userIDString)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to convert userID from string to int")
		return
	}

	err = cfg.DB.DeleteSpecificChirp(chirpID, userIDInt)
	if err != nil {
		// if strErr := fmt.Sprintf("%s\n", err); strErr == "403" {
		// 	fmt.Println("-=-")
		// 	respondWithError(w, 403, "User is not the author of this chirp")
		// 	return
		// }
		respondWithError(w, 403, "Failed to delete chirp")
		return
	}
	respondWithJSON(w, http.StatusOK, struct{}{})
}
