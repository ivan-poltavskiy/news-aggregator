package controller

import (
	"bytes"
	aggregatorv1 "com.teamdev/news-aggregator/api/v1"
	controller "com.teamdev/news-aggregator/internal/controller/mock_aggregator"
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	"io"
	"net/http"
	"testing"
)

//go:generate mockgen -destination=mock_aggregator/mock_client.go -package=controller  sigs.k8s.io/controller-runtime/pkg/client Client
//go:generate mockgen -destination=mock_aggregator/mock_status_client.go -package=mocks sigs.k8s.io/controller-runtime/pkg/client StatusClient
func TestFeedReconciler_addFeed(t *testing.T) {
	tests := []struct {
		name             string
		feed             aggregatorv1.Feed
		mockPostResponse *http.Response
		mockPostError    error
		expectedError    bool
	}{
		{
			name: "Success request",
			feed: aggregatorv1.Feed{
				Spec: aggregatorv1.FeedSpec{
					Url: "http://example.com/feed",
				},
			},
			mockPostResponse: &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBufferString(`{}`)),
			},
			mockPostError: nil,
			expectedError: false,
		},
		{
			name: "Failed request with error",
			feed: aggregatorv1.Feed{
				Spec: aggregatorv1.FeedSpec{
					Url: "http://example.com/feed",
				},
			},
			mockPostResponse: nil,
			mockPostError:    errors.New("failed to make POST request"),
			expectedError:    true,
		},
		{
			name: "Failed request with non-200 status",
			feed: aggregatorv1.Feed{
				Spec: aggregatorv1.FeedSpec{
					Url: "http://example.com/feed",
				},
			},
			mockPostResponse: &http.Response{
				StatusCode: http.StatusInternalServerError,
				Body:       io.NopCloser(bytes.NewBufferString(`{}`)),
			},
			mockPostError: nil,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockHttpClient := controller.NewMockHttpClient(ctrl)

			mockHttpClient.EXPECT().
				Post(gomock.Any(), "application/json", gomock.Any()).
				Return(tt.mockPostResponse, tt.mockPostError)

			reconciler := &FeedReconciler{
				HttpClient: mockHttpClient,
				HttpsLinks: HttpsClientData{
					EndpointForSourceManaging: "http://mock-server/create-feed",
				},
			}

			err := reconciler.addFeed(tt.feed)

			if (err != nil) != tt.expectedError {
				t.Errorf("addFeed() error = %v, expectedError %v", err, tt.expectedError)
			}
		})
	}
}
func TestFeedReconciler_deleteFeed(t *testing.T) {
	tests := []struct {
		name              string
		feed              aggregatorv1.Feed
		mockDeleteRequest *http.Request
		mockDeleteError   error
		mockStatusCode    int
		expectedError     bool
	}{
		{
			name: "Success delete request",
			feed: aggregatorv1.Feed{
				Spec: aggregatorv1.FeedSpec{
					Name: "test-feed",
				},
			},
			mockDeleteRequest: &http.Request{},
			mockDeleteError:   nil,
			mockStatusCode:    http.StatusOK,
			expectedError:     false,
		},
		{
			name: "Failed delete request with error",
			feed: aggregatorv1.Feed{
				Spec: aggregatorv1.FeedSpec{
					Name: "test-feed",
				},
			},
			mockDeleteRequest: &http.Request{},
			mockDeleteError:   fmt.Errorf("delete request failed"),
			mockStatusCode:    http.StatusInternalServerError,
			expectedError:     true,
		},
		{
			name: "Failed delete request with non-OK status code",
			feed: aggregatorv1.Feed{
				Spec: aggregatorv1.FeedSpec{
					Name: "test-feed",
				},
			},
			mockDeleteRequest: &http.Request{},
			mockDeleteError:   nil,
			mockStatusCode:    http.StatusInternalServerError,
			expectedError:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockHttpClient := controller.NewMockHttpClient(ctrl)

			mockHttpClient.EXPECT().
				Do(gomock.Any()).
				Return(&http.Response{
					StatusCode: tt.mockStatusCode,
					Body:       io.NopCloser(bytes.NewBufferString(`{}`)),
				}, tt.mockDeleteError).
				Times(1)

			reconciler := &FeedReconciler{
				HttpClient: mockHttpClient,
				HttpsLinks: HttpsClientData{
					EndpointForSourceManaging: "http://mock-delete-url",
				},
			}

			err := reconciler.deleteFeed(&tt.feed)
			if (err != nil) != tt.expectedError {
				t.Fatalf("expected error: %v, got: %v", tt.expectedError, err)
			}
		})
	}
}
