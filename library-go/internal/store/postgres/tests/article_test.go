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

func TestArticleStorage_GetOne(t *testing.T) {
	logger := logging.GetLogger("../../../../../logs", "test.log")
	db, mock, err := sqlmock.New()
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()
	r := postgres.NewArticleStorage(db, logger)

	testTable := []struct {
		name           string
		mock           func(uuid string)
		uuid           string
		expectedResult *domain.Article
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
					"A.edition_date",
					"A.rating",
					"A.description",
					"A.text",
					"A.local_url",
					"A.image_url",
					"A.web_url",
					"A.language",
					"A.download_count",
					"A.created_at",
					"Au.uuid",
					"Au.full_name",
					"D.uuid as direction_uuid",
					"D.name as direction_name",
					"array_agg(DISTINCT T) as tags").
					From("article AS A").
					LeftJoin("author AS Au ON Au.uuid = A.author_uuid").
					LeftJoin("direction AS D ON D.uuid = A.direction_uuid").
					LeftJoin("tag AS T ON  T.uuid = any (A.tags_uuids)").
					Where("A.uuid = ?", uuid).
					PlaceholderFormat(squirrel.Dollar).
					GroupBy("A.uuid", "A.title", "A.difficulty", "A.edition_date", "A.rating", "A.description", "A.text", "A.local_url", "A.image_url", "A.web_url", "A.language", "A.download_count", "A.created_at", "Au.uuid", "Au.full_name", "D.uuid", "D.name").
					ToSql()

				rows := sqlmock.NewRows([]string{"A.uuid", "A.title", "A.difficulty", "A.edition_date", "A.rating", "A.description", "A.text", "A.local_url", "A.image_url", "A.web_url", "A.language", "A.download_count", "A.created_at", "Au.uuid", "Au.full_name", "D.uuid", "D.name", `{"uuid", "name"}`}).
					AddRow("1", "Test Title", "Test Difficulty", 2000, 5.0, "Test Description", "Test Text", "Test LocalURL", "Test ImageURL", "Test WebURL", "Test Language", 10, time.Date(2000, 01, 01, 0, 0, 0, 0, time.UTC), "1", "Test Author", "1", "Test Direction", `{"1,test tag"}`)

				mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(uuid).WillReturnRows(rows)
			},
			expectedResult: domain.TestArticle(),

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
					"A.edition_date",
					"A.rating",
					"A.description",
					"A.text",
					"A.local_url",
					"A.image_url",
					"A.web_url",
					"A.language",
					"A.download_count",
					"A.created_at",
					"Au.uuid",
					"Au.full_name",
					"D.uuid as direction_uuid",
					"D.name as direction_name",
					"array_agg(DISTINCT T) as tags").
					From("article AS A").
					LeftJoin("author AS Au ON Au.uuid = A.author_uuid").
					LeftJoin("direction AS D ON D.uuid = A.direction_uuid").
					LeftJoin("tag AS T ON  T.uuid = any (A.tags_uuids)").
					Where("A.uuid = ?", uuid).
					PlaceholderFormat(squirrel.Dollar).
					GroupBy("A.uuid", "A.title", "A.difficulty", "A.edition_date", "A.rating", "A.description", "A.text", "A.local_url", "A.image_url", "A.web_url", "A.language", "A.download_count", "A.created_at", "Au.uuid", "Au.full_name", "D.uuid", "D.name").
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

