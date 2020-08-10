package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/DaivikDave/book-api/data"
	"github.com/DaivikDave/book-api/util"
	"github.com/gorilla/mux"
	"gopkg.in/go-playground/validator.v9"
)

//Book Handler
type Book struct {
	DB *sql.DB
}

//Create a Book Handler using DB
func CreateBookHandler(db *sql.DB) *Book {
	return &Book{
		DB: db,
	}
}

//Handler Function to Create a User
func (book Book) Create(rw http.ResponseWriter, req *http.Request) {
	var b data.Book

	// decode book from request body
	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&b); err != nil {
		util.RespondWithError(rw, http.StatusBadRequest, "Invalid Request")
		return
	}

	defer req.Body.Close()

	//validate request parameters
	validate := validator.New()
	err := validate.Struct(b)

	if err != nil {
		util.RespondWithError(rw, http.StatusBadRequest, err.Error())
		return
	}

	if err := b.Create(book.DB); err != nil {
		util.RespondWithError(rw, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	util.RespondWithJSON(rw, http.StatusOK, b)
}

func (book Book) Update(rw http.ResponseWriter, req *http.Request) {
	var b data.Book

	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&b); err != nil {
		util.RespondWithError(rw, http.StatusBadRequest, "Invalid Request")
		return
	}

	defer req.Body.Close()

	//validate request parameters
	validate := validator.New()
	err := validate.StructExcept(b, "Author")

	if err != nil {
		util.RespondWithError(rw, http.StatusBadRequest, err.Error())
		return
	}

	if err := b.Update(book.DB); err != nil {
		util.RespondWithError(rw, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	util.RespondWithJSON(rw, http.StatusOK, b)
}

// Handler function to get books by Author Name
func (book Book) GetByAuthor(rw http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)

	books, err := data.GetByAuthor(book.DB, vars["author"])

	if err != nil {
		util.RespondWithError(rw, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	util.RespondWithJSON(rw, http.StatusOK, books)

}

//Handler function to delete book by book ID
func (book Book) Delete(rw http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		util.RespondWithError(rw, http.StatusBadRequest, "Invalid Book ID")
		return
	}

	b := data.Book{
		ID: int64(id),
	}

	err = b.Delete(book.DB)

	if err != nil {
		util.RespondWithError(rw, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	util.RespondWithJSON(rw, http.StatusOK, nil)

}
