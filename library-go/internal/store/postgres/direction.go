package postgres

import (
	"database/sql"
	"github.com/Masterminds/squirrel"
	"library-go/internal/domain"
	"library-go/internal/store"
	"library-go/pkg/logging"
	"strings"
)

const (
	UpdateDirectionQuery = `UPDATE direction SET 
                   name = COALESCE(NULLIF($1, ''), name)
		WHERE uuid = $2 RETURNING *`
)

type directionStorage struct {
	logger *logging.Logger
	db     *sql.DB
}

func NewDirectionStorage(db *sql.DB, logger *logging.Logger) store.DirectionStorage {
	return &directionStorage{
		logger: logger,
		db:     db,
	}
}

func (ds *directionStorage) GetOne(UUID string) (*domain.Direction, error) {
	query, args, _ := squirrel.Select("uuid", "name").
		From("direction").
		Where(squirrel.Eq{
			"uuid": UUID,
		}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	var direction domain.Direction

	if err := ds.db.QueryRow(query, args...).Scan(
		&direction.UUID,
		&direction.Name,
	); err != nil {
		ds.logger.Errorf("error occurred while selecting direction from DB. err: %v", err)
		return nil, err
	}

	return &direction, nil
}

func (ds *directionStorage) GetAll(limit, offset int) ([]*domain.Direction, error) {

	query, _, _ := squirrel.Select("uuid", "name").
		From("direction").
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	rows, err := ds.db.Query(query)
	if err != nil {
		ds.logger.Errorf("error occurred while selecting all directions. err: %v", err)
		return nil, err
	}
	var directions []*domain.Direction

	for rows.Next() {
		direction := domain.Direction{}
		err := rows.Scan(
			&direction.UUID,
			&direction.Name,
		)
		if err != nil {
			ds.logger.Errorf("error occurred while selecting direction. err: %v", err)
			continue
		}
		directions = append(directions, &direction)
	}
	return directions, nil
}

func (ds *directionStorage) Create(directionCreateDTO *domain.CreateDirectionDTO) (string, error) {
	query, args, _ := squirrel.Insert("direction").
		Columns("name").
		Values(strings.Title(strings.ToLower(directionCreateDTO.Name))).
		PlaceholderFormat(squirrel.Dollar).
		Suffix("RETURNING  uuid").
		ToSql()

	tx, err := ds.db.Begin()
	if err != nil {
		ds.logger.Errorf("error occurred while creating transaction. err: %v", err)
		return "", err
	}

	var UUID string
	row := tx.QueryRow(query, args...)
	if err := row.Scan(&UUID); err != nil {
		tx.Rollback()
		ds.logger.Errorf("error occurred while creating direction. err: %v", err)
		return UUID, err
	}

	return UUID, tx.Commit()
}

func (ds *directionStorage) Delete(UUID string) error {
	query, args, _ := squirrel.Delete("direction").
		Where("uuid = ?", UUID).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	tx, err := ds.db.Begin()
	if err != nil {
		ds.logger.Errorf("error occurred while creating transaction. err: %v", err)
		return err
	}

	result, err := tx.Exec(query, args...)
	if err != nil {
		tx.Rollback()
		ds.logger.Errorf("error occurred while deleting direction. err: %v.", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		tx.Rollback()
		ds.logger.Errorf("error occurred while deleting direction (getting affected rows). err: %v", err)
		return err
	}

	if rowsAffected < 1 {
		tx.Rollback()
		ds.logger.Errorf("Direction with uuid %s was deleted.", UUID)
		return ErrNoRowsAffected
	}
	ds.logger.Infof("Direction with uuid %s wds deleted.", UUID)
	return tx.Commit()
}

func (ds *directionStorage) Update(directionUpdateDTO *domain.UpdateDirectionDTO) error {
	tx, err := ds.db.Begin()
	if err != nil {
		ds.logger.Errorf("error occurred while creating transaction. err: %v", err)
		return err
	}

	result, err := tx.Exec(UpdateDirectionQuery,
		strings.Title(strings.ToLower(directionUpdateDTO.Name)),
		directionUpdateDTO.UUID)

	if err != nil {
		tx.Rollback()
		ds.logger.Errorf("error occurred while updating direction. err: %v", err)
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		tx.Rollback()
		ds.logger.Errorf("Error occurred while updating direction. Err msg: %v.", err)
		return err
	}

	if rowsAffected < 1 {
		tx.Rollback()
		ds.logger.Errorf("Error occurred while updating direction. Err msg: %v.", ErrNoRowsAffected)
		return ErrNoRowsAffected
	}

	return tx.Commit()
}
