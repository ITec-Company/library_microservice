// Code generated by MockGen. DO NOT EDIT.
// Source: interface.go

// Package mock_service is a generated GoMock package.
package mock_service

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	domain "library-go/internal/domain"
	reflect "reflect"
)

// MockArticleService is a mock of ArticleService interface
type MockArticleService struct {
	ctrl     *gomock.Controller
	recorder *MockArticleServiceMockRecorder
}

// MockArticleServiceMockRecorder is the mock recorder for MockArticleService
type MockArticleServiceMockRecorder struct {
	mock *MockArticleService
}

// NewMockArticleService creates a new mock instance
func NewMockArticleService(ctrl *gomock.Controller) *MockArticleService {
	mock := &MockArticleService{ctrl: ctrl}
	mock.recorder = &MockArticleServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockArticleService) EXPECT() *MockArticleServiceMockRecorder {
	return m.recorder
}

// GetByUUID mocks base method
func (m *MockArticleService) GetByUUID(ctx context.Context, UUID string) (*domain.Article, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByUUID", ctx, UUID)
	ret0, _ := ret[0].(*domain.Article)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByUUID indicates an expected call of GetByUUID
func (mr *MockArticleServiceMockRecorder) GetByUUID(ctx, UUID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByUUID", reflect.TypeOf((*MockArticleService)(nil).GetByUUID), ctx, UUID)
}

// GetAll mocks base method
func (m *MockArticleService) GetAll(ctx context.Context, limit, offset int) ([]*domain.Article, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAll", ctx, limit, offset)
	ret0, _ := ret[0].([]*domain.Article)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAll indicates an expected call of GetAll
func (mr *MockArticleServiceMockRecorder) GetAll(ctx, limit, offset interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAll", reflect.TypeOf((*MockArticleService)(nil).GetAll), ctx, limit, offset)
}

// Delete mocks base method
func (m *MockArticleService) Delete(ctx context.Context, UUID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, UUID)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete
func (mr *MockArticleServiceMockRecorder) Delete(ctx, UUID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockArticleService)(nil).Delete), ctx, UUID)
}

// Create mocks base method
func (m *MockArticleService) Create(ctx context.Context, article *domain.CreateArticleDTO) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, article)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create
func (mr *MockArticleServiceMockRecorder) Create(ctx, article interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockArticleService)(nil).Create), ctx, article)
}

// Update mocks base method
func (m *MockArticleService) Update(ctx context.Context, article *domain.UpdateArticleDTO) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, article)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update
func (mr *MockArticleServiceMockRecorder) Update(ctx, article interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockArticleService)(nil).Update), ctx, article)
}

// LoadLocalFIle mocks base method
func (m *MockArticleService) LoadLocalFIle(ctx context.Context, path string) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LoadLocalFIle", ctx, path)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LoadLocalFIle indicates an expected call of LoadLocalFIle
func (mr *MockArticleServiceMockRecorder) LoadLocalFIle(ctx, path interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LoadLocalFIle", reflect.TypeOf((*MockArticleService)(nil).LoadLocalFIle), ctx, path)
}

// MockAudioService is a mock of AudioService interface
type MockAudioService struct {
	ctrl     *gomock.Controller
	recorder *MockAudioServiceMockRecorder
}

// MockAudioServiceMockRecorder is the mock recorder for MockAudioService
type MockAudioServiceMockRecorder struct {
	mock *MockAudioService
}

// NewMockAudioService creates a new mock instance
func NewMockAudioService(ctrl *gomock.Controller) *MockAudioService {
	mock := &MockAudioService{ctrl: ctrl}
	mock.recorder = &MockAudioServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockAudioService) EXPECT() *MockAudioServiceMockRecorder {
	return m.recorder
}

