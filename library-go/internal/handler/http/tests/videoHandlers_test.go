package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
	"io"
	"library-go/internal/domain"
	"library-go/internal/handler/http"
	mock_service "library-go/internal/service/mocks"
	"library-go/pkg/logging"
	"mime/multipart"
	"net/http/httptest"
	"os"
	"testing"
)

func TestVideoHandler_GetAll(t *testing.T) {
	type mockBehavior func(s *mock_service.MockVideoService, sortingOptions *domain.SortFilterPagination)

	testTable := []struct {
		name, input         string
		sortingOptions      domain.SortFilterPagination
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:  "OK",
			input: "",
			sortingOptions: domain.SortFilterPagination{
				SortBy:         "",
				Order:          "",
				FiltersAndArgs: nil,
				Limit:          0,
				Page:           0,
			},
			mockBehavior: func(s *mock_service.MockVideoService, sortingOptions *domain.SortFilterPagination) {
				s.EXPECT().GetAll(sortingOptions).Return([]*domain.Video{domain.TestVideo(), domain.TestVideo()}, 1, nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: "[{\"uuid\":\"1\",\"title\":\"Test Title\",\"direction\":{\"uuid\":\"1\",\"name\":\"Test Direction\"},\"rating\":5,\"difficulty\":\"Test Difficulty\",\"description\":\"Test Description\",\"local_url\":\"Test LocalURL\",\"language\":\"Test Language\",\"tags\":[{\"uuid\":\"1\",\"name\":\"Test Tag\"}],\"download_count\":10,\"created_at\":\"0001-01-01T00:00:00Z\"},{\"uuid\":\"1\",\"title\":\"Test Title\",\"direction\":{\"uuid\":\"1\",\"name\":\"Test Direction\"},\"rating\":5,\"difficulty\":\"Test Difficulty\",\"description\":\"Test Description\",\"local_url\":\"Test LocalURL\",\"language\":\"Test Language\",\"tags\":[{\"uuid\":\"1\",\"name\":\"Test Tag\"}],\"download_count\":10,\"created_at\":\"0001-01-01T00:00:00Z\"}]\n",
		},
		{
			name:  "no rows in result",
			input: "",
			sortingOptions: domain.SortFilterPagination{
				SortBy:         "",
				Order:          "",
				FiltersAndArgs: nil,
				Limit:          0,
				Page:           0,
			},
			mockBehavior: func(s *mock_service.MockVideoService, sortingOptions *domain.SortFilterPagination) {
				s.EXPECT().GetAll(sortingOptions).Return(nil, 0, nil)
			},
			expectedStatusCode:  404,
			expectedRequestBody: "{\"ErrorMsg\":\"no rows in result\"}\n",
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			service := mock_service.NewMockVideoService(c)
			testCase.mockBehavior(service, &testCase.sortingOptions)

			logger := logging.GetLogger("../../../../logs", "test.log")
			middleware := http.NewMiddlewares(logger)
			VideoHandler := http.NewVideoHandler(service, logger)

			router := httprouter.New()
			VideoHandler.Register(router, &middleware)

			w := httptest.NewRecorder()

			req := httptest.NewRequest("GET", fmt.Sprintf("/videos%s", testCase.input), nil)

			router.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedRequestBody, w.Body.String())
		})
	}
}

