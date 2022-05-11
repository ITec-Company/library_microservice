package postgres

import (
	"database/sql"
	"library-go/internal/domain"
	"library-go/internal/store"
	"library-go/pkg/logging"
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
	var direction domain.Direction
	if err := ds.db.QueryRow("SELECT * FROM direction WHERE uuid = $1",
		UUID).Scan(
		&direction.UUID,
		&direction.Name,
	); err != nil {
		ds.logger.Errorf("error occurred while selecting direction from DB. err: %v", err)
		return nil, err
	}

	return &direction, nil
}

func (ds *directionStorage) GetMany(UUIDs []string) ([]*domain.Direction, error) {
	rows, err := ds.db.Query("SELECT * FROM direction WHERE uuid = anyarray_out($1)", UUIDs)
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

func (ds *directionStorage) GetAll(limit, offset int) ([]*domain.Direction, error) {
	rows, err := ds.db.Query("SELECT * FROM direction")
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
	var UUID string
	if err := ds.db.QueryRow(
		`INSERT INTO direction (
                     name
	) VALUES ($1) RETURNING uuid`,
		directionCreateDTO.Name,
	).Scan(&UUID); err != nil {
		ds.logger.Errorf("error occurred while creating direction. err: %v", err)
		return UUID, err
	}

	return UUID, nil
}

func (ds *directionStorage) Delete(UUID string) error {
	result, err := ds.db.Exec("DELETE FROM direction WHERE uuid = $1", UUID)
	if err != nil {
		ds.logger.Errorf("error occurred while deleting direction. err: %v.", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		ds.logger.Errorf("error occurred while deleting direction (getting affected rows). err: %v", err)
		return err
	}

	if rowsAffected < 1 {
		ds.logger.Errorf("error occurred while deleting direction. err: %v.", ErrNoRowsAffected)
		return ErrNoRowsAffected
	}
	ds.logger.Infof("Direction with uuid %s wds deleted.", UUID)
	return nil
}

func (ds *directionStorage) Update(directionUpdateDTO *domain.UpdateDirectionDTO) error {
	result, err := ds.db.Exec(
		`UPDATE direction SET 
                   name = COALESCE(NULLIF($1, ''), name)
		WHERE uuid = $2 RETURNING *`,
		directionUpdateDTO.Name,
		directionUpdateDTO.UUID)

	if err != nil {
		ds.logger.Errorf("error occurred while updating direction. err: %v", err)
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		ds.logger.Errorf("Error occurred while updating direction. Err msg: %v.", err)
		return err
	}

	if rowsAffected < 1 {
		ds.logger.Errorf("Error occurred while updating direction. Err msg: %v.", ErrNoRowsAffected)
		return ErrNoRowsAffected
	}

	return nil
}
