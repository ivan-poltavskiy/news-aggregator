// Code generated by MockGen. DO NOT EDIT.
// Source: storage.go

// Package mock_aggregator is a generated GoMock package.
package storage

import (
	news "news-aggregator/entity/news"
	source "news-aggregator/entity/source"
	storage "news-aggregator/storage"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockStorage is a mock of Storage interface.
type MockStorage struct {
	ctrl     *gomock.Controller
	recorder *MockStorageMockRecorder
}

// MockStorageMockRecorder is the mock recorder for MockStorage.
type MockStorageMockRecorder struct {
	mock *MockStorage
}

// NewMockStorage creates a new mock instance.
func NewMockStorage(ctrl *gomock.Controller) *MockStorage {
	mock := &MockStorage{ctrl: ctrl}
	mock.recorder = &MockStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStorage) EXPECT() *MockStorageMockRecorder {
	return m.recorder
}

// DeleteSourceByName mocks base method.
func (m *MockStorage) DeleteSourceByName(arg0 source.Name) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteSourceByName", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteSourceByName indicates an expected call of DeleteSourceByName.
func (mr *MockStorageMockRecorder) DeleteSourceByName(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteSourceByName", reflect.TypeOf((*MockStorage)(nil).DeleteSourceByName), arg0)
}

// GetNews mocks base method.
func (m *MockStorage) GetNews(path string) ([]news.News, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetNews", path)
	ret0, _ := ret[0].([]news.News)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetNews indicates an expected call of GetNews.
func (mr *MockStorageMockRecorder) GetNews(path interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetNews", reflect.TypeOf((*MockStorage)(nil).GetNews), path)
}

// GetNewsBySourceName mocks base method.
func (m *MockStorage) GetNewsBySourceName(sourceName source.Name, sourceStorage storage.Source) ([]news.News, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetNewsBySourceName", sourceName, sourceStorage)
	ret0, _ := ret[0].([]news.News)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetNewsBySourceName indicates an expected call of GetNewsBySourceName.
func (mr *MockStorageMockRecorder) GetNewsBySourceName(sourceName, sourceStorage interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetNewsBySourceName", reflect.TypeOf((*MockStorage)(nil).GetNewsBySourceName), sourceName, sourceStorage)
}

// GetSourceByName mocks base method.
func (m *MockStorage) GetSourceByName(arg0 source.Name) (source.Source, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSourceByName", arg0)
	ret0, _ := ret[0].(source.Source)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSourceByName indicates an expected call of GetSourceByName.
func (mr *MockStorageMockRecorder) GetSourceByName(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSourceByName", reflect.TypeOf((*MockStorage)(nil).GetSourceByName), arg0)
}

// GetSources mocks base method.
func (m *MockStorage) GetSources() ([]source.Source, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSources")
	ret0, _ := ret[0].([]source.Source)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSources indicates an expected call of GetSources.
func (mr *MockStorageMockRecorder) GetSources() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSources", reflect.TypeOf((*MockStorage)(nil).GetSources))
}

// IsSourceExists mocks base method.
func (m *MockStorage) IsSourceExists(arg0 source.Name) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsSourceExists", arg0)
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsSourceExists indicates an expected call of IsSourceExists.
func (mr *MockStorageMockRecorder) IsSourceExists(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsSourceExists", reflect.TypeOf((*MockStorage)(nil).IsSourceExists), arg0)
}

// SaveNews mocks base method.
func (m *MockStorage) SaveNews(providedSource source.Source, news []news.News) (source.Source, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveNews", providedSource, news)
	ret0, _ := ret[0].(source.Source)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SaveNews indicates an expected call of SaveNews.
func (mr *MockStorageMockRecorder) SaveNews(providedSource, news interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveNews", reflect.TypeOf((*MockStorage)(nil).SaveNews), providedSource, news)
}

// SaveSource mocks base method.
func (m *MockStorage) SaveSource(source source.Source) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveSource", source)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveSource indicates an expected call of SaveSource.
func (mr *MockStorageMockRecorder) SaveSource(source interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveSource", reflect.TypeOf((*MockStorage)(nil).SaveSource), source)
}

// MockNews is a mock of News interface.
type MockNews struct {
	ctrl     *gomock.Controller
	recorder *MockNewsMockRecorder
}

// MockNewsMockRecorder is the mock recorder for MockNews.
type MockNewsMockRecorder struct {
	mock *MockNews
}

// NewMockNews creates a new mock instance.
func NewMockNews(ctrl *gomock.Controller) *MockNews {
	mock := &MockNews{ctrl: ctrl}
	mock.recorder = &MockNewsMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockNews) EXPECT() *MockNewsMockRecorder {
	return m.recorder
}

// GetNews mocks base method.
func (m *MockNews) GetNews(path string) ([]news.News, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetNews", path)
	ret0, _ := ret[0].([]news.News)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetNews indicates an expected call of GetNews.
func (mr *MockNewsMockRecorder) GetNews(path interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetNews", reflect.TypeOf((*MockNews)(nil).GetNews), path)
}

// GetNewsBySourceName mocks base method.
func (m *MockNews) GetNewsBySourceName(sourceName source.Name, sourceStorage storage.Source) ([]news.News, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetNewsBySourceName", sourceName, sourceStorage)
	ret0, _ := ret[0].([]news.News)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetNewsBySourceName indicates an expected call of GetNewsBySourceName.
func (mr *MockNewsMockRecorder) GetNewsBySourceName(sourceName, sourceStorage interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetNewsBySourceName", reflect.TypeOf((*MockNews)(nil).GetNewsBySourceName), sourceName, sourceStorage)
}

// SaveNews mocks base method.
func (m *MockNews) SaveNews(providedSource source.Source, news []news.News) (source.Source, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveNews", providedSource, news)
	ret0, _ := ret[0].(source.Source)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SaveNews indicates an expected call of SaveNews.
func (mr *MockNewsMockRecorder) SaveNews(providedSource, news interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveNews", reflect.TypeOf((*MockNews)(nil).SaveNews), providedSource, news)
}

// MockSource is a mock of Source interface.
type MockSource struct {
	ctrl     *gomock.Controller
	recorder *MockSourceMockRecorder
}

// MockSourceMockRecorder is the mock recorder for MockSource.
type MockSourceMockRecorder struct {
	mock *MockSource
}

// NewMockSource creates a new mock instance.
func NewMockSource(ctrl *gomock.Controller) *MockSource {
	mock := &MockSource{ctrl: ctrl}
	mock.recorder = &MockSourceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSource) EXPECT() *MockSourceMockRecorder {
	return m.recorder
}

// DeleteSourceByName mocks base method.
func (m *MockSource) DeleteSourceByName(arg0 source.Name) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteSourceByName", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteSourceByName indicates an expected call of DeleteSourceByName.
func (mr *MockSourceMockRecorder) DeleteSourceByName(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteSourceByName", reflect.TypeOf((*MockSource)(nil).DeleteSourceByName), arg0)
}

// GetSourceByName mocks base method.
func (m *MockSource) GetSourceByName(arg0 source.Name) (source.Source, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSourceByName", arg0)
	ret0, _ := ret[0].(source.Source)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSourceByName indicates an expected call of GetSourceByName.
func (mr *MockSourceMockRecorder) GetSourceByName(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSourceByName", reflect.TypeOf((*MockSource)(nil).GetSourceByName), arg0)
}

// GetSources mocks base method.
func (m *MockSource) GetSources() ([]source.Source, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSources")
	ret0, _ := ret[0].([]source.Source)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSources indicates an expected call of GetSources.
func (mr *MockSourceMockRecorder) GetSources() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSources", reflect.TypeOf((*MockSource)(nil).GetSources))
}

// IsSourceExists mocks base method.
func (m *MockSource) IsSourceExists(arg0 source.Name) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsSourceExists", arg0)
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsSourceExists indicates an expected call of IsSourceExists.
func (mr *MockSourceMockRecorder) IsSourceExists(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsSourceExists", reflect.TypeOf((*MockSource)(nil).IsSourceExists), arg0)
}

// SaveSource mocks base method.
func (m *MockSource) SaveSource(source source.Source) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveSource", source)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveSource indicates an expected call of SaveSource.
func (mr *MockSourceMockRecorder) SaveSource(source interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveSource", reflect.TypeOf((*MockSource)(nil).SaveSource), source)
}
