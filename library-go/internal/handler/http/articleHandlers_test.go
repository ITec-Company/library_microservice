package http

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
	"io"
	"library-go/internal/domain"
	mock_service "library-go/internal/service/mocks"
	"library-go/pkg/logging"
	"mime/multipart"
	"net/http/httptest"
	"os"
	"testing"
)

func TestArticleHandler_GetAll(t *testing.T) {
	type mockBehavior func(s *mock_service.MockArticleService, ctx context.Context, limit, offset int)

	testTable := []struct {
		name, input         string
		limit, offset       int
		ctx                 context.Context
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:   "OK",
			input:  "limit=0&offset=0",
			ctx:    context.Background(),
			limit:  0,
			offset: 0,
			mockBehavior: func(s *mock_service.MockArticleService, ctx context.Context, limit, offset int) {
				s.EXPECT().GetAll(ctx, limit, offset).Return([]*domain.Article{domain.TestArticle(), domain.TestArticle()}, nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: fmt.Sprintln(`[{"uuid":"1","title":"Test Title","direction":{"uuid":"1","name":"test direction"},"difficulty":"Test Difficulty","author":{"uuid":"1","full_name":"test Author"},"edition_date":"2000-01-01T00:00:00Z","rating":5,"description":"Test Description","url":"Test URL","language":"Test Language","tags":[{"uuid":"1","name":"test Tag"}],"download_count":10},{"uuid":"1","title":"Test Title","direction":{"uuid":"1","name":"test direction"},"difficulty":"Test Difficulty","author":{"uuid":"1","full_name":"test Author"},"edition_date":"2000-01-01T00:00:00Z","rating":5,"description":"Test Description","url":"Test URL","language":"Test Language","tags":[{"uuid":"1","name":"test Tag"}],"download_count":10}]`),
		},
		{
			name:   "OK empty input",
			input:  "",
			ctx:    context.Background(),
			limit:  0,
			offset: 0,
			mockBehavior: func(s *mock_service.MockArticleService, ctx context.Context, limit, offset int) {
				s.EXPECT().GetAll(ctx, limit, offset).Return([]*domain.Article{domain.TestArticle(), domain.TestArticle()}, nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: fmt.Sprintln(`[{"uuid":"1","title":"Test Title","direction":{"uuid":"1","name":"test direction"},"difficulty":"Test Difficulty","author":{"uuid":"1","full_name":"test Author"},"edition_date":"2000-01-01T00:00:00Z","rating":5,"description":"Test Description","url":"Test URL","language":"Test Language","tags":[{"uuid":"1","name":"test Tag"}],"download_count":10},{"uuid":"1","title":"Test Title","direction":{"uuid":"1","name":"test direction"},"difficulty":"Test Difficulty","author":{"uuid":"1","full_name":"test Author"},"edition_date":"2000-01-01T00:00:00Z","rating":5,"description":"Test Description","url":"Test URL","language":"Test Language","tags":[{"uuid":"1","name":"test Tag"}],"download_count":10}]`),
		},
		{
			name:   "OK invalid input",
			input:  "limit=one&offset=-10",
			ctx:    context.Background(),
			limit:  0,
			offset: 0,
			mockBehavior: func(s *mock_service.MockArticleService, ctx context.Context, limit, offset int) {
				s.EXPECT().GetAll(ctx, limit, offset).Return([]*domain.Article{domain.TestArticle(), domain.TestArticle()}, nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: fmt.Sprintln(`[{"uuid":"1","title":"Test Title","direction":{"uuid":"1","name":"test direction"},"difficulty":"Test Difficulty","author":{"uuid":"1","full_name":"test Author"},"edition_date":"2000-01-01T00:00:00Z","rating":5,"description":"Test Description","url":"Test URL","language":"Test Language","tags":[{"uuid":"1","name":"test Tag"}],"download_count":10},{"uuid":"1","title":"Test Title","direction":{"uuid":"1","name":"test direction"},"difficulty":"Test Difficulty","author":{"uuid":"1","full_name":"test Author"},"edition_date":"2000-01-01T00:00:00Z","rating":5,"description":"Test Description","url":"Test URL","language":"Test Language","tags":[{"uuid":"1","name":"test Tag"}],"download_count":10}]`),
		},
		{
			name:   "No rows in result",
			input:  "limit=0&offset=0",
			ctx:    context.Background(),
			limit:  0,
			offset: 0,
			mockBehavior: func(s *mock_service.MockArticleService, ctx context.Context, limit, offset int) {
				s.EXPECT().GetAll(ctx, limit, offset).Return(nil, errors.New("now rows in result"))
			},
			expectedStatusCode:  500,
			expectedRequestBody: fmt.Sprintf(`{"ErrorMsg":"error occurred while getting all articles. err: %v"}%v`, errors.New("now rows in result"), "\n"),
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			service := mock_service.NewMockArticleService(c)
			testCase.mockBehavior(service, testCase.ctx, testCase.limit, testCase.offset)

			ArticleHandler := NewArticleHandler(service, logging.GetLogger())

			router := httprouter.New()
			ArticleHandler.Register(router)

			w := httptest.NewRecorder()

			req := httptest.NewRequest("GET", fmt.Sprintf("/articles?%s", testCase.input), nil)

			router.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedRequestBody, w.Body.String())
		})
	}
}

func TestArticleHandler_GetByUUID(t *testing.T) {
	type mockBehavior func(s *mock_service.MockArticleService, ctx context.Context, uuid string)

	testTable := []struct {
		name, uuid, input   string
		ctx                 context.Context
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:  "OK",
			input: "1",
			uuid:  "1",
			ctx:   context.Background(),
			mockBehavior: func(s *mock_service.MockArticleService, ctx context.Context, uuid string) {
				s.EXPECT().GetByUUID(ctx, uuid).Return(domain.TestArticle(), nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: "{\"uuid\":\"1\",\"title\":\"Test Title\",\"direction\":{\"uuid\":\"1\",\"name\":\"test Direction\"},\"difficulty\":\"Test Difficulty\",\"author\":{\"uuid\":\"1\",\"full_name\":\"test Author\"},\"edition_date\":\"2000-01-01T00:00:00Z\",\"rating\":5,\"description\":\"Test Description\",\"url\":\"Test URL\",\"language\":\"Test Language\",\"tags\":[{\"uuid\":\"1\",\"name\":\"test Tag\"}],\"download_count\":10}\n",
		},
		{
			name:                "invalid uuid",
			input:               "-1",
			uuid:                "-1",
			ctx:                 context.Background(),
			mockBehavior:        func(s *mock_service.MockArticleService, ctx context.Context, uuid string) {},
			expectedStatusCode:  400,
			expectedRequestBody: "{\"ErrorMsg\":\"uuid need to be uint\"}\n",
		},
		{
			name:  "no rows in result",
			input: "1",
			uuid:  "1",
			ctx:   context.Background(),
			mockBehavior: func(s *mock_service.MockArticleService, ctx context.Context, uuid string) {
				s.EXPECT().GetByUUID(ctx, uuid).Return(nil, errors.New("now rows in result"))
			},
			expectedStatusCode:  404,
			expectedRequestBody: "{\"ErrorMsg\":\"error occurred while getting article from DB by UUID. err: now rows in result\"}\n",
		},
		{
			name:                "string input",
			input:               "one",
			uuid:                "one",
			ctx:                 context.Background(),
			mockBehavior:        func(s *mock_service.MockArticleService, ctx context.Context, uuid string) {},
			expectedStatusCode:  400,
			expectedRequestBody: "{\"ErrorMsg\":\"uuid need to be uint\"}\n",
		},
		// is it correct ?????
		{
			name:                "empty input",
			input:               "",
			uuid:                "",
			ctx:                 context.Background(),
			mockBehavior:        func(s *mock_service.MockArticleService, ctx context.Context, uuid string) {},
			expectedStatusCode:  301,
			expectedRequestBody: "<a href=\"/article\">Moved Permanently</a>.\n\n",
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			service := mock_service.NewMockArticleService(c)
			testCase.mockBehavior(service, testCase.ctx, testCase.uuid)

			ArticleHandler := NewArticleHandler(service, logging.GetLogger())

			router := httprouter.New()
			ArticleHandler.Register(router)

			w := httptest.NewRecorder()

			req := httptest.NewRequest("GET", fmt.Sprintf("/article/%s", testCase.input), nil)

			router.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedRequestBody, w.Body.String())
		})
	}
}

func TestArticleHandler_Create(t *testing.T) {
	type mockBehavior func(s *mock_service.MockArticleService, ctx context.Context, createArticleDTO domain.CreateArticleDTO)

	testTable := []struct {
		name                string
		ctx                 context.Context
		mockBehavior        mockBehavior
		createArticleDTO    *domain.CreateArticleDTO
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:             "OK",
			ctx:              context.Background(),
			createArticleDTO: domain.TestArticleCreateDTO(),
			mockBehavior: func(s *mock_service.MockArticleService, ctx context.Context, createArticleDTO domain.CreateArticleDTO) {
				s.EXPECT().Create(ctx, createArticleDTO).Return("1", nil)
			},
			expectedStatusCode:  201,
			expectedRequestBody: "Article created successfully. UUID: 1",
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			service := mock_service.NewMockArticleService(c)
			testCase.mockBehavior(service, testCase.ctx, *testCase.createArticleDTO)

			ArticleHandler := NewArticleHandler(service, logging.GetLogger())

			router := httprouter.New()
			ArticleHandler.Register(router)

			body := bytes.Buffer{}
			writer := multipart.NewWriter(&body)
			file, _ := os.Open("test_file.txt")

			part, _ := writer.CreateFormFile("file", "file")
			io.Copy(part, file)
			defer file.Close()

			writer.WriteField("title", "test title")
			writer.WriteField("direction_uuid", "1")
			writer.WriteField("author_uuid", "1")
			writer.WriteField("difficulty", "test difficulty")
			writer.WriteField("edition_date", "2000-01-01 00:00:00 +0000 UTC")
			writer.WriteField("description", "test description")
			writer.WriteField("language", "test language")
			writer.WriteField("tags_uuids", `[1]`)

			w := httptest.NewRecorder()

			writer.Close()
			req := httptest.NewRequest("POST", "/article", bytes.NewReader(body.Bytes()))
			req.Header.Set("Content-Type", writer.FormDataContentType())

			router.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedRequestBody, w.Body.String())
		})
	}

}

