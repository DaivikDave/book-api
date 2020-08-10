package data

import (
	"database/sql"
	"fmt"
)

//Author Model
type Author struct {
	ID          int64  `json:"-"`
	Name        string `json:"name" validate:"required"`
	Description string `json:"description" validate:"required"`
}

//Create an Author
func (author *Author) Create(db *sql.DB) error {
	res, err := db.Exec("insert into authors(name,description) values(?,?)", author.Name, author.Description)

	if err != nil {
		fmt.Printf("%s", err)
		return err
	}

	author.ID, err = res.LastInsertId()

	return err
}

//Get Author from name
func (author *Author) Get(db *sql.DB, name string) error {
	err := db.QueryRow("select id, name, description from  authors where name=?", name).Scan(&author.ID, &author.Name, &author.Description)
	return err
}
