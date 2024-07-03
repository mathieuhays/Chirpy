package database

import (
	"errors"
	"time"
)

type Session struct {
	Token      string
	Expiration time.Time
	User       int
}

func (db *DB) CreateSession(user int, token string, expiration time.Time) (Session, error) {
	structure, err := db.loadDB()
	if err != nil {
		return Session{}, err
	}

	if _, exists := structure.Sessions[token]; exists {
		return Session{}, errors.New("token already exist")
	}

	structure.Sessions[token] = Session{
		Token:      token,
		Expiration: expiration,
		User:       user,
	}

	err = db.writeDB(structure)
	if err != nil {
		return Session{}, err
	}

	return structure.Sessions[token], nil
}

func (db *DB) GetSession(token string) (Session, bool) {
	structure, err := db.loadDB()
	if err != nil {
		return Session{}, false
	}

	if session, ok := structure.Sessions[token]; ok {
		return session, true
	}

	return Session{}, false
}

func (db *DB) RevokeSession(token string) error {
	structure, err := db.loadDB()
	if err != nil {
		return err
	}

	delete(structure.Sessions, token)

	return db.writeDB(structure)
}
