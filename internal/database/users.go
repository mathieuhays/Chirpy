package database

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrEmailAlreadyExists = errors.New("email already exists")
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

	if _, exists := db.GetUserByEmail(email); exists {
		return DBUser{}, ErrEmailAlreadyExists
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
	structure.UserEmailIndex[email] = newIndex

	err = db.writeDB(structure)
	if err != nil {
		return DBUser{}, err
	}

	return structure.Users[newIndex], nil
}

func (db *DB) GetUser(id int) (DBUser, bool) {
	structure, err := db.loadDB()
	if err != nil {
		return DBUser{}, false
	}

	if user, ok := structure.Users[id]; ok {
		return user, true
	}

	return DBUser{}, false
}

func (db *DB) GetUserByEmail(email string) (DBUser, bool) {
	structure, err := db.loadDB()
	if err != nil {
		return DBUser{}, false
	}

	idx, exists := structure.UserEmailIndex[email]
	if !exists {
		return DBUser{}, false
	}

	if user, ok := structure.Users[idx]; ok {
		return user, true
	}

	return DBUser{}, false
}
