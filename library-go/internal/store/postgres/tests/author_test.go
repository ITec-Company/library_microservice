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
	"testing"
)

func TestAuthorStorage_GetOne(t *testing.T) {
	logger := logging.GetLogger("../../../../../logs", "test.log")
	db, mock, _ := sqlmock.New()
	defer db.Close()
	r := postgres.NewAuthorStorage(db, logger)

	testTable := []struct {
		name           string
		mock           func(uuid string)
		uuid           string
		expectedResult *domain.Author
		expectError    bool
	}{
		{
			name: "OK",
			uuid: "1",
			mock: func(uuid string) {
				query, _, _ := squirrel.Select("uuid", "full_name").
					From("author").
					Where(squirrel.Eq{
						"uuid": uuid,
					}).
					PlaceholderFormat(squirrel.Dollar).ToSql()

				rows := sqlmock.NewRows([]string{"uuid", "full_name"}).
					AddRow("1", "Test Author")

				mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(uuid).WillReturnRows(rows)
			},
			expectedResult: domain.TestAuthor(),
			expectError:    false,
		},
		{
			name: "no rows in result",
			uuid: "1",
			mock: func(uuid string) {
				query, _, _ := squirrel.Select("uuid", "full_name").
					From("author").
					Where(squirrel.Eq{
						"uuid": uuid,
					}).
					PlaceholderFormat(squirrel.Dollar).ToSql()

				mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(uuid).WillReturnError(errors.New("no rows in result"))
			},
			expectedResult: nil,
			expectError:    true,
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

func TestAuthorStorage_GetAll(t *testing.T) {
	logger := logging.GetLogger("../../../../../logs", "test.log")
	db, mock, _ := sqlmock.New()
	defer db.Close()
	r := postgres.NewAuthorStorage(db, logger)

	testTable := []struct {
		name           string
		mock           func(page, limit int)
		inputPage      int
		inputLimit     int
		expectedResult []*domain.Author
		expectError    bool
	}{
		{
			name: "OK",
			mock: func(page, limit int) {
				rows := sqlmock.NewRows([]string{"uuid", "full_name"}).
					AddRow("1", "Test Author").
					AddRow("1", "Test Author")
				query, _, _ := squirrel.Select("uuid", "full_name").
					From("author").
					PlaceholderFormat(squirrel.Dollar).
					ToSql()

				mock.ExpectQuery(regexp.QuoteMeta(query)).WillReturnRows(rows)
			},
			expectedResult: []*domain.Author{domain.TestAuthor(), domain.TestAuthor()},
			expectError:    false,
		},
		{
			name: "no rows in result",
			mock: func(page, limit int) {
				query, _, _ := squirrel.Select("uuid", "full_name").
					From("author").
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

func TestAuthorStorage_Create(t *testing.T) {
	logger := logging.GetLogger("../../../../../logs", "test.log")
	db, mock, err := sqlmock.New()
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()
	r := postgres.NewAuthorStorage(db, logger)

	testTable := []struct {
		name           string
		mock           func(dto *domain.CreateAuthorDTO)
		dto            *domain.CreateAuthorDTO
		expectedResult string
		expectError    bool
	}{
		{
			name: "OK",
			dto:  domain.TestAuthorCreateDTO(),
			mock: func(dto *domain.CreateAuthorDTO) {
				rows := sqlmock.NewRows([]string{"uuid"}).AddRow("1")
				query, _, _ := squirrel.Insert("author").
					Columns("full_name").
					Values(dto.FullName).
					Suffix("RETURNING  uuid").
					ToSql()

				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(dto.FullName).WillReturnRows(rows)
				mock.ExpectCommit()
			},
			expectedResult: "1",
			expectError:    false,
		},
		{
			name: "full_name duplicate",
			dto:  domain.TestAuthorCreateDTO(),
			mock: func(dto *domain.CreateAuthorDTO) {
				query, _, _ := squirrel.Insert("author").
					Columns("full_name").
					Values(dto.FullName).
					ToSql()

				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(dto.FullName).WillReturnError(errors.New("full_name duplicate"))
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

func TestAuthorStorage_Delete(t *testing.T) {
	logger := logging.GetLogger("../../../../../logs", "test.log")
	db, mock, err := sqlmock.New()
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()
	r := postgres.NewAuthorStorage(db, logger)

	testTable := []struct {
		name          string
		uuid          string
		mock          func(uuid string)
		expectedError bool
	}{
		{
			name: "OK",
			uuid: "1",
			mock: func(uuid string) {
				result := sqlmock.NewResult(1, 1)
				query, _, _ := squirrel.Delete("author").
					Where("uuid =?", uuid).
					PlaceholderFormat(squirrel.Dollar).
					ToSql()

				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(uuid).WillReturnResult(result)
				mock.ExpectCommit()
			},
			expectedError: false,
		},
		{
			name: "No rows affected",
			uuid: "1",
			mock: func(uuid string) {
				query, _, _ := squirrel.Delete("author").
					Where("uuid =?", uuid).
					PlaceholderFormat(squirrel.Dollar).
					ToSql()

				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(uuid).WillReturnError(errors.New("no rows affected"))
				mock.ExpectRollback()
			},
			expectedError: true,
		},
	}
	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock(tt.uuid)
			err := r.Delete(tt.uuid)
			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestAuthorStorage_Update(t *testing.T) {
	logger := logging.GetLogger("../../../../../logs", "test.log")
	db, mock, err := sqlmock.New()
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()
	r := postgres.NewAuthorStorage(db, logger)

	testTable := []struct {
		name          string
		mock          func(dto *domain.UpdateAuthorDTO)
		inputData     *domain.UpdateAuthorDTO
		expectedError bool
	}{
		{
			name:      "OK",
			inputData: domain.TestAuthorUpdateDTO(),
			mock: func(dto *domain.UpdateAuthorDTO) {
				result := sqlmock.NewResult(1, 1)

				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(postgres.UpdateAuthorQuery)).WithArgs(dto.FullName, dto.UUID).WillReturnResult(result)
				mock.ExpectCommit()
			},
			expectedError: false,
		},
		{
			name:      "Invalid UUID or DB error",
			inputData: domain.TestAuthorUpdateDTO(),
			mock: func(dto *domain.UpdateAuthorDTO) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(postgres.UpdateAuthorQuery)).WithArgs(dto.FullName, dto.UUID).WillReturnError(errors.New("DB error"))
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
