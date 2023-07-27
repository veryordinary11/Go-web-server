package database

import (
	"fmt"
	"os"
	"sort"
	"strconv"
)

type Chirp struct {
	ID       int    `json:"id"`
	Body     string `json:"body"`
	AuthorID int    `json:"authorId"`
}

// CreateChirp creates a new chirp and saves it to disk
func (db *DB) CreateChirp(userId, body string) (Chirp, error) {
	db.mux.Lock()
	defer db.mux.Unlock()

	// Load the current database into mermory
	dbStructure, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	// Convert the userId string to an int
	userIdInt, err := strconv.Atoi(userId)
	if err != nil {
		return Chirp{}, err
	}

	// Find the next available Id for the new chirp
	maxID := 0
	for _, chirp := range dbStructure.Chirps {
		if chirp.ID > maxID {
			maxID = chirp.ID
		}
	}
	newID := maxID + 1

	// Create the new chrip
	newChirp := Chirp{
		AuthorID: userIdInt,
		ID:       newID,
		Body:     body,
	}

	// Save the new chirp to the database
	dbStructure.Chirps[newID] = newChirp

	// Write and updated database back to disk
	err = db.writeDB(dbStructure)
	if err != nil {
		return Chirp{}, err
	}

	return newChirp, nil
}

// GetChirps returns all chirps in the database
func (db *DB) GetChirps() ([]Chirp, error) {
	db.mux.RLock()
	defer db.mux.RUnlock()

	// Load the current database into mermory
	dbStructure, err := db.loadDB()
	if err != nil {
		return nil, err
	}

	// Convert the map of chirps to a slice and sort by ID
	chirps := make([]Chirp, 0, len(dbStructure.Chirps))
	for _, chirp := range dbStructure.Chirps {
		chirps = append(chirps, chirp)
	}

	sort.Slice(chirps, func(i, j int) bool {
		return chirps[i].ID < chirps[j].ID
	})

	return chirps, nil
}

// GetChirpByID returns a chirp by ID
func (db *DB) GetChirpByID(id int) (Chirp, error) {
	db.mux.RLock()
	defer db.mux.RUnlock()

	// Read the databse file into mermory
	dbStructure, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	// Find the chirp with the matching ID
	chirp, ok := dbStructure.Chirps[id]
	if !ok {
		return Chirp{}, os.ErrNotExist
	}

	return chirp, nil
}

// DeleteChirp deletes a chirp from the database
func (db *DB) DeleteChirp(userId, chirpId string) error {
	db.mux.Lock()
	defer db.mux.Unlock()

	// Convert the chirpId string to an int
	id, err := strconv.Atoi(chirpId)
	if err != nil {
		return err
	}

	// Load the current database into mermory
	dbStructure, err := db.loadDB()
	if err != nil {
		return err
	}

	// Find the chirp with the matching ID
	chirp, ok := dbStructure.Chirps[id]
	if !ok {
		return os.ErrNotExist
	}

	// Check if the user is the author of the chirp
	if strconv.Itoa(chirp.AuthorID) != userId {
		return fmt.Errorf("user is not the author of the chirp")
	}

	// Delete the chirp from the database
	delete(dbStructure.Chirps, id)

	// Write the updated database back to disk
	err = db.writeDB(dbStructure)
	if err != nil {
		return err
	}

	return nil
}
