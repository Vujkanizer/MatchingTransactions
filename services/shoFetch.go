package services

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/go-resty/resty/v2"
)

type Customer struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type Order struct {
	ID                int      `json:"id"`
	Email             string   `json:"email"`
	CreatedAt         string   `json:"created_at"`
	TotalPrice        string   `json:"total_price"`
	FinancialStatus   string   `json:"financial_status"`
	FulfillmentStatus string   `json:"fulfillment_status"`
	Customer          Customer `json:"customer"`
}

type OrdersResponse struct {
	Orders []Order `json:"orders"`
}

type ShoUserKeys struct {
	ShoStoreName string
	ShoApi       string
}

func GetShopifyOrders(w http.ResponseWriter, r *http.Request) {

	db := r.Context().Value("db").(*sql.DB)

	userID, ok := r.Context().Value(UserIDKey).(int)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusUnauthorized)
		return
	}

	getShopifyKeys := func(db *sql.DB, userID int) (ShoUserKeys, error) {
		var keys ShoUserKeys
		query := `SELECT shoStoreName, shoApi FROM user WHERE id = ?`
		err := db.QueryRow(query, userID).Scan(&keys.ShoStoreName, &keys.ShoApi)
		if err != nil {
			return keys, err
		}
		return keys, nil
	}

	keys, err := getShopifyKeys(db, userID)
	if err != nil {
		http.Error(w, "Failed to get user keys", http.StatusInternalServerError)
		return
	}

	client := resty.New()

	resp, err := client.R().
		SetHeader("X-Shopify-Access-Token", keys.ShoApi).
		SetResult(&OrdersResponse{}).
		Get("https://" + keys.ShoStoreName + ".myshopify.com/admin/api/2023-01/orders.json")

	if err != nil {
		http.Error(w, "Failed to get Orders", http.StatusInternalServerError)
		return
	}

	if resp.StatusCode() == http.StatusNotFound {
		http.Error(w, "Orders not found", http.StatusNotFound)
		return
	}

	ordersResponse := resp.Result().(*OrdersResponse)

	insertQuery := `INSERT INTO shodata (user_id, orders) VALUES (?, ?)`
	for _, order := range ordersResponse.Orders {
		orderJSON, err := json.Marshal(order)
		if err != nil {
			http.Error(w, "Failed to marshal order data", http.StatusInternalServerError)
			return
		}
		_, err = db.Exec(insertQuery, userID, orderJSON)
		if err != nil {
			http.Error(w, "Failed to insert order data into database", http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ordersResponse.Orders)
}
