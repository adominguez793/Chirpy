package database

import "time"

type Revocation struct {
	Token     string    `json:"token"`
	RevokedAt time.Time `json:"revoked_at"`
}

func (db *DB) RevokeToken(token string) error {
	dbStructure, err := db.loadDB()
	if err != nil {
		return err
	}
	dbStructure.Revocations[token] = Revocation{
		Token:     token,
		RevokedAt: time.Now().UTC(),
	}

	err = db.writeDB(dbStructure)
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) IsTokenRevoked(token string) (bool, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return false, err
	}
	revocation, ok := dbStructure.Revocations[token]
	if !ok {
		return false, nil
	}
	if revocation.RevokedAt.IsZero() {
		return false, nil
	}
	return true, nil
}