// GetByUUID mocks base method
func (m *MockAudioService) GetByUUID(ctx context.Context, UUID string) (*domain.Audio, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByUUID", ctx, UUID)
	ret0, _ := ret[0].(*domain.Audio)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByUUID indicates an expected call of GetByUUID
func (mr *MockAudioServiceMockRecorder) GetByUUID(ctx, UUID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByUUID", reflect.TypeOf((*MockAudioService)(nil).GetByUUID), ctx, UUID)
}

// GetAll mocks base method
func (m *MockAudioService) GetAll(ctx context.Context, limit, offset int) ([]*domain.Audio, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAll", ctx, limit, offset)
	ret0, _ := ret[0].([]*domain.Audio)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAll indicates an expected call of GetAll
func (mr *MockAudioServiceMockRecorder) GetAll(ctx, limit, offset interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAll", reflect.TypeOf((*MockAudioService)(nil).GetAll), ctx, limit, offset)
}

// Delete mocks base method
func (m *MockAudioService) Delete(ctx context.Context, UUID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, UUID)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete
func (mr *MockAudioServiceMockRecorder) Delete(ctx, UUID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockAudioService)(nil).Delete), ctx, UUID)
}

// Create mocks base method
func (m *MockAudioService) Create(ctx context.Context, audio *domain.CreateAudioDTO) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, audio)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create
func (mr *MockAudioServiceMockRecorder) Create(ctx, audio interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockAudioService)(nil).Create), ctx, audio)
}

// Update mocks base method
func (m *MockAudioService) Update(ctx context.Context, audio *domain.UpdateAudioDTO) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, audio)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update
func (mr *MockAudioServiceMockRecorder) Update(ctx, audio interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockAudioService)(nil).Update), ctx, audio)
}

// LoadLocalFIle mocks base method
func (m *MockAudioService) LoadLocalFIle(ctx context.Context, path string) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LoadLocalFIle", ctx, path)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LoadLocalFIle indicates an expected call of LoadLocalFIle
func (mr *MockAudioServiceMockRecorder) LoadLocalFIle(ctx, path interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LoadLocalFIle", reflect.TypeOf((*MockAudioService)(nil).LoadLocalFIle), ctx, path)
}

// MockAuthorService is a mock of AuthorService interface
type MockAuthorService struct {
	ctrl     *gomock.Controller
	recorder *MockAuthorServiceMockRecorder
}

// MockAuthorServiceMockRecorder is the mock recorder for MockAuthorService
type MockAuthorServiceMockRecorder struct {
	mock *MockAuthorService
}

// NewMockAuthorService creates a new mock instance
func NewMockAuthorService(ctrl *gomock.Controller) *MockAuthorService {
	mock := &MockAuthorService{ctrl: ctrl}
	mock.recorder = &MockAuthorServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockAuthorService) EXPECT() *MockAuthorServiceMockRecorder {
	return m.recorder
}

// GetByUUID mocks base method
func (m *MockAuthorService) GetByUUID(ctx context.Context, UUID string) (*domain.Author, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByUUID", ctx, UUID)
	ret0, _ := ret[0].(*domain.Author)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByUUID indicates an expected call of GetByUUID
func (mr *MockAuthorServiceMockRecorder) GetByUUID(ctx, UUID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByUUID", reflect.TypeOf((*MockAuthorService)(nil).GetByUUID), ctx, UUID)
}

// GetAll mocks base method
func (m *MockAuthorService) GetAll(ctx context.Context, limit, offset int) ([]*domain.Author, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAll", ctx, limit, offset)
	ret0, _ := ret[0].([]*domain.Author)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAll indicates an expected call of GetAll
func (mr *MockAuthorServiceMockRecorder) GetAll(ctx, limit, offset interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAll", reflect.TypeOf((*MockAuthorService)(nil).GetAll), ctx, limit, offset)
}

// Delete mocks base method
func (m *MockAuthorService) Delete(ctx context.Context, UUID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, UUID)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete
func (mr *MockAuthorServiceMockRecorder) Delete(ctx, UUID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockAuthorService)(nil).Delete), ctx, UUID)
}

// Create mocks base method
func (m *MockAuthorService) Create(ctx context.Context, author *domain.CreateAuthorDTO) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, author)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create
func (mr *MockAuthorServiceMockRecorder) Create(ctx, author interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockAuthorService)(nil).Create), ctx, author)
}

// Update mocks base method
func (m *MockAuthorService) Update(ctx context.Context, author *domain.UpdateAuthorDTO) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, author)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update
func (mr *MockAuthorServiceMockRecorder) Update(ctx, author interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockAuthorService)(nil).Update), ctx, author)
}

// MockBookService is a mock of BookService interface
type MockBookService struct {
	ctrl     *gomock.Controller
	recorder *MockBookServiceMockRecorder
}

