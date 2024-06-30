package database

import "errors"

type DBUser struct {
	Id    int
	Email string
}

func (db *DB) CreateUser(email string) (DBUser, error) {
	structure, err := db.loadDB()
	if err != nil {
		return DBUser{}, err
	}

	if len(email) == 0 {
		return DBUser{}, errors.New("email is required")
	}

	newIndex := len(structure.Users) + 1
	structure.Users[newIndex] = DBUser{
		Id:    newIndex,
		Email: email,
	}

	err = db.writeDB(structure)
	if err != nil {
		return DBUser{}, err
	}

	return structure.Users[newIndex], nil
}
