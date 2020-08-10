package data

import (
	"database/sql"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

// User Struct
type User struct {
	ID       int64  `json:"-"`
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"` // encrypted password
}

//Token Response struct
type Token struct {
	AccessToken string `json:"access_token"`
	Expiry      int64  `json:"exp"`
}

// Create a User
func (user *User) Create(db *sql.DB) error {

	if err := user.hashPassword(); err != nil {
		return err
	}

	res, err := db.Exec("insert into users(name,email,password) values(?,?,?)", user.Name, user.Email, user.Password)

	if err != nil {
		return err
	}

	user.ID, err = res.LastInsertId()

	return err
}

//Get a user with matching credentials, returns error if not found
func (user *User) Get(db *sql.DB) error {
	password := user.Password
	err := db.QueryRow("select id, name, email, password from  users where email=? ", user.Email).Scan(&user.ID, &user.Name, &user.Email, &user.Password)
	if err != nil {
		return err
	}

	err = VerifyPassword(user.Password, password)
	return err
}

//Replaces Users password with it's hash
func (user *User) hashPassword() error {
	// Hash password
	hashedPassword, err := Hash(user.Password)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)

	return nil
}

//Create a jwt token after user has authenticated
func CreateTokenResponse(id int64) (*Token, error) {
	atClaims := jwt.MapClaims{}

	//update authorized and user_id claims
	atClaims["authorized"] = true
	atClaims["user_id"] = id

	exp := time.Now().Add(time.Minute * 15).Unix()
	atClaims["exp"] = exp
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)

	//get jwt token from access token
	tokenString, err := at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))

	if err != nil {
		return nil, err
	}

	token := &Token{
		AccessToken: tokenString,
		Expiry:      exp,
	}

	return token, nil
}

//Hash Password using bcrypt algorithm
func Hash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

//Verify Password
func VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