func TestVideoHandler_GetByUUID(t *testing.T) {
	type mockBehavior func(s *mock_service.MockVideoService, uuid string)

	testTable := []struct {
		name, uuid, input   string
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:  "OK",
			input: "1",
			uuid:  "1",
			mockBehavior: func(s *mock_service.MockVideoService, uuid string) {
				s.EXPECT().GetByUUID(uuid).Return(domain.TestVideo(), nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: "{\"uuid\":\"1\",\"title\":\"Test Title\",\"direction\":{\"uuid\":\"1\",\"name\":\"Test Direction\"},\"rating\":5,\"difficulty\":\"Test Difficulty\",\"description\":\"Test Description\",\"local_url\":\"Test LocalURL\",\"language\":\"Test Language\",\"tags\":[{\"uuid\":\"1\",\"name\":\"Test Tag\"}],\"download_count\":10,\"created_at\":\"0001-01-01T00:00:00Z\"}\n",
		},
		{
			name:                "invalid uuid",
			input:               "-1",
			uuid:                "-1",
			mockBehavior:        func(s *mock_service.MockVideoService, uuid string) {},
			expectedStatusCode:  400,
			expectedRequestBody: "{\"ErrorMsg\":\"uuid need to be uint\"}\n",
		},
		{
			name:  "no rows in result",
			input: "1",
			uuid:  "1",
			mockBehavior: func(s *mock_service.MockVideoService, uuid string) {
				s.EXPECT().GetByUUID(uuid).Return(nil, errors.New("now rows in result"))
			},
			expectedStatusCode:  404,
			expectedRequestBody: "{\"ErrorMsg\":\"error occurred while getting video from DB by UUID. err: now rows in result\"}\n",
		},
		{
			name:                "string input",
			input:               "one",
			uuid:                "one",
			mockBehavior:        func(s *mock_service.MockVideoService, uuid string) {},
			expectedStatusCode:  400,
			expectedRequestBody: "{\"ErrorMsg\":\"uuid need to be uint\"}\n",
		},
		// is it correct ?????
		{
			name:                "empty input",
			input:               "",
			uuid:                "",
			mockBehavior:        func(s *mock_service.MockVideoService, uuid string) {},
			expectedStatusCode:  301,
			expectedRequestBody: "<a href=\"/video\">Moved Permanently</a>.\n\n",
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			service := mock_service.NewMockVideoService(c)
			testCase.mockBehavior(service, testCase.uuid)

			logger := logging.GetLogger("../../../../logs", "test.log")
			middleware := http.NewMiddlewares(logger)
			VideoHandler := http.NewVideoHandler(service, logger)

			router := httprouter.New()
			VideoHandler.Register(router, &middleware)

			w := httptest.NewRecorder()

			req := httptest.NewRequest("GET", fmt.Sprintf("/video/%s", testCase.input), nil)

			router.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedRequestBody, w.Body.String())
		})
	}
}

func TestVideoHandler_Create(t *testing.T) {
	type mockBehavior func(s *mock_service.MockVideoService, createVideoDTO *domain.CreateVideoDTO)
	type mockSaveFile func(s *mock_service.MockVideoService, path, fileName string, file io.Reader)

	file, _ := os.Open("test_file_video.mp4")
	fileinfo, _ := file.Stat()
	fileBytes := make([]byte, fileinfo.Size())
	file.Read(fileBytes)

	testTable := []struct {
		name                string
		contentType         string
		mockBehavior        mockBehavior
		mockSaveFile        mockSaveFile
		createVideoDTO      domain.CreateVideoDTO
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:           "OK",
			contentType:    "MultipartFormData",
			createVideoDTO: *domain.TestVideoCreateDTO(),
			mockBehavior: func(s *mock_service.MockVideoService, createVideoDTO *domain.CreateVideoDTO) {
				s.EXPECT().Create(createVideoDTO).Return("1", nil)
			},
			mockSaveFile: func(s *mock_service.MockVideoService, path, fileName string, file io.Reader) {
				s.EXPECT().SaveFile("../store/videos/1/", "title(Test_Title).mp4", bytes.NewBuffer(fileBytes)).Return(nil)
			},
			expectedStatusCode:  201,
			expectedRequestBody: "{\"infoMsg\":\"Video created successfully. UUID: 1\"}\n",
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			pr, pw := io.Pipe()

			writer := multipart.NewWriter(pw)
			go func() {
				defer writer.Close()

				file, _ := os.Open("test_file_video.mp4")

				defer file.Close()

				part, _ := writer.CreateFormFile("file", "test_file_video.mp4")

				_, _ = io.Copy(part, file)

				writer.WriteField("title", "Test Title")
				writer.WriteField("direction_uuid", "1")
				writer.WriteField("difficulty", "Test Difficulty")
				writer.WriteField("description", "Test Description")
				writer.WriteField("language", "Test Language")
				writer.WriteField("tags_uuids", `1`)
				writer.WriteField("web_url", `text`)
			}()

			w := httptest.NewRecorder()

			req := httptest.NewRequest("POST", fmt.Sprintf("/video"), pr)
			if testCase.contentType == "MultipartFormData" {
				req.Header.Set("Content-Type", writer.FormDataContentType())
			} else {
				req.Header.Set("Content-Type", "application/json")
			}

			defer req.Body.Close()
			service := mock_service.NewMockVideoService(c)

			testCase.mockSaveFile(service, "../store/videos/1/", "title(Test_Title).mp4", file)
			testCase.mockBehavior(service, &testCase.createVideoDTO)

			logger := logging.GetLogger("../../../../logs", "test.log")
			middleware := http.NewMiddlewares(logger)
			VideoHandler := http.NewVideoHandler(service, logger)

			router := httprouter.New()
			VideoHandler.Register(router, &middleware)

			router.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedRequestBody, w.Body.String())
		})
	}

}