func TestArticleHandler_Delete(t *testing.T) {
	type mockBehavior func(s *mock_service.MockArticleService, ctx context.Context, uuid string)

	testTable := []struct {
		name, uuid, input   string
		ctx                 context.Context
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:  "OK",
			input: "1",
			uuid:  "1",
			ctx:   context.Background(),
			mockBehavior: func(s *mock_service.MockArticleService, ctx context.Context, uuid string) {
				s.EXPECT().Delete(ctx, uuid).Return(nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: "{\"infoMsg\":\"Article with UUID 1 was deleted\"}\n",
		},
		{
			name:  "no rows in result",
			input: "1",
			uuid:  "1",
			ctx:   context.Background(),
			mockBehavior: func(s *mock_service.MockArticleService, ctx context.Context, uuid string) {
				s.EXPECT().Delete(ctx, uuid).Return(errors.New("no rows affected"))
			},
			expectedStatusCode:  404,
			expectedRequestBody: "{\"ErrorMsg\":\"error occurred while deleting article from DB. err: no rows affected\"}\n",
		},
		{
			name:                "invalid uuid",
			input:               "-1",
			uuid:                "-1",
			ctx:                 context.Background(),
			mockBehavior:        func(s *mock_service.MockArticleService, ctx context.Context, uuid string) {},
			expectedStatusCode:  400,
			expectedRequestBody: "{\"ErrorMsg\":\"uuid need to be uint\"}\n",
		},
		{
			name:                "string uuid",
			input:               "one",
			uuid:                "one",
			ctx:                 context.Background(),
			mockBehavior:        func(s *mock_service.MockArticleService, ctx context.Context, uuid string) {},
			expectedStatusCode:  400,
			expectedRequestBody: "{\"ErrorMsg\":\"uuid need to be uint\"}\n",
		},
		// is it correct ?????
		{
			name:                "empty input",
			input:               "",
			uuid:                "",
			ctx:                 context.Background(),
			mockBehavior:        func(s *mock_service.MockArticleService, ctx context.Context, uuid string) {},
			expectedStatusCode:  404,
			expectedRequestBody: "404 page not found\n",
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			service := mock_service.NewMockArticleService(c)
			testCase.mockBehavior(service, testCase.ctx, testCase.uuid)

			ArticleHandler := NewArticleHandler(service, logging.GetLogger())

			router := httprouter.New()
			ArticleHandler.Register(router)

			w := httptest.NewRecorder()

			req := httptest.NewRequest("DELETE", fmt.Sprintf("/article/%s", testCase.input), nil)

			router.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedRequestBody, w.Body.String())
		})
	}
}

