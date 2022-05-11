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
	"library-go/internal/domain"
	mock_service "library-go/internal/service/mocks"
	"library-go/pkg/logging"
	"net/http/httptest"
	"testing"
)

func TestBookHandler_GetAll(t *testing.T) {
	type mockBehavior func(s *mock_service.MockBookService, ctx context.Context, limit, offset int)

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
			mockBehavior: func(s *mock_service.MockBookService, ctx context.Context, limit, offset int) {
				s.EXPECT().GetAll(ctx, limit, offset).Return([]*domain.Book{domain.TestBook(), domain.TestBook()}, nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: "[{\"uuid\":\"1\",\"title\":\"Test Title\",\"direction\":{\"uuid\":\"1\",\"name\":\"test direction\"},\"author\":{\"uuid\":\"1\",\"full_name\":\"test Author\"},\"difficulty\":\"Test Difficulty\",\"edition_date\":\"2000-01-01T00:00:00Z\",\"rating\":5,\"description\":\"Test Description\",\"url\":\"Test URL\",\"language\":\"Test Language\",\"tags\":[{\"uuid\":\"1\",\"name\":\"test Tag\"}],\"download_count\":10},{\"uuid\":\"1\",\"title\":\"Test Title\",\"direction\":{\"uuid\":\"1\",\"name\":\"test direction\"},\"author\":{\"uuid\":\"1\",\"full_name\":\"test Author\"},\"difficulty\":\"Test Difficulty\",\"edition_date\":\"2000-01-01T00:00:00Z\",\"rating\":5,\"description\":\"Test Description\",\"url\":\"Test URL\",\"language\":\"Test Language\",\"tags\":[{\"uuid\":\"1\",\"name\":\"test Tag\"}],\"download_count\":10}]\n",
		},
		{
			name:   "OK empty input",
			input:  "",
			ctx:    context.Background(),
			limit:  0,
			offset: 0,
			mockBehavior: func(s *mock_service.MockBookService, ctx context.Context, limit, offset int) {
				s.EXPECT().GetAll(ctx, limit, offset).Return([]*domain.Book{domain.TestBook(), domain.TestBook()}, nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: "[{\"uuid\":\"1\",\"title\":\"Test Title\",\"direction\":{\"uuid\":\"1\",\"name\":\"test direction\"},\"author\":{\"uuid\":\"1\",\"full_name\":\"test Author\"},\"difficulty\":\"Test Difficulty\",\"edition_date\":\"2000-01-01T00:00:00Z\",\"rating\":5,\"description\":\"Test Description\",\"url\":\"Test URL\",\"language\":\"Test Language\",\"tags\":[{\"uuid\":\"1\",\"name\":\"test Tag\"}],\"download_count\":10},{\"uuid\":\"1\",\"title\":\"Test Title\",\"direction\":{\"uuid\":\"1\",\"name\":\"test direction\"},\"author\":{\"uuid\":\"1\",\"full_name\":\"test Author\"},\"difficulty\":\"Test Difficulty\",\"edition_date\":\"2000-01-01T00:00:00Z\",\"rating\":5,\"description\":\"Test Description\",\"url\":\"Test URL\",\"language\":\"Test Language\",\"tags\":[{\"uuid\":\"1\",\"name\":\"test Tag\"}],\"download_count\":10}]\n",
		},
		{
			name:   "OK invalid input",
			input:  "limit=one&offset=-10",
			ctx:    context.Background(),
			limit:  0,
			offset: 0,
			mockBehavior: func(s *mock_service.MockBookService, ctx context.Context, limit, offset int) {
				s.EXPECT().GetAll(ctx, limit, offset).Return([]*domain.Book{domain.TestBook(), domain.TestBook()}, nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: "[{\"uuid\":\"1\",\"title\":\"Test Title\",\"direction\":{\"uuid\":\"1\",\"name\":\"test direction\"},\"author\":{\"uuid\":\"1\",\"full_name\":\"test Author\"},\"difficulty\":\"Test Difficulty\",\"edition_date\":\"2000-01-01T00:00:00Z\",\"rating\":5,\"description\":\"Test Description\",\"url\":\"Test URL\",\"language\":\"Test Language\",\"tags\":[{\"uuid\":\"1\",\"name\":\"test Tag\"}],\"download_count\":10},{\"uuid\":\"1\",\"title\":\"Test Title\",\"direction\":{\"uuid\":\"1\",\"name\":\"test direction\"},\"author\":{\"uuid\":\"1\",\"full_name\":\"test Author\"},\"difficulty\":\"Test Difficulty\",\"edition_date\":\"2000-01-01T00:00:00Z\",\"rating\":5,\"description\":\"Test Description\",\"url\":\"Test URL\",\"language\":\"Test Language\",\"tags\":[{\"uuid\":\"1\",\"name\":\"test Tag\"}],\"download_count\":10}]\n",
		},
		{
			name:   "No rows in result",
			input:  "limit=0&offset=0",
			ctx:    context.Background(),
			limit:  0,
			offset: 0,
			mockBehavior: func(s *mock_service.MockBookService, ctx context.Context, limit, offset int) {
				s.EXPECT().GetAll(ctx, limit, offset).Return(nil, errors.New("now rows in result"))
			},
			expectedStatusCode:  500,
			expectedRequestBody: fmt.Sprintf(`{"ErrorMsg":"error occurred while getting all books. err: %v"}%v`, errors.New("now rows in result"), "\n"),
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			service := mock_service.NewMockBookService(c)
			testCase.mockBehavior(service, testCase.ctx, testCase.limit, testCase.offset)

			BookHandler := NewBookHandler(service, logging.GetLogger())

			router := httprouter.New()
			BookHandler.Register(router)

			w := httptest.NewRecorder()

			req := httptest.NewRequest("GET", fmt.Sprintf("/books?%s", testCase.input), nil)

			router.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedRequestBody, w.Body.String())
		})
	}
}

func TestBookHandler_GetByUUID(t *testing.T) {
	type mockBehavior func(s *mock_service.MockBookService, ctx context.Context, uuid string)

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
			mockBehavior: func(s *mock_service.MockBookService, ctx context.Context, uuid string) {
				s.EXPECT().GetByUUID(ctx, uuid).Return(domain.TestBook(), nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: "{\"uuid\":\"1\",\"title\":\"Test Title\",\"direction\":{\"uuid\":\"1\",\"name\":\"test direction\"},\"author\":{\"uuid\":\"1\",\"full_name\":\"test Author\"},\"difficulty\":\"Test Difficulty\",\"edition_date\":\"2000-01-01T00:00:00Z\",\"rating\":5,\"description\":\"Test Description\",\"url\":\"Test URL\",\"language\":\"Test Language\",\"tags\":[{\"uuid\":\"1\",\"name\":\"test Tag\"}],\"download_count\":10}\n",
		},
		{
			name:                "invalid uuid",
			input:               "-1",
			uuid:                "-1",
			ctx:                 context.Background(),
			mockBehavior:        func(s *mock_service.MockBookService, ctx context.Context, uuid string) {},
			expectedStatusCode:  400,
			expectedRequestBody: "{\"ErrorMsg\":\"uuid need to be uint\"}\n",
		},
		{
			name:  "no rows in result",
			input: "1",
			uuid:  "1",
			ctx:   context.Background(),
			mockBehavior: func(s *mock_service.MockBookService, ctx context.Context, uuid string) {
				s.EXPECT().GetByUUID(ctx, uuid).Return(nil, errors.New("now rows in result"))
			},
			expectedStatusCode:  404,
			expectedRequestBody: "{\"ErrorMsg\":\"error occurred while getting book from DB by UUID. err: now rows in result\"}\n",
		},
		{
			name:                "string input",
			input:               "one",
			uuid:                "one",
			ctx:                 context.Background(),
			mockBehavior:        func(s *mock_service.MockBookService, ctx context.Context, uuid string) {},
			expectedStatusCode:  400,
			expectedRequestBody: "{\"ErrorMsg\":\"uuid need to be uint\"}\n",
		},
		// is it correct ?????
		{
			name:                "empty input",
			input:               "",
			uuid:                "",
			ctx:                 context.Background(),
			mockBehavior:        func(s *mock_service.MockBookService, ctx context.Context, uuid string) {},
			expectedStatusCode:  301,
			expectedRequestBody: "<a href=\"/book\">Moved Permanently</a>.\n\n",
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			service := mock_service.NewMockBookService(c)
			testCase.mockBehavior(service, testCase.ctx, testCase.uuid)

			BookHandler := NewBookHandler(service, logging.GetLogger())

			router := httprouter.New()
			BookHandler.Register(router)

			w := httptest.NewRecorder()

			req := httptest.NewRequest("GET", fmt.Sprintf("/book/%s", testCase.input), nil)

			router.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedRequestBody, w.Body.String())
		})
	}
}

