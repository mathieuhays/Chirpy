package database

type DBChirp struct {
	Id       int
	Body     string
	AuthorId int
}

func (db *DB) CreateChirp(body string, authorId int) (DBChirp, error) {
	structure, err := db.loadDB()
	if err != nil {
		return DBChirp{}, err
	}

	newIndex := structure.ChirpAutoIncrement + 1
	structure.Chirps[newIndex] = DBChirp{
		Id:       newIndex,
		Body:     body,
		AuthorId: authorId,
	}
	structure.ChirpAutoIncrement = newIndex
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

func (db *DB) DeleteChirp(id int) error {
	structure, err := db.loadDB()
	if err != nil {
		return err
	}

	delete(structure.Chirps, id)

	return db.writeDB(structure)
}