func TestVideoHandler_Delete(t *testing.T) {
	type mockBehavior func(s *mock_service.MockVideoService, uuid string, path string)

	testTable := []struct {
		name, uuid, path, input string
		mockBehavior            mockBehavior
		expectedStatusCode      int
		expectedRequestBody     string
	}{
		{
			name:  "OK",
			input: "1",
			uuid:  "1",
			mockBehavior: func(s *mock_service.MockVideoService, uuid string, path string) {
				s.EXPECT().Delete(uuid, fmt.Sprintf("%s%s/", "../store/videos/", uuid)).Return(nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: "{\"infoMsg\":\"Video with UUID 1 was deleted\"}\n",
		},
		{
			name:  "no rows in result",
			input: "1",
			uuid:  "1",
			mockBehavior: func(s *mock_service.MockVideoService, uuid string, path string) {
				s.EXPECT().Delete(uuid, fmt.Sprintf("%s%s/", "../store/videos/", uuid)).Return(errors.New("no rows affected"))
			},
			expectedStatusCode:  404,
			expectedRequestBody: "{\"ErrorMsg\":\"error occurred while deleting video from DB. err: no rows affected\"}\n",
		},
		{
			name:                "invalid uuid",
			input:               "-1",
			uuid:                "-1",
			mockBehavior:        func(s *mock_service.MockVideoService, uuid string, path string) {},
			expectedStatusCode:  400,
			expectedRequestBody: "{\"ErrorMsg\":\"uuid need to be uint\"}\n",
		},
		{
			name:                "string uuid",
			input:               "one",
			uuid:                "one",
			mockBehavior:        func(s *mock_service.MockVideoService, uuid string, path string) {},
			expectedStatusCode:  400,
			expectedRequestBody: "{\"ErrorMsg\":\"uuid need to be uint\"}\n",
		},
		// is it correct ?????
		{
			name:                "empty input",
			input:               "",
			uuid:                "",
			mockBehavior:        func(s *mock_service.MockVideoService, uuid string, path string) {},
			expectedStatusCode:  404,
			expectedRequestBody: "404 page not found\n",
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			service := mock_service.NewMockVideoService(c)
			testCase.mockBehavior(service, testCase.uuid, fmt.Sprintf("%s%s/", "../store/videos/", testCase.uuid))

			logger := logging.GetLogger("../../../../logs", "test.log")
			middleware := http.NewMiddlewares(logger)
			VideoHandler := http.NewVideoHandler(service, logger)

			router := httprouter.New()
			VideoHandler.Register(router, &middleware)

			w := httptest.NewRecorder()

			req := httptest.NewRequest("DELETE", fmt.Sprintf("/video/%s", testCase.input), nil)

			router.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedRequestBody, w.Body.String())
		})
	}
}

