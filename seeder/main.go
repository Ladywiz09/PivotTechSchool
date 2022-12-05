package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"log"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

type product struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Price       int    `json:"price"`
	Description string `json:"description"`
}

var products []product

func getProductsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatal(err)
	}

	p := productlookup(id)
	if p == nil {
		log.Printf("Product with id %d not found", id)
		return
	}
	if err := json.NewEncoder(w).Encode(p); err != nil {
		log.Println("Error encoding product to json", err)
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatal(err)
	}
}

func productlookup(id int) *product {
	for _, p := range products {
		if p.Id == id {
			return &p
		}
	}
	return nil
}

func deleteProductsHandler(w http.ResponseWriter, r *http.Request) {
	os.Remove("products.db")
	w.WriteHeader(http.StatusOK)
}

func addProductHandler(w http.ResponseWriter, r *http.Request) {
	var p product
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatal(err)
	}
	products = append(products, p)
	w.WriteHeader(http.StatusCreated)
}

func readProducts() {
	bs, err := os.ReadFile("products.json")
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(bs, &products)
	if err != nil {
		log.Fatal(err)
	}
}

func main() error {
	var productsDB string
	var jsonProducts string

	flag.StringVar(&productsDB, "db", "products.db", "Database file")
	flag.StringVar(&jsonProducts, "json", "products.json", "JSON file")
	flag.Parse()

	db, err := sql.Open("sqlite3", productsDB)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	prod, err := os.ReadFile(jsonProducts)
	if err := db; err != nil {
		log.Fatal(err)
	}

	var payload product

	if err := json.Unmarshal(prod, &payload); err != nil {
		return err
		}
	return nil
}
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare("INSERT INTO products(id, name, price, description) VALUES(?, ?, ?, ?)")
	if err != nil {
		return err
	}

	for _, p := range prod {
		_, err = stmt.Exec(p.Id, p.Name, p.Price, p.Description)
		if err != nil {
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}

	return nil

	rows, err := db.Query("SELECT id, name, price, description FROM products")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var name string
		var price int
		var description string
		err = rows.Scan(&id, &name, &price, &description)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(id, name, price, description)

		err = rows.Err()
		if err != nil {
			log.Fatal(err)
		}
	}

