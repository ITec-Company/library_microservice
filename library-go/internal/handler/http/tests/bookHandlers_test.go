package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
	"image"
	"image/jpeg"
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

func TestBookHandler_GetAll(t *testing.T) {
	type mockBehavior func(s *mock_service.MockBookService, sortingOptions *domain.SortFilterPagination)

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
			mockBehavior: func(s *mock_service.MockBookService, sortingOptions *domain.SortFilterPagination) {
				s.EXPECT().GetAll(sortingOptions).Return([]*domain.Book{domain.TestBook(), domain.TestBook()}, 1, nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: "[{\"uuid\":\"1\",\"title\":\"Test Title\",\"direction\":{\"uuid\":\"1\",\"name\":\"Test Direction\"},\"author\":{\"uuid\":\"1\",\"full_name\":\"Test Author\"},\"difficulty\":\"Test Difficulty\",\"edition_date\":2000,\"rating\":5,\"description\":\"Test Description\",\"local_url\":\"Test LocalURL\",\"language\":\"Test Language\",\"tags\":[{\"uuid\":\"1\",\"name\":\"Test Tag\"}],\"download_count\":10,\"image_url\":\"imageURL\",\"created_at\":\"0001-01-01T00:00:00Z\"},{\"uuid\":\"1\",\"title\":\"Test Title\",\"direction\":{\"uuid\":\"1\",\"name\":\"Test Direction\"},\"author\":{\"uuid\":\"1\",\"full_name\":\"Test Author\"},\"difficulty\":\"Test Difficulty\",\"edition_date\":2000,\"rating\":5,\"description\":\"Test Description\",\"local_url\":\"Test LocalURL\",\"language\":\"Test Language\",\"tags\":[{\"uuid\":\"1\",\"name\":\"Test Tag\"}],\"download_count\":10,\"image_url\":\"imageURL\",\"created_at\":\"0001-01-01T00:00:00Z\"}]\n",
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
			mockBehavior: func(s *mock_service.MockBookService, sortingOptions *domain.SortFilterPagination) {
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

			service := mock_service.NewMockBookService(c)
			testCase.mockBehavior(service, &testCase.sortingOptions)

			logger := logging.GetLogger("../../../../logs", "test.log")
			BookHandler := http.NewBookHandler(service, logger)

			router := httprouter.New()
			middleware := http.NewMiddlewares(logger)
			BookHandler.Register(router, &middleware)

			w := httptest.NewRecorder()

			req := httptest.NewRequest("GET", fmt.Sprintf("/books%s", testCase.input), nil)

			router.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedRequestBody, w.Body.String())
		})
	}
}

func TestBookHandler_GetByUUID(t *testing.T) {
	type mockBehavior func(s *mock_service.MockBookService, uuid string)

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
			mockBehavior: func(s *mock_service.MockBookService, uuid string) {
				s.EXPECT().GetByUUID(uuid).Return(domain.TestBook(), nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: "{\"uuid\":\"1\",\"title\":\"Test Title\",\"direction\":{\"uuid\":\"1\",\"name\":\"Test Direction\"},\"author\":{\"uuid\":\"1\",\"full_name\":\"Test Author\"},\"difficulty\":\"Test Difficulty\",\"edition_date\":2000,\"rating\":5,\"description\":\"Test Description\",\"local_url\":\"Test LocalURL\",\"language\":\"Test Language\",\"tags\":[{\"uuid\":\"1\",\"name\":\"Test Tag\"}],\"download_count\":10,\"image_url\":\"imageURL\",\"created_at\":\"0001-01-01T00:00:00Z\"}\n",
		},
		{
			name:                "invalid uuid",
			input:               "-1",
			uuid:                "-1",
			mockBehavior:        func(s *mock_service.MockBookService, uuid string) {},
			expectedStatusCode:  400,
			expectedRequestBody: "{\"ErrorMsg\":\"uuid need to be uint\"}\n",
		},
		{
			name:  "no rows in result",
			input: "1",
			uuid:  "1",
			mockBehavior: func(s *mock_service.MockBookService, uuid string) {
				s.EXPECT().GetByUUID(uuid).Return(nil, errors.New("now rows in result"))
			},
			expectedStatusCode:  404,
			expectedRequestBody: "{\"ErrorMsg\":\"error occurred while getting book from DB by UUID. err: now rows in result\"}\n",
		},
		{
			name:                "string input",
			input:               "one",
			uuid:                "one",
			mockBehavior:        func(s *mock_service.MockBookService, uuid string) {},
			expectedStatusCode:  400,
			expectedRequestBody: "{\"ErrorMsg\":\"uuid need to be uint\"}\n",
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			service := mock_service.NewMockBookService(c)
			testCase.mockBehavior(service, testCase.uuid)

			logger := logging.GetLogger("../../../../logs", "test.log")
			middleware := http.NewMiddlewares(logger)
			BookHandler := http.NewBookHandler(service, logger)

			router := httprouter.New()
			BookHandler.Register(router, &middleware)

			w := httptest.NewRecorder()

			req := httptest.NewRequest("GET", fmt.Sprintf("/book/%s", testCase.input), nil)

			router.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedRequestBody, w.Body.String())
		})
	}
}

