package services

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Transaction struct {
	MadeOn      string `json:"MadeOn"`
	Currency    string `json:"Currency"`
	AmountStr   string `json:"AmountStr"`
	Reference   string `json:"Reference"`
	PartnerName string `json:"PartnerName"`
}

var headerMappings = map[string]map[string]string{
	"BankA": {
		"Datum obdelave":  "MadeOn",
		"Valuta":          "Currency",
		"Znesek v dobro":  "AmountStr",
		"Tuja referenca":  "Reference",
		"Naziv partnerja": "PartnerName",
	},
	// "BankB": {
	// 	"Datum":          "MadeOn",
	// 	"Valuta":         "Currency",
	// 	"Znesek":         "AmountStr",
	// 	"Referenca":      "Reference",
	// 	"Partner":        "PartnerName",
	// },
	// Dodati banke
}

func ImportJSONData(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var rawTransactions []map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&rawTransactions)
	if err != nil {
		http.Error(w, "Failed to parse JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	for _, rawT := range rawTransactions {
		var t Transaction

		// Provera svakog hedera iz mape za sve banke
		for _, mappings := range headerMappings {
			for jsonHeader, structField := range mappings {
				if value, ok := rawT[jsonHeader]; ok {
					switch structField {
					case "MadeOn":
						t.MadeOn = value.(string)
					case "Currency":
						t.Currency = value.(string)
					case "AmountStr":
						t.AmountStr = value.(string)
					case "Reference":
						t.Reference = value.(string)
					case "PartnerName":
						t.PartnerName = value.(string)
					}
				}
			}
		}

		if t.MadeOn == "" || t.Currency == "" || t.AmountStr == "" || t.Reference == "" || t.PartnerName == "" {
			log.Printf("Incomplete transaction data: %+v", t)
			continue
		}

		amountStr := strings.ReplaceAll(t.AmountStr, ",", "")
		amount, err := strconv.ParseFloat(amountStr, 64)
		if err != nil {
			log.Printf("Invalid amount value: %s, error: %v", t.AmountStr, err)
			continue
		}

		madeOn := strings.ReplaceAll(t.MadeOn, ".", "-")
		madeOnTime, err := time.Parse("02-01-2006", madeOn) // Format: DD-MM-YYYY
		if err != nil {
			log.Printf("Invalid date value: %s, error: %v", t.MadeOn, err)
			continue
		}

		// Unos podataka u bazu
		_, err = db.Exec(`INSERT INTO SaltData (amount, currency, reference, madeOn, partnerName) VALUES (?, ?, ?, ?, ?)`,
			amount, t.Currency, t.Reference, madeOnTime, t.PartnerName)
		if err != nil {
			log.Printf("Failed to insert data into database: %v", err)
			http.Error(w, "Database insert error", http.StatusInternalServerError)
			return
		}
	}

	fmt.Fprintln(w, "Data imported successfully!")
}
