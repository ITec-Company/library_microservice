package postgres

import (
	"database/sql"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/lib/pq"
	"library-go/internal/domain"
	"library-go/internal/store"
	"library-go/pkg/logging"
	"strings"
)

const (
	getOneBookQuery = `SELECT 
		B.uuid,
		B.title,
		B.difficulty,
		B.edition_date,
		B.rating,
		B.description,
		B.url,
		B.language,
		B.download_count,
		Au.uuid,
		Au.full_name,
		D.uuid as direction_uuid,
		D.name as direction_name,
		array_agg(DISTINCT T) as tags
	FROM book AS B
	lEFT JOIN author AS Au ON Au.uuid = B.author_uuid
	LEFT JOIN direction AS D ON D.uuid = B.direction_uuid
	LEFT JOIN tag AS T ON  T.uuid = any (B.tags_uuids)
	WHERE  B.uuid = $1
	GROUP BY B.uuid, B.title, B.difficulty, B.edition_date, B.rating, B.description, B.url, B.language, B.download_count, Au.uuid, Au.full_name, D.uuid, D.name`
	getAllBooksQuery = `SELECT 
		B.uuid,
		B.title,
		B.difficulty,
		B.edition_date,
		B.rating,
		B.description,
		B.url,
		B.language,
		B.download_count,
		Au.uuid,
		Au.full_name,
		D.uuid as direction_uuid,
		D.name as direction_name,
		array_agg(DISTINCT T) as tags
	FROM book AS B
	lEFT JOIN author AS Au ON Au.uuid = B.author_uuid
	LEFT JOIN direction AS D ON D.uuid = B.direction_uuid
	LEFT JOIN tag AS T ON  T.uuid = any (B.tags_uuids)
	GROUP BY B.uuid, B.title, B.difficulty, B.edition_date, B.rating, B.description, B.url, B.language, B.download_count, Au.uuid, Au.full_name, D.uuid, D.name`
	createBookQuery = `INSERT INTO book (
                     title, 
                     direction_uuid, 
                     author_uuid,
                  	 difficulty,
                     edition_date,
                     rating, 
                     description,
                     url, 
                     language, 
                     tags_uuids, 
                     download_count
				) SELECT $1, $2 , $3, $4, $5, $6, $7, $8, $9, $10, $11
				WHERE EXISTS(SELECT uuid FROM author where $3 = author.uuid) AND
				EXISTS(SELECT uuid FROM direction where $2 = direction.uuid) AND
			    EXISTS(SELECT uuid FROM tag where tag.uuid = any($10)) RETURNING book.uuid`
	deleteBookQuery = `DELETE FROM book WHERE uuid = $1`
	updateBookQuery = `UPDATE book SET 
			title = COALESCE(NULLIF($1, ''), title),  
			direction_uuid = (CASE WHEN (EXISTS(SELECT uuid FROM direction where direction.uuid = $2)) THEN $2 ELSE COALESCE(NULLIF($2, 0), direction_uuid) END), 
			author_uuid = (CASE WHEN (EXISTS(SELECT uuid FROM author where author.uuid = $3)) THEN $3 ELSE COALESCE(NULLIF($3, 0), author_uuid) END), 
			difficulty = COALESCE($4, difficulty), 
			edition_date = COALESCE($5, edition_date), 
			rating = COALESCE(NULLIF($6, 0), rating), 
			description = COALESCE(NULLIF($7, ''), description), 
			url = COALESCE(NULLIF($8, ''), url), 
			language = COALESCE(NULLIF($9, ''), language), 
			tags_uuids = (CASE WHEN (EXISTS(SELECT uuid FROM tag where tag.uuid = any($10))) THEN $10 ELSE COALESCE($10, tags_uuids) END),
			download_count = COALESCE(NULLIF($11, 0), download_count)
		WHERE uuid = $12`
)

type bookStorage struct {
	logger *logging.Logger
	db     *sql.DB
}

func NewBookStorage(db *sql.DB, logger *logging.Logger) store.BookStorage {
	return &bookStorage{
		logger: logger,
		db:     db,
	}
}

func (bs *bookStorage) GetOne(UUID string) (*domain.Book, error) {
	var book domain.Book
	var tagsStr []string
	if err := bs.db.QueryRow(getOneBookQuery,
		UUID).Scan(
		&book.UUID,
		&book.Title,
		&book.Difficulty,
		&book.EditionDate,
		&book.Rating,
		&book.Description,
		&book.URL,
		&book.Language,
		&book.DownloadCount,
		&book.Author.UUID,
		&book.Author.FullName,
		&book.Direction.UUID,
		&book.Direction.Name,
		pq.Array(&tagsStr),
	); err != nil {
		bs.logger.Errorf("error occurred while selecting book from DB. err: %v", err)
		return nil, err
	}
	for _, t := range tagsStr {
		t = strings.Replace(t, "(", "", -1)
		t = strings.Replace(t, ")", "", -1)
		data := strings.Split(t, ",")
		var tag domain.Tag
		tag.UUID = data[0]
		tag.Name = data[1]
		book.Tags = append(book.Tags, tag)
	}
	return &book, nil
}

