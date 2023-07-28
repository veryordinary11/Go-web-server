package database

import (
	"fmt"
	"strconv"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID          int    `json:"id"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	IsChirpyRed bool   `json:"isChirpyRed"`
}

type UserWithoutPassword struct {
	ID          int    `json:"id"`
	Email       string `json:"email"`
	IsChirpyRed bool   `json:"isChirpyRed"`
}

type UserLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// CreateUser creates a new user and saves it to disk
func (db *DB) CreateUser(email, password string) (UserWithoutPassword, error) {
	db.mux.Lock()
	defer db.mux.Unlock()

	// Load the current database into mermory
	dbStructure, err := db.loadDB()
	if err != nil {
		return UserWithoutPassword{}, err
	}

	// Check if the email address is already in use
	for _, user := range dbStructure.Users {
		if user.Email == email {
			return UserWithoutPassword{}, fmt.Errorf("email address already in use")
		}
	}

	// Find the next available Id for the new user
	maxID := 0
	for _, user := range dbStructure.Users {
		if user.ID > maxID {
			maxID = user.ID
		}
	}
	newID := maxID + 1

	// Hash the password using bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return UserWithoutPassword{}, err
	}

	// Create the new user
	newUser := User{
		ID:          newID,
		Email:       email,
		Password:    string(hashedPassword),
		IsChirpyRed: false,
	}

	// Save the new user to the database
	dbStructure.Users[newID] = newUser

	// Write and updated database back to disk
	err = db.writeDB(dbStructure)
	if err != nil {
		return UserWithoutPassword{}, err
	}

	return UserWithoutPassword{ID: newUser.ID, Email: newUser.Email}, nil
}

// GetUserByEmail returns a user with email address given
func (db *DB) GetUserByEmail(email string) (User, error) {
	db.mux.Lock()
	defer db.mux.Unlock()

	// Load the current database into mermory
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	// Find the user with the given email address
	for _, user := range dbStructure.Users {
		if user.Email == email {
			return user, nil
		}
	}

	return User{}, fmt.Errorf("user not found")
}

// GetUserByID returns a user with ID given
func (db *DB) GetUserByID(userID string) (User, error) {
	db.mux.Lock()
	defer db.mux.Unlock()

	// Load the current database into mermory
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	// Find the user with the given ID
	for _, user := range dbStructure.Users {
		if strconv.Itoa(user.ID) == userID {
			return user, nil
		}
	}

	return User{}, fmt.Errorf("user not found")
}

// UpdateUser updates the user with the given ID, email and password
func (db *DB) UpdateUser(userID, email, password string) (UserWithoutPassword, error) {
	// Find the user by userID in the database
	user, err := db.GetUserByID(userID)
	if err != nil {
		return UserWithoutPassword{}, err
	}

	// Update the user's email and password
	user.Email = email
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return UserWithoutPassword{}, err
	}
	user.Password = string(hashedPassword)

	// Save the updated user to the database
	err = db.saveUser(user)
	if err != nil {
		return UserWithoutPassword{}, err
	}

	updatedUserWithoutPassword := UserWithoutPassword{
		ID:          user.ID,
		Email:       user.Email,
		IsChirpyRed: user.IsChirpyRed,
	}

	return updatedUserWithoutPassword, nil
}

// saveUser saves the user to the database
func (db *DB) saveUser(user User) error {
	db.mux.Lock()
	defer db.mux.Unlock()

	// Load the current database into mermory
	dbStructure, err := db.loadDB()
	if err != nil {
		return err
	}

	// Save the updated user to the database
	dbStructure.Users[user.ID] = user

	// Write and updated database back to disk
	err = db.writeDB(dbStructure)
	if err != nil {
		return err
	}

	return nil
}

// UpdateUserIsChirpyRed updates the user's IsChirpyRed status
func (db *DB) UpdateUserIsChirpyRed(user User) error {
	db.mux.Lock()
	defer db.mux.Unlock()

	// Load the current database into mermory
	dbStructure, err := db.loadDB()
	if err != nil {
		return err
	}

	// Update the user's IsChirpyRed status
	user.IsChirpyRed = true
	dbStructure.Users[user.ID] = user

	// Write and updated database back to disk
	err = db.writeDB(dbStructure)
	if err != nil {
		return err
	}

	return nil
}
