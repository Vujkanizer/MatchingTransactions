package services

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

// TransactionResponse je struktura koja odgovara JSON formatu za frontend
type TransactionResponse struct {
	ID              int     `json:"id"`
	Status          string  `json:"status"`
	Confidence      int     `json:"confidence"`
	Value           float64 `json:"value"`
	Name            string  `json:"name"`
	PaymentDate     string  `json:"paymentDate"`
	Address         string  `json:"address"`
	City            string  `json:"city"`
	Country         string  `json:"country"`
	PossibleOrderID string  `json:"possibleOrderId"`
	Ref             string  `json:"ref"`
}

// GetTransactions dohvaća transakcije iz baze i šalje ih kao JSON
func GetTransactions(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	// Ekstrakcija userID iz query parametara
	userID := r.URL.Query().Get("userID")
	if userID == "" {
		http.Error(w, "Missing userID parameter", http.StatusBadRequest)
		return
	}

	// Query za dohvaćanje podataka
	query := `SELECT id, madeOn, amount, reference, currency, partnerName FROM saltdata WHERE userID = ?`
	rows, err := db.Query(query, userID)
	if err != nil {
		log.Printf("Database query error: %v", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var transactions []TransactionResponse

	// Iteracija kroz rezultate
	for rows.Next() {
		var id int
		var madeOn string
		var amount float64
		var reference, currency, partnerName string

		err := rows.Scan(&id, &madeOn, &amount, &reference, &currency, &partnerName)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}

		// Parsiranje datuma u odgovarajući format za frontend
		parsedDate, err := time.Parse("2006-01-02", madeOn)
		if err != nil {
			log.Printf("Error parsing date: %v", err)
			parsedDate = time.Now()
		}
		formattedDate := parsedDate.Format("Jan 02, 2006, 15:04 PM")

		// Kreiranje transakcijskog JSON objekta
		transaction := TransactionResponse{
			ID:              id,
			Status:          "matched", //Sta cemo sve od podataka koristiti
			Confidence:      100,
			Value:           amount,
			Name:            partnerName,
			PaymentDate:     formattedDate,
			Address:         "123 Main St",
			City:            "Ljubljana",
			Country:         "Slovenia",
			PossibleOrderID: "2334",
			Ref:             reference,
		}

		transactions = append(transactions, transaction)
	}

	// Slanje JSON odgovora
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(transactions)
}