func TestBookHandler_Create(t *testing.T) {
	type mockBehavior func(s *mock_service.MockBookService, createBookDTO *domain.CreateBookDTO)
	type mockSaveFile func(s *mock_service.MockBookService, path, fileName string, file io.Reader)
	type mockSaveImage func(s *mock_service.MockBookService, path string, image *image.Image)

	file, _ := os.Open("test_file_book.pdf")
	defer file.Close()

	fileinfo, _ := file.Stat()
	fileBytes := make([]byte, fileinfo.Size())

	file.Read(fileBytes)

	imgFile, _ := os.Open("test_file_image.jpg")

	defer imgFile.Close()

	img, _ := jpeg.Decode(imgFile)

	testTable := []struct {
		name                string
		contentType         string
		mockBehavior        mockBehavior
		mockSaveFile        mockSaveFile
		mockSaveImage       mockSaveImage
		createBookDTO       domain.CreateBookDTO
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:          "OK",
			contentType:   "MultipartFormData",
			createBookDTO: *domain.TestBookCreateDTO(),
			mockBehavior: func(s *mock_service.MockBookService, createBookDTO *domain.CreateBookDTO) {
				s.EXPECT().Create(createBookDTO).Return("1", nil)
			},
			mockSaveFile: func(s *mock_service.MockBookService, path, fileName string, file io.Reader) {
				s.EXPECT().SaveFile("../store/books/1/", "author(1)-title(Test_Title).pdf", bytes.NewBuffer(fileBytes)).Return(nil)
			},
			mockSaveImage: func(s *mock_service.MockBookService, path string, img *image.Image) {
				s.EXPECT().SaveImage("../store/books/1/", img).Return(nil)
			},
			expectedStatusCode:  201,
			expectedRequestBody: "{\"infoMsg\":\"Book created successfully. UUID: 1\"}\n",
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

				file, _ := os.Open("test_file_book.pdf")
				part, _ := writer.CreateFormFile("file", "test_file_book.pdf")
				io.Copy(part, file)
				defer file.Close()

				imgFile, _ := os.Open("test_file_image.jpg")
				imgPart, _ := writer.CreateFormFile("image", "test_file_image.jpg")
				io.Copy(imgPart, imgFile)
				defer imgFile.Close()

				writer.WriteField("title", "Test Title")
				writer.WriteField("direction_uuid", "1")
				writer.WriteField("author_uuid", "1")
				writer.WriteField("difficulty", "Test Difficulty")
				writer.WriteField("edition_date", "2000")
				writer.WriteField("description", "Test Description")
				writer.WriteField("language", "Test Language")
				writer.WriteField("tags_uuids", `1`)
			}()

			w := httptest.NewRecorder()

			req := httptest.NewRequest("POST", fmt.Sprintf("/book"), pr)
			if testCase.contentType == "MultipartFormData" {
				req.Header.Set("Content-Type", writer.FormDataContentType())
			} else {
				req.Header.Set("Content-Type", "application/json")
			}

			defer req.Body.Close()
			service := mock_service.NewMockBookService(c)

			testCase.mockSaveFile(service, "../store/books/1/", "author(1)-title(Test_Title).pdf", file)
			testCase.mockSaveImage(service, "../store/books/1/", &img)
			testCase.mockBehavior(service, &testCase.createBookDTO)

			logger := logging.GetLogger("../../../../logs", "test.log")
			middleware := http.NewMiddlewares(logger)
			BookHandler := http.NewBookHandler(service, logger)

			router := httprouter.New()
			BookHandler.Register(router, &middleware)

			router.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedRequestBody, w.Body.String())
		})
	}

}

