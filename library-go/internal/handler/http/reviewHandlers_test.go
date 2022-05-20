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

func TestReviewHandler_GetAll(t *testing.T) {
	type mockBehavior func(s *mock_service.MockReviewService, ctx context.Context, limit, offset int)

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
			mockBehavior: func(s *mock_service.MockReviewService, ctx context.Context, limit, offset int) {
				s.EXPECT().GetAll(ctx, limit, offset).Return([]*domain.Review{domain.TestReview(), domain.TestReview()}, nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: "[{\"uuid\":\"1\",\"full_name\":\"Test Review\",\"text\":\"Test Text\",\"rating\":5.5,\"date\":\"2000-01-01T00:00:00Z\",\"source\":\"Test Source\",\"literature_uuid\":\"1\"},{\"uuid\":\"1\",\"full_name\":\"Test Review\",\"text\":\"Test Text\",\"rating\":5.5,\"date\":\"2000-01-01T00:00:00Z\",\"source\":\"Test Source\",\"literature_uuid\":\"1\"}]\n",
		},
		{
			name:   "OK empty input",
			input:  "",
			ctx:    context.Background(),
			limit:  0,
			offset: 0,
			mockBehavior: func(s *mock_service.MockReviewService, ctx context.Context, limit, offset int) {
				s.EXPECT().GetAll(ctx, limit, offset).Return([]*domain.Review{domain.TestReview(), domain.TestReview()}, nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: "[{\"uuid\":\"1\",\"full_name\":\"Test Review\",\"text\":\"Test Text\",\"rating\":5.5,\"date\":\"2000-01-01T00:00:00Z\",\"source\":\"Test Source\",\"literature_uuid\":\"1\"},{\"uuid\":\"1\",\"full_name\":\"Test Review\",\"text\":\"Test Text\",\"rating\":5.5,\"date\":\"2000-01-01T00:00:00Z\",\"source\":\"Test Source\",\"literature_uuid\":\"1\"}]\n",
		},
		{
			name:   "OK invalid input",
			input:  "limit=one&offset=-10",
			ctx:    context.Background(),
			limit:  0,
			offset: 0,
			mockBehavior: func(s *mock_service.MockReviewService, ctx context.Context, limit, offset int) {
				s.EXPECT().GetAll(ctx, limit, offset).Return([]*domain.Review{domain.TestReview(), domain.TestReview()}, nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: "[{\"uuid\":\"1\",\"full_name\":\"Test Review\",\"text\":\"Test Text\",\"rating\":5.5,\"date\":\"2000-01-01T00:00:00Z\",\"source\":\"Test Source\",\"literature_uuid\":\"1\"},{\"uuid\":\"1\",\"full_name\":\"Test Review\",\"text\":\"Test Text\",\"rating\":5.5,\"date\":\"2000-01-01T00:00:00Z\",\"source\":\"Test Source\",\"literature_uuid\":\"1\"}]\n",
		},
		{
			name:   "No rows in result",
			input:  "limit=0&offset=0",
			ctx:    context.Background(),
			limit:  0,
			offset: 0,
			mockBehavior: func(s *mock_service.MockReviewService, ctx context.Context, limit, offset int) {
				s.EXPECT().GetAll(ctx, limit, offset).Return(nil, errors.New("now rows in result"))
			},
			expectedStatusCode:  500,
			expectedRequestBody: "{\"ErrorMsg\":\"error occurred while getting all reviews. err: now rows in result\"}\n",
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			service := mock_service.NewMockReviewService(c)
			testCase.mockBehavior(service, testCase.ctx, testCase.limit, testCase.offset)

			logger := logging.GetLogger()
			middleware := NewMiddlewares(logger)
			ReviewHandler := NewReviewHandler(service, logger, &middleware)

			router := httprouter.New()
			ReviewHandler.Register(router)

			w := httptest.NewRecorder()

			req := httptest.NewRequest("GET", fmt.Sprintf("/reviews?%s", testCase.input), nil)

			router.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedRequestBody, w.Body.String())
		})
	}
}

func TestReviewHandler_GetByUUID(t *testing.T) {
	type mockBehavior func(s *mock_service.MockReviewService, ctx context.Context, uuid string)

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
			mockBehavior: func(s *mock_service.MockReviewService, ctx context.Context, uuid string) {
				s.EXPECT().GetByUUID(ctx, uuid).Return(domain.TestReview(), nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: "{\"uuid\":\"1\",\"full_name\":\"Test Review\",\"text\":\"Test Text\",\"rating\":5.5,\"date\":\"2000-01-01T00:00:00Z\",\"source\":\"Test Source\",\"literature_uuid\":\"1\"}\n",
		},
		{
			name:                "invalid uuid",
			input:               "-1",
			uuid:                "-1",
			ctx:                 context.Background(),
			mockBehavior:        func(s *mock_service.MockReviewService, ctx context.Context, uuid string) {},
			expectedStatusCode:  400,
			expectedRequestBody: "{\"ErrorMsg\":\"uuid need to be uint\"}\n",
		},
		{
			name:  "no rows in result",
			input: "1",
			uuid:  "1",
			ctx:   context.Background(),
			mockBehavior: func(s *mock_service.MockReviewService, ctx context.Context, uuid string) {
				s.EXPECT().GetByUUID(ctx, uuid).Return(nil, errors.New("now rows in result"))
			},
			expectedStatusCode:  404,
			expectedRequestBody: "{\"ErrorMsg\":\"error occurred while getting review from DB by UUID. err: now rows in result\"}\n",
		},
		{
			name:                "string input",
			input:               "one",
			uuid:                "one",
			ctx:                 context.Background(),
			mockBehavior:        func(s *mock_service.MockReviewService, ctx context.Context, uuid string) {},
			expectedStatusCode:  400,
			expectedRequestBody: "{\"ErrorMsg\":\"uuid need to be uint\"}\n",
		},
		// is it correct ?????
		{
			name:                "empty input",
			input:               "",
			uuid:                "",
			ctx:                 context.Background(),
			mockBehavior:        func(s *mock_service.MockReviewService, ctx context.Context, uuid string) {},
			expectedStatusCode:  301,
			expectedRequestBody: "<a href=\"/review\">Moved Permanently</a>.\n\n",
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			service := mock_service.NewMockReviewService(c)
			testCase.mockBehavior(service, testCase.ctx, testCase.uuid)

			logger := logging.GetLogger()
			middleware := NewMiddlewares(logger)
			ReviewHandler := NewReviewHandler(service, logger, &middleware)

			router := httprouter.New()
			ReviewHandler.Register(router)

			w := httptest.NewRecorder()

			req := httptest.NewRequest("GET", fmt.Sprintf("/review/%s", testCase.input), nil)

			router.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedRequestBody, w.Body.String())
		})
	}
}

