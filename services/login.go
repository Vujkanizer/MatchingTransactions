package services

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"ordermatch/config"
	"time"

	"github.com/golang-jwt/jwt"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

var db *sql.DB
var secretKey = []byte("secret-key")

func init() {
	db = config.DbConnect()
	if db == nil {
		log.Fatalf("Error connecting to database: database connection is nil")
	}
}

type User struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var u User
	json.NewDecoder(r.Body).Decode(&u)
	fmt.Printf("The user request value %v\n", u)

	valid, userID, err := validateUser(db, u.Email, u.Password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Error validating user")
		return
	}

	if valid {
		tokenString, err := createToken(userID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "Error generating token")
			return
		}

		fmt.Println("Generated token:", tokenString)

		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, tokenString)
		return
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, "Invalid credentials")
	}
}

// func ProtectedHandler(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")
// 	tokenString := r.Header.Get("Authorization")
// 	if tokenString == "" {
// 		w.WriteHeader(http.StatusUnauthorized)
// 		fmt.Fprint(w, "Missing authorization header")
// 		return
// 	}
// 	tokenString = tokenString[len("Bearer "):]

// 	err := VerifyToken(tokenString)
// 	if err != nil {
// 		w.WriteHeader(http.StatusUnauthorized)
// 		fmt.Fprint(w, "Invalid token")
// 		return
// 	}

// 	w.WriteHeader(http.StatusOK)
// 	fmt.Fprint(w, "Welcome to the protected area")
// }

func createToken(userID int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func VerifyToken(tokenString string) (int, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		return 0, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID := int(claims["user_id"].(float64))
		return userID, nil
	}

	return 0, fmt.Errorf("invalid token")
}

func validateUser(db *sql.DB, email, password string) (bool, int, error) {
	var dbEmail, dbPassword, dbSalt string
	var userID int
	query := "SELECT id, email, password, salt FROM user WHERE email = ?"
	err := db.QueryRow(query, email).Scan(&userID, &dbEmail, &dbPassword, &dbSalt)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, 0, nil
		}
		return false, 0, err
	}

	combined := password + dbSalt
	err = bcrypt.CompareHashAndPassword([]byte(dbPassword), []byte(combined))
	if err != nil {
		return false, 0, nil
	}

	return true, userID, nil
}