func TestArticleStorage_GetAll(t *testing.T) {
	logger := logging.GetLogger("../../../../../logs", "test.log")
	db, mock, err := sqlmock.New()
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()
	r := postgres.NewArticleStorage(db, logger)

	testTable := []struct {
		name           string
		mock           func(sortOptions *domain.SortFilterPagination)
		sortOptions    *domain.SortFilterPagination
		expectedResult []*domain.Article
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
					"A.uuid",
					"A.title",
					"A.difficulty",
					"A.edition_date",
					"A.rating",
					"A.description",
					"A.text",
					"A.local_url",
					"A.image_url",
					"A.web_url",
					"A.language",
					"A.download_count",
					"A.created_at",
					"Au.uuid",
					"Au.full_name",
					"D.uuid as direction_uuid",
					"D.name as direction_name",
					"array_agg(DISTINCT T) as tags",
					"count(*) OVER() AS full_count").
					From("article AS A").
					LeftJoin("author AS Au ON Au.uuid = A.author_uuid").
					LeftJoin("direction AS D ON D.uuid = A.direction_uuid").
					LeftJoin("tag AS T ON  T.uuid = any (A.tags_uuids)").
					GroupBy("A.uuid", "A.title", "A.difficulty", "A.edition_date", "A.rating", "A.description", "A.text", "A.local_url", "A.image_url", "A.web_url", "A.language", "A.download_count", "A.created_at", "Au.uuid", "Au.full_name", "D.uuid", "D.name").
					ToSql()

				rows := sqlmock.NewRows([]string{"A.uuid", "A.title", "A.difficulty", "A.edition_date", "A.rating", "A.description", "A.text", "A.local_url", "A.image_url", "A.web_url", "A.language", "A.download_count", "A.created_at", "Au.uuid", "Au.full_name", "D.uuid", "D.name", `{"uuid", "name"}`, "full_count"}).
					AddRow("1", "Test Title", "Test Difficulty", 2000, 5.0, "Test Description", "Test Text", "Test LocalURL", "Test ImageURL", "Test WebURL", "Test Language", 10, time.Date(2000, 01, 01, 0, 0, 0, 0, time.UTC), "1", "Test Author", "1", "Test Direction", `{"1,test tag"}`, 1).
					AddRow("1", "Test Title", "Test Difficulty", 2000, 5.0, "Test Description", "Test Text", "Test LocalURL", "Test ImageURL", "Test WebURL", "Test Language", 10, time.Date(2000, 01, 01, 0, 0, 0, 0, time.UTC), "1", "Test Author", "1", "Test Direction", `{"1,test tag"}`, 1)

				mock.ExpectQuery(regexp.QuoteMeta(query)).WillReturnRows(rows)
			},
			expectedResult: []*domain.Article{domain.TestArticle(), domain.TestArticle()},
			expectError:    false,
		},
		{
			name: "no rows in result",
			sortOptions: &domain.SortFilterPagination{
				SortBy:         "",
				Order:          "",
				FiltersAndArgs: nil,
				Limit:          0,
				Page:           0,
			},
			mock: func(sortOptions *domain.SortFilterPagination) {
				query, _, _ := squirrel.Select(
					"A.uuid",
					"A.title",
					"A.difficulty",
					"A.edition_date",
					"A.rating",
					"A.description",
					"A.text",
					"A.local_url",
					"A.image_url",
					"A.web_url",
					"A.language",
					"A.download_count",
					"A.created_at",
					"Au.uuid",
					"Au.full_name",
					"D.uuid as direction_uuid",
					"D.name as direction_name",
					"array_agg(DISTINCT T) as tags",
					"count(*) OVER() AS full_count").
					From("article AS A").
					LeftJoin("author AS Au ON Au.uuid = A.author_uuid").
					LeftJoin("direction AS D ON D.uuid = A.direction_uuid").
					LeftJoin("tag AS T ON  T.uuid = any (A.tags_uuids)").
					GroupBy("A.uuid", "A.title", "A.difficulty", "A.edition_date", "A.rating", "A.description", "A.text", "A.local_url", "A.image_url", "A.web_url", "A.language", "A.download_count", "A.created_at", "Au.uuid", "Au.full_name", "D.uuid", "D.name").
					ToSql()

				mock.ExpectQuery(regexp.QuoteMeta(query)).WillReturnError(errors.New("DB error"))
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

func TestArticleStorage_Create(t *testing.T) {
	logger := logging.GetLogger("../../../../../logs", "test.log")
	db, mock, err := sqlmock.New()
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()
	r := postgres.NewArticleStorage(db, logger)

	testTable := []struct {
		name           string
		mock           func(dto *domain.CreateArticleDTO)
		dto            *domain.CreateArticleDTO
		expectedResult string
		expectError    bool
	}{
		{
			name: "OK",
			dto:  domain.TestArticleCreateDTO(),
			mock: func(dto *domain.CreateArticleDTO) {
				rows := sqlmock.NewRows([]string{"uuid"}).AddRow("1")
				localURL := strings.Split(dto.LocalURL, "|split|")
				imageURL := strings.Split(dto.ImageURL, "|split|")

				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(postgres.CreateArticleQuery)).WithArgs(dto.Title, dto.DirectionUUID, dto.AuthorUUID, dto.Difficulty, dto.EditionDate, dto.Description, dto.Text, localURL[0], localURL[1], dto.WebURL, dto.Language, pq.Array(dto.TagsUUIDs), imageURL[0], imageURL[1]).WillReturnRows(rows)
				mock.ExpectCommit()
			},
			expectedResult: "1",
			expectError:    false,
		},
		{
			name: "Data base error",
			dto:  domain.TestArticleCreateDTO(),
			mock: func(dto *domain.CreateArticleDTO) {
				localURL := strings.Split(dto.LocalURL, "|split|")
				imageURL := strings.Split(dto.ImageURL, "|split|")

				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(postgres.CreateArticleQuery)).WithArgs(dto.Title, dto.DirectionUUID, dto.AuthorUUID, dto.Difficulty, dto.EditionDate, dto.Description, dto.Text, localURL[0], localURL[1], dto.WebURL, dto.Language, pq.Array(dto.TagsUUIDs), imageURL[0], imageURL[1]).WillReturnError(errors.New("DB error"))
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

func TestArticleStorage_Delete(t *testing.T) {
	logger := logging.GetLogger("../../../../../logs", "test.log")
	db, mock, err := sqlmock.New()
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()
	r := postgres.NewArticleStorage(db, logger)

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
				result := sqlmock.NewResult(1, 1)
				query, _, _ := squirrel.Delete("article").
					Where("uuid = ?", uuid).
					PlaceholderFormat(squirrel.Dollar).
					ToSql()

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
				query, _, _ := squirrel.Delete("article").
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

func TestArticleStorage_Update(t *testing.T) {
	logger := logging.GetLogger("../../../../../logs", "test.log")
	db, mock, err := sqlmock.New()
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()
	r := postgres.NewArticleStorage(db, logger)

	testTable := []struct {
		name          string
		mock          func(dto *domain.UpdateArticleDTO)
		inputData     *domain.UpdateArticleDTO
		expectedError bool
	}{
		{
			name:      "OK",
			inputData: domain.TestArticleUpdateDTO(),
			mock: func(dto *domain.UpdateArticleDTO) {
				result := sqlmock.NewResult(1, 1)

				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(postgres.UpdateArticleQuery)).WithArgs(dto.Title, dto.DirectionUUID, dto.AuthorUUID, dto.Difficulty, dto.EditionDate, dto.Description, dto.Text, dto.LocalURL, dto.WebURL, dto.Language, pq.Array(dto.TagsUUIDs), dto.UUID).WillReturnResult(result)
				mock.ExpectCommit()
			},
			expectedError: false,
		},
		{
			name:      "Invalid UUID or DB error",
			inputData: domain.TestArticleUpdateDTO(),
			mock: func(dto *domain.UpdateArticleDTO) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(postgres.UpdateArticleQuery)).WithArgs(dto.Title, dto.DirectionUUID, dto.AuthorUUID, dto.Difficulty, dto.EditionDate, dto.Description, dto.Text, dto.LocalURL, dto.WebURL, dto.Language, pq.Array(dto.TagsUUIDs), dto.UUID).WillReturnError(errors.New("DB error"))
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

func TestArticleStorage_Rate(t *testing.T) {
	logger := logging.GetLogger("../../../../../logs", "test.log")
	db, mock, err := sqlmock.New()
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()
	r := postgres.NewArticleStorage(db, logger)

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
				mock.ExpectExec(regexp.QuoteMeta(postgres.RateArticleQuery)).WithArgs(rating, uuid).WillReturnResult(result)
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
				mock.ExpectExec(regexp.QuoteMeta(postgres.RateArticleQuery)).WithArgs(rating, uuid).WillReturnError(errors.New("DB error"))
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

func TestArticleStorage_DownloadCountUp(t *testing.T) {
	logger := logging.GetLogger("../../../../../logs", "test.log")
	db, mock, err := sqlmock.New()
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()
	r := postgres.NewArticleStorage(db, logger)

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
				mock.ExpectExec(regexp.QuoteMeta(postgres.ArticleDownloadCountUpQuery)).WithArgs(uuid).WillReturnResult(result)
				mock.ExpectCommit()
			},
			expectedError: false,
		},
		{
			name: "Invalid UUID or DB error",
			uuid: "1",
			mock: func(uuid string) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(postgres.ArticleDownloadCountUpQuery)).WithArgs(uuid).WillReturnError(errors.New("DB error"))
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
