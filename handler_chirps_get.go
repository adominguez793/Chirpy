package main

import (
	"net/http"
	"sort"
	"strconv"

	"github.com/adominguez793/Chirpy/internal/database"
	"github.com/go-chi/chi/v5"
)

func (cfg *apiConfig) handlerGetChirp(w http.ResponseWriter, r *http.Request) {
	dbChirps, err := cfg.DB.GetChirps()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve chirps")
		return
	}
	chirps := []database.Chirp{}
	for _, dbChirp := range dbChirps {
		chirps = append(chirps, database.Chirp{
			ID:   dbChirp.ID,
			Body: dbChirp.Body,
		})
	}
	sort.Slice(chirps, func(i, j int) bool {
		return chirps[i].ID < chirps[j].ID
	})
	respondWithJSON(w, http.StatusOK, chirps)
}

func (cfg *apiConfig) handlerGetSpecificChirp(w http.ResponseWriter, r *http.Request) {
	userIDString := chi.URLParam(r, "userID")
	userID, err := strconv.Atoi(userIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Failed to use GET command with an int URL Parameter")
		return
	}

	dbChirp, err := cfg.DB.GetSpecificChirp(userID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't get Chirp")
		return
	}

	respondWithJSON(w, http.StatusOK, database.Chirp{
		Body: dbChirp.Body,
		ID:   dbChirp.ID,
	})
}
