package main

import (
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	type returnVals struct {
		Token string `json:"token"`
	}

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
	token, err := jwt.ParseWithClaims(
		strToken,
		&claimsStruct,
		func(token *jwt.Token) (interface{}, error) { return []byte(cfg.Secret), nil },
	)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Token failed to be validated")
		return
	}
	issuer, err := token.Claims.GetIssuer()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to retrieve issuer from token")
		return
	}
	if issuer != "chirpy-refresh" {
		respondWithError(w, http.StatusUnauthorized, "Token is not a refresh token")
		return
	}

	checker, err := cfg.DB.IsTokenRevoked(strToken)
	if err != nil || checker == true {
		respondWithError(w, http.StatusUnauthorized, "There are revocations for this token in the database")
		return
	}

	subject, err := token.Claims.GetSubject()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to retrieve subject from token")
		return
	}

	newAccessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "chirpy-access",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(time.Hour)),
		Subject:   subject,
	})
	signedNewAccessToken, err := newAccessToken.SignedString([]byte(cfg.Secret))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to sign token")
	}

	respondWithJSON(w, http.StatusOK, returnVals{
		Token: signedNewAccessToken,
	})

}
