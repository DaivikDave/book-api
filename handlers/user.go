package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/DaivikDave/book-api/data"
	"github.com/DaivikDave/book-api/util"
	"gopkg.in/go-playground/validator.v9"
)

// User Handler Struct
type User struct {
	DB *sql.DB
}

// Creates UserHandler with DB
func CreateUserHandler(db *sql.DB) *User {
	return &User{
		DB: db,
	}
}

//Handler Function for Creating a User
func (user User) Create(rw http.ResponseWriter, req *http.Request) {
	var u data.User

	//decode json from Request Body
	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&u); err != nil {
		util.RespondWithError(rw, http.StatusBadRequest, "Invalid Request")
		return
	}

	defer req.Body.Close()

	//validate request parameters
	validate := validator.New()
	err := validate.Struct(u)

	if err != nil {
		util.RespondWithError(rw, http.StatusBadRequest, err.Error())
		return
	}

	//create a user record in db
	if err := u.Create(user.DB); err != nil {
		util.RespondWithError(rw, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	//create an access token for created user
	token, err := data.CreateTokenResponse(u.ID)

	if err != nil {
		util.RespondWithError(rw, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	util.RespondWithJSON(rw, http.StatusOK, token)
}

//Login Handler for User
func (user User) Login(rw http.ResponseWriter, req *http.Request) {
	var u data.User
	//decode json in user struct
	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&u); err != nil {
		util.RespondWithError(rw, http.StatusBadRequest, "Invalid Request")
		return
	}

	defer req.Body.Close()

	//validate request parameters
	validate := validator.New()
	err := validate.StructExcept(u, "Name")

	if err != nil {
		util.RespondWithError(rw, http.StatusBadRequest, err.Error())
		return
	}

	if err := u.Get(user.DB); err != nil {
		util.RespondWithError(rw, http.StatusOK, "Invalid credentials")
		return
	}

	token, err := data.CreateTokenResponse(u.ID)

	if err != nil {
		util.RespondWithError(rw, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	util.RespondWithJSON(rw, http.StatusOK, token)
}
