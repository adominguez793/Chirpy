package database

import (
	"errors"
	"fmt"
)

type Chirp struct {
	ID   int    `json:"id"`
	Body string `json:"body"`
}

// GetChirps returns all chirps in the database
func (db *DB) GetChirps() ([]Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return []Chirp{}, err
	}
	chirps := make([]Chirp, 0, len(dbStructure.Chirps))
	for _, chirp := range dbStructure.Chirps {
		chirps = append(chirps, chirp)
	}
	return chirps, nil
}

func (db *DB) GetSpecificChirp(userID int) (Chirp, error) {
	if userID < 1 {
		fmt.Printf("Requested with a number less than one: %d\n", userID)
		return Chirp{}, errors.New("Bad number")
	}

	dbStructure, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}
	if userID > len(dbStructure.Chirps) {
		fmt.Printf("That chirp doesn't exist")
		return Chirp{}, errors.New("userID is too large")
	}

	chirp := dbStructure.Chirps[userID]
	return chirp, nil
}

func (db *DB) CreateChirp(body string) (Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}
	id := len(dbStructure.Chirps) + 1
	chirp := Chirp{
		ID:   id,
		Body: body,
	}
	dbStructure.Chirps[id] = chirp

	err = db.writeDB(dbStructure)
	if err != nil {
		return Chirp{}, err
	}
	return chirp, nil
}
