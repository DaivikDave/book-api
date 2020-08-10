package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/DaivikDave/book-api/handlers"
	"github.com/DaivikDave/book-api/middleware"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {

	// read configuration parameters from the .env file
	err := loadEnvironmentVariables()

	if err != nil {
		log.Fatalf("Error loading .env file")
		return
	}

	//get database connection
	db, err := getDatabase()

	if err != nil {
		log.Fatalf(fmt.Sprintf("Error connecting to the database : %s", err))
	}

	//create instances of handlers
	userHandler := handlers.CreateUserHandler(db)
	authorHandler := handlers.CreateAuthorHandler(db)
	bookHandler := handlers.CreateBookHandler(db)

	router := mux.NewRouter()

	//register routes
	router.HandleFunc("/user/create", userHandler.Create).Methods(http.MethodPost)
	router.HandleFunc("/user/login", userHandler.Login).Methods(http.MethodPost)

	//routes that require users to be authenticated
	authRouter := router.PathPrefix("/").Subrouter()
	authRouter.Use(middleware.AuthenticatedMiddleWare)

	authRouter.HandleFunc("/author/create", authorHandler.Create).Methods(http.MethodPost)

	authRouter.HandleFunc("/book/create", bookHandler.Create).Methods(http.MethodPost)
	authRouter.HandleFunc("/book/{author}", bookHandler.GetByAuthor).Methods(http.MethodGet)
	authRouter.HandleFunc("/book/update", bookHandler.Update).Methods(http.MethodPut)
	authRouter.HandleFunc("/book/delete/{id:[0-9]+}", bookHandler.Delete).Methods(http.MethodDelete)

	srv := &http.Server{
		Handler: router,
		Addr:    ":8000",
	}

	log.Fatal(srv.ListenAndServe())
}

// read configuration variables from the .env file
func loadEnvironmentVariables() error {
	// load .env file
	err := godotenv.Load(".env")

	return err
}

// set up Database Connection
func getDatabase() (*sql.DB, error) {

	//get parameters from database
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	database := os.Getenv("DB_DATABASE")
	username := os.Getenv("DB_USERNAME")
	pass := os.Getenv("DB_PASSWORD")

	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", username, pass, host, port, database)

	db, err := sql.Open("mysql",
		connectionString)

	if err != nil {
		return nil, err
	}

	//check the database connection
	err = db.Ping()

	if err != nil {
		return nil, err
	}

	return db, nil
}