// MockBookServiceMockRecorder is the mock recorder for MockBookService
type MockBookServiceMockRecorder struct {
	mock *MockBookService
}

// NewMockBookService creates a new mock instance
func NewMockBookService(ctrl *gomock.Controller) *MockBookService {
	mock := &MockBookService{ctrl: ctrl}
	mock.recorder = &MockBookServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockBookService) EXPECT() *MockBookServiceMockRecorder {
	return m.recorder
}

// GetByUUID mocks base method
func (m *MockBookService) GetByUUID(ctx context.Context, UUID string) (*domain.Book, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByUUID", ctx, UUID)
	ret0, _ := ret[0].(*domain.Book)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByUUID indicates an expected call of GetByUUID
func (mr *MockBookServiceMockRecorder) GetByUUID(ctx, UUID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByUUID", reflect.TypeOf((*MockBookService)(nil).GetByUUID), ctx, UUID)
}

// GetAll mocks base method
func (m *MockBookService) GetAll(ctx context.Context, limit, offset int) ([]*domain.Book, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAll", ctx, limit, offset)
	ret0, _ := ret[0].([]*domain.Book)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAll indicates an expected call of GetAll
func (mr *MockBookServiceMockRecorder) GetAll(ctx, limit, offset interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAll", reflect.TypeOf((*MockBookService)(nil).GetAll), ctx, limit, offset)
}

// Delete mocks base method
func (m *MockBookService) Delete(ctx context.Context, UUID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, UUID)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete
func (mr *MockBookServiceMockRecorder) Delete(ctx, UUID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockBookService)(nil).Delete), ctx, UUID)
}

// Create mocks base method
func (m *MockBookService) Create(ctx context.Context, book *domain.CreateBookDTO) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, book)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create
func (mr *MockBookServiceMockRecorder) Create(ctx, book interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockBookService)(nil).Create), ctx, book)
}

// Update mocks base method
func (m *MockBookService) Update(ctx context.Context, book *domain.UpdateBookDTO) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, book)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update
func (mr *MockBookServiceMockRecorder) Update(ctx, book interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockBookService)(nil).Update), ctx, book)
}

// LoadLocalFIle mocks base method
func (m *MockBookService) LoadLocalFIle(ctx context.Context, path string) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LoadLocalFIle", ctx, path)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LoadLocalFIle indicates an expected call of LoadLocalFIle
func (mr *MockBookServiceMockRecorder) LoadLocalFIle(ctx, path interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LoadLocalFIle", reflect.TypeOf((*MockBookService)(nil).LoadLocalFIle), ctx, path)
}

// MockDirectionService is a mock of DirectionService interface
type MockDirectionService struct {
	ctrl     *gomock.Controller
	recorder *MockDirectionServiceMockRecorder
}

// MockDirectionServiceMockRecorder is the mock recorder for MockDirectionService
type MockDirectionServiceMockRecorder struct {
	mock *MockDirectionService
}

// NewMockDirectionService creates a new mock instance
func NewMockDirectionService(ctrl *gomock.Controller) *MockDirectionService {
	mock := &MockDirectionService{ctrl: ctrl}
	mock.recorder = &MockDirectionServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockDirectionService) EXPECT() *MockDirectionServiceMockRecorder {
	return m.recorder
}

// GetByUUID mocks base method
func (m *MockDirectionService) GetByUUID(ctx context.Context, UUID string) (*domain.Direction, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByUUID", ctx, UUID)
	ret0, _ := ret[0].(*domain.Direction)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByUUID indicates an expected call of GetByUUID
func (mr *MockDirectionServiceMockRecorder) GetByUUID(ctx, UUID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByUUID", reflect.TypeOf((*MockDirectionService)(nil).GetByUUID), ctx, UUID)
}

// GetAll mocks base method
func (m *MockDirectionService) GetAll(ctx context.Context, limit, offset int) ([]*domain.Direction, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAll", ctx, limit, offset)
	ret0, _ := ret[0].([]*domain.Direction)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAll indicates an expected call of GetAll
func (mr *MockDirectionServiceMockRecorder) GetAll(ctx, limit, offset interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAll", reflect.TypeOf((*MockDirectionService)(nil).GetAll), ctx, limit, offset)
}

// Delete mocks base method
func (m *MockDirectionService) Delete(ctx context.Context, UUID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, UUID)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete
func (mr *MockDirectionServiceMockRecorder) Delete(ctx, UUID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockDirectionService)(nil).Delete), ctx, UUID)
}

