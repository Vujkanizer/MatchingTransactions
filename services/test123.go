package services

// import (
// 	"encoding/json"
// 	"net/http"

// 	"github.com/go-resty/resty/v2"
// )

// type tSUserKeys struct {
// 	WooKey    string
// 	WooSecret string
// 	WooUrl    string
// }

// // Struktura za meta podatke narudžbine
// type WooMetaData struct {
// 	Key   string `json:"key"`
// 	Value string `json:"value"`
// }

// // Struktura za narudžbine sa meta podacima
// type TestWoo struct {
// 	ID       int           `json:"id"`
// 	Status   string        `json:"status"`
// 	MetaData []WooMetaData `json:"meta_data"`
// }

// // Funkcija za fetchovanje WooCommerce narudžbina
// func testGetWooCommerceOrders(w http.ResponseWriter, r *http.Request) {
// 	// Hardkodovane vrednosti za testiranje (unesi svoje vrednosti)
// 	keys := tSUserKeys{
// 		WooKey:    "your_consumer_key",
// 		WooSecret: "your_consumer_secret",
// 		WooUrl:    "https://yourstore.com/wp-json/wc/v3/orders", // Unesi URL svog WooCommerce API-ja
// 	}

// 	// Kreiranje HTTP REST klijenta
// 	client := resty.New()

// 	// Kreiramo zahtev za fetchovanje narudžbina sa meta podacima
// 	resp, err := client.R().
// 		SetHeader("Accept", "application/json").
// 		SetHeader("Content-Type", "application/json").
// 		SetBasicAuth(keys.WooKey, keys.WooSecret). // WooCommerce koristi Basic Auth sa ključevima
// 		SetResult(&[]TestWoo{}).                   // Definišemo gde se smešta odgovor
// 		Get(keys.WooUrl)

// 	if err != nil {
// 		http.Error(w, "Failed to get WooCommerce orders", http.StatusInternalServerError)
// 		return
// 	}

// 	// Parsiramo odgovor iz WooCommerce API-ja
// 	orders := resp.Result().(*[]TestWoo)

// 	// Prikaz narudžbina u JSON formatu, uključujući meta podatke
// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(orders)
// }
