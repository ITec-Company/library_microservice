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

func TestVideoStorage_GetOne(t *testing.T) {
	logger := logging.GetLogger("../../../../../logs", "test.log")
	db, mock, err := sqlmock.New()
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()
	r := postgres.NewVideoStorage(db, logger)

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
				query, _, _ := squirrel.Select(
					"V.uuid",
					"V.title",
					"V.difficulty",
					"V.rating",
					"V.description",
					"V.local_url",
					"V.web_url",
					"V.language",
					"V.download_count",
					"V.created_at",
					"D.uuid as direction_uuid",
					"D.name as direction_name",
					"array_agg(DISTINCT T) as tags").
					From("video AS V").
					Where("V.uuid = ?", uuid).
					PlaceholderFormat(squirrel.Dollar).
					LeftJoin("direction AS D ON D.uuid = V.direction_uuid").
					LeftJoin("tag AS T ON  T.uuid = any (V.tags_uuids)").
					GroupBy("V.uuid, V.title, V.difficulty, V.rating, V.description, V.local_url, V.web_url, V.language, V.download_count, V.created_at, D.uuid, D.name").
					ToSql()

				rows := sqlmock.NewRows([]string{"V.uuid", "V.title", "V.difficulty", "V.rating", "V.description", "V.local_url", "V.web_url", "V.language", "V.download_count", "V.created_at", "D.uuid", "D.name", `{"uuid", "name"}`}).
					AddRow("1", "Test Title", "Test Difficulty", 5.0, "Test Description", "Test LocalURL", "Test WebURL", "Test Language", 10, time.Date(2000, 01, 01, 0, 0, 0, 0, time.UTC), "1", "Test Direction", `{"1,test tag"}`)

				mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(uuid).WillReturnRows(rows)
			},
			expectedResult: domain.TestVideo(),

			expectError: false,
		},
		{
			name: "no rows in result",
			uuid: "1",
			mock: func(uuid string) {
				query, _, _ := squirrel.Select(
					"V.uuid",
					"V.title",
					"V.difficulty",
					"V.rating",
					"V.description",
					"V.local_url",
					"V.web_url",
					"V.language",
					"V.download_count",
					"V.created_at",
					"D.uuid as direction_uuid",
					"D.name as direction_name",
					"array_agg(DISTINCT T) as tags").
					From("video AS V").
					Where("V.uuid = ?", uuid).
					PlaceholderFormat(squirrel.Dollar).
					LeftJoin("direction AS D ON D.uuid = V.direction_uuid").
					LeftJoin("tag AS T ON  T.uuid = any (V.tags_uuids)").
					GroupBy("V.uuid, V.title, V.difficulty, V.rating, V.description, V.local_url, V.web_url, V.language, V.download_count, V.created_at, D.uuid, D.name").
					ToSql()

				mock.ExpectQuery(regexp.QuoteMeta(query)).WillReturnError(errors.New("no rows in result"))
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
	logger := logging.GetLogger("../../../../../logs", "test.log")
	db, mock, err := sqlmock.New()
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()
	r := postgres.NewVideoStorage(db, logger)

	testTable := []struct {
		name           string
		sortOptions    *domain.SortFilterPagination
		mock           func(sortOptions *domain.SortFilterPagination)
		expectedResult []*domain.Video
		expectError    bool
	}{
		{
			name: "OK",
			sortOptions: &domain.SortFilterPagination{
				SortBy:         "",
				Order:          "",
				FiltersAndArgs: nil,
				Limit:          0,
				Page:           0,
			},
			mock: func(sortOptions *domain.SortFilterPagination) {
				query, _, _ := squirrel.Select(
					"V.uuid",
					"V.title",
					"V.difficulty",
					"V.rating",
					"V.description",
					"V.local_url",
					"V.web_url",
					"V.language",
					"V.download_count",
					"V.created_at",
					"D.uuid as direction_uuid",
					"D.name as direction_name",
					"array_agg(DISTINCT T) as tags",
					"count(*) OVER() AS full_count").
					From("video AS V").
					LeftJoin("direction AS D ON D.uuid = V.direction_uuid").
					LeftJoin("tag AS T ON  T.uuid = any (V.tags_uuids)").
					GroupBy("V.uuid, V.title, V.difficulty, V.rating, V.description, V.local_url, V.web_url, V.language, V.download_count, V.created_at, D.uuid, D.name").
					ToSql()

				rows := sqlmock.NewRows([]string{"V.uuid", "V.title", "V.difficulty", "V.rating", "V.description", "V.local_url", "V.web_url", "V.language", "V.download_count", "V.created_at", "D.uuid", "D.name", `{"uuid", "name"}`, "full_count"}).
					AddRow("1", "Test Title", "Test Difficulty", 5.0, "Test Description", "Test LocalURL", "Test WebURL", "Test Language", 10, time.Date(2000, 01, 01, 0, 0, 0, 0, time.UTC), "1", "Test Direction", `{"1,test tag"}`, 1).
					AddRow("1", "Test Title", "Test Difficulty", 5.0, "Test Description", "Test LocalURL", "Test WebURL", "Test Language", 10, time.Date(2000, 01, 01, 0, 0, 0, 0, time.UTC), "1", "Test Direction", `{"1,test tag"}`, 1)

				mock.ExpectQuery(regexp.QuoteMeta(query)).WillReturnRows(rows)
			},
			expectedResult: []*domain.Video{domain.TestVideo(), domain.TestVideo()},
			expectError:    false,
		},
		{
			name: "No rows in result",
			sortOptions: &domain.SortFilterPagination{
				SortBy:         "",
				Order:          "",
				FiltersAndArgs: nil,
				Limit:          0,
				Page:           0,
			},
			mock: func(sortOptions *domain.SortFilterPagination) {
				query, _, _ := squirrel.Select(
					"V.uuid",
					"V.title",
					"V.difficulty",
					"V.rating",
					"V.description",
					"V.local_url",
					"V.web_url",
					"V.language",
					"V.download_count",
					"V.created_at",
					"D.uuid as direction_uuid",
					"D.name as direction_name",
					"array_agg(DISTINCT T) as tags",
					"count(*) OVER() AS full_count").
					From("video AS V").
					LeftJoin("direction AS D ON D.uuid = V.direction_uuid").
					LeftJoin("tag AS T ON  T.uuid = any (V.tags_uuids)").
					GroupBy("V.uuid, V.title, V.difficulty, V.rating, V.description, V.local_url, V.web_url, V.language, V.download_count, V.created_at, D.uuid, D.name").
					ToSql()

				mock.ExpectQuery(regexp.QuoteMeta(query)).WillReturnError(errors.New("no rows in result"))
			},
			expectedResult: nil,
			expectError:    true,
		},
	}
	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock(tt.sortOptions)
			result, _, err := r.GetAll(tt.sortOptions)
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
	logger := logging.GetLogger("../../../../../logs", "test.log")
	db, mock, err := sqlmock.New()
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()
	r := postgres.NewVideoStorage(db, logger)

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
				rows := sqlmock.NewRows([]string{"uuid"}).AddRow("1")
				localURL := strings.Split(dto.LocalURL, "|split|")
				if len(localURL) < 2 {
					localURL = append(localURL, "")
				}

				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(postgres.CreateVideoQuery)).WithArgs(dto.Title, dto.Difficulty, dto.DirectionUUID, dto.Description, localURL[0], localURL[1], dto.WebURL, dto.Language, pq.Array(dto.TagsUUIDs)).WillReturnRows(rows)
				mock.ExpectCommit()
			},
			expectedResult: "1",
			expectError:    false,
		},
		{
			name: "no rows in result",
			dto:  domain.TestVideoCreateDTO(),
			mock: func(dto *domain.CreateVideoDTO) {
				localURL := strings.Split(dto.LocalURL, "|split|")
				if len(localURL) < 2 {
					localURL = append(localURL, "")
				}

				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(postgres.CreateVideoQuery)).WithArgs(dto.Title, dto.Difficulty, dto.DirectionUUID, dto.Description, localURL[0], localURL[1], dto.WebURL, dto.Language, pq.Array(dto.TagsUUIDs)).WillReturnError(errors.New("no rows in result"))
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
	logger := logging.GetLogger("../../../../../logs", "test.log")
	db, mock, err := sqlmock.New()
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()
	r := postgres.NewVideoStorage(db, logger)

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
				query, _, _ := squirrel.Delete("video").
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
				result := sqlmock.NewResult(0, 0)
				query, _, _ := squirrel.Delete("video").
					Where("uuid = ?", uuid).
					PlaceholderFormat(squirrel.Dollar).
					ToSql()

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

func TestVideoStorage_Update(t *testing.T) {
	logger := logging.GetLogger("../../../../../logs", "test.log")
	db, mock, err := sqlmock.New()
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()
	r := postgres.NewVideoStorage(db, logger)

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
				result := sqlmock.NewResult(1, 1)

				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(postgres.UpdateVideoQuery)).WithArgs(dto.Title, dto.Difficulty, dto.DirectionUUID, dto.Description, dto.LocalURL, dto.WebURL, dto.Language, pq.Array(dto.TagsUUIDs), dto.UUID).WillReturnResult(result)
				mock.ExpectCommit()
			},
			expectedError: false,
		},
		{
			name:      "Invalid UUID or DB error",
			inputData: domain.TestVideoUpdateDTO(),
			mock: func(dto *domain.UpdateVideoDTO) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(postgres.UpdateVideoQuery)).WithArgs(dto.Title, dto.Difficulty, dto.DirectionUUID, dto.Description, dto.LocalURL, dto.WebURL, dto.Language, pq.Array(dto.TagsUUIDs), dto.UUID).WillReturnError(errors.New("DB error"))
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

func TestVideoStorage_Rate(t *testing.T) {
	logger := logging.GetLogger("../../../../../logs", "test.log")
	db, mock, err := sqlmock.New()
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()
	r := postgres.NewVideoStorage(db, logger)

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
				mock.ExpectExec(regexp.QuoteMeta(postgres.RateVideoQuery)).WithArgs(rating, uuid).WillReturnResult(result)
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
				mock.ExpectExec(regexp.QuoteMeta(postgres.RateVideoQuery)).WithArgs(rating, uuid).WillReturnError(errors.New("DB error"))
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

func TestVideoStorage_DownloadCountUp(t *testing.T) {
	logger := logging.GetLogger("../../../../../logs", "test.log")
	db, mock, err := sqlmock.New()
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()
	r := postgres.NewVideoStorage(db, logger)

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
				mock.ExpectExec(regexp.QuoteMeta(postgres.VideoDownloadCountUpQuery)).WithArgs(uuid).WillReturnResult(result)
				mock.ExpectCommit()
			},
			expectedError: false,
		},
		{
			name: "Invalid UUID or DB error",
			uuid: "1",
			mock: func(uuid string) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(postgres.VideoDownloadCountUpQuery)).WithArgs(uuid).WillReturnError(errors.New("DB error"))
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
