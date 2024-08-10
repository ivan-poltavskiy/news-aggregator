package news

import (
	"errors"
	"net/http"
	"net/http/httptest"
	client "news-aggregator/client/mock_aggregator"
	"news-aggregator/entity/news"
	storage "news-aggregator/storage/mock_aggregator"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestFetchNewsHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := client.NewMockClient(ctrl)
	mockStorage := storage.NewMockStorage(ctrl)
	handler := NewNewsHandler(mockStorage)

	tests := []struct {
		name           string
		mockFetchNews  func()
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "Success Fetch",
			mockFetchNews: func() {
				expectedNews := []news.News{
					{Title: "Sample News 1", Description: "Description 1", Link: "http://example.com/1"},
					{Title: "Sample News 2", Description: "Description 2", Link: "http://example.com/2"},
				}
				mockClient.EXPECT().FetchNews().Return(expectedNews, nil).Times(1)
				mockClient.EXPECT().Print(expectedNews).Times(1)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   "",
		},
		{
			name: "Fetch News Error",
			mockFetchNews: func() {
				mockClient.EXPECT().FetchNews().Return(nil, errors.New("some error")).Times(1)
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "some error\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFetchNews()

			httptest.NewRequest(http.MethodGet, "/news", nil)
			rec := httptest.NewRecorder()

			handler.FetchNewsHandler(rec, mockClient)

			assert.Equal(t, tt.expectedStatus, rec.Code)
			assert.Equal(t, tt.expectedBody, rec.Body.String())
		})
	}
}