// Create mocks base method
func (m *MockDirectionService) Create(ctx context.Context, direction *domain.CreateDirectionDTO) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, direction)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create
func (mr *MockDirectionServiceMockRecorder) Create(ctx, direction interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockDirectionService)(nil).Create), ctx, direction)
}

// Update mocks base method
func (m *MockDirectionService) Update(ctx context.Context, direction *domain.UpdateDirectionDTO) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, direction)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update
func (mr *MockDirectionServiceMockRecorder) Update(ctx, direction interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockDirectionService)(nil).Update), ctx, direction)
}

// MockReviewService is a mock of ReviewService interface
type MockReviewService struct {
	ctrl     *gomock.Controller
	recorder *MockReviewServiceMockRecorder
}

// MockReviewServiceMockRecorder is the mock recorder for MockReviewService
type MockReviewServiceMockRecorder struct {
	mock *MockReviewService
}

// NewMockReviewService creates a new mock instance
func NewMockReviewService(ctrl *gomock.Controller) *MockReviewService {
	mock := &MockReviewService{ctrl: ctrl}
	mock.recorder = &MockReviewServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockReviewService) EXPECT() *MockReviewServiceMockRecorder {
	return m.recorder
}

// GetByUUID mocks base method
func (m *MockReviewService) GetByUUID(ctx context.Context, UUID string) (*domain.Review, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByUUID", ctx, UUID)
	ret0, _ := ret[0].(*domain.Review)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByUUID indicates an expected call of GetByUUID
func (mr *MockReviewServiceMockRecorder) GetByUUID(ctx, UUID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByUUID", reflect.TypeOf((*MockReviewService)(nil).GetByUUID), ctx, UUID)
}

// GetAll mocks base method
func (m *MockReviewService) GetAll(ctx context.Context, limit, offset int) ([]*domain.Review, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAll", ctx, limit, offset)
	ret0, _ := ret[0].([]*domain.Review)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAll indicates an expected call of GetAll
func (mr *MockReviewServiceMockRecorder) GetAll(ctx, limit, offset interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAll", reflect.TypeOf((*MockReviewService)(nil).GetAll), ctx, limit, offset)
}

// Delete mocks base method
func (m *MockReviewService) Delete(ctx context.Context, UUID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, UUID)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete
func (mr *MockReviewServiceMockRecorder) Delete(ctx, UUID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockReviewService)(nil).Delete), ctx, UUID)
}

// Create mocks base method
func (m *MockReviewService) Create(ctx context.Context, review *domain.CreateReviewDTO) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, review)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create
func (mr *MockReviewServiceMockRecorder) Create(ctx, review interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockReviewService)(nil).Create), ctx, review)
}

// Update mocks base method
func (m *MockReviewService) Update(ctx context.Context, review *domain.UpdateReviewDTO) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, review)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update
func (mr *MockReviewServiceMockRecorder) Update(ctx, review interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockReviewService)(nil).Update), ctx, review)
}

// MockTagService is a mock of TagService interface
type MockTagService struct {
	ctrl     *gomock.Controller
	recorder *MockTagServiceMockRecorder
}

// MockTagServiceMockRecorder is the mock recorder for MockTagService
type MockTagServiceMockRecorder struct {
	mock *MockTagService
}

// NewMockTagService creates a new mock instance
func NewMockTagService(ctrl *gomock.Controller) *MockTagService {
	mock := &MockTagService{ctrl: ctrl}
	mock.recorder = &MockTagServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockTagService) EXPECT() *MockTagServiceMockRecorder {
	return m.recorder
}

