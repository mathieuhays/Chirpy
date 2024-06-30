package database

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
)

type DB struct {
	path string
	mux  *sync.RWMutex
}

type DBStructure struct {
	Chirps         map[int]DBChirp
	Users          map[int]DBUser
	UserEmailIndex map[string]int
}

func NewDB(path string) (*DB, error) {
	db := &DB{
		path: path,
		mux:  &sync.RWMutex{},
	}

	err := db.ensureDB()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (db *DB) ensureDB() error {
	_, err := os.ReadFile(db.path)
	if errors.Is(err, os.ErrNotExist) {
		structure := DBStructure{
			Chirps:         make(map[int]DBChirp),
			Users:          make(map[int]DBUser),
			UserEmailIndex: map[string]int{},
		}
		return db.writeDB(structure)
	} else if err != nil {
		return err
	}

	return nil
}

func (db *DB) loadDB() (structure DBStructure, loadErr error) {
	db.mux.RLock()
	defer db.mux.RUnlock()

	data, err := os.ReadFile(db.path)
	if err != nil {
		loadErr = err
		return
	}

	err = json.Unmarshal(data, &structure)
	if err != nil {
		loadErr = err
		return
	}

	return
}

func (db *DB) writeDB(structure DBStructure) error {
	db.mux.Lock()
	defer db.mux.Unlock()

	data, err := json.Marshal(structure)
	if err != nil {
		return err
	}

	return os.WriteFile(db.path, data, 0666)
}