//func TestBookHandler_Create(t *testing.T) {
//	type mockBehavior func(s *mock_service.MockBookService, ctx context.Context, createBookDTO domain.CreateBookDTO)
//
//	testTable := []struct {
//		name                string
//		ctx                 context.Context
//		inputBody           func() *io.PipeReader
//		mockBehavior        mockBehavior
//		createBookDTO       domain.CreateBookDTO
//		expectedStatusCode  int
//		expectedRequestBody string
//	}{
//		{
//			name: "OK",
//			ctx:  context.Background(),
//			inputBody: func() *io.PipeReader {
//				pr, pw := io.Pipe()
//				file, _ := os.CreateTemp("./", "test_file")
//				defer file.Close()
//				fileBytes, _ := ioutil.ReadAll(file)
//				writer := multipart.NewWriter(pw)
//				part, _ := writer.CreateFormFile("file", "book.pdf")
//				part.Write(fileBytes)
//				writer.WriteField("file", "test file")
//				writer.WriteField("title", "test title")
//				writer.WriteField("direction_uuid", "1")
//				writer.WriteField("author_uuid", "1")
//				writer.WriteField("difficulty", "test difficulty")
//				writer.WriteField("edition_date", "test edition_date")
//				writer.WriteField("description", "test description")
//				writer.WriteField("language", "test language")
//				writer.WriteField("tags_uuids", `{"1"}`)
//				return pr
//			},
//			mockBehavior: func(s *mock_service.MockBookService, ctx context.Context, createBookDTO domain.CreateBookDTO) {
//				s.EXPECT().Create(ctx, createBookDTO).Return("1", nil)
//			},
//			expectedStatusCode:  201,
//			expectedRequestBody: "Book created successfully. UUID: 1",
//		},
//	}
//	for _, testCase := range testTable {
//		t.Run(testCase.name, func(t *testing.T) {
//			c := gomock.NewController(t)
//			defer c.Finish()
//
//			service := mock_service.NewMockBookService(c)
//			testCase.mockBehavior(service, testCase.ctx, testCase.createBookDTO)
//
//			BookHandler := NewBookHandler(service, logging.GetLogger())
//
//			router := httprouter.New()
//			BookHandler.Register(router)
//
//			w := httptest.NewRecorder()
//			req := httptest.NewRequest("POST", fmt.Sprintf("/book"), nil)
//
//			router.ServeHTTP(w, req)
//
//			assert.Equal(t, testCase.expectedStatusCode, w.Code)
//			assert.Equal(t, testCase.expectedRequestBody, w.Body.String())
//		})
//	}
//
//}

