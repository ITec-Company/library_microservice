package tests

import (
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Masterminds/squirrel"
	"github.com/stretchr/testify/assert"
	"library-go/internal/domain"
	"library-go/internal/store/postgres"
	"library-go/pkg/logging"
	"regexp"
	"strings"
	"testing"
	"time"
)

func TestReviewStorage_GetOne(t *testing.T) {
	logger := logging.GetLogger("../../../../../logs", "test.log")
	db, mock, err := sqlmock.New()
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()
	r := postgres.NewReviewStorage(db, logger)

	testTable := []struct {
		name           string
		mock           func(uuid string)
		uuid           string
		expectedResult *domain.Review
		expectError    bool
	}{
		{
			name: "OK",
			uuid: "1",
			mock: func(uuid string) {
				query, _, _ := squirrel.Select("uuid", "full_name", "text", "rating", "source", "date", "literature_uuid").
					From("review").
					Where(squirrel.Eq{
						"uuid": uuid,
					}).
					PlaceholderFormat(squirrel.Dollar).
					ToSql()
				rows := sqlmock.NewRows([]string{"uuid", "full_name", "text", "rating", "source", "date", "literature_uuid"}).
					AddRow("1", "Test Review", "Test Text", 5.5, "Test Source", time.Date(2000, 01, 01, 0, 0, 0, 0, time.UTC), "1")

				mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(uuid).WillReturnRows(rows)
			},
			expectedResult: domain.TestReview(),

			expectError: false,
		},
		{
			name: "invalid uuid",
			uuid: "1",
			mock: func(uuid string) {
				query, _, _ := squirrel.Select("uuid", "full_name", "text", "rating", "source", "date", "literature_uuid").
					From("review").
					Where(squirrel.Eq{
						"uuid": uuid,
					}).
					PlaceholderFormat(squirrel.Dollar).
					ToSql()

				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WillReturnError(errors.New("no rows un result"))
			},
			expectedResult: nil,

			expectError: true,
		},
	}
	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock(tt.uuid)
			result, err := r.GetOne(tt.uuid)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult, result)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestReviewStorage_GetAll(t *testing.T) {
	logger := logging.GetLogger("../../../../../logs", "test.log")
	db, mock, err := sqlmock.New()
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()
	r := postgres.NewReviewStorage(db, logger)

	testTable := []struct {
		name           string
		mock           func(page, limit int)
		inputPage      int
		inputLimit     int
		expectedResult []*domain.Review
		expectError    bool
	}{
		{
			name:       "OK",
			inputPage:  0,
			inputLimit: 0,
			mock: func(page, limit int) {
				query, _, _ := squirrel.Select("uuid", "full_name", "text", "rating", "source", "date", "literature_uuid").
					From("review").
					PlaceholderFormat(squirrel.Dollar).
					ToSql()
				rows := sqlmock.NewRows([]string{"uuid", "full_name", "text", "rating", "source", "date", "literature_uuid"}).
					AddRow("1", "Test Review", "Test Text", 5.5, "Test Source", time.Date(2000, 01, 01, 0, 0, 0, 0, time.UTC), "1").
					AddRow("1", "Test Review", "Test Text", 5.5, "Test Source", time.Date(2000, 01, 01, 0, 0, 0, 0, time.UTC), "1")

				mock.ExpectQuery(regexp.QuoteMeta(query)).WillReturnRows(rows)
			},
			expectedResult: []*domain.Review{domain.TestReview(), domain.TestReview()},
			expectError:    false,
		},
		{
			name:       "no rows in result",
			inputPage:  0,
			inputLimit: 0,
			mock: func(page, limit int) {
				query, _, _ := squirrel.Select("uuid", "full_name", "text", "rating", "source", "date", "literature_uuid").
					From("review").
					PlaceholderFormat(squirrel.Dollar).
					ToSql()

				mock.ExpectQuery(regexp.QuoteMeta(query)).WillReturnError(errors.New("no rows in result"))
			},
			expectedResult: nil,
			expectError:    true,
		},
	}
	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock(tt.inputPage, tt.inputLimit)
			result, err := r.GetAll(tt.inputPage, tt.inputLimit)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult, result)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestReviewStorage_Create(t *testing.T) {
	logger := logging.GetLogger("../../../../../logs", "test.log")
	db, mock, err := sqlmock.New()
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()
	r := postgres.NewReviewStorage(db, logger)

	testTable := []struct {
		name           string
		mock           func(dto *domain.CreateReviewDTO)
		dto            *domain.CreateReviewDTO
		expectedResult string
		expectError    bool
	}{
		{
			name: "OK",
			dto:  domain.TestReviewCreateDTO(),
			mock: func(dto *domain.CreateReviewDTO) {
				query, _, _ := squirrel.Insert("review").
					Columns("text", "full_name", "source", "literature_uuid").
					Values(dto.Text, dto.FullName, dto.Source, dto.LiteratureUUID).
					PlaceholderFormat(squirrel.Dollar).
					Suffix("RETURNING uuid").
					ToSql()
				rows := sqlmock.NewRows([]string{"uuid"}).AddRow("1")

				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(dto.Text, dto.FullName, dto.Source, dto.LiteratureUUID).WillReturnRows(rows)
				mock.ExpectCommit()
			},
			expectedResult: "1",
			expectError:    false,
		},
		{
			name: "No rows in result",
			dto:  domain.TestReviewCreateDTO(),
			mock: func(dto *domain.CreateReviewDTO) {
				query, _, _ := squirrel.Insert("review").
					Columns("text", "full_name", "source", "literature_uuid").
					Values(dto.Text, dto.FullName, dto.Source, dto.LiteratureUUID).
					PlaceholderFormat(squirrel.Dollar).
					Suffix("RETURNING uuid").
					ToSql()

				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(dto.Text, dto.FullName, dto.Source, dto.LiteratureUUID).WillReturnError(errors.New("no rows in result"))
				mock.ExpectRollback()
			},
			expectedResult: "",
			expectError:    true,
		},
	}
	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock(tt.dto)
			result, err := r.Create(tt.dto)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult, result)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestReviewStorage_Delete(t *testing.T) {
	logger := logging.GetLogger("../../../../../logs", "test.log")
	db, mock, err := sqlmock.New()
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()
	r := postgres.NewReviewStorage(db, logger)

	testTable := []struct {
		name          string
		mock          func(uuid string)
		inputData     string
		expectedError bool
	}{
		{
			name:      "OK",
			inputData: "1",
			mock: func(uuid string) {
				query, _, _ := squirrel.Delete("review").
					Where("uuid = ?", uuid).
					PlaceholderFormat(squirrel.Dollar).
					ToSql()
				result := sqlmock.NewResult(1, 1)

				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(uuid).WillReturnResult(result)
				mock.ExpectCommit()
			},
			expectedError: false,
		},
		{
			name:      "No rows affected",
			inputData: "1",
			mock: func(uuid string) {
				query, _, _ := squirrel.Delete("review").
					Where("uuid = ?", uuid).
					PlaceholderFormat(squirrel.Dollar).
					ToSql()
				result := sqlmock.NewResult(0, 0)

				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(uuid).WillReturnResult(result)
				mock.ExpectRollback()
			},
			expectedError: true,
		},
	}
	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock(tt.inputData)
			err := r.Delete(tt.inputData)
			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestReviewStorage_Update(t *testing.T) {
	logger := logging.GetLogger("../../../../../logs", "test.log")
	db, mock, err := sqlmock.New()
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()
	r := postgres.NewReviewStorage(db, logger)

	testTable := []struct {
		name          string
		mock          func(dto *domain.UpdateReviewDTO)
		inputData     *domain.UpdateReviewDTO
		expectedError bool
	}{
		{
			name:      "OK",
			inputData: domain.TestReviewUpdateDTO(),
			mock: func(dto *domain.UpdateReviewDTO) {
				mock.ExpectBegin()
				result := sqlmock.NewResult(1, 1)
				mock.ExpectExec(regexp.QuoteMeta(postgres.UpdateReviewQuery)).WithArgs(dto.Text, strings.Title(strings.ToLower(dto.FullName)), dto.UUID).WillReturnResult(result)
				mock.ExpectCommit()
			},
			expectedError: false,
		},
		{
			name:      "Invalid UUID or DB error",
			inputData: domain.TestReviewUpdateDTO(),
			mock: func(dto *domain.UpdateReviewDTO) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(postgres.UpdateReviewQuery)).WithArgs(dto.Text, strings.Title(strings.ToLower(dto.FullName)), dto.UUID).WillReturnError(errors.New("DB error"))
				mock.ExpectRollback()
			},
			expectedError: true,
		},
	}
	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock(tt.inputData)
			err := r.Update(tt.inputData)
			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestReviewStorage_Rate(t *testing.T) {
	logger := logging.GetLogger("../../../../../logs", "test.log")
	db, mock, err := sqlmock.New()
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()
	r := postgres.NewReviewStorage(db, logger)

	testTable := []struct {
		name          string
		mock          func(uuid string, rating float32)
		uuid          string
		rating        float32
		expectedError bool
	}{
		{
			name:   "OK",
			uuid:   "1",
			rating: 4.0,
			mock: func(uuid string, rating float32) {
				result := sqlmock.NewResult(1, 1)

				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(postgres.RateReviewQuery)).WithArgs(rating, uuid).WillReturnResult(result)
				mock.ExpectCommit()
			},
			expectedError: false,
		},
		{
			name:   "Invalid UUID or DB error",
			uuid:   "1",
			rating: 4.0,
			mock: func(uuid string, rating float32) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(postgres.RateReviewQuery)).WithArgs(rating, uuid).WillReturnError(errors.New("DB error"))
				mock.ExpectRollback()
			},
			expectedError: true,
		},
	}
	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock(tt.uuid, tt.rating)
			err := r.Rate(tt.uuid, tt.rating)
			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
