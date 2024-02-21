package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func (cfg *apiConfig) handlerPolka(w http.ResponseWriter, r *http.Request) {
	type parameter struct {
		Event string `json:"event"`
		Data  struct {
			UserID int `json:"user_id"`
		}
	}

	keyHead := r.Header.Get("Authorization")
	if keyHead == "" {
		respondWithError(w, http.StatusUnauthorized, "Authorization header is required")
		return
	}
	header := strings.Split(keyHead, " ")
	if len(header) != 2 || header[0] != "ApiKey" {
		respondWithError(w, http.StatusUnauthorized, "Invalid ApiKey header format")
		return
	}

	key := header[1]
	if cfg.PolkaKey != key {
		respondWithError(w, http.StatusUnauthorized, "Invalid ApiKey")
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := parameter{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to decode parameter")
		return
	}

	if params.Event != "user.upgraded" {
		respondWithJSON(w, http.StatusOK, struct{}{})
		return
	}

	_, err = cfg.DB.UpgradeChirpyRed(params.Data.UserID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to upgrade user to chirpy red")
		return
	}
	respondWithJSON(w, http.StatusOK, struct{}{})
}
