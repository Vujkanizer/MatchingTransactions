package services

type TransactionsResponseSalt struct {
	Transactions []TransactionSalt `json:"transactions"`
}

type SUserKeys struct {
	SaltKey     string
	SaltPrivate string
	SaltUrl     string
}

type contextKey string

const UserIDKey contextKey = "userID"

// func GetSaltEdgeTransactions(w http.ResponseWriter, r *http.Request) {

// 	db := r.Context().Value("db").(*sql.DB)

// 	userID, ok := r.Context().Value(UserIDKey).(int)
// 	if !ok {
// 		http.Error(w, "User ID not found in context", http.StatusUnauthorized)
// 		return
// 	}

// 	getSaltEdgeKeys := func(db *sql.DB, userID int) (SUserKeys, error) {
// 		var keys SUserKeys
// 		query := `SELECT saltKey, saltPrivate, saltUrl FROM user WHERE id = ?`
// 		err := db.QueryRow(query, userID).Scan(&keys.SaltKey, &keys.SaltPrivate, &keys.SaltUrl)
// 		if err != nil {
// 			return keys, err
// 		}
// 		return keys, nil
// 	}

// 	keys, err := getSaltEdgeKeys(db, userID)
// 	if err != nil {
// 		http.Error(w, "Failed to get user keys", http.StatusInternalServerError)
// 		return
// 	}

// 	client := resty.New()

// 	resp, err := client.R().
// 		SetHeader("Accept", "application/json").
// 		SetHeader("Content-Type", "application/json").
// 		SetHeader("App-id", keys.SaltKey).
// 		SetHeader("Secret", keys.SaltPrivate).
// 		SetResult(&TransactionsResponseSalt{}).
// 		Get(keys.SaltUrl)

// 	if err != nil {
// 		http.Error(w, "Failed to get transactions", http.StatusInternalServerError)
// 		return
// 	}

// 	transactions := resp.Result().(*TransactionsResponseSalt)

// 	tx, err := db.Begin()
// 	if err != nil {
// 		http.Error(w, "Failed to start database transaction", http.StatusInternalServerError)
// 		return
// 	}

// 	stmt, err := tx.Prepare(`INSERT INTO saltdata (user_id, transactions) VALUES (?, ?)`)
// 	if err != nil {
// 		tx.Rollback()
// 		http.Error(w, "Failed to prepare SQL statement", http.StatusInternalServerError)
// 		return
// 	}
// 	defer stmt.Close()

// 	for _, transaction := range transactions.Transactions {

// 		transactionJSON, err := json.Marshal(transaction)
// 		if err != nil {
// 			tx.Rollback()
// 			http.Error(w, "Failed to encode transaction to JSON", http.StatusInternalServerError)
// 			return
// 		}

// 		_, err = stmt.Exec(userID, transactionJSON)
// 		if err != nil {
// 			tx.Rollback()
// 			http.Error(w, "Failed to save transaction to the database", http.StatusInternalServerError)
// 			return
// 		}
// 	}

// 	err = tx.Commit()
// 	if err != nil {
// 		tx.Rollback()
// 		http.Error(w, "Failed to commit transaction", http.StatusInternalServerError)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(transactions)
// }