func TestReviewHandler_Create(t *testing.T) {
	type mockBehavior func(s *mock_service.MockReviewService, ctx context.Context, dto domain.CreateReviewDTO)

	testTable := []struct {
		name                string
		inputBodyJSON       map[string]interface{}
		ctx                 context.Context
		dto                 domain.CreateReviewDTO
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name: "OK",
			inputBodyJSON: map[string]interface{}{
				"text":            "Test Text",
				"full_name":       "Test Review",
				"literature_uuid": "1",
				"source":          "Test Source",
			},
			ctx: context.Background(),
			dto: *domain.TestReviewCreateDTO(),
			mockBehavior: func(s *mock_service.MockReviewService, ctx context.Context, dto domain.CreateReviewDTO) {
				s.EXPECT().Create(ctx, &dto).Return("1", nil)
			},
			expectedStatusCode:  201,
			expectedRequestBody: "{\"infoMsg\":\"Review created successfully. UUID: 1\"}\n",
		},
		{
			name: "service error",
			inputBodyJSON: map[string]interface{}{
				"text":            "Test Text",
				"full_name":       "Test Review",
				"literature_uuid": "1",
				"source":          "Test Source",
			},
			ctx: context.Background(),
			dto: *domain.TestReviewCreateDTO(),
			mockBehavior: func(s *mock_service.MockReviewService, ctx context.Context, dto domain.CreateReviewDTO) {
				s.EXPECT().Create(ctx, &dto).Return("", errors.New("some service error"))
			},
			expectedStatusCode:  500,
			expectedRequestBody: "{\"ErrorMsg\":\"error occurred while creating review into DB. err: some service error\"}\n",
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			service := mock_service.NewMockReviewService(c)
			testCase.mockBehavior(service, testCase.ctx, testCase.dto)

			logger := logging.GetLogger()
			middleware := NewMiddlewares(logger)
			ReviewHandler := NewReviewHandler(service, logger, &middleware)

			router := httprouter.New()
			ReviewHandler.Register(router)

			body, _ := json.Marshal(testCase.inputBodyJSON)
			req := httptest.NewRequest("POST", fmt.Sprintf("/review"), bytes.NewReader(body))

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedRequestBody, w.Body.String())
		})
	}

}

