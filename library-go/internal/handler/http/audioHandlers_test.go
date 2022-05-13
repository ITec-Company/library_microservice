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

func TestAudioHandler_GetAll(t *testing.T) {
	type mockBehavior func(s *mock_service.MockAudioService, ctx context.Context, limit, offset int)

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
			mockBehavior: func(s *mock_service.MockAudioService, ctx context.Context, limit, offset int) {
				s.EXPECT().GetAll(ctx, limit, offset).Return([]*domain.Audio{domain.TestAudio(), domain.TestAudio()}, nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: "[{\"uuid\":\"1\",\"title\":\"Test Title\",\"direction\":{\"uuid\":\"1\",\"name\":\"test direction\"},\"difficulty\":\"Test Difficulty\",\"rating\":5,\"url\":\"Test URL\",\"language\":\"Test Language\",\"tags\":[{\"uuid\":\"1\",\"name\":\"test Tag\"}],\"download_count\":10},{\"uuid\":\"1\",\"title\":\"Test Title\",\"direction\":{\"uuid\":\"1\",\"name\":\"test direction\"},\"difficulty\":\"Test Difficulty\",\"rating\":5,\"url\":\"Test URL\",\"language\":\"Test Language\",\"tags\":[{\"uuid\":\"1\",\"name\":\"test Tag\"}],\"download_count\":10}]\n",
		},
		{
			name:   "OK empty input",
			input:  "",
			ctx:    context.Background(),
			limit:  0,
			offset: 0,
			mockBehavior: func(s *mock_service.MockAudioService, ctx context.Context, limit, offset int) {
				s.EXPECT().GetAll(ctx, limit, offset).Return([]*domain.Audio{domain.TestAudio(), domain.TestAudio()}, nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: "[{\"uuid\":\"1\",\"title\":\"Test Title\",\"direction\":{\"uuid\":\"1\",\"name\":\"test direction\"},\"difficulty\":\"Test Difficulty\",\"rating\":5,\"url\":\"Test URL\",\"language\":\"Test Language\",\"tags\":[{\"uuid\":\"1\",\"name\":\"test Tag\"}],\"download_count\":10},{\"uuid\":\"1\",\"title\":\"Test Title\",\"direction\":{\"uuid\":\"1\",\"name\":\"test direction\"},\"difficulty\":\"Test Difficulty\",\"rating\":5,\"url\":\"Test URL\",\"language\":\"Test Language\",\"tags\":[{\"uuid\":\"1\",\"name\":\"test Tag\"}],\"download_count\":10}]\n",
		},
		{
			name:   "OK invalid input",
			input:  "limit=one&offset=-10",
			ctx:    context.Background(),
			limit:  0,
			offset: 0,
			mockBehavior: func(s *mock_service.MockAudioService, ctx context.Context, limit, offset int) {
				s.EXPECT().GetAll(ctx, limit, offset).Return([]*domain.Audio{domain.TestAudio(), domain.TestAudio()}, nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: "[{\"uuid\":\"1\",\"title\":\"Test Title\",\"direction\":{\"uuid\":\"1\",\"name\":\"test direction\"},\"difficulty\":\"Test Difficulty\",\"rating\":5,\"url\":\"Test URL\",\"language\":\"Test Language\",\"tags\":[{\"uuid\":\"1\",\"name\":\"test Tag\"}],\"download_count\":10},{\"uuid\":\"1\",\"title\":\"Test Title\",\"direction\":{\"uuid\":\"1\",\"name\":\"test direction\"},\"difficulty\":\"Test Difficulty\",\"rating\":5,\"url\":\"Test URL\",\"language\":\"Test Language\",\"tags\":[{\"uuid\":\"1\",\"name\":\"test Tag\"}],\"download_count\":10}]\n",
		},
		{
			name:   "No rows in result",
			input:  "limit=0&offset=0",
			ctx:    context.Background(),
			limit:  0,
			offset: 0,
			mockBehavior: func(s *mock_service.MockAudioService, ctx context.Context, limit, offset int) {
				s.EXPECT().GetAll(ctx, limit, offset).Return(nil, errors.New("now rows in result"))
			},
			expectedStatusCode:  500,
			expectedRequestBody: "{\"ErrorMsg\":\"error occurred while getting all audios. err: now rows in result\"}\n",
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			service := mock_service.NewMockAudioService(c)
			testCase.mockBehavior(service, testCase.ctx, testCase.limit, testCase.offset)

			AudioHandler := NewAudioHandler(service, logging.GetLogger())

			router := httprouter.New()
			AudioHandler.Register(router)

			w := httptest.NewRecorder()

			req := httptest.NewRequest("GET", fmt.Sprintf("/audios?%s", testCase.input), nil)

			router.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedRequestBody, w.Body.String())
		})
	}
}

func TestAudioHandler_GetByUUID(t *testing.T) {
	type mockBehavior func(s *mock_service.MockAudioService, ctx context.Context, uuid string)

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
			mockBehavior: func(s *mock_service.MockAudioService, ctx context.Context, uuid string) {
				s.EXPECT().GetByUUID(ctx, uuid).Return(domain.TestAudio(), nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: "{\"uuid\":\"1\",\"title\":\"Test Title\",\"direction\":{\"uuid\":\"1\",\"name\":\"test direction\"},\"difficulty\":\"Test Difficulty\",\"rating\":5,\"url\":\"Test URL\",\"language\":\"Test Language\",\"tags\":[{\"uuid\":\"1\",\"name\":\"test Tag\"}],\"download_count\":10}\n",
		},
		{
			name:                "invalid uuid",
			input:               "-1",
			uuid:                "-1",
			ctx:                 context.Background(),
			mockBehavior:        func(s *mock_service.MockAudioService, ctx context.Context, uuid string) {},
			expectedStatusCode:  400,
			expectedRequestBody: "{\"ErrorMsg\":\"uuid need to be uint\"}\n",
		},
		{
			name:  "no rows in result",
			input: "1",
			uuid:  "1",
			ctx:   context.Background(),
			mockBehavior: func(s *mock_service.MockAudioService, ctx context.Context, uuid string) {
				s.EXPECT().GetByUUID(ctx, uuid).Return(nil, errors.New("now rows in result"))
			},
			expectedStatusCode:  404,
			expectedRequestBody: "{\"ErrorMsg\":\"error occurred while getting audio from DB by UUID. err: now rows in result\"}\n",
		},
		{
			name:                "string input",
			input:               "one",
			uuid:                "one",
			ctx:                 context.Background(),
			mockBehavior:        func(s *mock_service.MockAudioService, ctx context.Context, uuid string) {},
			expectedStatusCode:  400,
			expectedRequestBody: "{\"ErrorMsg\":\"uuid need to be uint\"}\n",
		},
		// is it correct ?????
		{
			name:                "empty input",
			input:               "",
			uuid:                "",
			ctx:                 context.Background(),
			mockBehavior:        func(s *mock_service.MockAudioService, ctx context.Context, uuid string) {},
			expectedStatusCode:  301,
			expectedRequestBody: "<a href=\"/audio\">Moved Permanently</a>.\n\n",
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			service := mock_service.NewMockAudioService(c)
			testCase.mockBehavior(service, testCase.ctx, testCase.uuid)

			AudioHandler := NewAudioHandler(service, logging.GetLogger())

			router := httprouter.New()
			AudioHandler.Register(router)

			w := httptest.NewRecorder()

			req := httptest.NewRequest("GET", fmt.Sprintf("/audio/%s", testCase.input), nil)

			router.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedRequestBody, w.Body.String())
		})
	}
}

