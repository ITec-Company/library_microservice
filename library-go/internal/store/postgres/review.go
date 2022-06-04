package postgres

import (
	"database/sql"
	"github.com/Masterminds/squirrel"
	"library-go/internal/domain"
	"library-go/internal/store"
	"library-go/pkg/logging"
	"time"
)

const (
	getOneReviewQuery = `SELECT 
		uuid,
		text,
		rating,
		date,
		literature_uuid
	FROM review WHERE  uuid = $1`
	getAllReviewsQuery = `SELECT 
		uuid,
		full_name,
		text,
		source,
		rating,
		date,
		literature_uuid
	FROM review 
	GROUP BY uuid, full_name, source,text, rating, date, literature_uuid`
	createReviewQuery = `INSERT INTO review (
			text,
			full_name,
			source,
			date,
			rating,
			literature_uuid
		) SELECT $1, $2 , $3, $4, $5, $6 RETURNING uuid`
	deleteReviewQuery = `DELETE FROM review WHERE uuid = $1`
	updateReviewQuery = `UPDATE review SET 
			text = COALESCE(NULLIF($1, ''), text),
			full_name = COALESCE(NULLIF($2, ''), full_name),
			rating = COALESCE(NULLIF($3, 0), rating)
		WHERE uuid = $4`
)

type reviewStorage struct {
	logger *logging.Logger
	db     *sql.DB
}

func NewReviewStorage(db *sql.DB, logger *logging.Logger) store.ReviewStorage {
	return &reviewStorage{
		logger: logger,
		db:     db,
	}
}

func (rs *reviewStorage) GetOne(UUID string) (*domain.Review, error) {
	query, args, _ := squirrel.Select("uuid", "full_name", "text", "rating", "source", "date", "literature_uuid").
		From("review").
		Where(squirrel.Eq{
			"uuid": UUID,
		}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	var review domain.Review

	if err := rs.db.QueryRow(query, args...).Scan(
		&review.UUID,
		&review.FullName,
		&review.Text,
		&review.Rating,
		&review.Source,
		&review.Date,
		&review.LiteratureUUID,
	); err != nil {
		rs.logger.Errorf("error occurred while selecting review from DB. err: %v", err)
		return nil, err
	}

	return &review, nil
}

func (rs *reviewStorage) GetAll(limit, offset int) ([]*domain.Review, error) {
	query, _, _ := squirrel.Select("uuid", "full_name", "text", "rating", "source", "date", "literature_uuid").
		From("review").
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	rows, err := rs.db.Query(query)
	if err != nil {
		rs.logger.Errorf("error occurred while selecting all reviews. err: %v", err)
		return nil, err
	}
	var reviews []*domain.Review

	for rows.Next() {
		review := domain.Review{}
		err := rows.Scan(
			&review.UUID,
			&review.FullName,
			&review.Text,
			&review.Rating,
			&review.Source,
			&review.Date,
			&review.LiteratureUUID,
		)
		if err != nil {
			rs.logger.Errorf("error occurred while selecting review. err: %v", err)
			continue
		}
		reviews = append(reviews, &review)
	}
	return reviews, nil
}

func (rs *reviewStorage) Create(reviewCreateDTO *domain.CreateReviewDTO) (string, error) {
	tx, err := rs.db.Begin()
	if err != nil {
		rs.logger.Errorf("error occurred while creating transaction. err: %v", err)
		return "", err
	}

	var UUID string
	row := tx.QueryRow(createReviewQuery,
		reviewCreateDTO.Text,
		reviewCreateDTO.FullName,
		reviewCreateDTO.Source,
		time.Now(),
		0,
		reviewCreateDTO.LiteratureUUID,
	)
	if err := row.Scan(&UUID); err != nil {
		tx.Rollback()
		rs.logger.Errorf("error occurred while creating review. err: %v", err)
		return UUID, err
	}

	return UUID, tx.Commit()
}

func (rs *reviewStorage) Delete(UUID string) error {
	tx, err := rs.db.Begin()
	if err != nil {
		rs.logger.Errorf("error occurred while creating transaction. err: %v", err)
		return err
	}

	result, err := tx.Exec(deleteReviewQuery, UUID)
	if err != nil {
		tx.Rollback()
		rs.logger.Errorf("error occurred while deleting review. err: %v.", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		tx.Rollback()
		rs.logger.Errorf("error occurred while deleting review (getting affected rows). err: %v", err)
		return err
	}

	if rowsAffected < 1 {
		tx.Rollback()
		rs.logger.Errorf("Direction with uuid %s wds deleted.", UUID)
		return ErrNoRowsAffected
	}
	rs.logger.Infof("Review with uuid %s was deleted.", UUID)
	return tx.Commit()
}

func (rs *reviewStorage) Update(reviewUpdateDTO *domain.UpdateReviewDTO) error {
	tx, err := rs.db.Begin()
	if err != nil {
		rs.logger.Errorf("error occurred while creating transaction. err: %v", err)
		return err
	}

	result, err := tx.Exec(updateReviewQuery,
		reviewUpdateDTO.Text,
		reviewUpdateDTO.FullName,
		reviewUpdateDTO.Rating,
		reviewUpdateDTO.UUID,
	)
	if err != nil {
		tx.Rollback()
		rs.logger.Errorf("error occurred while updating review. err: %v", err)
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		tx.Rollback()
		rs.logger.Errorf("error occurred while updating review (getting affected rows). err: %v", err)
		return err
	}
	if rowsAffected < 1 {
		tx.Rollback()
		rs.logger.Errorf("error occurred while updating review. err: %v.", ErrNoRowsAffected)
		return ErrNoRowsAffected
	}

	rs.logger.Infof("Review with uuid %s was updated.", reviewUpdateDTO.UUID)

	return tx.Commit()
}
