package postgres

import (
	"database/sql"
	"github.com/lib/pq"
	"library-go/internal/domain"
	"library-go/internal/store"
	"library-go/pkg/logging"
)

type tagStorage struct {
	logger *logging.Logger
	db     *sql.DB
}

func NewTagStorage(db *sql.DB, logger *logging.Logger) store.TagStorage {
	return &tagStorage{
		logger: logger,
		db:     db,
	}
}

func (ts *tagStorage) GetOne(UUID string) (*domain.Tag, error) {
	var tag domain.Tag
	if err := ts.db.QueryRow("SELECT * FROM tag WHERE uuid = $1",
		UUID).Scan(
		&tag.UUID,
		&tag.Name,
	); err != nil {
		ts.logger.Errorf("error occurred while selecting tag from DB. err: %v", err)
		return nil, err
	}

	return &tag, nil
}

func (ts *tagStorage) GetMany(UUIDs []string) ([]*domain.Tag, error) {
	rows, err := ts.db.Query("SELECT * FROM tag WHERE uuid = any($1)", pq.Array(UUIDs))
	if err != nil {
		ts.logger.Errorf("error occurred while selecting all tags. err: %v", err)
		return nil, err
	}
	var tags []*domain.Tag

	for rows.Next() {
		tag := domain.Tag{}
		err := rows.Scan(
			&tag.UUID,
			&tag.Name,
		)
		if err != nil {
			ts.logger.Errorf("error occurred while selecting tag. err: %v", err)
			continue
		}
		tags = append(tags, &tag)
	}
	return tags, nil
}

func (ts *tagStorage) GetAll(limit, offset int) ([]*domain.Tag, error) {
	rows, err := ts.db.Query("SELECT * FROM tag")
	if err != nil {
		ts.logger.Errorf("error occurred while selecting all tags. err: %v", err)
		return nil, err
	}
	var tags []*domain.Tag

	for rows.Next() {
		tag := domain.Tag{}
		err := rows.Scan(
			&tag.UUID,
			&tag.Name,
		)
		if err != nil {
			ts.logger.Errorf("error occurred while selecting tag. err: %v", err)
			continue
		}
		tags = append(tags, &tag)
	}
	return tags, nil
}

func (ts *tagStorage) Create(tagCreateDTO *domain.CreateTagDTO) (string, error) {
	var UUID string
	if err := ts.db.QueryRow(
		`INSERT INTO tag (
                     name
	) VALUES ($1) RETURNING uuid`,
		tagCreateDTO.Name,
	).Scan(&UUID); err != nil {
		ts.logger.Errorf("error occurred while creating tag. err: %v", err)
		return UUID, err
	}

	return UUID, nil
}

func (ts *tagStorage) Delete(UUID string) error {
	result, err := ts.db.Exec("DELETE FROM tag WHERE uuid = $1", UUID)
	if err != nil {
		ts.logger.Errorf("error occurred while deleting tag. err: %v.", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		ts.logger.Errorf("error occurred while deleting tag (getting affected rows). err: %v", err)
		return err
	}

	if rowsAffected < 1 {
		ts.logger.Errorf("error occurred while deleting tag. err: %v.", ErrNoRowsAffected)
		return ErrNoRowsAffected
	}
	ts.logger.Infof("Tag with uuid %s wts deleted.", UUID)
	return nil
}

func (ts *tagStorage) Update(tagUpdateDTO *domain.UpdateTagDTO) error {
	var tag domain.Tag
	if err := ts.db.QueryRow(
		`UPDATE tag SET 
                   name = $1 
		WHERE uuid = $2 RETURNING *`,
		tagUpdateDTO.Name,
		tagUpdateDTO.UUID,
	).Scan(
		&tag.UUID,
		&tag.Name,
	); err != nil {
		ts.logger.Errorf("error occurred while updating tag. err: %v", err)
		return err
	}

	return nil
}
