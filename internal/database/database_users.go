package database

import (
	"errors"
	"fmt"
	"strconv"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID          int    `json:"id"`
	Email       string `json:"email"`
	Password    []byte `json:"password"`
	IsChirpyRed bool   `json:"is_chirpy_red"`
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
		ID:          id,
		Email:       email,
		Password:    encryptedPassword,
		IsChirpyRed: false,
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

func (db *DB) UpdateUser(userID, newEmail, newPassword string) error {
	dbStructure, err := db.loadDB()
	if err != nil {
		return err
	}
	encryptedNewPassword, err := encryptPassword(newPassword)
	if err != nil {
		return err
	}
	userIDInt, err := strconv.Atoi(userID)
	if err != nil {
		return err
	}

	_, existence := dbStructure.Users[userIDInt]
	if existence {
		dbStructure.Users[userIDInt] = User{
			ID:       userIDInt,
			Email:    newEmail,
			Password: encryptedNewPassword,
		}
		err = db.writeDB(dbStructure)
		if err != nil {
			return err
		}
		return nil
	} else {
		fmt.Println("Provided user ID is not valid")
		return errors.New("False ID")
	}
}

func encryptPassword(password string) ([]byte, error) {
	bytePassword := []byte(password)
	encryptedPassword, err := bcrypt.GenerateFromPassword(bytePassword, bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return encryptedPassword, nil
}

func (db *DB) UpgradeChirpyRed(userID int) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	oldUser, ok := dbStructure.Users[userID]
	if !ok {
		return User{}, errors.New("User with this user ID doesn't exist")
	}
	dbStructure.Users[userID] = User{
		ID:          oldUser.ID,
		Email:       oldUser.Email,
		Password:    oldUser.Password,
		IsChirpyRed: true,
	}
	upgradedUser := dbStructure.Users[userID]

	db.writeDB(dbStructure)

	return upgradedUser, nil
}
