package store

import (
	"fmt"
	"io"
	"library/domain/model"
	"os"
)

// BookRepository ...
type BookRepository struct {
	Store *Store
}

// SaveBook to local store and to DB
func (r *BookRepository) SaveBook(book *model.Book, bookSrc io.Reader) error {

	if err := os.MkdirAll(fmt.Sprintf("store/books/%s/%s/", book.SubDirection.SubDirection, book.Diffuculty), os.ModePerm); err != nil {
		r.Store.Logger.Errorf("Error occured while creating directory for file. Err msg: %v.", err)
		return err
	}

	bookFile, err := os.Create(fmt.Sprintf("store/books/%s/%s/%s%s", book.SubDirection.SubDirection, book.Diffuculty, book.Author.FullName, book.Title))
	if err != nil {
		r.Store.Logger.Errorf("Error occured while creating file. Err msg: %v.", err)
		return err
	}

	defer bookFile.Close()

	_, err = io.Copy(bookFile, bookSrc)
	if err != nil {
		r.Store.Logger.Errorf("Error occured while writing book file. Err msg: %v.", err)
		return err
	}

	if err := r.Store.Db.QueryRow(
		`INSERT INTO book
		(author_id, sub_direction_id, title, edition_date, diffuculty, rating, description, language, URL, DowloadCount)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id`,
		book.Author.ID,
		book.SubDirection.ID,
		book.Title,
		book.EditionDate,
		string(book.Diffuculty),
		book.Rating,
		book.Description,
		book.Language,
		book.URL,
		book.DownloadCount,
	).Scan(&book.ID); err != nil {
		r.Store.Logger.Errorf("Error occured while saving books data in db. Err msg: %v.", err)
		return err
	}

	return nil
}
