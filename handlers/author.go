package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/DaivikDave/book-api/data"
	"github.com/DaivikDave/book-api/util"
	"gopkg.in/go-playground/validator.v9"
)

//Author Handler Struct
type Author struct {
	DB *sql.DB
}

//Create an Author Handler with the DB
func CreateAuthorHandler(db *sql.DB) *Author {
	return &Author{
		DB: db,
	}
}

//Handler function for creating Author
func (author Author) Create(rw http.ResponseWriter, req *http.Request) {
	var a data.Author

	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&a); err != nil {
		util.RespondWithError(rw, http.StatusBadRequest, "Invalid Request")
		return
	}

	defer req.Body.Close()

	//validate request parameters
	validate := validator.New()
	err := validate.Struct(a)

	if err != nil {
		util.RespondWithError(rw, http.StatusBadRequest, err.Error())
		return
	}

	if err := a.Create(author.DB); err != nil {
		util.RespondWithError(rw, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	util.RespondWithJSON(rw, http.StatusOK, a)
}
