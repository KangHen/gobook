package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

/**
* Connection Database
**/
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

/*
* Book Type
*/

type Book struct{
    ID int
    Name string
    CategoryId int
    CreatedAt string
}

type BookData struct {
    PageTitle string
    Books []Book
}

var tmpl = template.Must(template.ParseGlob("templates/*"))

/**
* Action
**/
func bookIndex(w http.ResponseWriter, r *http.Request) {
    db := dbConn()
    query := `SELECT id, name, category_id, created_at FROM books;`
    
    rows, err := db.Query(query)
    defer rows.Close()

    if err != nil {
        log.Println(err)
        http.Error(w, err.Error(), http.StatusInternalServerError)

        return
    }

    books := []Book{}
    
    for rows.Next() {
        book := Book{}
        rows.Scan(&book.ID, &book.Name, &book.CategoryId, &book.CreatedAt)

        books = append(books, book)
    }

    data := BookData{
        PageTitle: "Book",
        Books: books,
    }

    tmpl.ExecuteTemplate(w, "index.html", data)
    defer db.Close()
}

func bookShow(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    bookId := vars["id"]

    db := dbConn()
    query := "SELECT id, name, category_id, created_at FROM books WHERE id = ?"

    var b Book

    err := db.QueryRow(query, bookId).Scan(&b.ID, &b.Name, &b.CategoryId, &b.CreatedAt)

    if err != nil {
        fmt.Fprintf(w, "Book with %s not found and has some error %s", bookId, err)

        return
    }

    fmt.Fprintf(w, "Found Book with name : %s , category : %d , created at : %s and by id : %d", b.Name, b.CategoryId, b.CreatedAt, b.ID)

    defer db.Close()
}

func bookCreate(w http.ResponseWriter, r *http.Request) {
    name := "Kala Senjana"
    categoryId := 1
    createdAt := time.Now()

    db := dbConn()
    result, err := db.Exec(`INSERT INTO books (name, category_id, created_at) values (?, ?, ?)`, name, categoryId, createdAt)

    if err != nil {
        fmt.Fprintf(w, "Error , Failed Stored the Book")
    }

    bookId, err := result.LastInsertId()

    if err != nil {
        fmt.Fprintf(w, "Error , Lat Id not found")
        return
    }
    

    defer db.Close()
    
    if (bookId > 0) {
        http.Redirect(w, r, "/books", http.StatusSeeOther)
    }
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
    defer dbConn().Close()

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