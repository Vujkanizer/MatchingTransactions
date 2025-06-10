package services

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
)

// OrderResponse - Struktura odgovora za pojedinačne narudžbine
type OrderResponse struct {
	ID         int     `json:"id"`
	OrderID    string  `json:"orderId"`
	Email      string  `json:"orderEmail"`
	CreatedAt  string  `json:"createdAt"`
	TotalPrice float64 `json:"totalPrice"`
	Status     string  `json:"status"`
	Name       string  `json:"name"`
	Surname    string  `json:"surname"`
}

// GetOrdersByTransactionID dohvaća sve povezane narudžbine za datu transakciju
func GetOrdersByTransactionID(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	transactionID := r.URL.Query().Get("transactionID")
	if transactionID == "" {
		http.Error(w, "Missing transactionID parameter", http.StatusBadRequest)
		return
	}

	// SQL upit za dohvaćanje narudžbina povezanih sa transakcijom
	query := `SELECT id, orderId, orderEmail, createdAt, totalPrice, status, name, surname FROM woodata WHERE saltDataId = ?`
	rows, err := db.Query(query, transactionID)
	if err != nil {
		log.Printf("Database query error: %v", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var orders []OrderResponse

	// Iteracija kroz rezultate
	for rows.Next() {
		var order OrderResponse
		err := rows.Scan(&order.ID, &order.OrderID, &order.Email, &order.CreatedAt, &order.TotalPrice, &order.Status, &order.Name, &order.Surname)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}

		orders = append(orders, order)
	}

	// Ako nema povezanih narudžbina, vraćamo prazan niz
	if len(orders) == 0 {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode([]OrderResponse{})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orders)
}