func TestBookHandler_Delete(t *testing.T) {
	type mockBehavior func(s *mock_service.MockBookService, uuid string, path string)

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
			mockBehavior: func(s *mock_service.MockBookService, uuid string, path string) {
				s.EXPECT().Delete(uuid, fmt.Sprintf("%s%s/", "../store/books/", uuid)).Return(nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: "{\"infoMsg\":\"Book with UUID 1 was deleted\"}\n",
		},
		{
			name:  "no rows in result",
			input: "1",
			uuid:  "1",
			mockBehavior: func(s *mock_service.MockBookService, uuid string, path string) {
				s.EXPECT().Delete(uuid, fmt.Sprintf("%s%s/", "../store/books/", uuid)).Return(errors.New("no rows affected"))
			},
			expectedStatusCode:  404,
			expectedRequestBody: "{\"ErrorMsg\":\"error occurred while deleting book from DB. err: no rows affected\"}\n",
		},
		{
			name:                "invalid uuid",
			input:               "-1",
			uuid:                "-1",
			mockBehavior:        func(s *mock_service.MockBookService, uuid string, path string) {},
			expectedStatusCode:  400,
			expectedRequestBody: "{\"ErrorMsg\":\"uuid need to be uint\"}\n",
		},
		{
			name:                "string uuid",
			input:               "one",
			uuid:                "one",
			mockBehavior:        func(s *mock_service.MockBookService, uuid string, path string) {},
			expectedStatusCode:  400,
			expectedRequestBody: "{\"ErrorMsg\":\"uuid need to be uint\"}\n",
		},
		// is it correct ?????
		{
			name:                "empty input",
			input:               "",
			uuid:                "",
			mockBehavior:        func(s *mock_service.MockBookService, uuid string, path string) {},
			expectedStatusCode:  404,
			expectedRequestBody: "404 page not found\n",
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			service := mock_service.NewMockBookService(c)
			testCase.mockBehavior(service, testCase.uuid, fmt.Sprintf("%s%s/", "../store/books/", testCase.uuid))

			logger := logging.GetLogger("../../../../logs", "test.log")
			middleware := http.NewMiddlewares(logger)
			BookHandler := http.NewBookHandler(service, logger)

			router := httprouter.New()
			BookHandler.Register(router, &middleware)

			w := httptest.NewRecorder()

			req := httptest.NewRequest("DELETE", fmt.Sprintf("/book/%s", testCase.input), nil)

			router.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedRequestBody, w.Body.String())
		})
	}
}

func TestBookHandler_Update(t *testing.T) {
	type mockBehavior func(s *mock_service.MockBookService, dto *domain.UpdateBookDTO)

	testTable := []struct {
		name                string
		inputBodyJSON       map[string]interface{}
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
				"edition_date":   2000,
				"rating":         5.5,
				"description":    "Test Description",
				"language":       "Test Language",
				"tags_uuids":     []string{"1"},
				"download_count": 10,
			},
			dto: *domain.TestBookUpdateDTO(),
			mockBehavior: func(s *mock_service.MockBookService, dto *domain.UpdateBookDTO) {
				s.EXPECT().Update(dto).Return(nil)
			},
			expectedStatusCode:  200,
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
				"edition_date":   2000,
				"rating":         5.5,
				"description":    "Test Description",
				"language":       "Test Language",
				"tags_uuids":     []string{"1"},
				"download_count": 10,
			},
			dto: *domain.TestBookUpdateDTO(),
			mockBehavior: func(s *mock_service.MockBookService, dto *domain.UpdateBookDTO) {
				s.EXPECT().Update(dto).Return(errors.New("no rows affected"))
			},
			expectedStatusCode:  404,
			expectedRequestBody: "{\"ErrorMsg\":\"error occurred while updating book into DB. err: no rows affected\"}\n",
		},
		{
			name:                "empty input body JSON or nil UUID",
			inputBodyJSON:       map[string]interface{}{},
			dto:                 *domain.TestBookUpdateDTO(),
			mockBehavior:        func(s *mock_service.MockBookService, dto *domain.UpdateBookDTO) {},
			expectedStatusCode:  400,
			expectedRequestBody: "{\"ErrorMsg\":\"error occurred while decoding JSON request. err: nil UUID\"}\n",
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			service := mock_service.NewMockBookService(c)
			testCase.mockBehavior(service, &testCase.dto)

			logger := logging.GetLogger("../../../../logs", "test.log")
			middleware := http.NewMiddlewares(logger)
			BookHandler := http.NewBookHandler(service, logger)

			router := httprouter.New()
			BookHandler.Register(router, &middleware)

			body, _ := json.Marshal(testCase.inputBodyJSON)
			req := httptest.NewRequest("PUT", fmt.Sprintf("/book"), bytes.NewReader(body))

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedRequestBody, w.Body.String())
		})
	}
}

