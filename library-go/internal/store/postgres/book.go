package postgres

import (
	"database/sql"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/lib/pq"
	"library-go/internal/domain"
	"library-go/internal/store"
	"library-go/pkg/logging"
	"math"
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
		B.local_url,
		B.language,
		B.download_count,
		B.image_url,
		B.created_at,
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
	GROUP BY B.uuid, B.title, B.difficulty, B.edition_date, B.rating, B.description, B.local_url, B.language, B.download_count, B.image_url, B.created_at, Au.uuid, Au.full_name, D.uuid, D.name`

	getAllBooksQuery = `SELECT 
		B.uuid,
		B.title,
		B.difficulty,
		B.edition_date,
		B.rating,
		B.description,
		B.local_url,
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
	GROUP BY B.uuid, B.title, B.difficulty, B.edition_date, B.rating, B.description, B.local_url, B.language, B.download_count, Au.uuid, Au.full_name, D.uuid, D.name`

	createBookQuery = `INSERT INTO book (
                     title, 
                     direction_uuid, 
                     author_uuid,
                     difficulty,
                     edition_date,
                     description,
                     local_url, 
                     language, 
                     tags_uuids, 
                     image_url
				) SELECT 
				      $1, 
				      $2, 
				      $3, 
				      $4, 
				      $5, 
				      $6,  
				      $7 || (SELECT last_value from book_uuid_seq) || $8, 
				      $9, 
				      $10, 
				      $11 || (SELECT last_value from book_uuid_seq) || $12
				WHERE EXISTS(SELECT uuid FROM author where $3 = author.uuid) AND
				EXISTS(SELECT uuid FROM direction where $2 = direction.uuid) AND
			    EXISTS(SELECT uuid FROM tag where tag.uuid = any($10)) RETURNING book.uuid`

	deleteBookQuery = `DELETE FROM book WHERE uuid = $1`

	updateBookQuery = `UPDATE book SET 
			title = COALESCE(NULLIF($1, ''), title),  
			direction_uuid = (CASE WHEN (EXISTS(SELECT uuid FROM direction where direction.uuid = $2)) THEN $2 ELSE direction_uuid END), 
			author_uuid = (CASE WHEN (EXISTS(SELECT uuid FROM author where author.uuid = $3)) THEN $3 ELSE author_uuid END), 
			difficulty = (CASE WHEN ($4 = any(enum_range(difficulty))) THEN $4 ELSE difficulty END), 
			edition_date = (CASE WHEN ($5 != date('0001-01-01 00:00:00')) THEN $5 ELSE edition_date END),
			description = COALESCE(NULLIF($6, ''), description), 
			local_url = COALESCE(NULLIF($7, ''), local_url), 
			language = COALESCE(NULLIF($8, ''), language), 
			tags_uuids = (CASE WHEN (EXISTS(SELECT uuid FROM tag where tag.uuid = any($9))) THEN $9 ELSE COALESCE($9, tags_uuids) END)
		WHERE uuid = $10`

	rateBookQuery = `WITH grades AS (
   		 SELECT avg((select avg(a) from unnest(array_append(all_grades, $1)) as a)) AS avg
   		 FROM book
		)
		UPDATE book SET
    	    all_grades = (CASE WHEN (0.0 < $1 AND $1 < 5.1) THEN array_append(all_grades, $1) ELSE all_grades END),
    	    rating = (CASE WHEN (0.0 < $1 AND $1 < 5.1) THEN grades.avg  ELSE rating END)
		FROM grades
		WHERE uuid = $2`

	bookDownloadCountUpQuery = `UPDATE book SET
			download_count = (download_count + 1)
			WHERE uuid = $1`
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
	query, args, _ := squirrel.Select(
		"B.uuid",
		"B.title",
		"B.difficulty",
		"B.edition_date",
		"B.rating",
		"B.description",
		"B.local_url",
		"B.image_url",
		"B.language",
		"B.download_count",
		"B.created_at",
		"Au.uuid",
		"Au.full_name",
		"D.uuid as direction_uuid",
		"D.name as direction_name",
		"array_agg(DISTINCT T) as tags").
		From("book AS B").
		Where("B.uuid = ?", UUID).
		PlaceholderFormat(squirrel.Dollar).
		LeftJoin("author AS Au ON Au.uuid = B.author_uuid").
		LeftJoin("direction AS D ON D.uuid = B.direction_uuid").
		LeftJoin("tag AS T ON  T.uuid = any (B.tags_uuids)").
		GroupBy("B.uuid, B.title, B.difficulty, B.edition_date, B.rating, B.description, B.local_url, B.image_url, B.language, B.download_count, B.created_at, Au.uuid, Au.full_name, D.uuid, D.name").
		ToSql()

	var book domain.Book
	var tagsStr []string
	if err := bs.db.QueryRow(query, args...).Scan(
		&book.UUID,
		&book.Title,
		&book.Difficulty,
		&book.EditionDate,
		&book.Rating,
		&book.Description,
		&book.LocalURL,
		&book.ImageURL,
		&book.Language,
		&book.DownloadCount,
		&book.CreatedAt,
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

func (bs *bookStorage) GetAll(sortOptions *domain.SortFilterPagination) ([]*domain.Book, int, error) {
	s := squirrel.Select(
		"B.uuid",
		"B.title",
		"B.difficulty",
		"B.edition_date",
		"B.rating",
		"B.description",
		"B.local_url",
		"B.image_url",
		"B.language",
		"B.download_count",
		"B.created_at",
		"Au.uuid",
		"Au.full_name",
		"D.uuid as direction_uuid",
		"D.name as direction_name",
		"array_agg(DISTINCT T) as tags",
		"count(*) OVER() AS full_count").
		From("book AS B").
		LeftJoin("author AS Au ON Au.uuid = B.author_uuid").
		LeftJoin("direction AS D ON D.uuid = B.direction_uuid").
		LeftJoin("tag AS T ON  T.uuid = any (B.tags_uuids)").
		GroupBy("B.uuid, B.title, B.difficulty, B.edition_date, B.rating, B.description, B.local_url, B.image_url, B.language, B.download_count, B.created_at, Au.uuid, Au.full_name, D.uuid, D.name")

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
		return nil, 0, err
	}

	var books []*domain.Book
	var fullCount int

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
			&book.LocalURL,
			&book.ImageURL,
			&book.Language,
			&book.DownloadCount,
			&book.CreatedAt,
			&book.Author.UUID,
			&book.Author.FullName,
			&book.Direction.UUID,
			&book.Direction.Name,
			pq.Array(&tagsStr),
			&fullCount,
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

	var pagesCount int

	if sortOptions.Limit != 0 {
		pagesCount = int(math.Ceil(float64(fullCount) / float64(sortOptions.Limit)))
	}

	return books, pagesCount, nil
}

func (bs *bookStorage) Create(bookCreateDTO *domain.CreateBookDTO) (string, error) {

	//s, _, _ := squirrel.Insert("book").
	//	Columns("title", "direction_uuid", "author_uuid", "difficulty", "edition_date", "rating", "description", "local_url", "language", "tags_uuids", "download_count", "image_url", "created_at").
	//	Values(bookCreateDTO.Title, bookCreateDTO.DirectionUUID, bookCreateDTO.AuthorUUID, bookCreateDTO.Difficulty, bookCreateDTO.EditionDate, 0, bookCreateDTO.Description, bookCreateDTO.URL, bookCreateDTO.Language, pq.Array(bookCreateDTO.TagsUUIDs), 0, bookCreateDTO.ImageURL, bookCreateDTO.CreatedAt).
	//	Suffix("WHERE EXISTS(SELECT uuid FROM author where author_uuid = author.uuid)", bookCreateDTO.AuthorUUID).
	//	Suffix("AND EXISTS(SELECT uuid FROM direction where direction_uuid = direction.uuid)", bookCreateDTO.DirectionUUID).
	//	Suffix("AND EXISTS(SELECT uuid FROM tag where tag.uuid = ?", pq.Array(bookCreateDTO.TagsUUIDs)).
	//	Suffix("RETURNING uuid").
	//	PlaceholderFormat(squirrel.Dollar).
	//	ToSql()

	tx, err := bs.db.Begin()
	if err != nil {
		bs.logger.Errorf("error occurred while creating transaction. err: %v", err)
		return "", err
	}

	var UUID string

	localURL := strings.Split(bookCreateDTO.LocalURL, "|split|")
	if len(localURL) < 2 {
		localURL = append(localURL, "")
	}
	imageURL := strings.Split(bookCreateDTO.ImageURL, "|split|")
	if len(imageURL) < 2 {
		imageURL = append(imageURL, "")
	}

	row := tx.QueryRow(createBookQuery,
		bookCreateDTO.Title,
		bookCreateDTO.DirectionUUID,
		bookCreateDTO.AuthorUUID,
		bookCreateDTO.Difficulty,
		bookCreateDTO.EditionDate,
		bookCreateDTO.Description,
		localURL[0],
		localURL[1],
		bookCreateDTO.Language,
		pq.Array(bookCreateDTO.TagsUUIDs),
		imageURL[0],
		imageURL[1],
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
		bs.logger.Errorf("No book with UUID %s wbs found", UUID)
		return ErrNoRowsAffected
	}
	bs.logger.Infof("Book with uuid %s wbs deleted.", UUID)
	return tx.Commit()
}

func (bs *bookStorage) Update(bookUpdateDTO *domain.UpdateBookDTO) error {
	if bookUpdateDTO.DirectionUUID == "" {
		bookUpdateDTO.DirectionUUID = "0"
	}
	if bookUpdateDTO.AuthorUUID == "" {
		bookUpdateDTO.AuthorUUID = "0"
	}

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
		bookUpdateDTO.Description,
		bookUpdateDTO.LocalURL,
		bookUpdateDTO.Language,
		pq.Array(bookUpdateDTO.TagsUUIDs),
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

func (bs *bookStorage) Rate(UUID string, rating float32) error {
	tx, err := bs.db.Begin()
	if err != nil {
		bs.logger.Errorf("error occurred while creating transaction. err: %v", err)
		return err
	}

	result, err := tx.Exec(rateBookQuery,
		rating,
		UUID,
	)
	if err != nil {
		tx.Rollback()
		bs.logger.Errorf("error occurred while rating book. err: %v", err)
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		tx.Rollback()
		bs.logger.Errorf("error occurred while raing book (getting affected rows). err: %v", err)
		return err
	}
	if rowsAffected < 1 {
		tx.Rollback()
		bs.logger.Errorf("error occurred while raing book. err: %v.", ErrNoRowsAffected)
		return ErrNoRowsAffected
	}

	bs.logger.Infof("book with uuid %s was rated.", UUID)

	return tx.Commit()
}

func (bs *bookStorage) DownloadCountUp(UUID string) error {
	tx, err := bs.db.Begin()
	if err != nil {
		bs.logger.Errorf("error occurred while creating transaction. err: %v", err)
		return err
	}

	result, err := tx.Exec(bookDownloadCountUpQuery,
		UUID,
	)
	if err != nil {
		tx.Rollback()
		bs.logger.Errorf("error occurred while rating book. err: %v", err)
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		tx.Rollback()
		bs.logger.Errorf("error occurred while raing book (getting affected rows). err: %v", err)
		return err
	}
	if rowsAffected < 1 {
		tx.Rollback()
		bs.logger.Errorf("error occurred while raing book. err: %v.", ErrNoRowsAffected)
		return ErrNoRowsAffected
	}

	bs.logger.Infof("book with uuid %s was rated.", UUID)

	return tx.Commit()
}
