package services

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"ordermatch/config"

	"golang.org/x/crypto/bcrypt"
)

func init() {
	db = config.DbConnect()
	if db == nil {
		log.Fatalf("Error connecting to database: database connection is nil")
	}
}

type RegisterUser struct {
	Name     string `json:"name"`
	Surname  string `json:"surname"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Username string `json:"username"`
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var ru RegisterUser
	if err := json.NewDecoder(r.Body).Decode(&ru); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid request payload: %v", err)
		return
	}

	salt, err := generateSalt()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error generating salt: %v", err)
		return
	}

	hashedPassword, err := hashPassword(ru.Password, salt)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error hashing password: %v", err)
		return
	}

	err = createUser(db, ru.Name, ru.Surname, ru.Email, hashedPassword, salt, ru.Username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error creating user: %v", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, "User registered successfully")
}

func generateSalt() (string, error) {
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(salt), nil
}

func hashPassword(password, salt string) (string, error) {
	combined := password + salt
	hash, err := bcrypt.GenerateFromPassword([]byte(combined), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func createUser(db *sql.DB, name string, surname string, email string, password string, salt string, username string) error {
	query := "INSERT INTO user (name, surname, email, password, salt, username) VALUES (?, ?, ?, ?, ?, ?)"
	_, err := db.Exec(query, name, surname, email, password, salt, username)
	if err != nil {
		fmt.Printf("Error executing query: %v\n", err)
	}
	return err
}
