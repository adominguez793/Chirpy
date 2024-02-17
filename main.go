package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/adominguez793/Chirpy/internal/database"
	"github.com/go-chi/chi/v5"
)

type apiConfig struct {
	fileserverHits int
	DB             *database.DB
}

func main() {
	const filepathRoot = "."
	const port = "8080"

	db, err := database.NewDB("database.json")
	if err != nil {
		log.Fatal(err)
	}

	apiCfg := apiConfig{
		fileserverHits: 0,
		DB:             db,
	}

	router := chi.NewRouter()

	fsHandler := apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))
	router.Handle("/app/*", fsHandler)
	router.Handle("/app", fsHandler)

	apiRouter := chi.NewRouter()
	apiRouter.Get("/healthz", handlerReadiness)
	apiRouter.Get("/reset", apiCfg.handlerReset)
	apiRouter.Post("/users", apiCfg.handlerCreateUser)
	apiRouter.Post("/chirps", apiCfg.handlerCreateChirp)
	apiRouter.Get("/chirps", apiCfg.handlerGetChirp)
	apiRouter.Get("/chirps/{userID}", apiCfg.handlerGetSpecificChirp)

	adminRouter := chi.NewRouter()
	adminRouter.Get("/metrics", apiCfg.handlerMetrics)

	router.Mount("/api", apiRouter)
	router.Mount("/admin", adminRouter)

	corsMux := middlewareCors(router)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: corsMux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}

// func respondWithError(w http.ResponseWriter, code int, msg string) {
// 	if code > 499 {
// 		log.Printf("Responding with 5XX error: %s", msg)
// 	}
// 	type errorResponse struct {
// 		Error string `json:"error"`
// 	}
// 	respondWithJSON(w, code, errorResponse{
// 		Error: msg,
// 	})
// }

// func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
// 	w.Header().Set("Content-Type", "application/json")
// 	dat, err := json.Marshal(payload)
// 	if err != nil {
// 		log.Printf("Error marshalling JSON: %s", err)
// 		w.WriteHeader(500)
// 		return
// 	}
// 	w.WriteHeader(code)
// 	w.Write(dat)
// }

func naughtyChecker(chirp string) string {
	cleanChirp := strings.ToLower(chirp)
	splitCleanChirp := strings.Split(cleanChirp, " ")

	splitRegularChirp := strings.Split(chirp, " ")

	naughtyWords := []string{"kerfuffle", "sharbert", "fornax"}
	for i := 0; i < len(splitCleanChirp); i++ {
		for _, word := range naughtyWords {
			if splitCleanChirp[i] == word {
				splitRegularChirp[i] = "****"
			}
		}
	}
	profaneFreeChirp := strings.Join(splitRegularChirp, " ")
	return profaneFreeChirp
}
