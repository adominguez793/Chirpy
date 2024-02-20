package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/adominguez793/Chirpy/internal/database"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

type apiConfig struct {
	fileserverHits int
	DB             *database.DB
	Secret         string
}

func main() {
	const filepathRoot = "."
	const port = "8080"

	db, err := database.NewDB("database.json")
	if err != nil {
		log.Fatal(err)
	}

	err = godotenv.Load() // by default, godotenv will look for a file named .env in the current directory
	if err != nil {
		log.Fatalf("error loading .env file: %s\n", err)
	}
	jwtSecret := os.Getenv("JWT_SECRET") // os.Getenv in this case loads the JWT_SECRET var in .env

	apiCfg := apiConfig{
		fileserverHits: 0,
		DB:             db,
		Secret:         jwtSecret,
	}

	router := chi.NewRouter()

	fsHandler := apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))
	router.Handle("/app/*", fsHandler)
	router.Handle("/app", fsHandler)

	apiRouter := chi.NewRouter()
	apiRouter.Get("/healthz", handlerReadiness)
	apiRouter.Get("/reset", apiCfg.handlerReset)
	apiRouter.Post("/users", apiCfg.handlerCreateUser)
	apiRouter.Put("/users", apiCfg.handlerUpdateUser)
	apiRouter.Post("/login", apiCfg.handlerLogin)
	apiRouter.Post("/chirps", apiCfg.handlerCreateChirp)
	apiRouter.Get("/chirps", apiCfg.handlerGetChirp)
	apiRouter.Get("/chirps/{userID}", apiCfg.handlerGetSpecificChirp)
	apiRouter.Post("/refresh", apiCfg.handlerRefresh)
	apiRouter.Post("/revoke", apiCfg.handlerRevoke)
	apiRouter.Delete("/chirps/{chirpID}", apiCfg.handlerDeleteChirp)

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
