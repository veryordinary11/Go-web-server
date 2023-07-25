package database

import (
	"encoding/json"
	"os"
	"sync"
)

type DB struct {
	path string
	mux  *sync.RWMutex
}

type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
	Users  map[int]User  `json:"users"`
}

// NewDB creates a new database connection
// and creates the database file if it doesn't exist
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

// ensureDB creates the database file if it doesn't exist
func (db *DB) ensureDB() error {
	_, err := os.Stat(db.path)
	if os.IsNotExist(err) {
		dbStructure := DBStructure{
			Chirps: make(map[int]Chirp),
			Users:  make(map[int]User),
		}

		// Write the empty database to disk
		err := db.writeDB(dbStructure)
		if err != nil {
			return err
		}

	} else if err != nil {
		return err
	}

	// Load the existing database
	dbStructure, err := db.loadDB()
	if err != nil {
		return err
	}

	// Initialize the Users map if it is nil
	if dbStructure.Users == nil {
		dbStructure.Users = make(map[int]User)
	}

	// Write the updated database back to disk
	err = db.writeDB(dbStructure)
	if err != nil {
		return err
	}

	return nil
}

// loadDB loads the database from disk
func (db *DB) loadDB() (DBStructure, error) {
	data, err := os.ReadFile(db.path)
	if err != nil {
		return DBStructure{}, err
	}

	dbStructure := DBStructure{}
	err = json.Unmarshal(data, &dbStructure)
	if err != nil {
		return DBStructure{}, err
	}

	return dbStructure, nil
}

// writeDB writes the database to disk
func (db *DB) writeDB(dbStructure DBStructure) error {
	data, err := json.Marshal(dbStructure)
	if err != nil {
		return err
	}

	err = os.WriteFile(db.path, data, 0644)
	if err != nil {
		return err
	}

	return nil
}
