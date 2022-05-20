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

func TestTagHandler_GetAll(t *testing.T) {
	type mockBehavior func(s *mock_service.MockTagService, ctx context.Context, limit, offset int)

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
			mockBehavior: func(s *mock_service.MockTagService, ctx context.Context, limit, offset int) {
				s.EXPECT().GetAll(ctx, limit, offset).Return([]*domain.Tag{domain.TestTag(), domain.TestTag()}, nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: "[{\"uuid\":\"1\",\"name\":\"Test Tag\"},{\"uuid\":\"1\",\"name\":\"Test Tag\"}]\n",
		},
		{
			name:   "OK empty input",
			input:  "",
			ctx:    context.Background(),
			limit:  0,
			offset: 0,
			mockBehavior: func(s *mock_service.MockTagService, ctx context.Context, limit, offset int) {
				s.EXPECT().GetAll(ctx, limit, offset).Return([]*domain.Tag{domain.TestTag(), domain.TestTag()}, nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: "[{\"uuid\":\"1\",\"name\":\"Test Tag\"},{\"uuid\":\"1\",\"name\":\"Test Tag\"}]\n",
		},
		{
			name:   "OK invalid input",
			input:  "limit=one&offset=-10",
			ctx:    context.Background(),
			limit:  0,
			offset: 0,
			mockBehavior: func(s *mock_service.MockTagService, ctx context.Context, limit, offset int) {
				s.EXPECT().GetAll(ctx, limit, offset).Return([]*domain.Tag{domain.TestTag(), domain.TestTag()}, nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: "[{\"uuid\":\"1\",\"name\":\"Test Tag\"},{\"uuid\":\"1\",\"name\":\"Test Tag\"}]\n",
		},
		{
			name:   "No rows in result",
			input:  "limit=0&offset=0",
			ctx:    context.Background(),
			limit:  0,
			offset: 0,
			mockBehavior: func(s *mock_service.MockTagService, ctx context.Context, limit, offset int) {
				s.EXPECT().GetAll(ctx, limit, offset).Return(nil, errors.New("now rows in result"))
			},
			expectedStatusCode:  500,
			expectedRequestBody: "{\"ErrorMsg\":\"error occurred while getting all tags. err: now rows in result\"}\n",
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			service := mock_service.NewMockTagService(c)
			testCase.mockBehavior(service, testCase.ctx, testCase.limit, testCase.offset)

			logger := logging.GetLogger("../../../../logs", "test.log")
			middleware := NewMiddlewares(logger)
			TagHandler := NewTagHandler(service, logger, &middleware)

			router := httprouter.New()
			TagHandler.Register(router)

			w := httptest.NewRecorder()

			req := httptest.NewRequest("GET", fmt.Sprintf("/tags?%s", testCase.input), nil)

			router.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedRequestBody, w.Body.String())
		})
	}
}

func TestTagHandler_GetByUUID(t *testing.T) {
	type mockBehavior func(s *mock_service.MockTagService, ctx context.Context, uuid string)

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
			mockBehavior: func(s *mock_service.MockTagService, ctx context.Context, uuid string) {
				s.EXPECT().GetByUUID(ctx, uuid).Return(domain.TestTag(), nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: "{\"uuid\":\"1\",\"name\":\"Test Tag\"}\n",
		},
		{
			name:                "invalid uuid",
			input:               "-1",
			uuid:                "-1",
			ctx:                 context.Background(),
			mockBehavior:        func(s *mock_service.MockTagService, ctx context.Context, uuid string) {},
			expectedStatusCode:  400,
			expectedRequestBody: "{\"ErrorMsg\":\"uuid need to be uint\"}\n",
		},
		{
			name:  "no rows in result",
			input: "1",
			uuid:  "1",
			ctx:   context.Background(),
			mockBehavior: func(s *mock_service.MockTagService, ctx context.Context, uuid string) {
				s.EXPECT().GetByUUID(ctx, uuid).Return(nil, errors.New("now rows in result"))
			},
			expectedStatusCode:  404,
			expectedRequestBody: "{\"ErrorMsg\":\"error occurred while getting tag from DB by UUID. err: now rows in result\"}\n",
		},
		{
			name:                "string input",
			input:               "one",
			uuid:                "one",
			ctx:                 context.Background(),
			mockBehavior:        func(s *mock_service.MockTagService, ctx context.Context, uuid string) {},
			expectedStatusCode:  400,
			expectedRequestBody: "{\"ErrorMsg\":\"uuid need to be uint\"}\n",
		},
		// is it correct ?????
		{
			name:                "empty input",
			input:               "",
			uuid:                "",
			ctx:                 context.Background(),
			mockBehavior:        func(s *mock_service.MockTagService, ctx context.Context, uuid string) {},
			expectedStatusCode:  301,
			expectedRequestBody: "<a href=\"/tag\">Moved Permanently</a>.\n\n",
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			service := mock_service.NewMockTagService(c)
			testCase.mockBehavior(service, testCase.ctx, testCase.uuid)

			logger := logging.GetLogger("../../../../logs", "test.log")
			middleware := NewMiddlewares(logger)
			TagHandler := NewTagHandler(service, logger, &middleware)

			router := httprouter.New()
			TagHandler.Register(router)

			w := httptest.NewRecorder()

			req := httptest.NewRequest("GET", fmt.Sprintf("/tag/%s", testCase.input), nil)

			router.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedRequestBody, w.Body.String())
		})
	}
}