func TestBookHandler_Delete(t *testing.T) {
	type mockBehavior func(s *mock_service.MockBookService, ctx context.Context, uuid string)

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
			mockBehavior: func(s *mock_service.MockBookService, ctx context.Context, uuid string) {
				s.EXPECT().Delete(ctx, uuid).Return(nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: "{\"infoMsg\":\"Book with UUID 1 was deleted\"}\n",
		},
		{
			name:  "no rows in result",
			input: "1",
			uuid:  "1",
			ctx:   context.Background(),
			mockBehavior: func(s *mock_service.MockBookService, ctx context.Context, uuid string) {
				s.EXPECT().Delete(ctx, uuid).Return(errors.New("no rows affected"))
			},
			expectedStatusCode:  404,
			expectedRequestBody: "{\"ErrorMsg\":\"error occurred while deleting book from DB. err: no rows affected\"}\n",
		},
		{
			name:                "invalid uuid",
			input:               "-1",
			uuid:                "-1",
			ctx:                 context.Background(),
			mockBehavior:        func(s *mock_service.MockBookService, ctx context.Context, uuid string) {},
			expectedStatusCode:  400,
			expectedRequestBody: "{\"ErrorMsg\":\"uuid need to be uint\"}\n",
		},
		{
			name:                "string uuid",
			input:               "one",
			uuid:                "one",
			ctx:                 context.Background(),
			mockBehavior:        func(s *mock_service.MockBookService, ctx context.Context, uuid string) {},
			expectedStatusCode:  400,
			expectedRequestBody: "{\"ErrorMsg\":\"uuid need to be uint\"}\n",
		},
		// is it correct ?????
		{
			name:                "empty input",
			input:               "",
			uuid:                "",
			ctx:                 context.Background(),
			mockBehavior:        func(s *mock_service.MockBookService, ctx context.Context, uuid string) {},
			expectedStatusCode:  404,
			expectedRequestBody: "404 page not found\n",
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			service := mock_service.NewMockBookService(c)
			testCase.mockBehavior(service, testCase.ctx, testCase.uuid)

			BookHandler := NewBookHandler(service, logging.GetLogger())

			router := httprouter.New()
			BookHandler.Register(router)

			w := httptest.NewRecorder()

			req := httptest.NewRequest("DELETE", fmt.Sprintf("/book/%s", testCase.input), nil)

			router.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedRequestBody, w.Body.String())
		})
	}
}