func TestVideoHandler_Update(t *testing.T) {
	type mockBehavior func(s *mock_service.MockVideoService, dto *domain.UpdateVideoDTO)

	testTable := []struct {
		name                string
		inputBodyJSON       map[string]interface{}
		dto                 domain.UpdateVideoDTO
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
				"description":    "Test Description",
				"web_url":        "Test URL",
				"language":       "Test Language",
				"tags_uuids":     []string{"1"},
				"download_count": 10,
			},
			dto: *domain.TestVideoUpdateDTO(),
			mockBehavior: func(s *mock_service.MockVideoService, dto *domain.UpdateVideoDTO) {
				s.EXPECT().Update(dto).Return(nil)
			},
			expectedStatusCode:  201,
			expectedRequestBody: "{\"infoMsg\":\"Video updated successfully\"}\n",
		},
		{
			name: "no rows in result",
			inputBodyJSON: map[string]interface{}{
				"uuid":           "1",
				"title":          "Test Title",
				"direction_uuid": "1",
				"difficulty":     "Test Difficulty",
				"rating":         5.5,
				"description":    "Test Description",
				"web_url":        "Test URL",
				"language":       "Test Language",
				"tags_uuids":     []string{"1"},
				"download_count": 10,
			},
			dto: *domain.TestVideoUpdateDTO(),
			mockBehavior: func(s *mock_service.MockVideoService, dto *domain.UpdateVideoDTO) {
				s.EXPECT().Update(dto).Return(errors.New("no rows affected"))
			},
			expectedStatusCode:  500,
			expectedRequestBody: "{\"ErrorMsg\":\"error occurred while updating video into DB. err: no rows affected\"}\n",
		},
		{
			name:                "empty input body JSON or nil UUID",
			inputBodyJSON:       map[string]interface{}{},
			dto:                 *domain.TestVideoUpdateDTO(),
			mockBehavior:        func(s *mock_service.MockVideoService, dto *domain.UpdateVideoDTO) {},
			expectedStatusCode:  400,
			expectedRequestBody: "{\"ErrorMsg\":\"error occurred while decoding JSON request. err: nil UUID\"}\n",
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			service := mock_service.NewMockVideoService(c)
			testCase.mockBehavior(service, &testCase.dto)

			logger := logging.GetLogger("../../../../logs", "test.log")
			middleware := http.NewMiddlewares(logger)
			VideoHandler := http.NewVideoHandler(service, logger)

			router := httprouter.New()
			VideoHandler.Register(router, &middleware)

			body, _ := json.Marshal(testCase.inputBodyJSON)
			req := httptest.NewRequest("PUT", fmt.Sprintf("/video"), bytes.NewReader(body))

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedRequestBody, w.Body.String())
		})
	}
}

func TestVideoHandler_UpdateFile(t *testing.T) {
	type mockBehavior func(s *mock_service.MockVideoService, dto *domain.UpdateVideoFileDTO)

	d, _ := os.Open("test_file_video.mp4")
	dtoFile := new(bytes.Buffer)
	dtoFile.ReadFrom(d)

	testTable := []struct {
		name                string
		localURL            string
		inputBody           *io.Reader
		dto                 *domain.UpdateVideoFileDTO
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:     "OK",
			localURL: "/videos/1/title(software_testing).mp4",
			dto: &domain.UpdateVideoFileDTO{
				UUID:        "1",
				NewFileName: "title(software_testing).mp4",
				OldFileName: "title(software_testing).mp4",
				File:        dtoFile,
				LocalURL:    fmt.Sprintf("%s|split|/%s", "/videos/", "title(software_testing).mp4"),
				LocalPath:   fmt.Sprintf("%s%s/", "../store/videos/", "1"),
			},
			mockBehavior: func(s *mock_service.MockVideoService, dto *domain.UpdateVideoFileDTO) {
				s.EXPECT().UpdateFile(dto).Return(nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: "{\"infoMsg\":\"File updated successfully\"}\n",
		},
		{
			name:     "invalid UUID",
			localURL: "/videos/0/title(software_testing).mp4",
			dto: &domain.UpdateVideoFileDTO{
				UUID:        "0",
				NewFileName: "title(software_testing).mp4",
				OldFileName: "title(software_testing).mp4",
				File:        dtoFile,
				LocalURL:    fmt.Sprintf("%s|split|/%s", "/videos/", "title(software_testing).mp4"),
				LocalPath:   fmt.Sprintf("%s%s/", "../store/videos/", "0"),
			},
			mockBehavior: func(s *mock_service.MockVideoService, dto *domain.UpdateVideoFileDTO) {
				s.EXPECT().UpdateFile(dto).Return(errors.New("no video with such UUID was found"))
			},
			expectedStatusCode:  404,
			expectedRequestBody: "{\"ErrorMsg\":\"error occurred while saving video into local store. err: no video with such UUID was found.\"}\n",
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			service := mock_service.NewMockVideoService(c)
			testCase.mockBehavior(service, testCase.dto)

			logger := logging.GetLogger("../../../../logs", "test.log")
			middleware := http.NewMiddlewares(logger)
			VideoHandler := http.NewVideoHandler(service, logger)

			router := httprouter.New()
			VideoHandler.Register(router, &middleware)

			w := httptest.NewRecorder()
			file, _ := os.Open("test_file_video.mp4")
			req := httptest.NewRequest("PUT", fmt.Sprintf("/file/video?localurl=%s", testCase.localURL), file)

			router.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedRequestBody, w.Body.String())
		})
	}
}

