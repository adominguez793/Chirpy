module github.com/adominguez793/Chirpy

go 1.22.0

require github.com/go-chi/chi/v5 v5.0.11

//indirect  <<-- that text was to the right of the line above this one
// and it was causing an error for some reason

require (
	github.com/golang-jwt/jwt/v5 v5.2.0
	github.com/joho/godotenv v1.5.1
	golang.org/x/crypto v0.19.0
)