func (bs *bookStorage) GetAll(sortOptions *domain.SortFilterPagination) ([]*domain.Book, error) {
	s := squirrel.Select("B.uuid, B.title, B.difficulty, B.edition_date, B.rating, B.description, B.url, B.language, B.download_count, Au.uuid, Au.full_name, D.uuid as direction_uuid, D.name as direction_name, array_agg(DISTINCT T) as tags").
		From("book AS B").
		LeftJoin("author AS Au ON Au.uuid = B.author_uuid").
		LeftJoin("direction AS D ON D.uuid = B.direction_uuid").
		LeftJoin("tag AS T ON  T.uuid = any (B.tags_uuids)").
		GroupBy("B.uuid, B.title, B.difficulty, B.edition_date, B.rating, B.description, B.url, B.language, B.download_count, Au.uuid, Au.full_name, D.uuid, D.name")

	if sortOptions.Limit != 0 {
		s = s.Limit(sortOptions.Limit)
		if sortOptions.Page != 0 {
			offset := (sortOptions.Page - 1) * sortOptions.Limit
			s = s.Offset(offset)
		}
	}

	if sortOptions.FiltersAndArgs != nil {
		s = s.Where(sortOptions.FiltersAndArgs).PlaceholderFormat(squirrel.Dollar)
	}

	if sortOptions.SortBy != "" {
		s = s.OrderBy(fmt.Sprintf("%s %s", sortOptions.SortBy, sortOptions.Order))
	}

	query, args, _ := s.ToSql()

	rows, err := bs.db.Query(query, args...)
	if err != nil {
		bs.logger.Errorf("error occurred while selecting all books. err: %v", err)
		return nil, err
	}
	var books []*domain.Book

	for rows.Next() {
		book := domain.Book{}
		var tagsStr []string
		err := rows.Scan(
			&book.UUID,
			&book.Title,
			&book.Difficulty,
			&book.EditionDate,
			&book.Rating,
			&book.Description,
			&book.URL,
			&book.Language,
			&book.DownloadCount,
			&book.Author.UUID,
			&book.Author.FullName,
			&book.Direction.UUID,
			&book.Direction.Name,
			pq.Array(&tagsStr),
		)
		if err != nil {
			bs.logger.Errorf("error occurred while selecting book. err: %v", err)
			continue
		}

		for _, t := range tagsStr {
			t = strings.Replace(t, "(", "", -1)
			t = strings.Replace(t, ")", "", -1)
			data := strings.Split(t, ",")
			var tag domain.Tag
			tag.UUID = data[0]
			tag.Name = data[1]
			book.Tags = append(book.Tags, tag)
		}

		books = append(books, &book)
	}
	return books, nil
}

func (bs *bookStorage) Create(bookCreateDTO *domain.CreateBookDTO) (string, error) {
	tx, err := bs.db.Begin()
	if err != nil {
		bs.logger.Errorf("error occurred while creating transaction. err: %v", err)
		return "", err
	}

	var UUID string

	row := tx.QueryRow(createBookQuery,
		bookCreateDTO.Title,
		bookCreateDTO.DirectionUUID,
		bookCreateDTO.AuthorUUID,
		bookCreateDTO.Difficulty,
		bookCreateDTO.EditionDate,
		0,
		bookCreateDTO.Description,
		bookCreateDTO.URL,
		bookCreateDTO.Language,
		pq.Array(bookCreateDTO.TagsUUIDs),
		0,
	)

	if err := row.Scan(&UUID); err != nil {
		tx.Rollback()
		bs.logger.Errorf("error occurred while creating book. err: %v", err)
		return UUID, err
	}

	return UUID, tx.Commit()
}

func (bs *bookStorage) Delete(UUID string) error {
	tx, err := bs.db.Begin()
	if err != nil {
		bs.logger.Errorf("error occurred while creating transaction. err: %v", err)
		return err
	}

	result, err := tx.Exec(deleteBookQuery, UUID)
	if err != nil {
		tx.Rollback()
		bs.logger.Errorf("error occurred while deleting book. err: %v.", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		tx.Rollback()
		bs.logger.Errorf("error occurred while deleting book (getting affected rows). err: %v", err)
		return err
	}

	if rowsAffected < 1 {
		tx.Rollback()
		bs.logger.Errorf("No book with UUID %s was found", UUID)
		return ErrNoRowsAffected
	}
	bs.logger.Infof("Book with uuid %s was deleted.", UUID)
	return tx.Commit()
}

func (bs *bookStorage) Update(bookUpdateDTO *domain.UpdateBookDTO) error {
	tx, err := bs.db.Begin()
	if err != nil {
		bs.logger.Errorf("error occurred while creating transaction. err: %v", err)
		return err
	}

	result, err := tx.Exec(updateBookQuery,
		bookUpdateDTO.Title,
		bookUpdateDTO.DirectionUUID,
		bookUpdateDTO.AuthorUUID,
		bookUpdateDTO.Difficulty,
		bookUpdateDTO.EditionDate,
		bookUpdateDTO.Rating,
		bookUpdateDTO.Description,
		bookUpdateDTO.URL,
		bookUpdateDTO.Language,
		pq.Array(bookUpdateDTO.TagsUUIDs),
		bookUpdateDTO.DownloadCount,
		bookUpdateDTO.UUID,
	)
	if err != nil {
		tx.Rollback()
		bs.logger.Errorf("error occurred while updating book. err: %v", err)
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		tx.Rollback()
		bs.logger.Errorf("error occurred while updating book (getting affected rows). err: %v", err)
		return err
	}
	if rowsAffected < 1 {
		tx.Rollback()
		bs.logger.Errorf("error occurred while updating book. err: %v.", ErrNoRowsAffected)
		return ErrNoRowsAffected
	}

	bs.logger.Infof("book with uuid %s was updated.", bookUpdateDTO.UUID)

	return tx.Commit()
}