func TestTagHandler_GetManyByUUIDs(t *testing.T) {
	type mockBehavior func(s *mock_service.MockTagService, ctx context.Context, uuids []string)

	testTable := []struct {
		name, input         string
		uuids               []string
		ctx                 context.Context
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:  "OK",
			input: "1,2",
			uuids: []string{"1", "2"},
			ctx:   context.Background(),
			mockBehavior: func(s *mock_service.MockTagService, ctx context.Context, uuids []string) {
				s.EXPECT().GetManyByUUIDs(ctx, uuids).Return([]*domain.Tag{domain.TestTag(), &domain.Tag{UUID: "2", Name: "Test Tag2"}}, nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: "[{\"uuid\":\"1\",\"name\":\"Test Tag\"},{\"uuid\":\"2\",\"name\":\"Test Tag2\"}]\n",
		},
		{
			name:                "invalid uuid",
			input:               "-1",
			uuids:               []string{"-1"},
			ctx:                 context.Background(),
			mockBehavior:        func(s *mock_service.MockTagService, ctx context.Context, uuids []string) {},
			expectedStatusCode:  400,
			expectedRequestBody: "{\"ErrorMsg\":\"uuid need to be uint\"}\n",
		},
		{
			name:  "no rows in result",
			input: "1",
			uuids: []string{"1"},
			ctx:   context.Background(),
			mockBehavior: func(s *mock_service.MockTagService, ctx context.Context, uuids []string) {
				s.EXPECT().GetManyByUUIDs(ctx, uuids).Return(nil, errors.New("now rows in result"))
			},
			expectedStatusCode:  404,
			expectedRequestBody: "{\"ErrorMsg\":\"error occurred while getting tags from DB by UUIDs. err: now rows in result\"}\n",
		},
		{
			name:                "string input",
			input:               "one",
			uuids:               []string{"one"},
			ctx:                 context.Background(),
			mockBehavior:        func(s *mock_service.MockTagService, ctx context.Context, uuids []string) {},
			expectedStatusCode:  400,
			expectedRequestBody: "{\"ErrorMsg\":\"uuid need to be uint\"}\n",
		},
		// is it correct ?????
		{
			name:                "empty input",
			input:               "",
			uuids:               []string{},
			ctx:                 context.Background(),
			mockBehavior:        func(s *mock_service.MockTagService, ctx context.Context, uuids []string) {},
			expectedStatusCode:  301,
			expectedRequestBody: "<a href=\"/tags\">Moved Permanently</a>.\n\n",
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			service := mock_service.NewMockTagService(c)
			testCase.mockBehavior(service, testCase.ctx, testCase.uuids)

			logger := logging.GetLogger("../../../../logs", "test.log")
			middleware := NewMiddlewares(logger)
			TagHandler := NewTagHandler(service, logger, &middleware)

			router := httprouter.New()
			TagHandler.Register(router)

			w := httptest.NewRecorder()

			req := httptest.NewRequest("GET", fmt.Sprintf("/tags/%s", testCase.input), nil)

			router.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedRequestBody, w.Body.String())
		})
	}
}

