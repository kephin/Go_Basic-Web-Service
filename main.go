package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type Product struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	Manufacturer   string `json:"manufacturer"`
	Sku            string `json:"sku"`
	Upc            string `json:"upc"`
	PricePerUnit   string `json:"pricePerUnit"`
	QuantityOnHand int    `json:"quantityOnHand"`
}

var productList []Product

func init() {
	productsJSON := `[{
		"id": 1,
		"manufacturer": "Johns-Jenkins",
		"sku": "p5z343vdS",
		"upc": "939581000000",
		"pricePerUnit": "497.45",
		"quantityOnHand": 9703,
		"name": "sticky note"
	}, {
		"id": 2,
		"manufacturer": "Hessel, Schimmel and Feeney",
		"sku": "i7v300kmx",
		"upc": "740979000000",
		"pricePerUnit": "282.29",
		"quantityOnHand": 9217,
		"name": "leg warmers"
	}]`
	err := json.Unmarshal([]byte(productsJSON), &productList)
	if err != nil {
		log.Fatal(err)
	}
}

func getNextID() int {
	highestID := -1
	for _, product := range productList {
		if highestID < product.ID {
			highestID = product.ID
		}
	}
	return highestID + 1
}

func productHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		productJSON, err := json.Marshal(productList)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(productJSON)
	case http.MethodPost:
		var newProduct Product
		bodyBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err = json.Unmarshal(bodyBytes, &newProduct)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if newProduct.ID != 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		newProduct.ID = getNextID()
		productList = append(productList, newProduct)
		w.WriteHeader(http.StatusCreated)
	}
}

func main() {
	http.HandleFunc("/products", productHandler)
	http.ListenAndServe(":5000", nil)
}