func TestAudioHandler_Create(t *testing.T) {
	type mockBehavior func(s *mock_service.MockAudioService, ctx context.Context, createAudioDTO *domain.CreateAudioDTO)

	testTable := []struct {
		name                string
		ctx                 context.Context
		mockBehavior        mockBehavior
		createAudioDTO      *domain.CreateAudioDTO
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:           "OK",
			ctx:            context.Background(),
			createAudioDTO: domain.TestAudioCreateDTO(),
			mockBehavior: func(s *mock_service.MockAudioService, ctx context.Context, createAudioDTO *domain.CreateAudioDTO) {
				createAudioDTO.URL = "load/audio/1-Test Difficulty-file"
				s.EXPECT().Create(ctx, createAudioDTO).Return("1", nil)
			},
			expectedStatusCode:  201,
			expectedRequestBody: "{\"infoMsg\":\"Audio created successfully. UUID: 1\"}\n",
		},
		{
			name:           "service error",
			ctx:            context.Background(),
			createAudioDTO: domain.TestAudioCreateDTO(),
			mockBehavior: func(s *mock_service.MockAudioService, ctx context.Context, createAudioDTO *domain.CreateAudioDTO) {
				createAudioDTO.URL = "load/audio/1-Test Difficulty-file"
				s.EXPECT().Create(ctx, createAudioDTO).Return("", errors.New("some service error"))
			},
			expectedStatusCode:  500,
			expectedRequestBody: "{\"ErrorMsg\":\"error occurred while creating audio into DB. err: some service error\"}\n",
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			service := mock_service.NewMockAudioService(c)
			testCase.mockBehavior(service, testCase.ctx, testCase.createAudioDTO)

			AudioHandler := NewAudioHandler(service, logging.GetLogger())

			router := httprouter.New()
			AudioHandler.Register(router)

			body := bytes.Buffer{}
			writer := multipart.NewWriter(&body)
			file, _ := os.Open("test_file.txt")

			part, _ := writer.CreateFormFile("file", "file")
			io.Copy(part, file)
			defer file.Close()

			writer.WriteField("title", "Test Title")
			writer.WriteField("direction_uuid", "1")
			writer.WriteField("author_uuid", "1")
			writer.WriteField("difficulty", "Test Difficulty")
			writer.WriteField("language", "Test Language")
			writer.WriteField("tags_uuids", `1`)

			w := httptest.NewRecorder()

			writer.Close()
			req := httptest.NewRequest("POST", "/audio", bytes.NewReader(body.Bytes()))
			req.Header.Set("Content-Type", writer.FormDataContentType())

			router.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedRequestBody, w.Body.String())
		})
	}

}