func TestVideoHandler_Rate(t *testing.T) {
	type mockBehavior func(s *mock_service.MockVideoService, uuid string, rating float32)

	testTable := []struct {
		name, uuid, input   string
		rating              float32
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:   "OK",
			uuid:   "1",
			rating: 1,
			input:  "?uuid=1&rating=1",
			mockBehavior: func(s *mock_service.MockVideoService, uuid string, rating float32) {
				s.EXPECT().Rate(uuid, rating).Return(nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: "{\"infoMsg\":\"Video rated successfully. UUID: 1\"}\n",
		},
		{
			name:   "service error",
			uuid:   "1",
			rating: 1,
			input:  "?uuid=1&rating=1",
			mockBehavior: func(s *mock_service.MockVideoService, uuid string, rating float32) {
				s.EXPECT().Rate(uuid, rating).Return(errors.New("service error"))
			},
			expectedStatusCode:  400,
			expectedRequestBody: "{\"ErrorMsg\":\"error occurred while rating video into local store. err: service error.\"}\n",
		},
		{
			name:   "invalid UUID",
			uuid:   "0",
			rating: 1,
			input:  "?uuid=0&rating=1",
			mockBehavior: func(s *mock_service.MockVideoService, uuid string, rating float32) {
				s.EXPECT().Rate(uuid, rating).Return(errors.New("invalid UUID"))
			},
			expectedStatusCode:  400,
			expectedRequestBody: "{\"ErrorMsg\":\"error occurred while rating video into local store. err: invalid UUID.\"}\n",
		},
		{
			name:   "empty UUID",
			rating: 1,
			input:  "?rating=1",
			mockBehavior: func(s *mock_service.MockVideoService, uuid string, rating float32) {
			},
			expectedStatusCode:  400,
			expectedRequestBody: "{\"ErrorMsg\":\"uuid can't be empty\"}\n",
		},
		{
			name:  "empty rating",
			uuid:  "1",
			input: "?uuid=1",
			mockBehavior: func(s *mock_service.MockVideoService, uuid string, rating float32) {
			},
			expectedStatusCode:  400,
			expectedRequestBody: "{\"ErrorMsg\":\"rating query can't be empty\"}\n",
		},
		{
			name:  "invalid rating",
			uuid:  "1",
			input: "?uuid=1&rating=five",
			mockBehavior: func(s *mock_service.MockVideoService, uuid string, rating float32) {
			},
			expectedStatusCode:  400,
			expectedRequestBody: "{\"ErrorMsg\":\"error occurred while parsing rating. Should be float32. err: strconv.ParseFloat: parsing \\\"five\\\": invalid syntax.\"}\n",
		},
		{
			name:   "bad rating value",
			uuid:   "1",
			rating: -1,
			input:  "?uuid=1&rating=-1",
			mockBehavior: func(s *mock_service.MockVideoService, uuid string, rating float32) {
			},
			expectedStatusCode:  400,
			expectedRequestBody: "{\"ErrorMsg\":\"rating should be from 1.0 to 5.0\"}\n",
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			service := mock_service.NewMockVideoService(c)
			testCase.mockBehavior(service, testCase.uuid, testCase.rating)

			logger := logging.GetLogger("../../../../logs", "test.log")
			middleware := http.NewMiddlewares(logger)
			VideoHandler := http.NewVideoHandler(service, logger)

			router := httprouter.New()
			VideoHandler.Register(router, &middleware)

			w := httptest.NewRecorder()

			req := httptest.NewRequest("PUT", fmt.Sprintf("/rate/video%s", testCase.input), nil)

			router.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedRequestBody, w.Body.String())
		})
	}
}
