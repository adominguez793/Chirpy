package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/adominguez793/Chirpy/internal/database"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email               string `json:"email"`
		Password            string `json:"password"`
		ExpirationInSeconds *int   `json:"expires_in_seconds"`
	}
	type returnVals struct {
		database.User
		// ID           int    `json:"id"`
		// Email        string `json:"email"`
		// Token        string `json:"token"`
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
		// IsChirpyRed  bool   `json:"is_chirpy_red"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	user, err := cfg.DB.GetUserByEmail(params.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Bad email")
	}

	err = bcrypt.CompareHashAndPassword(user.Password, []byte(params.Password))
	if err != nil {
		respondWithError(w, 401, "Unauthorized: Passwords do not match")
	}

	var expiration int
	if params.ExpirationInSeconds != nil {
		expiration = *params.ExpirationInSeconds
	}

	if expiration == 0 || expiration > 86400 {
		expiration = 86400
	}

	strID := strconv.Itoa(user.ID)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "chirpy-access",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(time.Hour)),
		Subject:   strID,
	})
	signedToken, err := token.SignedString([]byte(cfg.Secret))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to sign token")
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "chirpy-refresh",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(time.Duration(1440) * time.Hour)),
		Subject:   strID,
	})
	signedRefreshToken, err := refreshToken.SignedString([]byte(cfg.Secret))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to sign refresh token")
	}

	respondWithJSON(w, http.StatusOK, returnVals{
		User: database.User{
			ID:          user.ID,
			Email:       user.Email,
			IsChirpyRed: user.IsChirpyRed,
		},
		Token:        signedToken,
		RefreshToken: signedRefreshToken,
	})
}
