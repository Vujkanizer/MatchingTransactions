package services

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
)

// WooCommerceKeys predstavlja strukturu za unos API ključeva
type WooCommerceKeys struct {
	WoKey      string `json:"woKey"`
	WoSecret   string `json:"woSecret"`
	WoStoreUrl string `json:"woStoreUrl"`
}

// UpdateWooKeys ažurira WooCommerce API podatke u bazi za korisnika
func UpdateWooKeys(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	// Dobijanje userID iz konteksta
	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Parsiranje JSON zahteva
	var keys WooCommerceKeys
	if err := json.NewDecoder(r.Body).Decode(&keys); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Ažuriranje podataka u bazi
	query := `UPDATE user SET woKey = ?, woSecret = ?, woStoreUrl = ? WHERE id = ?`
	_, err := db.Exec(query, keys.WoKey, keys.WoSecret, keys.WoStoreUrl, userID)
	if err != nil {
		log.Printf("Database update error: %v", err)
		http.Error(w, "Failed to update data", http.StatusInternalServerError)
		return
	}

	// Odgovor ka klijentu
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "WooCommerce keys updated successfully"}`))
}
