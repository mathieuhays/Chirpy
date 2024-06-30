package database

type DBChirp struct {
	Id   int
	Body string
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