func TestBookHandler_UpdateFile(t *testing.T) {
	type mockBehavior func(s *mock_service.MockBookService, dto *domain.UpdateBookFileDTO)

	d, _ := os.Open("test_file_book.pdf")
	dtoFile := new(bytes.Buffer)
	dtoFile.ReadFrom(d)

	testTable := []struct {
		name                string
		localURL            string
		inputBody           *io.Reader
		dto                 *domain.UpdateBookFileDTO
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:     "OK",
			localURL: "/books/1/title(software_testing).pdf",
			dto: &domain.UpdateBookFileDTO{
				UUID:        "1",
				NewFileName: "title(software_testing).pdf",
				OldFileName: "title(software_testing).pdf",
				File:        dtoFile,
				LocalURL:    fmt.Sprintf("%s|split|/%s", "/books/", "title(software_testing).pdf"),
				LocalPath:   fmt.Sprintf("%s%s/", "../store/books/", "1"),
			},
			mockBehavior: func(s *mock_service.MockBookService, dto *domain.UpdateBookFileDTO) {
				s.EXPECT().UpdateFile(dto).Return(nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: "{\"infoMsg\":\"File updated successfully\"}\n",
		},
		{
			name:     "invalid UUID",
			localURL: "/books/0/title(software_testing).pdf",
			dto: &domain.UpdateBookFileDTO{
				UUID:        "0",
				NewFileName: "title(software_testing).pdf",
				OldFileName: "title(software_testing).pdf",
				File:        dtoFile,
				LocalURL:    fmt.Sprintf("%s|split|/%s", "/books/", "title(software_testing).pdf"),
				LocalPath:   fmt.Sprintf("%s%s/", "../store/books/", "0"),
			},
			mockBehavior: func(s *mock_service.MockBookService, dto *domain.UpdateBookFileDTO) {
				s.EXPECT().UpdateFile(dto).Return(errors.New("no book with such UUID was found"))
			},
			expectedStatusCode:  404,
			expectedRequestBody: "{\"ErrorMsg\":\"error occurred while saving book into local store. err: no book with such UUID was found.\"}\n",
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			service := mock_service.NewMockBookService(c)
			testCase.mockBehavior(service, testCase.dto)

			logger := logging.GetLogger("../../../../logs", "test.log")
			middleware := http.NewMiddlewares(logger)
			BookHandler := http.NewBookHandler(service, logger)

			router := httprouter.New()
			BookHandler.Register(router, &middleware)

			w := httptest.NewRecorder()
			file, _ := os.Open("test_file_book.pdf")
			req := httptest.NewRequest("PUT", fmt.Sprintf("/file/book?localurl=%s", testCase.localURL), file)

			router.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedRequestBody, w.Body.String())
		})
	}
}

func TestBookHandler_UpdateImage(t *testing.T) {
	type mockBehavior func(s *mock_service.MockBookService, path string, image *image.Image)

	imgUpdFile, _ := os.Open("test_file_image.jpg")
	defer imgUpdFile.Close()
	img, _ := jpeg.Decode(imgUpdFile)

	testTable := []struct {
		name                string
		uuid, path, input   string
		img                 *image.Image
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:  "OK",
			uuid:  "1",
			path:  "../store/books/1/",
			input: "?uuid=1",
			img:   &img,
			mockBehavior: func(s *mock_service.MockBookService, path string, image *image.Image) {
				s.EXPECT().SaveImage(path, &img).Return(nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: "{\"infoMsg\":\"Book image updated successfully\"}\n",
		},
		{
			name:  "service error",
			uuid:  "0",
			path:  "../store/books/0/",
			input: "?uuid=0",
			img:   &img,
			mockBehavior: func(s *mock_service.MockBookService, path string, image *image.Image) {
				s.EXPECT().SaveImage(path, &img).Return(errors.New("service error"))
			},
			expectedStatusCode:  500,
			expectedRequestBody: "{\"ErrorMsg\":\"error occurred while saving book into local store. err: service error.\"}\n",
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			service := mock_service.NewMockBookService(c)
			testCase.mockBehavior(service, testCase.path, testCase.img)

			logger := logging.GetLogger("../../../../logs", "test.log")
			middleware := http.NewMiddlewares(logger)
			BookHandler := http.NewBookHandler(service, logger)

			router := httprouter.New()
			BookHandler.Register(router, &middleware)

			w := httptest.NewRecorder()

			imgFile, _ := os.Open("test_file_image.jpg")
			defer imgFile.Close()

			req := httptest.NewRequest("PUT", fmt.Sprintf("/image/book%s", testCase.input), imgFile)

			router.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedRequestBody, w.Body.String())
		})
	}
}

