package postgres

import (
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"library-go/internal/domain"
	"library-go/pkg/logging"
	"regexp"
	"testing"
)

func TestTagStorage_GetOne(t *testing.T) {
	logger := logging.GetLogger()
	db, mock, err := sqlmock.New()
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()
	r := NewTagStorage(db, logger)

	testTable := []struct {
		name           string
		mock           func(uuid string)
		uuid           string
		expectedResult *domain.Tag
		expectError    bool
	}{
		{
			name: "OK",
			uuid: "1",
			mock: func(uuid string) {
				rows := sqlmock.NewRows([]string{"uuid", "name"}).
					AddRow("1", "Test Tag")
				mock.ExpectQuery(regexp.QuoteMeta(getOneTagQuery)).WithArgs(uuid).WillReturnRows(rows)
			},
			expectedResult: domain.TestTag(),

			expectError: false,
		},
		{
			name: "DB error",
			uuid: "1",
			mock: func(uuid string) {
				mock.ExpectQuery(regexp.QuoteMeta(getOneTagQuery)).
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

func TestTagStorage_GetMany(t *testing.T) {
	logger := logging.GetLogger()
	db, mock, err := sqlmock.New()
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()
	r := NewTagStorage(db, logger)

	testTable := []struct {
		name           string
		mock           func(uuids []string)
		uuids          []string
		expectedResult []*domain.Tag
		expectError    bool
	}{
		{
			name:  "OK",
			uuids: []string{"1", "2"},
			mock: func(uuids []string) {
				rows := sqlmock.NewRows([]string{"uuid", "name"}).
					AddRow("1", "Test Tag").
					AddRow("1", "Test Tag")
				mock.ExpectQuery(regexp.QuoteMeta(getManyTagsQuery)).WithArgs(pq.Array(uuids)).WillReturnRows(rows)
			},
			expectedResult: []*domain.Tag{domain.TestTag(), domain.TestTag()},

			expectError: false,
		},
		{
			name:  "DB error",
			uuids: []string{"1", "2"},
			mock: func(uuids []string) {
				mock.ExpectQuery(regexp.QuoteMeta(getManyTagsQuery)).
					WillReturnError(errors.New("DB error"))
			},
			expectedResult: nil,

			expectError: true,
		},
	}
	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock(tt.uuids)
			result, err := r.GetMany(tt.uuids)
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

func TestTagStorage_GetAll(t *testing.T) {
	logger := logging.GetLogger()
	db, mock, err := sqlmock.New()
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()
	r := NewTagStorage(db, logger)

	testTable := []struct {
		name           string
		mock           func(page, limit int)
		inputPage      int
		inputLimit     int
		expectedResult []*domain.Tag
		expectError    bool
	}{
		{
			name:       "OK",
			inputPage:  0,
			inputLimit: 0,
			mock: func(page, limit int) {
				rows := sqlmock.NewRows([]string{"uuid", "name"}).
					AddRow("1", "Test Tag").
					AddRow("1", "Test Tag")
				mock.ExpectQuery(regexp.QuoteMeta(getAllTagsQuery)).WillReturnRows(rows)
			},
			expectedResult: []*domain.Tag{domain.TestTag(), domain.TestTag()},
			expectError:    false,
		},
		{
			name:       "Data base error",
			inputPage:  0,
			inputLimit: 0,
			mock: func(page, limit int) {
				mock.ExpectQuery(regexp.QuoteMeta(getAllTagsQuery)).WillReturnError(errors.New("DB error"))
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

func TestTagStorage_Create(t *testing.T) {
	logger := logging.GetLogger()
	db, mock, err := sqlmock.New()
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()
	r := NewTagStorage(db, logger)

	testTable := []struct {
		name           string
		mock           func(dto *domain.CreateTagDTO)
		dto            *domain.CreateTagDTO
		expectedResult string
		expectError    bool
	}{
		{
			name: "OK",
			dto:  domain.TestTagCreateDTO(),
			mock: func(dto *domain.CreateTagDTO) {
				mock.ExpectBegin()
				rows := sqlmock.NewRows([]string{"uuid"}).AddRow("1")
				mock.ExpectQuery(regexp.QuoteMeta(createTagQuery)).WithArgs(dto.Name).WillReturnRows(rows)
				mock.ExpectCommit()
			},
			expectedResult: "1",
			expectError:    false,
		},
		{
			name: "Data base error",
			dto:  domain.TestTagCreateDTO(),
			mock: func(dto *domain.CreateTagDTO) {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(createTagQuery)).WithArgs(dto.Name).WillReturnError(errors.New("DB error"))
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

func TestTagStorage_Delete(t *testing.T) {
	logger := logging.GetLogger()
	db, mock, err := sqlmock.New()
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()
	r := NewTagStorage(db, logger)

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
				mock.ExpectExec(regexp.QuoteMeta(deleteTagQuery)).WithArgs(uuid).WillReturnResult(result)
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
				mock.ExpectExec(regexp.QuoteMeta(deleteTagQuery)).WithArgs(uuid).WillReturnResult(result)
				mock.ExpectRollback()
			},
			expectedError: true,
		},
		{
			name:      "DB error",
			inputData: "1",
			mock: func(uuid string) {
				mock.ExpectBegin()
				mock.ExpectExec("DELETE FROM tag WHERE").WithArgs(uuid).WillReturnError(errors.New("DB error"))
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

func TestTagStorage_Update(t *testing.T) {
	logger := logging.GetLogger()
	db, mock, err := sqlmock.New()
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()
	r := NewTagStorage(db, logger)

	testTable := []struct {
		name          string
		mock          func(dto *domain.UpdateTagDTO)
		inputData     *domain.UpdateTagDTO
		expectedError bool
	}{
		{
			name:      "OK",
			inputData: domain.TestTagUpdateDTO(),
			mock: func(dto *domain.UpdateTagDTO) {
				mock.ExpectBegin()
				result := sqlmock.NewResult(1, 1)
				mock.ExpectExec(regexp.QuoteMeta(updateTagQuery)).WithArgs(dto.Name, dto.UUID).WillReturnResult(result)
				mock.ExpectCommit()
			},
			expectedError: false,
		},
		{
			name:      "Invalid UUID or DB error",
			inputData: domain.TestTagUpdateDTO(),
			mock: func(dto *domain.UpdateTagDTO) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(updateTagQuery)).WithArgs(dto.Name, dto.UUID).WillReturnError(errors.New("DB error"))
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
