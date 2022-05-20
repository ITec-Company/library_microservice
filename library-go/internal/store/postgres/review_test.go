package postgres

import (
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"library-go/internal/domain"
	"library-go/pkg/logging"
	"library-go/pkg/utils"
	"regexp"
	"testing"
	"time"
)

func TestReviewStorage_GetOne(t *testing.T) {
	logger := logging.GetLogger()
	db, mock, err := sqlmock.New()
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()
	r := NewReviewStorage(db, logger)

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
				rows := sqlmock.NewRows([]string{"uuid", "full_name", "text", "rating", "source", "date", "literature_uuid"}).
					AddRow("1", "Test Review", "Test Text", 5.5, "Test Source", time.Date(2000, 01, 01, 0, 0, 0, 0, time.UTC), "1")
				mock.ExpectQuery(regexp.QuoteMeta(getOneReviewQuery)).WithArgs(uuid).WillReturnRows(rows)
			},
			expectedResult: domain.TestReview(),

			expectError: false,
		},
		{
			name: "DB error",
			uuid: "1",
			mock: func(uuid string) {
				mock.ExpectQuery(regexp.QuoteMeta(getOneReviewQuery)).
					WillReturnError(errors.New("DB error"))
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
	logger := logging.GetLogger()
	db, mock, err := sqlmock.New()
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()
	r := NewReviewStorage(db, logger)

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
				rows := sqlmock.NewRows([]string{"uuid", "full_name", "text", "rating", "source", "date", "literature_uuid"}).
					AddRow("1", "Test Review", "Test Text", 5.5, "Test Source", time.Date(2000, 01, 01, 0, 0, 0, 0, time.UTC), "1").
					AddRow("1", "Test Review", "Test Text", 5.5, "Test Source", time.Date(2000, 01, 01, 0, 0, 0, 0, time.UTC), "1")
				mock.ExpectQuery(regexp.QuoteMeta(getAllReviewsQuery)).WillReturnRows(rows)
			},
			expectedResult: []*domain.Review{domain.TestReview(), domain.TestReview()},
			expectError:    false,
		},
		{
			name:       "Data base error",
			inputPage:  0,
			inputLimit: 0,
			mock: func(page, limit int) {
				mock.ExpectQuery(regexp.QuoteMeta(getAllReviewsQuery)).WillReturnError(errors.New("DB error"))
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
	logger := logging.GetLogger()
	db, mock, err := sqlmock.New()
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()
	r := NewReviewStorage(db, logger)

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
				mock.ExpectBegin()
				rows := sqlmock.NewRows([]string{"uuid"}).AddRow("1")
				mock.ExpectQuery(regexp.QuoteMeta(createReviewQuery)).WithArgs(dto.Text, dto.FullName, dto.Source, utils.AnyTime{}, 0, dto.LiteratureUUID).WillReturnRows(rows)
				mock.ExpectCommit()
			},
			expectedResult: "1",
			expectError:    false,
		},
		{
			name: "Data base error",
			dto:  domain.TestReviewCreateDTO(),
			mock: func(dto *domain.CreateReviewDTO) {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(createReviewQuery)).WithArgs(dto.Text, dto.FullName, dto.Source, utils.AnyTime{}, 0, dto.LiteratureUUID).WillReturnError(errors.New("DB error"))
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
	logger := logging.GetLogger()
	db, mock, err := sqlmock.New()
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()
	r := NewReviewStorage(db, logger)

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
				mock.ExpectBegin()
				result := sqlmock.NewResult(1, 1)
				mock.ExpectExec(regexp.QuoteMeta(deleteReviewQuery)).WithArgs(uuid).WillReturnResult(result)
				mock.ExpectCommit()
			},
			expectedError: false,
		},
		{
			name:      "No rows affected",
			inputData: "1",
			mock: func(uuid string) {
				mock.ExpectBegin()
				result := sqlmock.NewResult(0, 0)
				mock.ExpectExec(regexp.QuoteMeta(deleteReviewQuery)).WithArgs(uuid).WillReturnResult(result)
				mock.ExpectRollback()
			},
			expectedError: true,
		},
		{
			name:      "DB error",
			inputData: "1",
			mock: func(uuid string) {
				mock.ExpectBegin()
				mock.ExpectExec("DELETE FROM review WHERE").WithArgs(uuid).WillReturnError(errors.New("DB error"))
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
	logger := logging.GetLogger()
	db, mock, err := sqlmock.New()
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()
	r := NewReviewStorage(db, logger)

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
				mock.ExpectExec(regexp.QuoteMeta(updateReviewQuery)).WithArgs(dto.Text, dto.FullName, dto.Rating, dto.UUID).WillReturnResult(result)
				mock.ExpectCommit()
			},
			expectedError: false,
		},
		{
			name:      "Invalid UUID or DB error",
			inputData: domain.TestReviewUpdateDTO(),
			mock: func(dto *domain.UpdateReviewDTO) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(updateReviewQuery)).WithArgs(dto.Text, dto.FullName, dto.Rating, dto.UUID).WillReturnError(errors.New("DB error"))
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
