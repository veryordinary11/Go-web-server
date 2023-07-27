package database

func (db *DB) RevokedRefreshToken(refreshToken string) error {
	db.mux.Lock()
	defer db.mux.Unlock()

	// Load the current database into mermory
	dbStructure, err := db.loadDB()
	if err != nil {
		return err
	}

	// If there is no revokedTokens slice, create one
	if dbStructure.RevokedTokens == nil {
		dbStructure.RevokedTokens = make([]string, 0)
	}

	// Add the refresh token to the revokedTokens slice
	dbStructure.RevokedTokens = append(dbStructure.RevokedTokens, refreshToken)

	// Write the updated database to disk
	err = db.writeDB(dbStructure)
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) IsRefreshTokenRevoked(refreshToken string) (bool, error) {
	db.mux.Lock()
	defer db.mux.Unlock()

	// Load the current database into mermory
	dbStructure, err := db.loadDB()
	if err != nil {
		return false, err
	}

	// Check if the refresh token is in the revokedTokens slice
	for _, revokedToken := range dbStructure.RevokedTokens {
		if revokedToken == refreshToken {
			return true, nil
		}
	}

	return false, nil
}
