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
	UpdateReviewQuery = `UPDATE review SET 
			text = COALESCE(NULLIF($1, ''), text),
			full_name = COALESCE(NULLIF($2, ''), full_name)
		WHERE uuid = $3`

	RateReviewQuery = `WITH grades AS (
   		 SELECT avg((select avg(a) from unnest(array_append(all_grades, $1)) as a)) AS avg
   		 FROM review
		)
		UPDATE review SET
    	    all_grades = (CASE WHEN (0.0 < $1 AND $1 < 5.1) THEN array_append(all_grades, $1) ELSE all_grades END),
    	    rating = (CASE WHEN (0.0 < $1 AND $1 < 5.1) THEN grades.avg  ELSE rating END)
		FROM grades
		WHERE uuid = $2`
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
		&review.CreatedAt,
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
			&review.CreatedAt,
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
	query, args, _ := squirrel.Insert("review").
		Columns("text", "full_name", "source", "literature_uuid").
		Values(reviewCreateDTO.Text, reviewCreateDTO.FullName, reviewCreateDTO.Source, reviewCreateDTO.LiteratureUUID).
		PlaceholderFormat(squirrel.Dollar).
		Suffix("RETURNING uuid").
		ToSql()

	tx, err := rs.db.Begin()
	if err != nil {
		rs.logger.Errorf("error occurred while creating transaction. err: %v", err)
		return "", err
	}

	var UUID string
	row := tx.QueryRow(query, args...)
	if err := row.Scan(&UUID); err != nil {
		tx.Rollback()
		rs.logger.Errorf("error occurred while creating review. err: %v", err)
		return UUID, err
	}

	return UUID, tx.Commit()
}

func (rs *reviewStorage) Delete(UUID string) error {
	query, args, _ := squirrel.Delete("review").
		Where("uuid = ?", UUID).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	tx, err := rs.db.Begin()
	if err != nil {
		rs.logger.Errorf("error occurred while creating transaction. err: %v", err)
		return err
	}

	result, err := tx.Exec(query, args...)
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

	result, err := tx.Exec(UpdateReviewQuery,
		reviewUpdateDTO.Text,
		strings.Title(strings.ToLower(reviewUpdateDTO.FullName)),
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

func (rs *reviewStorage) Rate(UUID string, rating float32) error {
	tx, err := rs.db.Begin()
	if err != nil {
		rs.logger.Errorf("error occurred while creating transaction. err: %v", err)
		return err
	}

	result, err := tx.Exec(RateReviewQuery,
		rating,
		UUID,
	)
	if err != nil {
		tx.Rollback()
		rs.logger.Errorf("error occurred while rating review. err: %v", err)
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		tx.Rollback()
		rs.logger.Errorf("error occurred while raing review (getting affected rows). err: %v", err)
		return err
	}
	if rowsAffected < 1 {
		tx.Rollback()
		rs.logger.Errorf("error occurred while raing review. err: %v.", ErrNoRowsAffected)
		return ErrNoRowsAffected
	}

	rs.logger.Infof("review with uuid %s was rated.", UUID)

	return tx.Commit()
}
