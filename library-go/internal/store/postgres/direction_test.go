package postgres

import (
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"library-go/internal/domain"
	"library-go/pkg/logging"
	"regexp"
	"testing"
)

func TestDirectionStorage_GetOne(t *testing.T) {
	logger := logging.GetLogger()
	db, mock, err := sqlmock.New()
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()
	r := NewDirectionStorage(db, logger)

	testTable := []struct {
		name           string
		mock           func(uuid string)
		uuid           string
		expectedResult *domain.Direction
		expectError    bool
	}{
		{
			name: "OK",
			uuid: "1",
			mock: func(uuid string) {
				rows := sqlmock.NewRows([]string{"uuid", "name"}).
					AddRow("1", "Test Direction")
				mock.ExpectQuery(regexp.QuoteMeta(getOneDirectionQuery)).WithArgs(uuid).WillReturnRows(rows)
			},
			expectedResult: domain.TestDirection(),

			expectError: false,
		},
		{
			name: "DB error",
			uuid: "1",
			mock: func(uuid string) {
				mock.ExpectQuery(regexp.QuoteMeta(getOneDirectionQuery)).
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

func TestDirectionStorage_GetAll(t *testing.T) {
	logger := logging.GetLogger()
	db, mock, err := sqlmock.New()
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()
	r := NewDirectionStorage(db, logger)

	testTable := []struct {
		name           string
		mock           func(page, limit int)
		inputPage      int
		inputLimit     int
		expectedResult []*domain.Direction
		expectError    bool
	}{
		{
			name:       "OK",
			inputPage:  0,
			inputLimit: 0,
			mock: func(page, limit int) {
				rows := sqlmock.NewRows([]string{"uuid", "name"}).
					AddRow("1", "Test Direction").
					AddRow("1", "Test Direction")
				mock.ExpectQuery(regexp.QuoteMeta(getAllDirectionsQuery)).WillReturnRows(rows)
			},
			expectedResult: []*domain.Direction{domain.TestDirection(), domain.TestDirection()},
			expectError:    false,
		},
		{
			name:       "Data base error",
			inputPage:  0,
			inputLimit: 0,
			mock: func(page, limit int) {
				mock.ExpectQuery(regexp.QuoteMeta(getAllDirectionsQuery)).WillReturnError(errors.New("DB error"))
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

func TestDirectionStorage_Create(t *testing.T) {
	logger := logging.GetLogger()
	db, mock, err := sqlmock.New()
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()
	r := NewDirectionStorage(db, logger)

	testTable := []struct {
		name           string
		mock           func(dto *domain.CreateDirectionDTO)
		dto            *domain.CreateDirectionDTO
		expectedResult string
		expectError    bool
	}{
		{
			name: "OK",
			dto:  domain.TestDirectionCreateDTO(),
			mock: func(dto *domain.CreateDirectionDTO) {
				mock.ExpectBegin()
				rows := sqlmock.NewRows([]string{"uuid"}).AddRow("1")
				mock.ExpectQuery(regexp.QuoteMeta(createDirectionQuery)).WithArgs(dto.Name).WillReturnRows(rows)
				mock.ExpectCommit()
			},
			expectedResult: "1",
			expectError:    false,
		},
		{
			name: "Data base error",
			dto:  domain.TestDirectionCreateDTO(),
			mock: func(dto *domain.CreateDirectionDTO) {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(createDirectionQuery)).WithArgs(dto.Name).WillReturnError(errors.New("DB error"))
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

func TestDirectionStorage_Delete(t *testing.T) {
	logger := logging.GetLogger()
	db, mock, err := sqlmock.New()
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()
	r := NewDirectionStorage(db, logger)

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
				mock.ExpectExec(regexp.QuoteMeta(deleteDirectionQuery)).WithArgs(uuid).WillReturnResult(result)
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
				mock.ExpectExec(regexp.QuoteMeta(deleteDirectionQuery)).WithArgs(uuid).WillReturnResult(result)
				mock.ExpectRollback()
			},
			expectedError: true,
		},
		{
			name:      "DB error",
			inputData: "1",
			mock: func(uuid string) {
				mock.ExpectBegin()
				mock.ExpectExec("DELETE FROM direction WHERE").WithArgs(uuid).WillReturnError(errors.New("DB error"))
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

func TestDirectionStorage_Update(t *testing.T) {
	logger := logging.GetLogger()
	db, mock, err := sqlmock.New()
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()
	r := NewDirectionStorage(db, logger)

	testTable := []struct {
		name          string
		mock          func(dto *domain.UpdateDirectionDTO)
		inputData     *domain.UpdateDirectionDTO
		expectedError bool
	}{
		{
			name:      "OK",
			inputData: domain.TestDirectionUpdateDTO(),
			mock: func(dto *domain.UpdateDirectionDTO) {
				mock.ExpectBegin()
				result := sqlmock.NewResult(1, 1)
				mock.ExpectExec(regexp.QuoteMeta(updateDirectionQuery)).WithArgs(dto.Name, dto.UUID).WillReturnResult(result)
				mock.ExpectCommit()
			},
			expectedError: false,
		},
		{
			name:      "Invalid UUID or DB error",
			inputData: domain.TestDirectionUpdateDTO(),
			mock: func(dto *domain.UpdateDirectionDTO) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(updateDirectionQuery)).WithArgs(dto.Name, dto.UUID).WillReturnError(errors.New("DB error"))
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