func TestReviewHandler_Delete(t *testing.T) {
	type mockBehavior func(s *mock_service.MockReviewService, ctx context.Context, uuid string)

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
			mockBehavior: func(s *mock_service.MockReviewService, ctx context.Context, uuid string) {
				s.EXPECT().Delete(ctx, uuid).Return(nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: "{\"infoMsg\":\"Review with UUID 1 was deleted\"}\n",
		},
		{
			name:  "no rows in result",
			input: "1",
			uuid:  "1",
			ctx:   context.Background(),
			mockBehavior: func(s *mock_service.MockReviewService, ctx context.Context, uuid string) {
				s.EXPECT().Delete(ctx, uuid).Return(errors.New("no rows affected"))
			},
			expectedStatusCode:  404,
			expectedRequestBody: "{\"ErrorMsg\":\"error occurred while deleting review from DB. err: no rows affected\"}\n",
		},
		{
			name:                "invalid uuid",
			input:               "-1",
			uuid:                "-1",
			ctx:                 context.Background(),
			mockBehavior:        func(s *mock_service.MockReviewService, ctx context.Context, uuid string) {},
			expectedStatusCode:  400,
			expectedRequestBody: "{\"ErrorMsg\":\"uuid need to be uint\"}\n",
		},
		{
			name:                "string uuid",
			input:               "one",
			uuid:                "one",
			ctx:                 context.Background(),
			mockBehavior:        func(s *mock_service.MockReviewService, ctx context.Context, uuid string) {},
			expectedStatusCode:  400,
			expectedRequestBody: "{\"ErrorMsg\":\"uuid need to be uint\"}\n",
		},
		// is it correct ?????
		{
			name:                "empty input",
			input:               "",
			uuid:                "",
			ctx:                 context.Background(),
			mockBehavior:        func(s *mock_service.MockReviewService, ctx context.Context, uuid string) {},
			expectedStatusCode:  404,
			expectedRequestBody: "404 page not found\n",
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			service := mock_service.NewMockReviewService(c)
			testCase.mockBehavior(service, testCase.ctx, testCase.uuid)

			logger := logging.GetLogger()
			middleware := NewMiddlewares(logger)
			ReviewHandler := NewReviewHandler(service, logger, &middleware)

			router := httprouter.New()
			ReviewHandler.Register(router)

			w := httptest.NewRecorder()

			req := httptest.NewRequest("DELETE", fmt.Sprintf("/review/%s", testCase.input), nil)

			router.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedRequestBody, w.Body.String())
		})
	}
}

func TestReviewHandler_Update(t *testing.T) {
	type mockBehavior func(s *mock_service.MockReviewService, ctx context.Context, dto *domain.UpdateReviewDTO)

	testTable := []struct {
		name                string
		inputBodyJSON       map[string]interface{}
		ctx                 context.Context
		dto                 domain.UpdateReviewDTO
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name: "OK",
			inputBodyJSON: map[string]interface{}{
				"uuid":      "1",
				"full_name": "Test Review",
				"text":      "Test Text",
				"rating":    5.5,
			},
			ctx: context.Background(),
			dto: *domain.TestReviewUpdateDTO(),
			mockBehavior: func(s *mock_service.MockReviewService, ctx context.Context, dto *domain.UpdateReviewDTO) {
				s.EXPECT().Update(ctx, dto).Return(nil)
			},
			expectedStatusCode:  201,
			expectedRequestBody: "{\"infoMsg\":\"Review updated successfully\"}\n",
		},
		{
			name: "no rows in result",
			inputBodyJSON: map[string]interface{}{
				"uuid":      "1",
				"full_name": "Test Review",
				"text":      "Test Text",
				"rating":    5.5,
			},
			ctx: context.Background(),
			dto: *domain.TestReviewUpdateDTO(),
			mockBehavior: func(s *mock_service.MockReviewService, ctx context.Context, dto *domain.UpdateReviewDTO) {
				s.EXPECT().Update(ctx, dto).Return(errors.New("no rows affected"))
			},
			expectedStatusCode:  500,
			expectedRequestBody: "{\"ErrorMsg\":\"error occurred while updating review into DB. err: no rows affected\"}\n",
		},
		{
			name:                "empty input body JSON or nil UUID",
			inputBodyJSON:       map[string]interface{}{},
			ctx:                 context.Background(),
			dto:                 *domain.TestReviewUpdateDTO(),
			mockBehavior:        func(s *mock_service.MockReviewService, ctx context.Context, dto *domain.UpdateReviewDTO) {},
			expectedStatusCode:  400,
			expectedRequestBody: "{\"ErrorMsg\":\"error occurred while decoding JSON request. err: nil UUID\"}\n",
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			service := mock_service.NewMockReviewService(c)
			testCase.mockBehavior(service, testCase.ctx, &testCase.dto)

			logger := logging.GetLogger()
			middleware := NewMiddlewares(logger)
			ReviewHandler := NewReviewHandler(service, logger, &middleware)

			router := httprouter.New()
			ReviewHandler.Register(router)

			body, _ := json.Marshal(testCase.inputBodyJSON)
			req := httptest.NewRequest("PUT", fmt.Sprintf("/review"), bytes.NewReader(body))

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedRequestBody, w.Body.String())
		})
	}
}
