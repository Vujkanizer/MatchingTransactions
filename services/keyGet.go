package services

import (
	"database/sql"
	"log"
)

// Struktura za korisničke podatke
type UserKeys struct {
	WoKey       string
	WoSecret    string
	WoStoreUrl  string
	SaltKey     string
	SaltPrivate string
	SaltUrl     string
}

func GetWooCommerceKeys(db *sql.DB, userID int) (UserKeys, error) {
	var keys UserKeys
	query := `SELECT woKey, woSecret, woStoreUrl FROM user WHERE id = ?`
	err := db.QueryRow(query, userID).Scan(&keys.WoKey, &keys.WoSecret, &keys.WoStoreUrl)
	if err != nil {
		log.Println("Error fetching WooCommerce keys:", err)
		return keys, err
	}
	return keys, nil
}

func GetSaltEdgeKeys(db *sql.DB, userID int) (UserKeys, error) {
	var keys UserKeys
	query := `SELECT saltKey, saltPrivate, saltUrl FROM user WHERE id = ?`
	err := db.QueryRow(query, userID).Scan(&keys.SaltKey, &keys.SaltPrivate, &keys.SaltUrl)
	if err != nil {
		log.Println("Error fetching Salt Edge keys:", err)
		return keys, err
	}
	return keys, nil
}