func TestArticleHandler_Update(t *testing.T) {
	type mockBehavior func(s *mock_service.MockArticleService, ctx context.Context, dto *domain.UpdateArticleDTO)

	testTable := []struct {
		name                string
		inputBodyJSON       map[string]interface{}
		ctx                 context.Context
		dto                 domain.UpdateArticleDTO
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name: "OK",
			inputBodyJSON: map[string]interface{}{
				"uuid":           "1",
				"title":          "Test Title",
				"direction_uuid": "1",
				"author_uuid":    "1",
				"difficulty":     "Test Difficulty",
				"edition_date":   "2000-01-01T00:00:00Z",
				"rating":         5.5,
				"description":    "Test Description",
				"url":            "Test URL",
				"language":       "Test Language",
				"tags_uuids":     []string{"1"},
				"download_count": 10,
			},
			ctx: context.Background(),
			dto: *domain.TestArticleUpdateDTO(),
			mockBehavior: func(s *mock_service.MockArticleService, ctx context.Context, dto *domain.UpdateArticleDTO) {
				s.EXPECT().Update(ctx, dto).Return(nil)
			},
			expectedStatusCode:  201,
			expectedRequestBody: "{\"infoMsg\":\"Article updated successfully\"}\n",
		},
		{
			name: "no rows in result",
			inputBodyJSON: map[string]interface{}{
				"uuid":           "1",
				"title":          "Test Title",
				"direction_uuid": "1",
				"author_uuid":    "1",
				"difficulty":     "Test Difficulty",
				"edition_date":   "2000-01-01T00:00:00Z",
				"rating":         5.5,
				"description":    "Test Description",
				"url":            "Test URL",
				"language":       "Test Language",
				"tags_uuids":     []string{"1"},
				"download_count": 10,
			},
			ctx: context.Background(),
			dto: *domain.TestArticleUpdateDTO(),
			mockBehavior: func(s *mock_service.MockArticleService, ctx context.Context, dto *domain.UpdateArticleDTO) {
				s.EXPECT().Update(ctx, dto).Return(errors.New("no rows affected"))
			},
			expectedStatusCode:  500,
			expectedRequestBody: "{\"ErrorMsg\":\"error occurred while updating article into DB. err: no rows affected\"}\n",
		},
		{
			name:                "empty input body JSON or nil UUID",
			inputBodyJSON:       map[string]interface{}{},
			ctx:                 context.Background(),
			dto:                 *domain.TestArticleUpdateDTO(),
			mockBehavior:        func(s *mock_service.MockArticleService, ctx context.Context, dto *domain.UpdateArticleDTO) {},
			expectedStatusCode:  400,
			expectedRequestBody: "{\"ErrorMsg\":\"error occurred while decoding JSON request. err: nil UUID\"}\n",
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			service := mock_service.NewMockArticleService(c)
			testCase.mockBehavior(service, testCase.ctx, &testCase.dto)

			ArticleHandler := NewArticleHandler(service, logging.GetLogger())

			router := httprouter.New()
			ArticleHandler.Register(router)

			body, _ := json.Marshal(testCase.inputBodyJSON)
			req := httptest.NewRequest("PUT", fmt.Sprintf("/article"), bytes.NewReader(body))

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedRequestBody, w.Body.String())
		})
	}
}

