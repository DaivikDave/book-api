package data

import (
	"database/sql"
	"log"
)

// Model for Book
type Book struct {
	ID          int64  `json:"-"`
	Name        string `json:"name" validate:"required"`
	Description string `json:"description" validate:"required"`
	Author      string `json:"author_name" validate:"required"`
}

// Create a Book
func (book Book) Create(db *sql.DB) error {

	var author Author
	err := author.Get(db, book.Author)

	if err != nil {
		author.Name = book.Author
		author.Create(db)
	}

	//Insert in books table
	res, err := db.Exec("insert into books(name,description) values(?,?)", book.Name, book.Description)

	if err != nil {
		return err
	}

	//Get the id of the inserted row
	book.ID, err = res.LastInsertId()

	if err != nil {
		return err
	}

	//make an entry with the corresponding author_id,book_id in author_books
	_, err = db.Exec("insert into author_books(author_id,book_id) values(?,?)", author.ID, book.ID)
	return err
}

//Delete Book by ID
func (book *Book) Delete(db *sql.DB) error {
	_, err := db.Exec("delete from books where id=?", book.ID)
	if err != nil {
		return err
	}
	_, err = db.Exec("delete from author_books where book_id=?", book.ID)

	return err
}

//Get all books by an Author
func GetByAuthor(db *sql.DB, authorName string) (*[]Book, error) {
	books := []Book{}

	var author Author
	err := author.Get(db, authorName)

	if err != nil {
		return &books, nil
	}

	//select books of by author id
	rows, err := db.Query(
		"SELECT books.name,books.description from books INNER JOIN author_books on books.id = author_books.book_id where author_books.author_id =? ",
		author.ID)

	if err != nil {
		log.Printf("%s", err)
		return nil, err
	}

	defer rows.Close()

	//iterate through the result rows and populate the books array
	for rows.Next() {
		var book Book

		if err := rows.Scan(&book.Name, &book.Description); err != nil {

			return nil, err
		}

		books = append(books, book)
	}

	return &books, nil
}

//Update Book description by name
func (book *Book) Update(db *sql.DB) error {
	_, err := db.Exec("UPDATE books SET description=? WHERE name=?", book.Description, book.Name)
	return err
}
