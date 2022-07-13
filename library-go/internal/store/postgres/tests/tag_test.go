package tests

import (
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Masterminds/squirrel"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"library-go/internal/domain"
	"library-go/internal/store/postgres"
	"library-go/pkg/logging"
	"regexp"
	"strings"
	"testing"
)

func TestTagStorage_GetOne(t *testing.T) {
	logger := logging.GetLogger("../../../../../logs", "test.log")
	db, mock, err := sqlmock.New()
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()
	r := postgres.NewTagStorage(db, logger)

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
				query, _, _ := squirrel.Select("uuid", "name").
					From("tag").
					Where(squirrel.Eq{
						"uuid": uuid,
					}).
					PlaceholderFormat(squirrel.Dollar).
					ToSql()
				rows := sqlmock.NewRows([]string{"uuid", "name"}).
					AddRow("1", "Test Tag")

				mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(uuid).WillReturnRows(rows)
			},
			expectedResult: domain.TestTag(),

			expectError: false,
		},
		{
			name: "no rows in result",
			uuid: "1",
			mock: func(uuid string) {
				query, _, _ := squirrel.Select("uuid", "name").
					From("tag").
					Where(squirrel.Eq{
						"uuid": uuid,
					}).
					PlaceholderFormat(squirrel.Dollar).
					ToSql()

				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WillReturnError(errors.New("no rows in result"))
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
	logger := logging.GetLogger("../../../../../logs", "test.log")
	db, mock, err := sqlmock.New()
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()
	r := postgres.NewTagStorage(db, logger)

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
				query, _, _ := squirrel.Select("uuid", "name").
					From("tag").
					Where("uuid = any(?)", pq.Array(uuids)).
					PlaceholderFormat(squirrel.Dollar).
					ToSql()

				rows := sqlmock.NewRows([]string{"uuid", "name"}).
					AddRow("1", "Test Tag").
					AddRow("1", "Test Tag")
				mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(pq.Array(uuids)).WillReturnRows(rows)
			},
			expectedResult: []*domain.Tag{domain.TestTag(), domain.TestTag()},

			expectError: false,
		},
		{
			name:  "no rows in result",
			uuids: []string{"1", "2"},
			mock: func(uuids []string) {
				query, _, _ := squirrel.Select("uuid", "name").
					From("tag").
					Where("uuid = any(?)", pq.Array(uuids)).
					PlaceholderFormat(squirrel.Dollar).
					ToSql()

				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WillReturnError(errors.New("no rows in result"))
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
	logger := logging.GetLogger("../../../../../logs", "test.log")
	db, mock, err := sqlmock.New()
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()
	r := postgres.NewTagStorage(db, logger)

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
				query, _, _ := squirrel.Select("uuid", "name").
					From("tag").
					PlaceholderFormat(squirrel.Dollar).
					ToSql()
				rows := sqlmock.NewRows([]string{"uuid", "name"}).
					AddRow("1", "Test Tag").
					AddRow("1", "Test Tag")

				mock.ExpectQuery(regexp.QuoteMeta(query)).WillReturnRows(rows)
			},
			expectedResult: []*domain.Tag{domain.TestTag(), domain.TestTag()},
			expectError:    false,
		},
		{
			name:       "no rows in result",
			inputPage:  0,
			inputLimit: 0,
			mock: func(page, limit int) {
				query, _, _ := squirrel.Select("uuid", "name").
					From("tag").
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

func TestTagStorage_Create(t *testing.T) {
	logger := logging.GetLogger("../../../../../logs", "test.log")
	db, mock, err := sqlmock.New()
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()
	r := postgres.NewTagStorage(db, logger)

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
				query, _, _ := squirrel.Insert("tag").
					Columns("name").
					Values(strings.ToLower(dto.Name)).
					Suffix("RETURNING uuid").
					ToSql()
				rows := sqlmock.NewRows([]string{"uuid"}).AddRow("1")

				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(dto.Name).WillReturnRows(rows)
				mock.ExpectCommit()
			},
			expectedResult: "1",
			expectError:    false,
		},
		{
			name: "no rows in result",
			dto:  domain.TestTagCreateDTO(),
			mock: func(dto *domain.CreateTagDTO) {
				query, _, _ := squirrel.Insert("tag").
					Columns("name").
					Values(strings.ToLower(dto.Name)).
					Suffix("RETURNING uuid").
					ToSql()

				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(dto.Name).WillReturnError(errors.New("no rows in result"))
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
	logger := logging.GetLogger("../../../../../logs", "test.log")
	db, mock, err := sqlmock.New()
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()
	r := postgres.NewTagStorage(db, logger)

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
				query, _, _ := squirrel.Delete("tag").
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
				query, _, _ := squirrel.Delete("tag").
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
	logger := logging.GetLogger("../../../../../logs", "test.log")
	db, mock, err := sqlmock.New()
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()
	r := postgres.NewTagStorage(db, logger)

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
				result := sqlmock.NewResult(1, 1)

				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(postgres.UpdateTagQuery)).WithArgs(dto.Name, dto.UUID).WillReturnResult(result)
				mock.ExpectCommit()
			},
			expectedError: false,
		},
		{
			name:      "Invalid UUID or DB error",
			inputData: domain.TestTagUpdateDTO(),
			mock: func(dto *domain.UpdateTagDTO) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(postgres.UpdateTagQuery)).WithArgs(dto.Name, dto.UUID).WillReturnError(errors.New("DB error"))
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
