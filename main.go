package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

// Connect to the database
func dbConn() (db *sql.DB) {
    dbDriver := "mysql"
    dbUser := "root"
    dbPass := ""
    dbName := "go_books_store"
    db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
    
    if err != nil {
        panic(err.Error())
    }

    err = db.Ping()
    if err != nil {
        panic(err.Error())
    }

    return db
}

func bookIndex(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Welcome to the index page")
}

func bookShow(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Welcome to the show page")
}

func bookCreate(w http.ResponseWriter, r *http.Request) {
    name := "Atomic Habbit"
    categoryId := 1
    createdAt := time.Now()

    result, err := dbConn().Exec(`INSERT INTO books (name, category_id, created_at) values (?, ?, ?)`, name, categoryId, createdAt)

    if err != nil {
        fmt.Fprintf(w, "Error , Failed Stored the Book")
    }

    bookId, err := result.LastInsertId()

    if err != nil {
        fmt.Fprintf(w, "Error , Lat Id not found")
    }
    
    fmt.Fprintf(w, "Book Id : %d", bookId)
}

func bookStore(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Welcome to the store page")
}

func bookUpdate(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Welcome to the update page")
}

func bookDelete(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Welcome to the delete page")
}

func main() {
    r := mux.NewRouter()

    r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Welcome to the home page")
    })
    /**
    * Create a new subrouter for books
    * Define the routes for the books
    */
    //bookRouter := r.PathPrefix("/books").Subrouter()
    r.HandleFunc("/books", bookIndex)
    r.HandleFunc("/books/create", bookCreate)
    r.HandleFunc("/books/show/{id}", bookShow)
    r.HandleFunc("/books/store", bookStore).Methods("POST")
    r.HandleFunc("/books/update/{id}", bookUpdate).Methods("PUT")
    r.HandleFunc("/books/delete/{id}", bookDelete).Methods("DELETE")

    log.Fatal(http.ListenAndServe(":8000", r))
}