package services

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

// HashAndSave hashira API ključ in skrivnost ter ju shrani v bazo podatkov
func HashAndSave(apiKey string, apiSecret *string) error {
	// Hashira API ključ
	hashedKey, err := bcrypt.GenerateFromPassword([]byte(apiKey), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("neuspešno hashiranje API ključa: %v", err)
	}

	var hashedSecret []byte
	if apiSecret != nil {
		// Hashira API skrivnost
		hashedSecret, err = bcrypt.GenerateFromPassword([]byte(*apiSecret), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("neuspešno hashiranje API skrivnosti: %v", err)
		}
	}

	// Odpre povezavo z bazo podatkov
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_NAME"),
	)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("neuspešna povezava z bazo podatkov: %v", err)
	}
	defer db.Close()

	// Ping bazo podatkov za preverjanje povezave
	if err := db.Ping(); err != nil {
		return fmt.Errorf("neuspešen ping na bazo podatkov: %v", err)
	}

	var stmt *sql.Stmt

	// Pripravi SQL stavek samo z API ključem
	stmt, err = db.Prepare("INSERT INTO api_credentials (api_key) VALUES (?)")
	if err != nil {
		return fmt.Errorf("neuspešna priprava SQL stavka: %v", err)
	}
	defer stmt.Close()

	// Izvede SQL stavek
	_, err = stmt.Exec(hashedKey)
	if err != nil {
		return fmt.Errorf("neuspešna izvedba SQL stavka: %v", err)
	}

	if apiSecret != nil {
		// Pripravi SQL stavek z API ključem in skrivnostjo
		stmt, err = db.Prepare("INSERT INTO api_credentials (api_key, api_secret) VALUES (?, ?)")
		if err != nil {
			return fmt.Errorf("neuspešna priprava SQL stavka: %v", err)
		}
		defer stmt.Close()

		// Izvede SQL stavek
		_, err = stmt.Exec(hashedKey, hashedSecret)
		if err != nil {
			return fmt.Errorf("neuspešna izvedba SQL stavka: %v", err)
		}
	}

	return nil
}
