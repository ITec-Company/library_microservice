package postgres

import (
	"database/sql"
	"library-go/internal/domain"
	"library-go/internal/store"
	"library-go/pkg/logging"
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
	if err := as.db.QueryRow("SELECT * FROM author WHERE uuid = $1",
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
	rows, err := as.db.Query("SELECT * FROM author")
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

	var UUID string

	if err := as.db.QueryRow(
		`INSERT INTO author (
                     full_name
	) VALUES ($1) RETURNING uuid`,
		authorCreateDTO.FullName,
	).Scan(
		&UUID,
	); err != nil {
		as.logger.Errorf("error occurred while creating author. err: %v", err)
		return UUID, err
	}
	as.logger.Errorf("%v", UUID)

	return UUID, nil
}

func (as *authorStorage) Delete(UUID string) error {
	result, err := as.db.Exec("DELETE FROM author WHERE uuid = $1", UUID)
	if err != nil {
		as.logger.Errorf("error occurred while deleting author. err: %v.", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		as.logger.Errorf("error occurred while deleting author (getting affected rows). err: %v", err)
		return err
	}

	if rowsAffected < 1 {
		as.logger.Errorf("error occurred while deleting author. err: %v.", ErrNoRowsAffected)
		return ErrNoRowsAffected
	}
	as.logger.Infof("Author with uuid %s was deleted.", UUID)
	return nil
}

func (as *authorStorage) Update(authorUpdateDTO *domain.UpdateAuthorDTO) error {
	result, err := as.db.Exec(
		`UPDATE author SET
	              full_name = COALESCE(NULLIF($1, ''), full_name)
		WHERE uuid = $2`,
		authorUpdateDTO.FullName,
		authorUpdateDTO.UUID)

	if err != nil {
		as.logger.Errorf("error occurred while updating author. err: %v", err)
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		as.logger.Errorf("Error occurred while updating author. Err msg: %v.", err)
		return err
	}

	if rowsAffected < 1 {
		as.logger.Errorf("Error occurred while updating author. Err msg: %v.", ErrNoRowsAffected)
		return ErrNoRowsAffected
	}

	return nil
}
