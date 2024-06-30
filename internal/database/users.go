package database

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
)

type DBUser struct {
	Id       int
	Email    string
	Password string
}

func (db *DB) CreateUser(email string, password string) (DBUser, error) {
	structure, err := db.loadDB()
	if err != nil {
		return DBUser{}, err
	}

	if len(email) == 0 {
		return DBUser{}, errors.New("email is required")
	}

	if len(password) == 0 {
		return DBUser{}, errors.New("password is required")
	}

	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return DBUser{}, err
	}

	newIndex := len(structure.Users) + 1
	structure.Users[newIndex] = DBUser{
		Id:       newIndex,
		Email:    email,
		Password: string(encryptedPassword),
	}

	err = db.writeDB(structure)
	if err != nil {
		return DBUser{}, err
	}

	return structure.Users[newIndex], nil
}
