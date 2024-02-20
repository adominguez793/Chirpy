package main

import (
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		respondWithError(w, http.StatusUnauthorized, "Authorization header is required")
		return
	}
	bearerToken := strings.Split(authHeader, " ")
	if len(bearerToken) != 2 || bearerToken[0] != "Bearer" {
		respondWithError(w, http.StatusUnauthorized, "Invalid token format")
		return
	}
	strToken := bearerToken[1]
	claimsStruct := jwt.RegisteredClaims{}
	refreshToken, err := jwt.ParseWithClaims(
		strToken,
		&claimsStruct,
		func(token *jwt.Token) (interface{}, error) { return []byte(cfg.Secret), nil },
	)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Token failed to be validated")
		return
	}

	issuer, err := refreshToken.Claims.GetIssuer()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to retrieve issuer from token")
		return
	}
	if issuer != "chirpy-refresh" {
		respondWithError(w, http.StatusUnauthorized, "Not a refresh token")
		return
	}

	err = cfg.DB.RevokeToken(strToken)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't revoke session")
		return
	}
	respondWithJSON(w, http.StatusOK, struct{}{})
}