func TestTagHandler_Create(t *testing.T) {
	type mockBehavior func(s *mock_service.MockTagService, ctx context.Context, dto domain.CreateTagDTO)

	testTable := []struct {
		name                string
		inputBodyJSON       map[string]interface{}
		ctx                 context.Context
		dto                 domain.CreateTagDTO
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name: "OK",
			inputBodyJSON: map[string]interface{}{
				"name": "Test Tag",
			},
			ctx: context.Background(),
			dto: *domain.TestTagCreateDTO(),
			mockBehavior: func(s *mock_service.MockTagService, ctx context.Context, dto domain.CreateTagDTO) {
				s.EXPECT().Create(ctx, &dto).Return("1", nil)
			},
			expectedStatusCode:  201,
			expectedRequestBody: "{\"infoMsg\":\"Tag created successfully. UUID: 1\"}\n",
		},
		{
			name: "service error",
			inputBodyJSON: map[string]interface{}{
				"name": "Test Tag",
			},
			ctx: context.Background(),
			dto: *domain.TestTagCreateDTO(),
			mockBehavior: func(s *mock_service.MockTagService, ctx context.Context, dto domain.CreateTagDTO) {
				s.EXPECT().Create(ctx, &dto).Return("", errors.New("some service error"))
			},
			expectedStatusCode:  500,
			expectedRequestBody: "{\"ErrorMsg\":\"error occurred while creating tag into DB. err: some service error\"}\n",
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			service := mock_service.NewMockTagService(c)
			testCase.mockBehavior(service, testCase.ctx, testCase.dto)

			logger := logging.GetLogger("../../../../logs", "test.log")
			middleware := NewMiddlewares(logger)
			TagHandler := NewTagHandler(service, logger, &middleware)

			router := httprouter.New()
			TagHandler.Register(router)

			body, _ := json.Marshal(testCase.inputBodyJSON)
			req := httptest.NewRequest("POST", fmt.Sprintf("/tag"), bytes.NewReader(body))

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedRequestBody, w.Body.String())
		})
	}

}

