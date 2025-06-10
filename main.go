package main

import (
	"log"
	"net/http"
	"ordermatch/config" // Import the config package
	"ordermatch/middleware"
	"ordermatch/services"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	db := config.DbConnect() // Get the database connection

	// Glavni router za sve rute
	mainRouter := mux.NewRouter()

	// Router za rute koje zahtevaju CORS
	corsRouter := mux.NewRouter()
	corsRouter.HandleFunc("/register", services.RegisterHandler).Methods("POST") // Register
	corsRouter.HandleFunc("/login", services.LoginHandler).Methods("POST")       // Login

	// CORS konfiguracija
	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	originsOk := handlers.AllowedOrigins([]string{"http://localhost:5173"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS", "DELETE"})
	corsHandler := handlers.CORS(originsOk, headersOk, methodsOk)(corsRouter)

	// Rute koje zahtevaju autentifikaciju
	mainRouter.Handle("/import", middleware.Authenticate(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { //Upload
		services.ImportJSONData(w, r, db)
	}))).Methods("POST")

	mainRouter.Handle("/woocommerce/orders", middleware.Authenticate(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { //WOO Fetch
		services.GetWooCommerceOrders(w, r, db)
	}))).Methods("GET")

	mainRouter.Handle("/match-orders", middleware.Authenticate(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { //Match
		services.MatchOrdersToTransactionsHandler(w, r, db)
	}))).Methods("GET")

	mainRouter.Handle("/gettransactions", middleware.Authenticate(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { //GetTransactions
		services.GetTransactions(w, r, db)
	}))).Methods("GET")

	mainRouter.Handle("/getmatchedorders", middleware.Authenticate(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { //GetOrderForTransactions
		services.GetOrdersByTransactionID(w, r, db)
	}))).Methods("GET")

	mainRouter.Handle("/savekeys", middleware.Authenticate(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { //GetOrderForTransactions
		services.UpdateWooKeys(w, r, db)
	}))).Methods("POST")

	// Dodajte rute iz corsRouter u mainRouter
	mainRouter.PathPrefix("/register").Handler(corsHandler)
	mainRouter.PathPrefix("/login").Handler(corsHandler)

	// Pokretanje servera
	log.Println("Server started on :8082")
	if err := http.ListenAndServe(":8082", handlers.CombinedLoggingHandler(os.Stdout, mainRouter)); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
