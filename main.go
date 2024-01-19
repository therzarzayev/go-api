package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

// Book model
type Book struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
}

var db *sql.DB

func main() {
	var err error
	db, err = sql.Open("mysql", "root:rza3639@tcp(172.203.217.12:3306)/booksdb")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	router := mux.NewRouter()

	router.HandleFunc("/", getHome).Methods("GET")
	router.HandleFunc("/books", getBooks).Methods("GET")
	router.HandleFunc("/books/{id}", getBook).Methods("GET")
	router.HandleFunc("/books", addBook).Methods("POST")
	router.HandleFunc("/books/{id}", updateBook).Methods("PUT")
	router.HandleFunc("/books/{id}", deleteBook).Methods("DELETE")

	fmt.Println("Server :8001 portunda çalışıyor...")
	log.Fatal(http.ListenAndServe(":8001", router))
}

func getHome(w http.ResponseWriter, r *http.Request) {

}

func getBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var books []Book

	rows, err := db.Query("SELECT * FROM books")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var book Book
		if err := rows.Scan(&book.ID, &book.Title, &book.Author); err != nil {
			log.Fatal(err)
		}
		books = append(books, book)
	}

	if err := json.NewEncoder(w).Encode(books); err != nil {
		log.Fatal(err)
	}
}

func getBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	var book Book

	row := db.QueryRow("SELECT * FROM books WHERE id = ?", params["id"])
	if err := row.Scan(&book.ID, &book.Title, &book.Author); err != nil {
		log.Fatal(err)
	}

	if err := json.NewEncoder(w).Encode(book); err != nil {
		log.Fatal(err)
	}
}

func addBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var book Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		log.Fatal(err)
	}

	result, err := db.Exec("INSERT INTO books (title, author) VALUES (?, ?)", book.Title, book.Author)
	if err != nil {
		log.Fatal(err)
	}

	newID, _ := result.LastInsertId()
	book.ID = int(newID)

	if err := json.NewEncoder(w).Encode(book); err != nil {
		log.Fatal(err)
	}
}

func updateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	var book Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		log.Fatal(err)
	}

	_, err := db.Exec("UPDATE books SET title = ?, author = ? WHERE id = ?", book.Title, book.Author, params["id"])
	if err != nil {
		log.Fatal(err)
	}

	if err := json.NewEncoder(w).Encode(book); err != nil {
		log.Fatal(err)
	}
}

func deleteBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)

	_, err := db.Exec("DELETE FROM books WHERE id = ?", params["id"])
	if err != nil {
		log.Fatal(err)
	}

	w.Write([]byte("Kitap silindi"))
}