func TestBookHandler_Rate(t *testing.T) {
	type mockBehavior func(s *mock_service.MockBookService, uuid string, rating float32)

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
			mockBehavior: func(s *mock_service.MockBookService, uuid string, rating float32) {
				s.EXPECT().Rate(uuid, rating).Return(nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: "{\"infoMsg\":\"Book rated successfully. UUID: 1\"}\n",
		},
		{
			name:   "service error",
			uuid:   "1",
			rating: 1,
			input:  "?uuid=1&rating=1",
			mockBehavior: func(s *mock_service.MockBookService, uuid string, rating float32) {
				s.EXPECT().Rate(uuid, rating).Return(errors.New("service error"))
			},
			expectedStatusCode:  400,
			expectedRequestBody: "{\"ErrorMsg\":\"error occurred while rating book into local store. err: service error.\"}\n",
		},
		{
			name:   "invalid UUID",
			uuid:   "0",
			rating: 1,
			input:  "?uuid=0&rating=1",
			mockBehavior: func(s *mock_service.MockBookService, uuid string, rating float32) {
				s.EXPECT().Rate(uuid, rating).Return(errors.New("invalid UUID"))
			},
			expectedStatusCode:  400,
			expectedRequestBody: "{\"ErrorMsg\":\"error occurred while rating book into local store. err: invalid UUID.\"}\n",
		},
		{
			name:   "empty UUID",
			rating: 1,
			input:  "?rating=1",
			mockBehavior: func(s *mock_service.MockBookService, uuid string, rating float32) {
			},
			expectedStatusCode:  400,
			expectedRequestBody: "{\"ErrorMsg\":\"uuid can't be empty\"}\n",
		},
		{
			name:  "empty rating",
			uuid:  "1",
			input: "?uuid=1",
			mockBehavior: func(s *mock_service.MockBookService, uuid string, rating float32) {
			},
			expectedStatusCode:  400,
			expectedRequestBody: "{\"ErrorMsg\":\"rating query can't be empty\"}\n",
		},
		{
			name:  "invalid rating",
			uuid:  "1",
			input: "?uuid=1&rating=five",
			mockBehavior: func(s *mock_service.MockBookService, uuid string, rating float32) {
			},
			expectedStatusCode:  400,
			expectedRequestBody: "{\"ErrorMsg\":\"error occurred while parsing rating. Should be float32. err: strconv.ParseFloat: parsing \\\"five\\\": invalid syntax.\"}\n",
		},
		{
			name:   "bad rating value",
			uuid:   "1",
			rating: -1,
			input:  "?uuid=1&rating=-1",
			mockBehavior: func(s *mock_service.MockBookService, uuid string, rating float32) {
			},
			expectedStatusCode:  400,
			expectedRequestBody: "{\"ErrorMsg\":\"rating should be from 1.0 to 5.0\"}\n",
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			service := mock_service.NewMockBookService(c)
			testCase.mockBehavior(service, testCase.uuid, testCase.rating)

			logger := logging.GetLogger("../../../../logs", "test.log")
			middleware := http.NewMiddlewares(logger)
			BookHandler := http.NewBookHandler(service, logger)

			router := httprouter.New()
			BookHandler.Register(router, &middleware)

			w := httptest.NewRecorder()

			req := httptest.NewRequest("PUT", fmt.Sprintf("/rate/book%s", testCase.input), nil)

			router.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedRequestBody, w.Body.String())
		})
	}
}
