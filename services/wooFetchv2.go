package services

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-resty/resty/v2"
)

type WooMetaData struct {
	ID    int         `json:"id"`
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

type Billing struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

type OrderWoo struct {
	ID         int           `json:"id"`
	Billing    Billing       `json:"billing"`
	CreatedAt  string        `json:"date_created"`
	TotalPrice string        `json:"total"`
	Status     string        `json:"status"`
	MetaData   []WooMetaData `json:"meta_data"`
	Customer   string        `json:"customer"`
}

func GetWooCommerceOrders(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	// Uzmi userID iz konteksta
	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		fmt.Println("Failed to get userID from context")
		http.Error(w, "Failed to get userID from context", http.StatusInternalServerError)
		return
	}

	// Dohvati API podatke za korisnika
	var woKey, woSecret, woStoreUrl string
	err := db.QueryRow(`
        SELECT woKey, woSecret, woStoreUrl 
        FROM user 
        WHERE id = ?`, userID).Scan(&woKey, &woSecret, &woStoreUrl)

	if err != nil {
		fmt.Println("Failed to fetch API credentials:", err)
		http.Error(w, "Failed to fetch API credentials", http.StatusInternalServerError)
		return
	}

	client := resty.New()
	page := 1
	perPage := 100 // Broj narudžbina po stranici
	var allOrders []OrderWoo

	for {
		// Fetch orders per page
		resp, err := client.R().
			SetBasicAuth(woKey, woSecret).
			SetQueryParam("per_page", strconv.Itoa(perPage)).
			SetQueryParam("page", strconv.Itoa(page)).
			SetResult(&[]OrderWoo{}).
			Get(woStoreUrl + "/wp-json/wc/v3/orders")

		if err != nil {
			fmt.Println("Error occurred while making the request:", err)
			http.Error(w, "Failed to get orders", http.StatusInternalServerError)
			return
		}

		if resp.IsError() {
			fmt.Println("Error status code:", resp.StatusCode())
			fmt.Println("Response body:", resp.String())
			http.Error(w, "Failed to get orders: "+resp.Status(), http.StatusInternalServerError)
			return
		}

		orders := resp.Result().(*[]OrderWoo)

		if len(*orders) == 0 {
			break // Prekini ako više nema narudžbina
		}

		allOrders = append(allOrders, *orders...)
		page++ // Idi na sledeću stranicu
	}

	// Sačuvaj narudžbine u bazu
	for _, order := range allOrders {
		createdAtDate := parseDate(order.CreatedAt)

		_, err := db.Exec(`
            INSERT INTO woodata (userID, orderId, orderEmail, createdAt, totalPrice, status, name, surname) 
            VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
			userID, order.ID, order.Billing.Email, createdAtDate, order.TotalPrice, order.Status,
			order.Billing.FirstName, order.Billing.LastName,
		)
		if err != nil {
			fmt.Println("Failed to save order:", err)
			http.Error(w, "Failed to save order", http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(allOrders)
}

// Helper funkcija za dobijanje vrednosti iz meta podataka
// func getMetaDataValue(metaData []WooMetaData, key string) string {
// 	for _, meta := range metaData {
// 		if meta.Key == key {
// 			if val, ok := meta.Value.(string); ok {
// 				return val
// 			}
// 		}
// 	}
// 	return ""
// }

// Helper funkcija za parsiranje datuma iz stringa
func parseDate(dateTimeStr string) string {
	// Pretvori datum u format YYYY-MM-DD koristeći prilagođeni format za WooCommerce (bez vremenske zone)
	parsedTime, err := time.Parse("2006-01-02T15:04:05", dateTimeStr)
	if err != nil {
		fmt.Println("Error parsing date:", err)
		return ""
	}
	// Vrati samo datum u formatu YYYY-MM-DD
	return parsedTime.Format("2006-01-02")
}