func TestTagHandler_Delete(t *testing.T) {
	type mockBehavior func(s *mock_service.MockTagService, ctx context.Context, uuid string)

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
			mockBehavior: func(s *mock_service.MockTagService, ctx context.Context, uuid string) {
				s.EXPECT().Delete(ctx, uuid).Return(nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: "{\"infoMsg\":\"Tag with UUID 1 was deleted\"}\n",
		},
		{
			name:  "no rows in result",
			input: "1",
			uuid:  "1",
			ctx:   context.Background(),
			mockBehavior: func(s *mock_service.MockTagService, ctx context.Context, uuid string) {
				s.EXPECT().Delete(ctx, uuid).Return(errors.New("no rows affected"))
			},
			expectedStatusCode:  404,
			expectedRequestBody: "{\"ErrorMsg\":\"error occurred while deleting tag from DB. err: no rows affected\"}\n",
		},
		{
			name:                "invalid uuid",
			input:               "-1",
			uuid:                "-1",
			ctx:                 context.Background(),
			mockBehavior:        func(s *mock_service.MockTagService, ctx context.Context, uuid string) {},
			expectedStatusCode:  400,
			expectedRequestBody: "{\"ErrorMsg\":\"uuid need to be uint\"}\n",
		},
		{
			name:                "string uuid",
			input:               "one",
			uuid:                "one",
			ctx:                 context.Background(),
			mockBehavior:        func(s *mock_service.MockTagService, ctx context.Context, uuid string) {},
			expectedStatusCode:  400,
			expectedRequestBody: "{\"ErrorMsg\":\"uuid need to be uint\"}\n",
		},
		// is it correct ?????
		{
			name:                "empty input",
			input:               "",
			uuid:                "",
			ctx:                 context.Background(),
			mockBehavior:        func(s *mock_service.MockTagService, ctx context.Context, uuid string) {},
			expectedStatusCode:  404,
			expectedRequestBody: "404 page not found\n",
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			service := mock_service.NewMockTagService(c)
			testCase.mockBehavior(service, testCase.ctx, testCase.uuid)

			logger := logging.GetLogger("../../../../logs", "test.log")
			middleware := NewMiddlewares(logger)
			TagHandler := NewTagHandler(service, logger, &middleware)

			router := httprouter.New()
			TagHandler.Register(router)

			w := httptest.NewRecorder()

			req := httptest.NewRequest("DELETE", fmt.Sprintf("/tag/%s", testCase.input), nil)

			router.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedRequestBody, w.Body.String())
		})
	}
}

func TestTagHandler_Update(t *testing.T) {
	type mockBehavior func(s *mock_service.MockTagService, ctx context.Context, dto *domain.UpdateTagDTO)

	testTable := []struct {
		name                string
		inputBodyJSON       map[string]interface{}
		ctx                 context.Context
		dto                 domain.UpdateTagDTO
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name: "OK",
			inputBodyJSON: map[string]interface{}{
				"uuid": "1",
				"name": "Test Tag",
			},
			ctx: context.Background(),
			dto: *domain.TestTagUpdateDTO(),
			mockBehavior: func(s *mock_service.MockTagService, ctx context.Context, dto *domain.UpdateTagDTO) {
				s.EXPECT().Update(ctx, dto).Return(nil)
			},
			expectedStatusCode:  201,
			expectedRequestBody: "{\"infoMsg\":\"Tag updated successfully\"}\n",
		},
		{
			name: "no rows in result",
			inputBodyJSON: map[string]interface{}{
				"uuid": "1",
				"name": "Test Tag",
			},
			ctx: context.Background(),
			dto: *domain.TestTagUpdateDTO(),
			mockBehavior: func(s *mock_service.MockTagService, ctx context.Context, dto *domain.UpdateTagDTO) {
				s.EXPECT().Update(ctx, dto).Return(errors.New("no rows affected"))
			},
			expectedStatusCode:  500,
			expectedRequestBody: "{\"ErrorMsg\":\"error occurred while updating tag into DB. err: no rows affected\"}\n",
		},
		{
			name:                "empty input body JSON or nil UUID",
			inputBodyJSON:       map[string]interface{}{},
			ctx:                 context.Background(),
			dto:                 *domain.TestTagUpdateDTO(),
			mockBehavior:        func(s *mock_service.MockTagService, ctx context.Context, dto *domain.UpdateTagDTO) {},
			expectedStatusCode:  400,
			expectedRequestBody: "{\"ErrorMsg\":\"error occurred while decoding JSON request. err: nil UUID\"}\n",
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			service := mock_service.NewMockTagService(c)
			testCase.mockBehavior(service, testCase.ctx, &testCase.dto)

			logger := logging.GetLogger("../../../../logs", "test.log")
			middleware := NewMiddlewares(logger)
			TagHandler := NewTagHandler(service, logger, &middleware)

			router := httprouter.New()
			TagHandler.Register(router)

			body, _ := json.Marshal(testCase.inputBodyJSON)
			req := httptest.NewRequest("PUT", fmt.Sprintf("/tag"), bytes.NewReader(body))

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedRequestBody, w.Body.String())
		})
	}
}
