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
	"time"
)

func TestAudioStorage_GetOne(t *testing.T) {
	logger := logging.GetLogger("../../../../../logs", "test.log")
	db, mock, err := sqlmock.New()
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()
	r := postgres.NewAudioStorage(db, logger)

	testTable := []struct {
		name           string
		mock           func(uuid string)
		uuid           string
		expectedResult *domain.Audio
		expectError    bool
	}{
		{
			name: "OK",
			uuid: "1",
			mock: func(uuid string) {
				query, _, _ := squirrel.Select(
					"A.uuid",
					"A.title",
					"A.difficulty",
					"A.rating",
					"A.local_url",
					"A.language",
					"A.download_count",
					"A.created_at",
					"D.uuid as direction_uuid",
					"D.name as direction_name",
					"array_agg(DISTINCT T) as tags").
					From("audio AS A").
					LeftJoin("direction AS D ON D.uuid = A.direction_uuid").
					LeftJoin("tag AS T ON  T.uuid = any (A.tags_uuids)").
					Where("A.uuid = ?", uuid).
					PlaceholderFormat(squirrel.Dollar).
					GroupBy("A.uuid", "A.title", "A.difficulty", "A.rating", "A.local_url", "A.language", "A.download_count", "A.created_at", "D.uuid", "D.name").
					ToSql()
				rows := sqlmock.NewRows([]string{"A.uuid", "A.title", "A.difficulty", "A.rating", "A.url", "A.language", "A.download_count", "A.created_at", "D.uuid", "D.name", `{"uuid", "name"}`}).
					AddRow("1", "Test Title", "Test Difficulty", 5.0, "Test URL", "Test Language", 10, time.Date(2000, 01, 01, 0, 0, 0, 0, time.UTC), "1", "Test Direction", `{"1,test tag"}`)

				mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(uuid).WillReturnRows(rows)
			},
			expectedResult: domain.TestAudio(),

			expectError: false,
		},
		{
			name: "no rows in result",
			uuid: "1",
			mock: func(uuid string) {
				query, _, _ := squirrel.Select(
					"A.uuid",
					"A.title",
					"A.difficulty",
					"A.rating",
					"A.local_url",
					"A.language",
					"A.download_count",
					"A.created_at",
					"D.uuid as direction_uuid",
					"D.name as direction_name",
					"array_agg(DISTINCT T) as tags").
					From("audio AS A").
					LeftJoin("direction AS D ON D.uuid = A.direction_uuid").
					LeftJoin("tag AS T ON  T.uuid = any (A.tags_uuids)").
					Where("A.uuid = ?", uuid).
					PlaceholderFormat(squirrel.Dollar).
					GroupBy("A.uuid", "A.title", "A.difficulty", "A.rating", "A.local_url", "A.language", "A.download_count", "A.created_at", "D.uuid", "D.name").
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

func TestAudioStorage_GetAll(t *testing.T) {
	logger := logging.GetLogger("../../../../../logs", "test.log")
	db, mock, err := sqlmock.New()
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()
	r := postgres.NewAudioStorage(db, logger)

	testTable := []struct {
		name           string
		mock           func(sortOptions *domain.SortFilterPagination)
		pages          int
		input          *domain.SortFilterPagination
		expectedResult []*domain.Audio
		expectError    bool
	}{
		{
			name:  "OK",
			input: &domain.SortFilterPagination{},
			mock: func(sortOptions *domain.SortFilterPagination) {
				query, _, _ := squirrel.Select(
					"A.uuid",
					"A.title",
					"A.difficulty",
					"A.rating",
					"A.local_url",
					"A.language",
					"A.download_count",
					"A.created_at",
					"D.uuid as direction_uuid",
					"D.name as direction_name",
					"array_agg(DISTINCT T) as tags",
					"count(*) OVER() AS full_count").
					From("audio AS A").
					LeftJoin("direction AS D ON D.uuid = A.direction_uuid").
					LeftJoin("tag AS T ON  T.uuid = any (A.tags_uuids)").
					GroupBy("A.uuid, A.title, A.difficulty, A.rating, A.local_url, A.language, A.download_count, A.created_at, D.uuid, D.name").
					ToSql()
				rows := sqlmock.NewRows([]string{"A.uuid", "A.title", "A.difficulty", "A.rating", "A.url", "A.language", "A.download_count", "A.created_at", "D.uuid", "D.name", `{"uuid", "name"}`, "full_count"}).
					AddRow("1", "Test Title", "Test Difficulty", 5.0, "Test URL", "Test Language", 10, time.Date(2000, 01, 01, 0, 0, 0, 0, time.UTC), "1", "Test Direction", `{"1,test tag"}`, 1).
					AddRow("1", "Test Title", "Test Difficulty", 5.0, "Test URL", "Test Language", 10, time.Date(2000, 01, 01, 0, 0, 0, 0, time.UTC), "1", "Test Direction", `{"1,test tag"}`, 1)

				mock.ExpectQuery(regexp.QuoteMeta(query)).WillReturnRows(rows)
			},
			expectedResult: []*domain.Audio{domain.TestAudio(), domain.TestAudio()},
			expectError:    false,
		},
		{
			name:  "no rows in result",
			input: &domain.SortFilterPagination{},
			mock: func(sortOptions *domain.SortFilterPagination) {
				query, _, _ := squirrel.Select(
					"A.uuid",
					"A.title",
					"A.difficulty",
					"A.rating",
					"A.local_url",
					"A.language",
					"A.download_count",
					"A.created_at",
					"D.uuid as direction_uuid",
					"D.name as direction_name",
					"array_agg(DISTINCT T) as tags",
					"count(*) OVER() AS full_count").
					From("audio AS A").
					LeftJoin("direction AS D ON D.uuid = A.direction_uuid").
					LeftJoin("tag AS T ON  T.uuid = any (A.tags_uuids)").
					GroupBy("A.uuid, A.title, A.difficulty, A.rating, A.local_url, A.language, A.download_count, A.created_at, D.uuid, D.name").
					ToSql()

				mock.ExpectQuery(regexp.QuoteMeta(query)).WillReturnError(errors.New("no rows in result"))
			},
			expectedResult: nil,
			expectError:    true,
		},
	}
	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock(tt.input)
			result, _, err := r.GetAll(tt.input)
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

func TestAudioStorage_Create(t *testing.T) {
	logger := logging.GetLogger("../../../../../logs", "test.log")
	db, mock, err := sqlmock.New()
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()
	r := postgres.NewAudioStorage(db, logger)

	testTable := []struct {
		name           string
		mock           func(dto *domain.CreateAudioDTO)
		dto            *domain.CreateAudioDTO
		expectedResult string
		expectError    bool
	}{
		{
			name: "OK",
			dto:  domain.TestAudioCreateDTO(),
			mock: func(dto *domain.CreateAudioDTO) {
				rows := sqlmock.NewRows([]string{"uuid"}).AddRow("1")

				localURL := strings.Split(dto.LocalURL, "|split|")
				if len(localURL) < 2 {
					localURL = append(localURL, "")
				}

				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(postgres.CreateAudioQuery)).WithArgs(dto.Title, dto.Difficulty, dto.DirectionUUID, localURL[0], localURL[1], dto.Language, pq.Array(dto.TagsUUIDs)).WillReturnRows(rows)
				mock.ExpectCommit()
			},
			expectedResult: "1",
			expectError:    false,
		},
		{
			name: "no rows in result",
			dto:  domain.TestAudioCreateDTO(),
			mock: func(dto *domain.CreateAudioDTO) {
				localURL := strings.Split(dto.LocalURL, "|split|")
				if len(localURL) < 2 {
					localURL = append(localURL, "")
				}
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(postgres.CreateAudioQuery)).WithArgs(dto.Title, dto.Difficulty, dto.DirectionUUID, localURL[0], localURL[1], dto.Language, pq.Array(dto.TagsUUIDs)).WillReturnError(errors.New("no rows in result"))
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

func TestAudioStorage_Delete(t *testing.T) {
	logger := logging.GetLogger("../../../../../logs", "test.log")
	db, mock, err := sqlmock.New()
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()
	r := postgres.NewAudioStorage(db, logger)

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
				query, _, _ := squirrel.Delete("audio").
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
				query, _, _ := squirrel.Delete("audio").
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
				mock.ExpectExec("DELETE FROM audio WHERE").WithArgs(uuid).WillReturnError(errors.New("DB error"))
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

func TestAudioStorage_Update(t *testing.T) {
	logger := logging.GetLogger("../../../../../logs", "test.log")
	db, mock, err := sqlmock.New()
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()
	r := postgres.NewAudioStorage(db, logger)

	testTable := []struct {
		name          string
		mock          func(dto *domain.UpdateAudioDTO)
		inputData     *domain.UpdateAudioDTO
		expectedError bool
	}{
		{
			name:      "OK",
			inputData: domain.TestAudioUpdateDTO(),
			mock: func(dto *domain.UpdateAudioDTO) {
				mock.ExpectBegin()
				result := sqlmock.NewResult(1, 1)
				mock.ExpectExec(regexp.QuoteMeta(postgres.UpdateAudioQuery)).WithArgs(dto.Title, dto.Difficulty, dto.DirectionUUID, dto.LocalURL, dto.Language, pq.Array(dto.TagsUUIDs), dto.UUID).WillReturnResult(result)
				mock.ExpectCommit()
			},
			expectedError: false,
		},
		{
			name:      "Invalid UUID or DB error",
			inputData: domain.TestAudioUpdateDTO(),
			mock: func(dto *domain.UpdateAudioDTO) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(postgres.UpdateAudioQuery)).WithArgs(dto.Title, dto.Difficulty, dto.DirectionUUID, dto.LocalURL, dto.Language, pq.Array(dto.TagsUUIDs), dto.UUID).WillReturnError(errors.New("DB error"))
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

func TestAudioStorage_Rate(t *testing.T) {
	logger := logging.GetLogger("../../../../../logs", "test.log")
	db, mock, err := sqlmock.New()
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()
	r := postgres.NewAudioStorage(db, logger)

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
				mock.ExpectExec(regexp.QuoteMeta(postgres.RateAudioQuery)).WithArgs(rating, uuid).WillReturnResult(result)
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
				mock.ExpectExec(regexp.QuoteMeta(postgres.RateAudioQuery)).WithArgs(rating, uuid).WillReturnError(errors.New("DB error"))
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

func TestAudioStorage_DownloadCountUp(t *testing.T) {
	logger := logging.GetLogger("../../../../../logs", "test.log")
	db, mock, err := sqlmock.New()
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()
	r := postgres.NewAudioStorage(db, logger)

	testTable := []struct {
		name          string
		mock          func(uuid string)
		uuid          string
		expectedError bool
	}{
		{
			name: "OK",
			uuid: "1",
			mock: func(uuid string) {
				result := sqlmock.NewResult(1, 1)

				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(postgres.AudioDownloadCountUpQuery)).WithArgs(uuid).WillReturnResult(result)
				mock.ExpectCommit()
			},
			expectedError: false,
		},
		{
			name: "Invalid UUID or DB error",
			uuid: "1",
			mock: func(uuid string) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(postgres.AudioDownloadCountUpQuery)).WithArgs(uuid).WillReturnError(errors.New("DB error"))
				mock.ExpectRollback()
			},
			expectedError: true,
		},
	}
	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock(tt.uuid)
			err := r.DownloadCountUp(tt.uuid)
			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
