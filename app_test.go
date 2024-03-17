package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

var a App

func TestMain(m *testing.M) {
	err := a.Initialise(DbUser, DbPassword, "test")
	if err != nil {
		log.Fatal("Error occurred while initialising the database")
	}
	createDatabase()
	m.Run()
}

func createDatabase() {
	createTableQuery := `CREATE TABLE IF NOT EXISTS products(
		id int NOT NULL AUTO_INCREMENT,
		name varchar(255) NOT NULL,
		quantity int,
		price float(10,7),
		PRIMARY KEY (id)
	);`

	_, err := a.DB.Exec(createTableQuery)

	if err != nil {
		log.Fatal(err)
	}
}

func clearTable() {
	_, err := a.DB.Exec("DELETE FROM products")
	if err != nil {
		log.Println(err)
	}
	_, err = a.DB.Exec("ALTER TABLE products AUTO_INCREMENT=1")
	if err != nil {
		log.Println(err)
	}
}

func addProduct(name string, quantity int, price float64) {
	query := fmt.Sprintf("INSERT INTO products(name, quantity, price) values ('%v', '%v', '%v')", name, quantity, price)
	_, err := a.DB.Exec(query)
	if err != nil {
		log.Println(err)
	}
}

func TestGetProduct(t *testing.T) {
	clearTable()
	addProduct("keyboard", 100, 500)
	request, _ := http.NewRequest("GET", "/products/1", nil)
	response := sendRequest(request)
	checkStatusCode(t, http.StatusOK, response.Code)
}

func checkStatusCode(t *testing.T, expectedStatusCode int, actualStatusCode int) {
	if expectedStatusCode != actualStatusCode {
		t.Errorf("Expected status: %v, Received: %v", expectedStatusCode, actualStatusCode)
	}
}

func sendRequest(request *http.Request) *httptest.ResponseRecorder {
	recorder := httptest.NewRecorder()
	a.Router.ServeHTTP(recorder, request)

	return recorder
}

func TestCreateProduct(t *testing.T) {
	clearTable()
	var product = []byte(`{"name":"chair", "quantity":1, "price": 100}`)
	req, _ := http.NewRequest("POST", "/products", bytes.NewBuffer(product))
	req.Header.Set("Content-Type", "application/json")

	response := sendRequest(req)
	checkStatusCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	_ = json.Unmarshal(response.Body.Bytes(), &m)

	if m["name"] != "chair" {
		t.Errorf("Expected name: %v, Got: %v", "chair", m["chair"])
	}
	if m["quantity"] != 1.0 {
		t.Errorf("Expected name: %v, Got: %v", 1.0, m["quantity"])
	}
	if m["price"] != 100.0 {
		t.Errorf("Expected name: %v, Got: %v", 100.0, m["price"])
	}
}

func TestUpdateProduct(t *testing.T) {
	clearTable()
	addProduct("connector", 10, 10)
	req, _ := http.NewRequest("GET", "/products/1", nil)
	response := sendRequest(req)
	checkStatusCode(t, http.StatusOK, response.Code)

	var oldValue map[string]interface{}
	_ = json.Unmarshal(response.Body.Bytes(), &oldValue)

	var product = []byte(`{"name":"connector", "quantity":1, "price": 10}`)
	req, _ = http.NewRequest("PUT", "/products/1", bytes.NewBuffer(product))
	req.Header.Set("Content-Type", "application/json")
	response = sendRequest(req)
	var newValue map[string]interface{}
	_ = json.Unmarshal(response.Body.Bytes(), &newValue)

	if oldValue["id"] != newValue["id"] {
		t.Errorf("Expected id: %v, Got: %v", oldValue["id"], newValue["id"])
	}

	if oldValue["quantity"] == newValue["quantity"] {
		t.Errorf("Expected id: %v, Got: %v", newValue["quantity"], oldValue["quantity"])
	}

	if oldValue["price"] != newValue["price"] {
		t.Errorf("Expected id: %v, Got: %v", oldValue["price"], newValue["price"])
	}
}

func TestDeleteProduct(t *testing.T) {
	clearTable()
	addProduct("connector", 10, 10)

	req, _ := http.NewRequest("GET", "/products/1", nil)
	response := sendRequest(req)
	checkStatusCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("DELETE", "/products/1", nil)
	response = sendRequest(req)
	checkStatusCode(t, http.StatusNoContent, response.Code)

	req, _ = http.NewRequest("GET", "/products/1", nil)
	response = sendRequest(req)
	checkStatusCode(t, http.StatusNotFound, response.Code)
}
