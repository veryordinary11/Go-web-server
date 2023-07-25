package database

type User struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
}

// CreateUser creates a new user and saves it to disk
func (db *DB) CreateUser(email string) (User, error) {
	db.mux.Lock()
	defer db.mux.Unlock()

	// Load the current database into mermory
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	// Find the next available Id for the new user
	maxID := 0
	for _, user := range dbStructure.Users {
		if user.ID > maxID {
			maxID = user.ID
		}
	}
	newID := maxID + 1

	// Create the new user
	newUser := User{
		ID:    newID,
		Email: email,
	}

	// Save the new user to the database
	dbStructure.Users[newID] = newUser

	// Write and updated database back to disk
	err = db.writeDB(dbStructure)
	if err != nil {
		return User{}, err
	}

	return newUser, nil
}
