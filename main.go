package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

//Product type data sebagai ganti database
type Product struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Quantity int    `json:"quantity"`
}

var (
	database = make(map[string]Product)
)

//SetJSONRes kdjfkjdkfjdkfj
func SetJSONRes(res http.ResponseWriter, message []byte, httpCode int) {
	res.Header().Set("Content-type", "application/json")
	res.WriteHeader(httpCode)
	res.Write(message)
}
func main() {

	//init db
	database["001"] = Product{ID: "001", Name: "John", Quantity: 10}
	database["002"] = Product{ID: "002", Name: "Yang ke dua", Quantity: 4}

	http.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		message := []byte(`{"message": "server is up"}`)
		SetJSONRes(res, message, http.StatusOK)
	})

	http.HandleFunc("/get-products", func(res http.ResponseWriter, req *http.Request) {
		if req.Method != "GET" {
			message := []byte(`{"message": "invalid http method"}`)
			SetJSONRes(res, message, http.StatusMethodNotAllowed)
			return
		}

		var products []Product

		for _, product := range database {
			products = append(products, product)
		}

		productJSON, err := json.Marshal(&products)
		if err != nil {
			message := []byte(`{"message": "error when marshal"}`)
			SetJSONRes(res, message, http.StatusFailedDependency)
			return
		}
		SetJSONRes(res, productJSON, http.StatusOK)
	})

	http.HandleFunc("/add-product", func(res http.ResponseWriter, req *http.Request) {
		if req.Method != "POST" {
			message := []byte(`{"message" :"harus METHOD POST"}`)
			SetJSONRes(res, message, http.StatusMethodNotAllowed)
			return

		}
		var product Product
		payload := req.Body

		defer req.Body.Close()
		err := json.NewDecoder(payload).Decode(&product)
		if err != nil {
			message := []byte(`{"message" :"gagal parsing product"}`)
			SetJSONRes(res, message, http.StatusInternalServerError)
			return
		}
		database[product.ID] = product
		message := []byte(`{"message": "success create product"}`)
		SetJSONRes(res, message, http.StatusCreated)
	})

	http.HandleFunc("/get-product", func(res http.ResponseWriter, req *http.Request) {
		if req.Method != "GET" {
			message := []byte(`{"message": "invalid http method"}`)
			SetJSONRes(res, message, http.StatusMethodNotAllowed)
			return
		}

		if _, ok := req.URL.Query()["id"]; !ok {
			message := []byte(`{"message" :"gagal parsing product"}`)
			SetJSONRes(res, message, http.StatusBadRequest)
			return
		}

		id := req.URL.Query()["id"][0]

		product, ok := database[id]
		if !ok {
			message := []byte(`{"message" :"product tidak ditemukan"}`)
			SetJSONRes(res, message, http.StatusAccepted)
			return
		}

		productJSON, err := json.Marshal(&product)
		if err != nil {
			message := []byte(`{"message" :"error when parsing data"}`)
			SetJSONRes(res, message, http.StatusInternalServerError)
			return
		}

		SetJSONRes(res, productJSON, http.StatusOK)
	})

	err := http.ListenAndServe(":9000", nil)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
