package postgres

import (
	"database/sql"
	"library-go/internal/domain"
	"library-go/internal/store"
	"library-go/pkg/logging"
)

const (
	getOneAuthorQuery  = `SELECT * FROM author WHERE uuid = $1`
	getAllAuthorsQuery = `SELECT * FROM author`
	createAuthorQuery  = `INSERT INTO author (
                     full_name
	) VALUES ($1) RETURNING uuid`
	deleteAuthorQuery = `DELETE FROM author WHERE uuid = $1`
	updateAuthorQuery = `UPDATE author SET
						full_name = COALESCE(NULLIF($1, ''), full_name)
						WHERE uuid = $2`
)

type authorStorage struct {
	logger *logging.Logger
	db     *sql.DB
}

func NewAuthorStorage(db *sql.DB, logger *logging.Logger) store.AuthorStorage {
	return &authorStorage{
		logger: logger,
		db:     db,
	}
}

func (as *authorStorage) GetOne(UUID string) (*domain.Author, error) {
	var author domain.Author
	if err := as.db.QueryRow(getOneAuthorQuery,
		UUID).Scan(
		&author.UUID,
		&author.FullName,
	); err != nil {
		as.logger.Errorf("error occurred while selecting author from DB. err: %v", err)
		return nil, err
	}

	return &author, nil
}

func (as *authorStorage) GetAll(limit, offset int) ([]*domain.Author, error) {
	rows, err := as.db.Query(getAllAuthorsQuery)
	if err != nil {
		as.logger.Errorf("error occurred while selecting all authors. err: %v", err)
		return nil, err
	}
	var authors []*domain.Author

	for rows.Next() {
		author := domain.Author{}
		err := rows.Scan(
			&author.UUID,
			&author.FullName,
		)
		if err != nil {
			as.logger.Errorf("error occurred while selecting author. err: %v", err)
			continue
		}
		authors = append(authors, &author)
	}
	return authors, nil
}

func (as *authorStorage) Create(authorCreateDTO *domain.CreateAuthorDTO) (string, error) {
	tx, err := as.db.Begin()
	if err != nil {
		as.logger.Errorf("error occurred while creating transaction. err: %v", err)
		return "", err
	}

	var UUID string

	row := tx.QueryRow(createAuthorQuery,
		authorCreateDTO.FullName,
	)
	if err := row.Scan(
		&UUID,
	); err != nil {
		tx.Rollback()
		as.logger.Errorf("error occurred while creating author. err: %v", err)
		return UUID, err
	}

	return UUID, tx.Commit()
}

func (as *authorStorage) Delete(UUID string) error {
	tx, err := as.db.Begin()
	if err != nil {
		as.logger.Errorf("error occurred while creating transaction. err: %v", err)
		return err
	}

	result, err := tx.Exec(deleteAuthorQuery, UUID)
	if err != nil {
		tx.Rollback()
		as.logger.Errorf("error occurred while deleting author. err: %v.", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		tx.Rollback()
		as.logger.Errorf("error occurred while deleting author (getting affected rows). err: %v", err)
		return err
	}

	if rowsAffected < 1 {
		tx.Rollback()
		as.logger.Errorf("Author with uuid %s was deleted.", UUID)
		return ErrNoRowsAffected
	}
	as.logger.Infof("Author with uuid %s was deleted.", UUID)
	return tx.Commit()
}

func (as *authorStorage) Update(authorUpdateDTO *domain.UpdateAuthorDTO) error {
	tx, err := as.db.Begin()
	if err != nil {
		as.logger.Errorf("error occurred while creating transaction. err: %v", err)
		return err
	}

	result, err := tx.Exec(updateAuthorQuery,
		authorUpdateDTO.FullName,
		authorUpdateDTO.UUID)

	if err != nil {
		tx.Rollback()
		as.logger.Errorf("error occurred while updating author. err: %v", err)
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		tx.Rollback()
		as.logger.Errorf("Error occurred while updating author. Err msg: %v.", err)
		return err
	}

	if rowsAffected < 1 {
		tx.Rollback()
		as.logger.Errorf("Error occurred while updating author. Err msg: %v.", ErrNoRowsAffected)
		return ErrNoRowsAffected
	}

	return tx.Commit()
}
