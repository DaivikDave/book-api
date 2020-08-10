package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/DaivikDave/book-api/util"

	"github.com/dgrijalva/jwt-go"
)

// Middleware to check if a User is Authenticated
func AuthenticatedMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		//get the bearer token
		bearerToken := req.Header.Get("Authorization")

		//split the token type and get the access token
		if len(bearerToken) == 0 || len(strings.Split(bearerToken, " ")) != 2 {
			util.RespondWithError(rw, http.StatusBadRequest, "Invalid Access Token")
			return
		}

		tokenString := strings.Split(bearerToken, " ")[1]

		//parse the token string
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(os.Getenv("ACCESS_SECRET")), nil
		})

		// terminate if token is invalid
		if err != nil {
			util.RespondWithError(rw, http.StatusBadRequest, "Invalid Access Token")
			return
		}
		// forward the request if token is valid
		if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			next.ServeHTTP(rw, req)
			return
		}

		util.RespondWithError(rw, http.StatusUnauthorized, "Unauthorised")
	})

}