func TestAudioHandler_Delete(t *testing.T) {
	type mockBehavior func(s *mock_service.MockAudioService, ctx context.Context, uuid string)

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
			mockBehavior: func(s *mock_service.MockAudioService, ctx context.Context, uuid string) {
				s.EXPECT().Delete(ctx, uuid).Return(nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: "{\"infoMsg\":\"Audio with UUID 1 was deleted\"}\n",
		},
		{
			name:  "no rows in result",
			input: "1",
			uuid:  "1",
			ctx:   context.Background(),
			mockBehavior: func(s *mock_service.MockAudioService, ctx context.Context, uuid string) {
				s.EXPECT().Delete(ctx, uuid).Return(errors.New("no rows affected"))
			},
			expectedStatusCode:  404,
			expectedRequestBody: "{\"ErrorMsg\":\"error occurred while deleting audio from DB. err: no rows affected\"}\n",
		},
		{
			name:                "invalid uuid",
			input:               "-1",
			uuid:                "-1",
			ctx:                 context.Background(),
			mockBehavior:        func(s *mock_service.MockAudioService, ctx context.Context, uuid string) {},
			expectedStatusCode:  400,
			expectedRequestBody: "{\"ErrorMsg\":\"uuid need to be uint\"}\n",
		},
		{
			name:                "string uuid",
			input:               "one",
			uuid:                "one",
			ctx:                 context.Background(),
			mockBehavior:        func(s *mock_service.MockAudioService, ctx context.Context, uuid string) {},
			expectedStatusCode:  400,
			expectedRequestBody: "{\"ErrorMsg\":\"uuid need to be uint\"}\n",
		},
		// is it correct ?????
		{
			name:                "empty input",
			input:               "",
			uuid:                "",
			ctx:                 context.Background(),
			mockBehavior:        func(s *mock_service.MockAudioService, ctx context.Context, uuid string) {},
			expectedStatusCode:  404,
			expectedRequestBody: "404 page not found\n",
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			service := mock_service.NewMockAudioService(c)
			testCase.mockBehavior(service, testCase.ctx, testCase.uuid)

			AudioHandler := NewAudioHandler(service, logging.GetLogger())

			router := httprouter.New()
			AudioHandler.Register(router)

			w := httptest.NewRecorder()

			req := httptest.NewRequest("DELETE", fmt.Sprintf("/audio/%s", testCase.input), nil)

			router.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedRequestBody, w.Body.String())
		})
	}
}

