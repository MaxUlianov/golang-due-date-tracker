package main

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type userModel struct {
	Id       string
	Username string
	Password string
}

func authUser(db DBInstance, username string, password string) bool {
	var hashedPassword string

	err := db.self.QueryRow("SELECT password FROM app_users WHERE username = $1", username).Scan(&hashedPassword)
	if err != nil {
		return false
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

func createUser(db DBInstance, user userModel) (string, error) {

	result, err := db.self.Exec("INSERT INTO app_users (username, password) VALUES ($1, $2)", user.Username, user.Password)
	if err != nil {
		return "0", fmt.Errorf("addUser: %v", err)
	}
	id, err := result.RowsAffected()
	if err != nil {
		return "0", fmt.Errorf("addUser: %v", err)
	}
	return fmt.Sprintf("%d", id), nil

}

func getUserId(db DBInstance, username string) (string, error) {
	var userId string

	err := db.self.QueryRow("SELECT id FROM app_users WHERE username = $1", username).Scan(&userId)
	if err != nil {
		return "", fmt.Errorf("addUser: %v", err)
	}
	return userId, nil
}
