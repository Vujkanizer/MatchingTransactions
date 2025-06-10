package services

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type TransactionSalt struct {
	ID            string `json:"id"`
	Amount        string `json:"amount"`
	Currency      string `json:"currency"`
	TransactionID string `json:"transaction_id"`
	Status        string `json:"status"`
	UPNReference  string `json:"upn_reference"`
	Date          string `json:"made_on"`
	PartnerName   string `json:"partner_name"`
}

func MatchOrdersToTransactionsHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	// Extract the userID from the context
	userID, ok := r.Context().Value("userID").(int)
	if !ok || userID == 0 {
		http.Error(w, "User ID not found in context", http.StatusUnauthorized)
		return
	}

	// Call the service function
	err := MatchOrdersToTransactions(r.Context(), db)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error matching orders to transactions: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Matching process completed successfully"))
}

func MatchOrdersToTransactions(ctx context.Context, db *sql.DB) error {
	// Extract the userID from the context
	userID, ok := ctx.Value(UserIDKey).(int)
	if !ok || userID == 0 {
		return fmt.Errorf("user ID not found in context")
	}

	// 1. Retrieve all transactions for the given userID
	rows, err := db.Query(`SELECT id, madeOn, amount, reference, currency, partnerName FROM saltdata WHERE userID = ?`, userID)
	if err != nil {
		return fmt.Errorf("failed to get transactions: %v", err)
	}
	defer rows.Close()

	var transactions []TransactionSalt

	// 2. Iterate through the transactions
	for rows.Next() {
		var transaction TransactionSalt
		err := rows.Scan(&transaction.ID, &transaction.Date, &transaction.Amount, &transaction.Currency, &transaction.TransactionID, &transaction.PartnerName)
		if err != nil {
			return fmt.Errorf("failed to scan transaction: %v", err)
		}
		transactions = append(transactions, transaction)
	}

	// 3. Process each transaction
	for _, transaction := range transactions {
		transDate, err := time.Parse("2006-01-02T15:04:05-07:00", transaction.Date)
		if err != nil {
			return fmt.Errorf("failed to parse transaction date: %v", err)
		}

		// 4. Query to get orders created within the last 3 days of the transaction date for the given userID
		query := `SELECT id, createdAt, totalPrice, status, customer FROM woodata WHERE userID = ? AND DATE(createdAt) BETWEEN DATE_SUB(?, INTERVAL 3 DAY) AND ?`
		orderRows, err := db.Query(query, userID, transDate.Format("2006-01-02"), transDate.Format("2006-01-02"))
		if err != nil {
			return fmt.Errorf("failed to get orders: %v", err)
		}
		defer orderRows.Close()

		var orders []OrderWoo
		for orderRows.Next() {
			var order OrderWoo
			err := orderRows.Scan(&order.ID, &order.CreatedAt, &order.TotalPrice, &order.Status, &order.Customer)
			if err != nil {
				return fmt.Errorf("failed to scan order: %v", err)
			}
			orders = append(orders, order)
		}

		var matched bool

		// Convert partner name to lowercase for comparison
		lowercasePartnerName := strings.ToLower(transaction.PartnerName)

		// 5. Check if the transaction amount matches any order's total price
		for _, order := range orders {
			transAmount, err := strconv.ParseFloat(transaction.Amount, 64)
			if err != nil {
				return fmt.Errorf("failed to parse transaction amount: %v", err)
			}

			orderTotal, err := strconv.ParseFloat(order.TotalPrice, 64)
			if err != nil {
				return fmt.Errorf("failed to parse order total: %v", err)
			}

			if transAmount == orderTotal {
				lowercaseCustomer := strings.ToLower(order.Customer)

				if lowercaseCustomer == lowercasePartnerName || areNamesReversed(lowercaseCustomer, lowercasePartnerName) {
					fmt.Printf("Matched order ID: %v with transaction ID: %v and amount: %v\n", order.ID, transaction.ID, transaction.Amount)
					matched = true

					// Update the woodata table with the saltDataId
					_, err := db.Exec(`UPDATE woodata SET saltDataId = ? WHERE id = ?`, transaction.ID, order.ID) //Testirati ovo!
					if err != nil {
						return fmt.Errorf("failed to update wodata with saltDataId: %v", err)
					}
				}
			}
		}

		// If no matches were found, print all orders and the transaction details
		if !matched {
			fmt.Printf("No matches found for transaction ID: %v. Here are the details:\n", transaction.ID)
			fmt.Printf("Transaction Amount: %v, Partner Name: %s\n", transaction.Amount, transaction.PartnerName)
			fmt.Println("All orders for the same date:")
			for _, order := range orders {
				fmt.Printf("Order ID: %v, Total Price: %v, Customer: %s\n", order.ID, order.TotalPrice, order.Customer)
			}
		}
	}

	return nil
}

func areNamesReversed(name1, name2 string) bool {
	// Split names by space and check if they are reversed
	names1 := strings.Fields(name1)
	names2 := strings.Fields(name2)

	if len(names1) != 2 || len(names2) != 2 {
		return false
	}
	return names1[0] == names2[1] && names1[1] == names2[0]
}