func TestAudioHandler_Update(t *testing.T) {
	type mockBehavior func(s *mock_service.MockAudioService, ctx context.Context, dto *domain.UpdateAudioDTO)

	testTable := []struct {
		name                string
		inputBodyJSON       map[string]interface{}
		ctx                 context.Context
		dto                 domain.UpdateAudioDTO
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
				"difficulty":     "Test Difficulty",
				"rating":         5.5,
				"url":            "Test URL",
				"language":       "Test Language",
				"tags_uuids":     []string{"1"},
				"download_count": 10,
			},
			ctx: context.Background(),
			dto: *domain.TestAudioUpdateDTO(),
			mockBehavior: func(s *mock_service.MockAudioService, ctx context.Context, dto *domain.UpdateAudioDTO) {
				s.EXPECT().Update(ctx, dto).Return(nil)
			},
			expectedStatusCode:  201,
			expectedRequestBody: "{\"infoMsg\":\"Audio updated successfully\"}\n",
		},
		{
			name: "no rows in result",
			inputBodyJSON: map[string]interface{}{
				"uuid":           "1",
				"title":          "Test Title",
				"direction_uuid": "1",
				"difficulty":     "Test Difficulty",
				"rating":         5.5,
				"url":            "Test URL",
				"language":       "Test Language",
				"tags_uuids":     []string{"1"},
				"download_count": 10,
			},
			ctx: context.Background(),
			dto: *domain.TestAudioUpdateDTO(),
			mockBehavior: func(s *mock_service.MockAudioService, ctx context.Context, dto *domain.UpdateAudioDTO) {
				s.EXPECT().Update(ctx, dto).Return(errors.New("no rows affected"))
			},
			expectedStatusCode:  500,
			expectedRequestBody: "{\"ErrorMsg\":\"error occurred while updating audio into DB. err: no rows affected\"}\n",
		},
		{
			name:                "empty input body JSON or nil UUID",
			inputBodyJSON:       map[string]interface{}{},
			ctx:                 context.Background(),
			dto:                 *domain.TestAudioUpdateDTO(),
			mockBehavior:        func(s *mock_service.MockAudioService, ctx context.Context, dto *domain.UpdateAudioDTO) {},
			expectedStatusCode:  400,
			expectedRequestBody: "{\"ErrorMsg\":\"error occurred while decoding JSON request. err: nil UUID\"}\n",
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			service := mock_service.NewMockAudioService(c)
			testCase.mockBehavior(service, testCase.ctx, &testCase.dto)

			AudioHandler := NewAudioHandler(service, logging.GetLogger())

			router := httprouter.New()
			AudioHandler.Register(router)

			body, _ := json.Marshal(testCase.inputBodyJSON)
			req := httptest.NewRequest("PUT", fmt.Sprintf("/audio"), bytes.NewReader(body))

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedRequestBody, w.Body.String())
		})
	}
}

//func TestAudioHandler_Load(t *testing.T) {
//	type mockBehavior func(s *mock_service.MockAudioService, ctx context.Context, uuid string)
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
//			mockBehavior: func(s *mock_service.MockAudioService, ctx context.Context, uuid string) {
//				s.EXPECT().GetByUUID(ctx, uuid).Return(domain.TestAudio(), nil)
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
//			service := mock_service.NewMockAudioService(c)
//			testCase.mockBehavior(service, testCase.ctx, testCase.uuid)
//
//			AudioHandler := NewAudioHandler(service, logging.GetLogger())
//
//			router := httprouter.New()
//			AudioHandler.Register(router)
//
//			w := httptest.NewRecorder()
//
//			req := httptest.NewRequest("GET", fmt.Sprintf("/audio/%s", testCase.input), nil)
//
//			router.ServeHTTP(w, req)
//
//			assert.Equal(t, testCase.expectedStatusCode, w.Code)
//			assert.Equal(t, testCase.expectedRequestBody, w.Body.String())
//		})
//	}
//}
