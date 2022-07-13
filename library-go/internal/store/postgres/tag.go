package postgres

import (
	"database/sql"
	"github.com/Masterminds/squirrel"
	"github.com/lib/pq"
	"library-go/internal/domain"
	"library-go/internal/store"
	"library-go/pkg/logging"
	"strings"
)

const (
	UpdateTagQuery = `UPDATE tag SET 
                   name = $1 
		WHERE uuid = $2 RETURNING *`
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
	query, args, _ := squirrel.Select("uuid", "name").
		From("tag").
		Where(squirrel.Eq{
			"uuid": UUID,
		}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	var tag domain.Tag
	if err := ts.db.QueryRow(query, args...).Scan(
		&tag.UUID,
		&tag.Name,
	); err != nil {
		ts.logger.Errorf("error occurred while selecting tag from DB. err: %v", err)
		return nil, err
	}

	return &tag, nil
}

func (ts *tagStorage) GetMany(UUIDs []string) ([]*domain.Tag, error) {
	query, args, _ := squirrel.Select("uuid", "name").
		From("tag").
		Where("uuid = any(?)", pq.Array(UUIDs)).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	rows, err := ts.db.Query(query, args...)
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
	query, _, _ := squirrel.Select("uuid", "name").
		From("tag").
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	rows, err := ts.db.Query(query)
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
	query, args, _ := squirrel.Insert("tag").
		Columns("name").
		Values(strings.ToLower(tagCreateDTO.Name)).
		Suffix("RETURNING uuid").
		ToSql()

	tx, err := ts.db.Begin()
	if err != nil {
		ts.logger.Errorf("error occurred while creating transaction. err: %v", err)
		return "", err
	}

	var UUID string
	row := tx.QueryRow(query, args...)
	if err := row.Scan(&UUID); err != nil {
		tx.Rollback()
		ts.logger.Errorf("error occurred while creating tag. err: %v", err)
		return UUID, err
	}

	return UUID, tx.Commit()
}

func (ts *tagStorage) Delete(UUID string) error {
	query, args, _ := squirrel.Delete("tag").
		Where("uuid = ?", UUID).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	tx, err := ts.db.Begin()
	if err != nil {
		ts.logger.Errorf("error occurred while creating transaction. err: %v", err)
		return err
	}

	result, err := tx.Exec(query, args...)
	if err != nil {
		tx.Rollback()
		ts.logger.Errorf("error occurred while deleting tag. err: %v.", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		tx.Rollback()
		ts.logger.Errorf("error occurred while deleting tag (getting affected rows). err: %v", err)
		return err
	}

	if rowsAffected < 1 {
		tx.Rollback()
		ts.logger.Errorf("Tag with uuid %s wds deleted.", UUID)
		return ErrNoRowsAffected
	}
	ts.logger.Infof("Tag with uuid %s wts deleted.", UUID)
	return tx.Commit()
}

func (ts *tagStorage) Update(tagUpdateDTO *domain.UpdateTagDTO) error {
	tx, err := ts.db.Begin()
	if err != nil {
		ts.logger.Errorf("error occurred while creating transaction. err: %v", err)
		return err
	}

	result, err := tx.Exec(UpdateTagQuery,
		strings.ToLower(tagUpdateDTO.Name),
		tagUpdateDTO.UUID)

	if err != nil {
		tx.Rollback()
		ts.logger.Errorf("error occurred while updating tag. err: %v", err)
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		tx.Rollback()
		ts.logger.Errorf("Error occurred while updating tag. Err msg: %v.", err)
		return err
	}

	if rowsAffected < 1 {
		tx.Rollback()
		ts.logger.Errorf("Error occurred while updating tag. Err msg: %v.", ErrNoRowsAffected)
		return ErrNoRowsAffected
	}

	return tx.Commit()
}
