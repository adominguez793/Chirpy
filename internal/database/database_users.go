package database

import (
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password []byte `json:"password"`
}

func (db *DB) CreateUser(email, password string) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}
	encryptedPassword, err := encryptPassword(password)
	if err != nil {
		fmt.Printf("Failed to encrypt password: %s\n", err)
		return User{}, err
	}

	for _, userInfo := range dbStructure.Users {
		if userInfo.Email == email {
			fmt.Printf("Email (((%s))) is already taken", email)
			return User{}, errors.New("Taken email")
		}
	}

	id := len(dbStructure.Users) + 1
	user := User{
		ID:       id,
		Email:    email,
		Password: encryptedPassword,
	}
	dbStructure.Users[id] = user

	err = db.writeDB(dbStructure)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func (db *DB) GetUserByEmail(email string) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	for _, user := range dbStructure.Users {
		if user.Email == email {
			return user, nil
		}
	}
	fmt.Println("Provided email does not belong to any user.")
	return User{}, errors.New("Bad email")
}

func encryptPassword(password string) ([]byte, error) {
	bytePassword := []byte(password)
	encryptedPassword, err := bcrypt.GenerateFromPassword(bytePassword, bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return encryptedPassword, nil
}