//func TestArticleHandler_Load(t *testing.T) {
//	type mockBehavior func(s *mock_service.MockArticleService, ctx context.Context, uuid string)
//
//	testTable := []struct {
//		name, uuid, input   string
//		ctx                 context.Context
//		mockBehavior        mockBehavior
//		expectedStatusCode  int
//		expectedRequestBody string
//	}{
//		{
//			name:  "OK",
//			input: "1",
//			uuid:  "1",
//			ctx:   context.Background(),
//			mockBehavior: func(s *mock_service.MockArticleService, ctx context.Context, uuid string) {
//				s.EXPECT().GetByUUID(ctx, uuid).Return(domain.TestArticle(), nil)
//			},
//			expectedStatusCode:  200,
//			expectedRequestBody: fmt.Sprintln(`{"uuid":"1","title":"Test Title","direction":{"uuid":"1","name":"test direction"},"difficulty":"Test Difficulty","author":{"uuid":"1","full_name":"test Author"},"edition_date":"2000-01-01T00:00:00Z","rating":5,"description":"Test Description","url":"Test URL","language":"Test Language","tags":[{"uuid":"1","name":"test Tag"}],"download_count":10}`),
//		},
//	}
//	for _, testCase := range testTable {
//		t.Run(testCase.name, func(t *testing.T) {
//			c := gomock.NewController(t)
//			defer c.Finish()
//
//			service := mock_service.NewMockArticleService(c)
//			testCase.mockBehavior(service, testCase.ctx, testCase.uuid)
//
//			ArticleHandler := NewArticleHandler(service, logging.GetLogger())
//
//			router := httprouter.New()
//			ArticleHandler.Register(router)
//
//			w := httptest.NewRecorder()
//
//			req := httptest.NewRequest("GET", fmt.Sprintf("/article/%s", testCase.input), nil)
//
//			router.ServeHTTP(w, req)
//
//			assert.Equal(t, testCase.expectedStatusCode, w.Code)
//			assert.Equal(t, testCase.expectedRequestBody, w.Body.String())
//		})
//	}
//}