// GetByUUID mocks base method
func (m *MockTagService) GetByUUID(ctx context.Context, UUID string) (*domain.Tag, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByUUID", ctx, UUID)
	ret0, _ := ret[0].(*domain.Tag)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByUUID indicates an expected call of GetByUUID
func (mr *MockTagServiceMockRecorder) GetByUUID(ctx, UUID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByUUID", reflect.TypeOf((*MockTagService)(nil).GetByUUID), ctx, UUID)
}

// GetManyByUUIDs mocks base method
func (m *MockTagService) GetManyByUUIDs(ctx context.Context, UUIDs []string) ([]*domain.Tag, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetManyByUUIDs", ctx, UUIDs)
	ret0, _ := ret[0].([]*domain.Tag)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetManyByUUIDs indicates an expected call of GetManyByUUIDs
func (mr *MockTagServiceMockRecorder) GetManyByUUIDs(ctx, UUIDs interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetManyByUUIDs", reflect.TypeOf((*MockTagService)(nil).GetManyByUUIDs), ctx, UUIDs)
}

// GetAll mocks base method
func (m *MockTagService) GetAll(ctx context.Context, limit, offset int) ([]*domain.Tag, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAll", ctx, limit, offset)
	ret0, _ := ret[0].([]*domain.Tag)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAll indicates an expected call of GetAll
func (mr *MockTagServiceMockRecorder) GetAll(ctx, limit, offset interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAll", reflect.TypeOf((*MockTagService)(nil).GetAll), ctx, limit, offset)
}

// Delete mocks base method
func (m *MockTagService) Delete(ctx context.Context, UUID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, UUID)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete
func (mr *MockTagServiceMockRecorder) Delete(ctx, UUID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockTagService)(nil).Delete), ctx, UUID)
}

// Create mocks base method
func (m *MockTagService) Create(ctx context.Context, tag *domain.CreateTagDTO) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, tag)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create
func (mr *MockTagServiceMockRecorder) Create(ctx, tag interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockTagService)(nil).Create), ctx, tag)
}

// Update mocks base method
func (m *MockTagService) Update(ctx context.Context, tag *domain.UpdateTagDTO) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, tag)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update
func (mr *MockTagServiceMockRecorder) Update(ctx, tag interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockTagService)(nil).Update), ctx, tag)
}

// MockVideoService is a mock of VideoService interface
type MockVideoService struct {
	ctrl     *gomock.Controller
	recorder *MockVideoServiceMockRecorder
}

// MockVideoServiceMockRecorder is the mock recorder for MockVideoService
type MockVideoServiceMockRecorder struct {
	mock *MockVideoService
}

// NewMockVideoService creates a new mock instance
func NewMockVideoService(ctrl *gomock.Controller) *MockVideoService {
	mock := &MockVideoService{ctrl: ctrl}
	mock.recorder = &MockVideoServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockVideoService) EXPECT() *MockVideoServiceMockRecorder {
	return m.recorder
}

// GetByUUID mocks base method
func (m *MockVideoService) GetByUUID(ctx context.Context, UUID string) (*domain.Video, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByUUID", ctx, UUID)
	ret0, _ := ret[0].(*domain.Video)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByUUID indicates an expected call of GetByUUID
func (mr *MockVideoServiceMockRecorder) GetByUUID(ctx, UUID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByUUID", reflect.TypeOf((*MockVideoService)(nil).GetByUUID), ctx, UUID)
}

// GetAll mocks base method
func (m *MockVideoService) GetAll(ctx context.Context, limit, offset int) ([]*domain.Video, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAll", ctx, limit, offset)
	ret0, _ := ret[0].([]*domain.Video)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAll indicates an expected call of GetAll
func (mr *MockVideoServiceMockRecorder) GetAll(ctx, limit, offset interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAll", reflect.TypeOf((*MockVideoService)(nil).GetAll), ctx, limit, offset)
}

// Delete mocks base method
func (m *MockVideoService) Delete(ctx context.Context, UUID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, UUID)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete
func (mr *MockVideoServiceMockRecorder) Delete(ctx, UUID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockVideoService)(nil).Delete), ctx, UUID)
}

// Create mocks base method
func (m *MockVideoService) Create(ctx context.Context, video *domain.CreateVideoDTO) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, video)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create
func (mr *MockVideoServiceMockRecorder) Create(ctx, video interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockVideoService)(nil).Create), ctx, video)
}

// Update mocks base method
func (m *MockVideoService) Update(ctx context.Context, video *domain.UpdateVideoDTO) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, video)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update
func (mr *MockVideoServiceMockRecorder) Update(ctx, video interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockVideoService)(nil).Update), ctx, video)
}
