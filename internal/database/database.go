package database

import (
	"encoding/json"
	"errors"
	"github.com/mathieuhays/Chirpy/internal/chirp"
	"os"
	"sync"
)

type DB struct {
	path string
	mux  *sync.RWMutex
}

type DBStructure struct {
	Chirps map[int]chirp.Chirp
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

func (db *DB) CreateChirp(body string) (chirp.Chirp, error) {
	c, err := chirp.NewChirp(body)
	if err != nil {
		return chirp.Chirp{}, err
	}

	structure, err := db.loadDB()
	if err != nil {
		return chirp.Chirp{}, err
	}

	newIndex := len(structure.Chirps) + 1
	c.Id = newIndex
	structure.Chirps[newIndex] = c
	err = db.writeDB(structure)
	if err != nil {
		return chirp.Chirp{}, err
	}

	return c, nil
}

func (db *DB) GetChirps() ([]chirp.Chirp, error) {
	var chirps = make([]chirp.Chirp, 0)
	structure, err := db.loadDB()

	if err != nil {
		return chirps, err
	}

	for _, c := range structure.Chirps {
		chirps = append(chirps, c)
	}

	return chirps, nil
}

func (db *DB) ensureDB() error {
	_, err := os.ReadFile(db.path)
	if errors.Is(err, os.ErrNotExist) {
		structure := DBStructure{Chirps: make(map[int]chirp.Chirp)}
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
