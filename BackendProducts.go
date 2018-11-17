package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "secret"
	dbname   = "programmingdb"
)

type Product struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Price       int64  `json:"price"`
	Description string `json:"description"`
	Category    string `json:"category"`
}

var productsMap map[int64]Product
var psqlinfo string

func GetAllProducts(w http.ResponseWriter, r *http.Request) {
	psqlinfo = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlinfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	sqlStatement := `SELECT * FROM products;`
	rows, err := db.Query(sqlStatement)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	var product Product
	var products []Product
	for rows.Next() {
		err = rows.Scan(&product.ID, &product.Name, &product.Price, &product.Description, &product.Category)
		if err != nil {
			panic(err)
		}
		products = append(products, product)
		productsMap[product.ID] = product
	}

	json.NewEncoder(w).Encode(products)
}

func GetAllProductsByCategory(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	psqlinfo = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlinfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	sqlStatement := `SELECT * FROM products WHERE category=$1;`
	rows, err := db.Query(sqlStatement, params["name"])
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	var product Product
	var products []Product
	for rows.Next() {
		err = rows.Scan(&product.ID, &product.Name, &product.Price, &product.Description, &product.Category)
		if err != nil {
			panic(err)
		}
		products = append(products, product)
		productsMap[product.ID] = product
	}

	json.NewEncoder(w).Encode(products)
}

func GetAllProductsByName(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	psqlinfo = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlinfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	sqlStatement := `SELECT * FROM products WHERE name like $1;`
	rows, err := db.Query(sqlStatement, "%"+params["name"]+"%")
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	var product Product
	var products []Product
	for rows.Next() {
		err = rows.Scan(&product.ID, &product.Name, &product.Price, &product.Description, &product.Category)
		if err != nil {
			panic(err)
		}
		products = append(products, product)
		productsMap[product.ID] = product
	}

	json.NewEncoder(w).Encode(products)
}

func GetAProduct(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.ParseInt(params["id"], 10, 64)
	if err != nil {
		log.Fatal(err)
	}
	psqlinfo = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlinfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	sqlStatement := `SELECT * FROM products WHERE id=$1`
	rows, err := db.Query(sqlStatement, id)
	if err != nil {
		panic(err)
	}
	var product Product
	for rows.Next() {
		err = rows.Scan(&product.ID, &product.Name, &product.Price, &product.Description, &product.Category)
		if err != nil {
			panic(err)
		}
	}
	// product := productsMap[int64(id)]
	productsMap[id] = product
	fmt.Printf("id: %d", product.ID)
	json.NewEncoder(w).Encode(product)
}

func CreateAProduct(w http.ResponseWriter, r *http.Request) {
	product := Product{}
	psqlinfo = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		log.Fatal("Error decoding the body: ", err)
	}

	//Process the decoded data
	db, err := sql.Open("postgres", psqlinfo)
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	sqlStatement := `INSERT INTO products VALUES ($1, $2, $3, $4, $5)`
	_, err = db.Exec(sqlStatement, product.ID, product.Name, product.Price, product.Description, product.Category)
	if err != nil {
		panic(err)
	}

	var products []Product
	products = append(products, product)
	productsMap[product.ID] = product

	//Return information of all products, including the new created one
	json.NewEncoder(w).Encode(products)
}

func DeleteAProduct(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var products []Product
	for idx, product := range products {
		id, err := strconv.Atoi(params["id"])
		if err != nil {
			log.Fatal(err)
		}

		if product.ID == int64(id) {
			products = append(products[:idx], products[idx+1:]...)
			delete(productsMap, product.ID)
			break
		}
	}
}

func main() {
	productsMap = make(map[int64]Product)

	router := mux.NewRouter()
	router.HandleFunc("/products", GetAllProducts).Methods("GET")
	router.HandleFunc("/products/category/{name}", GetAllProductsByCategory).Methods("GET")
	router.HandleFunc("/products/names/{name}", GetAllProductsByName).Methods("GET")
	router.HandleFunc("/products/{id}", GetAProduct).Methods("GET")
	router.HandleFunc("/products/{id}", CreateAProduct).Methods("POST")
	router.HandleFunc("/products/{id}", DeleteAProduct).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":9000", router))
}
