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

func TestVideoStorage_GetOne(t *testing.T) {
	logger := logging.GetLogger("../../../../logs", "test.log")
	db, mock, err := sqlmock.New()
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()
	r := NewVideoStorage(db, logger)

	testTable := []struct {
		name           string
		mock           func(uuid string)
		uuid           string
		expectedResult *domain.Video
		expectError    bool
	}{
		{
			name: "OK",
			uuid: "1",
			mock: func(uuid string) {
				rows := sqlmock.NewRows([]string{"V.uuid", "V.title", "V.difficulty", "V.rating", "V.url", "V.language", "V.download_count", "D.uuid", "D.name", `{"uuid", "name"}`}).
					AddRow("1", "Test Title", "Test Difficulty", 5.0, "Test URL", "Test Language", 10, "1", "Test Direction", `{"1,Test Tag"}`)
				mock.ExpectQuery(regexp.QuoteMeta(getOneVideoQuery)).WithArgs(uuid).WillReturnRows(rows)
			},
			expectedResult: domain.TestVideo(),

			expectError: false,
		},
		{
			name: "DB error",
			uuid: "1",
			mock: func(uuid string) {
				mock.ExpectQuery(regexp.QuoteMeta(getOneVideoQuery)).
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

func TestVideoStorage_GetAll(t *testing.T) {
	logger := logging.GetLogger("../../../../logs", "test.log")
	db, mock, err := sqlmock.New()
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()
	r := NewVideoStorage(db, logger)

	testTable := []struct {
		name           string
		mock           func(page, limit int)
		inputPage      int
		inputLimit     int
		expectedResult []*domain.Video
		expectError    bool
	}{
		{
			name:       "OK",
			inputPage:  0,
			inputLimit: 0,
			mock: func(page, limit int) {
				rows := sqlmock.NewRows([]string{"V.uuid", "V.title", "V.difficulty", "V.rating", "V.url", "V.language", "V.download_count", "D.uuid", "D.name", `{"uuid", "name"}`}).
					AddRow("1", "Test Title", "Test Difficulty", 5.0, "Test URL", "Test Language", 10, "1", "Test Direction", `{"1,Test Tag"}`).
					AddRow("1", "Test Title", "Test Difficulty", 5.0, "Test URL", "Test Language", 10, "1", "Test Direction", `{"1,Test Tag"}`)
				mock.ExpectQuery(regexp.QuoteMeta(getAllVideosQuery)).WillReturnRows(rows)
			},
			expectedResult: []*domain.Video{domain.TestVideo(), domain.TestVideo()},
			expectError:    false,
		},
		{
			name:       "Data base error",
			inputPage:  0,
			inputLimit: 0,
			mock: func(page, limit int) {
				mock.ExpectQuery(regexp.QuoteMeta(getAllVideosQuery)).WillReturnError(errors.New("DB error"))
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

func TestVideoStorage_Create(t *testing.T) {
	logger := logging.GetLogger("../../../../logs", "test.log")
	db, mock, err := sqlmock.New()
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()
	r := NewVideoStorage(db, logger)

	testTable := []struct {
		name           string
		mock           func(dto *domain.CreateVideoDTO)
		dto            *domain.CreateVideoDTO
		expectedResult string
		expectError    bool
	}{
		{
			name: "OK",
			dto:  domain.TestVideoCreateDTO(),
			mock: func(dto *domain.CreateVideoDTO) {
				mock.ExpectBegin()
				rows := sqlmock.NewRows([]string{"uuid"}).AddRow("1")
				mock.ExpectQuery(regexp.QuoteMeta(createVideoQuery)).WithArgs(dto.Title, dto.Difficulty, dto.DirectionUUID, 0, dto.URL, dto.Language, pq.Array(dto.TagsUUIDs), 0).WillReturnRows(rows)
				mock.ExpectCommit()
			},
			expectedResult: "1",
			expectError:    false,
		},
		{
			name: "Data base error",
			dto:  domain.TestVideoCreateDTO(),
			mock: func(dto *domain.CreateVideoDTO) {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(createVideoQuery)).WithArgs(dto.Title, dto.Difficulty, dto.DirectionUUID, 0, dto.URL, dto.Language, pq.Array(dto.TagsUUIDs), 0).WillReturnError(errors.New("DB error"))
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

func TestVideoStorage_Delete(t *testing.T) {
	logger := logging.GetLogger("../../../../logs", "test.log")
	db, mock, err := sqlmock.New()
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()
	r := NewVideoStorage(db, logger)

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
				mock.ExpectExec(regexp.QuoteMeta(deleteVideoQuery)).WithArgs(uuid).WillReturnResult(result)
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
				mock.ExpectExec(regexp.QuoteMeta(deleteVideoQuery)).WithArgs(uuid).WillReturnResult(result)
				mock.ExpectRollback()
			},
			expectedError: true,
		},
		{
			name:      "DB error",
			inputData: "1",
			mock: func(uuid string) {
				mock.ExpectBegin()
				mock.ExpectExec("DELETE FROM video WHERE").WithArgs(uuid).WillReturnError(errors.New("DB error"))
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

func TestVideoStorage_Update(t *testing.T) {
	logger := logging.GetLogger("../../../../logs", "test.log")
	db, mock, err := sqlmock.New()
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()
	r := NewVideoStorage(db, logger)

	testTable := []struct {
		name          string
		mock          func(dto *domain.UpdateVideoDTO)
		inputData     *domain.UpdateVideoDTO
		expectedError bool
	}{
		{
			name:      "OK",
			inputData: domain.TestVideoUpdateDTO(),
			mock: func(dto *domain.UpdateVideoDTO) {
				mock.ExpectBegin()
				result := sqlmock.NewResult(1, 1)
				mock.ExpectExec(regexp.QuoteMeta(updateVideoQuery)).WithArgs(dto.Title, dto.Difficulty, dto.DirectionUUID, dto.Rating, dto.URL, dto.Language, pq.Array(dto.TagsUUIDs), dto.DownloadCount, dto.UUID).WillReturnResult(result)
				mock.ExpectCommit()
			},
			expectedError: false,
		},
		{
			name:      "Invalid UUID or DB error",
			inputData: domain.TestVideoUpdateDTO(),
			mock: func(dto *domain.UpdateVideoDTO) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(updateVideoQuery)).WithArgs(dto.Title, dto.Difficulty, dto.DirectionUUID, dto.Rating, dto.URL, dto.Language, pq.Array(dto.TagsUUIDs), dto.DownloadCount, dto.UUID).WillReturnError(errors.New("DB error"))
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
