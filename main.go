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

// dbConn returns a SQL database connection object.
//
// The returned connection is pinged to verify the connection is valid.
// If the connection is invalid, the function panics.
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

var tmpl = template.Must(template.New("").Funcs(template.FuncMap{
    "add": func(a, b int) int {
        return a + b
    },
}).ParseGlob("templates/*"))


// bookIndex responds to GET requests to "/" and shows all books in the database.
func bookIndex(w http.ResponseWriter, r *http.Request) {
    db := dbConn()
    query := `SELECT id, name, category_id, created_at FROM books;`
    
    rows, err := db.Query(query)

    if err != nil {
        log.Println(err)
        http.Error(w, err.Error(), http.StatusInternalServerError)

        return
    }

    defer rows.Close()

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

    tmpl.ExecuteTemplate(w, "book.html", data)
    defer db.Close()
}

// bookShow responds to GET requests to "/books/{id}" and shows a book with matching id
// from the database.
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

// bookEdit responds to GET requests to "/books/edit/{id}" and shows a book with matching id
// from the database in a form ready to be edited.
func bookEdit(w http.ResponseWriter, r *http.Request) {
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

    tmpl.ExecuteTemplate(w, "edit.html", b)

    defer db.Close()
}

// bookCreate responds to GET requests to "/books/create" and shows a form for creating a new book.
func bookCreate(w http.ResponseWriter, r *http.Request) {
    tmpl.ExecuteTemplate(w, "create.html", BookData{
        PageTitle: "Create Book",
    })
}

// bookStore responds to POST requests to "/books/store" and stores the book in the database.
// It will redirect to "/books" if the book is successfully stored.
func bookStore(w http.ResponseWriter, r *http.Request) {
    var (
        name = r.FormValue("name")
        category_id = r.FormValue("category_id")
        created_at = time.Now().Format("2006-01-02 15:04:05")
    )

    db := dbConn()
    result, err := db.Exec(`INSERT INTO books (name, category_id, created_at) values (?, ?, ?)`, name, category_id, created_at)

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

// bookUpdate responds to POST requests to "/books/update/{id}" and updates the book
// with matching id in the database.
//
// It will redirect to "/books" if the book is successfully updated.
func bookUpdate(w http.ResponseWriter, r *http.Request) {
    var bookId = mux.Vars(r)["id"]

    db := dbConn()
    query := "UPDATE books SET name = ?, category_id = ?, updated_at = ? WHERE id = ?"

    _, err := db.Exec(query, r.FormValue("name"), r.FormValue("category_id"), time.Now().Format("2006-01-02 15:04:05"), bookId)

    if err != nil {
        fmt.Fprintf(w, "Book with %s not found and has some error %s", bookId, err)

        return
    }

    http.Redirect(w, r, "/books", http.StatusSeeOther)
}

// bookDelete responds to POST requests to "/books/delete/{id}" and deletes the book
// with matching id in the database.
//
// It will redirect to "/books" if the book is successfully deleted.
func bookDelete(w http.ResponseWriter, r *http.Request) {
    var bookId = mux.Vars(r)["id"]

    db := dbConn()
    query := "DELETE FROM books WHERE id = ?"

    _, err := db.Exec(query, bookId)

    if err != nil {
        fmt.Fprintf(w, "Book with %s not found and has some error %s", bookId, err)

        return
    }

    http.Redirect(w, r, "/books", http.StatusSeeOther)
}

// main is the main entry point for the application.
// It creates a new router and sets up the routes for the books. It then
// starts the server and listens on port 8000.
func main() {
    defer dbConn().Close()

    r := mux.NewRouter()

    r.PathPrefix("/static/").
        Handler(http.StripPrefix("/static/", 
            http.FileServer(http.Dir("./assets/"))))

    r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        tmpl.ExecuteTemplate(w, "index.html", BookData{
            PageTitle: "Book",
        })
    })
    /**
    * Create a new subrouter for books
    * Define the routes for the books
    */
    r.HandleFunc("/books", bookIndex)
    bookRouter := r.PathPrefix("/books").Subrouter()
    bookRouter.HandleFunc("/create", bookCreate)
    bookRouter.HandleFunc("/show/{id}", bookShow)
    bookRouter.HandleFunc("/edit/{id}", bookEdit)
    bookRouter.HandleFunc("/store", bookStore).Methods("POST")
    bookRouter.HandleFunc("/update/{id}", bookUpdate).Methods("POST")
    bookRouter.HandleFunc("/delete/{id}", bookDelete)

    log.Fatal(http.ListenAndServe(":8000", r))
}