func TestBookHandler_Update(t *testing.T) {
	type mockBehavior func(s *mock_service.MockBookService, ctx context.Context, dto *domain.UpdateBookDTO)

	testTable := []struct {
		name                string
		inputBodyJSON       map[string]interface{}
		ctx                 context.Context
		dto                 domain.UpdateBookDTO
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
			dto: *domain.TestBookUpdateDTO(),
			mockBehavior: func(s *mock_service.MockBookService, ctx context.Context, dto *domain.UpdateBookDTO) {
				s.EXPECT().Update(ctx, dto).Return(nil)
			},
			expectedStatusCode:  201,
			expectedRequestBody: "{\"infoMsg\":\"Book updated successfully\"}\n",
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
			dto: *domain.TestBookUpdateDTO(),
			mockBehavior: func(s *mock_service.MockBookService, ctx context.Context, dto *domain.UpdateBookDTO) {
				s.EXPECT().Update(ctx, dto).Return(errors.New("no rows affected"))
			},
			expectedStatusCode:  500,
			expectedRequestBody: "{\"ErrorMsg\":\"error occurred while updating book into DB. err: no rows affected\"}\n",
		},
		{
			name:                "empty input body JSON or nil UUID",
			inputBodyJSON:       map[string]interface{}{},
			ctx:                 context.Background(),
			dto:                 *domain.TestBookUpdateDTO(),
			mockBehavior:        func(s *mock_service.MockBookService, ctx context.Context, dto *domain.UpdateBookDTO) {},
			expectedStatusCode:  400,
			expectedRequestBody: "{\"ErrorMsg\":\"error occurred while decoding JSON request. err: nil UUID\"}\n",
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			service := mock_service.NewMockBookService(c)
			testCase.mockBehavior(service, testCase.ctx, &testCase.dto)

			BookHandler := NewBookHandler(service, logging.GetLogger())

			router := httprouter.New()
			BookHandler.Register(router)

			body, _ := json.Marshal(testCase.inputBodyJSON)
			req := httptest.NewRequest("PUT", fmt.Sprintf("/book"), bytes.NewReader(body))

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedRequestBody, w.Body.String())
		})
	}
}

//func TestBookHandler_Load(t *testing.T) {
//	type mockBehavior func(s *mock_service.MockBookService, ctx context.Context, uuid string)
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
//			mockBehavior: func(s *mock_service.MockBookService, ctx context.Context, uuid string) {
//				s.EXPECT().GetByUUID(ctx, uuid).Return(domain.TestBook(), nil)
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
//			service := mock_service.NewMockBookService(c)
//			testCase.mockBehavior(service, testCase.ctx, testCase.uuid)
//
//			BookHandler := NewBookHandler(service, logging.GetLogger())
//
//			router := httprouter.New()
//			BookHandler.Register(router)
//
//			w := httptest.NewRecorder()
//
//			req := httptest.NewRequest("GET", fmt.Sprintf("/book/%s", testCase.input), nil)
//
//			router.ServeHTTP(w, req)
//
//			assert.Equal(t, testCase.expectedStatusCode, w.Code)
//			assert.Equal(t, testCase.expectedRequestBody, w.Body.String())
//		})
//	}
//}
