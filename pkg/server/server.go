package server

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"strconv"

	"api-golang-kubernetes/pkg/log"

	"github.com/gorilla/mux"
)

// Book Struct (Model)
type Book struct {
	ID     string  `json:"id"`
	Isbn   string  `json:"isbn"`
	Title  string  `json:"title"`
	Author *Author `json:"author"`
}

// Author struct
type Author struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

// Init books var as a slice Book struct
var books []Book

// CreatBook is the api to create a Item of book
func CreatBook(res http.ResponseWriter, req *http.Request) {
	log.D("Write the book information")

	res.Header().Set("Content-Type", "application/json")
	var book Book
	_ = json.NewDecoder(req.Body).Decode(&book)
	book.ID = strconv.Itoa(rand.Intn(100000000)) // Mock ID
	books = append(books, book)
	json.NewEncoder(res).Encode(book)

	// for debugging
	jsondata, err := json.Marshal(book)
	if err != nil {
		log.E("Cannot encode to Json", err)
	}
	log.D("%v", string(jsondata))
}

// GetBooks is the api to get all Items of book
func GetBooks(w http.ResponseWriter, r *http.Request) {
	log.D("Get all list of books")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)

	// for debugging
	jsondata, err := json.Marshal(books)
	if err != nil {
		log.E("Cannot encode to Json", err)
	}
	log.D("%v", string(jsondata))
}

// GetBook is the api to get an item of book
func GetBook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r) // get params

	log.D("Get the infomation of book for id:%v", params["id"])

	// Loop through books and find with id
	for _, item := range books {
		if item.ID == params["id"] {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(item)

			// for debugging
			jsondata, err := json.Marshal(item)
			if err != nil {
				log.E("Cannot encode to Json", err)
			}
			log.D("%v", string(jsondata))

			return
		}
	}

	// if the id is not in the database, the answer is "404 Not Found"
	w.WriteHeader(404)
	log.D("Fail to get the book information since it is invalid (" + params["id"] + ")")
}

// UpdateBook is the api to update an item of book
func UpdateBook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r) // get params
	log.D("Update the book information of %v", params["id"])

	for index, item := range books {
		if item.ID == params["id"] {
			w.Header().Set("Content-Type", "application/json")
			books = append(books[:index], books[index+1:]...)

			// same with createBook()
			var book Book
			_ = json.NewDecoder(r.Body).Decode(&book)
			book.ID = params["id"]
			books = append(books, book)
			json.NewEncoder(w).Encode(book)

			// for debugging
			jsondata, err := json.Marshal(book)
			if err != nil {
				log.E("Cannot encode to Json", err)
			}
			log.E("%v", string(jsondata))

			return
		}
	}

	w.WriteHeader(400)
	log.D("Fail to update the book information since it is invalid (" + params["id"] + ")")
}

// DeleteBook is the api to remove an item of book
func DeleteBook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r) // get params
	log.D("Delete the book information of %v", params["id"])

	isDeleted := false
	for index, item := range books {
		if item.ID == params["id"] {
			w.Header().Set("Content-Type", "application/json")
			books = append(books[:index], books[index+1:]...)
			isDeleted = true

			log.D("Successfully deleted the information of Book for id (" + params["id"] + ")")
			break
		}
	}

	if isDeleted == false {
		w.WriteHeader(404)
		log.D("Fail to delete the book information since it is invalid (" + params["id"] + ")")
	}
}

// InitServer initializes the REST api server
func InitServer() error {
	// Init Router
	r := mux.NewRouter()

	// Mock Data - @todo - implement DB
	books = append(books, Book{ID: "1", Isbn: "448743", Title: "Book One", Author: &Author{Firstname: "John", Lastname: "Park"}})
	books = append(books, Book{ID: "2", Isbn: "312234", Title: "Book Two", Author: &Author{Firstname: "Steave", Lastname: "Smith"}})

	// Route Handler / Endpoints
	r.HandleFunc("/api/books", CreatBook).Methods("POST")
	r.HandleFunc("/api/books", GetBooks).Methods("GET")
	r.HandleFunc("/api/books/{id}", GetBook).Methods("GET")
	r.HandleFunc("/api/books/{id}", UpdateBook).Methods("PUT")
	r.HandleFunc("/api/books/{id}", DeleteBook).Methods("DELETE")

	err := http.ListenAndServe(":8000", r)

	return err
}
