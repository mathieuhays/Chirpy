package database

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
)

type DBChirp struct {
	Id   int
	Body string
}

type DB struct {
	path string
	mux  *sync.RWMutex
}

type DBStructure struct {
	Chirps map[int]DBChirp
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

func (db *DB) CreateChirp(body string) (DBChirp, error) {
	structure, err := db.loadDB()
	if err != nil {
		return DBChirp{}, err
	}

	newIndex := len(structure.Chirps) + 1
	structure.Chirps[newIndex] = DBChirp{
		Id:   newIndex,
		Body: body,
	}
	err = db.writeDB(structure)
	if err != nil {
		return DBChirp{}, err
	}

	return structure.Chirps[newIndex], nil
}

func (db *DB) GetChirps() ([]DBChirp, error) {
	var chirps = make([]DBChirp, 0)
	structure, err := db.loadDB()

	if err != nil {
		return chirps, err
	}

	for _, c := range structure.Chirps {
		chirps = append(chirps, c)
	}

	return chirps, nil
}

func (db *DB) GetChirp(id int) (DBChirp, bool) {
	structure, err := db.loadDB()
	if err != nil {
		return DBChirp{}, false
	}

	if val, exists := structure.Chirps[id]; exists {
		return val, true
	}

	return DBChirp{}, false
}

func (db *DB) ensureDB() error {
	_, err := os.ReadFile(db.path)
	if errors.Is(err, os.ErrNotExist) {
		structure := DBStructure{Chirps: make(map[int]DBChirp)}
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
