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
			ID:       dbChirp.ID,
			Body:     dbChirp.Body,
			AuthorID: dbChirp.AuthorID,
		})
	}
	trueOrFalse := true
	optionalQuerySort := r.URL.Query().Get("sort")
	if optionalQuerySort == "" || optionalQuerySort == "asc" {
		trueOrFalse = true
	}
	if optionalQuerySort == "desc" {
		trueOrFalse = false
	}

	sortedChirps := SortChirps(chirps, trueOrFalse)

	strOptionalQueryAuthorID := r.URL.Query().Get("author_id")
	if strOptionalQueryAuthorID == "" {
		respondWithJSON(w, http.StatusOK, sortedChirps)
		return
	}
	intOptionalQueryAuthorID, err := strconv.Atoi(strOptionalQueryAuthorID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to convert query from string to int")
		return
	}

	specificAuthorChirps := []database.Chirp{}
	for _, chirp := range sortedChirps {
		if chirp.AuthorID == intOptionalQueryAuthorID {
			specificAuthorChirps = append(specificAuthorChirps, chirp)
		}
	}
	respondWithJSON(w, http.StatusOK, specificAuthorChirps)
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

func SortChirps(chirps []database.Chirp, trueOrFalse bool) []database.Chirp {
	if trueOrFalse == true {
		sort.Slice(chirps, func(i, j int) bool {
			return chirps[i].ID < chirps[j].ID
		})
		return chirps
	} else {
		sort.Slice(chirps, func(i, j int) bool {
			return chirps[i].ID > chirps[j].ID
		})
		return chirps
	}
}
