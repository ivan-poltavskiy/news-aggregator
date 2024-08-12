package controller

import (
	"bytes"
	aggregatorv1 "com.teamdev/news-aggregator/api/v1"
	controller "com.teamdev/news-aggregator/internal/controller/mock_aggregator"
	"context"
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"io"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"net/http"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
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

func TestFeedReconciler_updateFeed(t *testing.T) {
	tests := []struct {
		name            string
		feed            aggregatorv1.Feed
		mockPutResponse *http.Response
		mockPutError    error
		expectedError   bool
	}{
		{
			name: "Success update request",
			feed: aggregatorv1.Feed{
				Spec: aggregatorv1.FeedSpec{
					Url:  "http://example.com/feed",
					Name: "test-feed",
				},
			},
			mockPutResponse: &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBufferString(`{}`)),
			},
			mockPutError:  nil,
			expectedError: false,
		},
		{
			name: "Failed update request with error",
			feed: aggregatorv1.Feed{
				Spec: aggregatorv1.FeedSpec{
					Url:  "http://example.com/feed",
					Name: "test-feed",
				},
			},
			mockPutResponse: nil,
			mockPutError:    errors.New("failed to make PUT request"),
			expectedError:   true,
		},
		{
			name: "Failed update request with non-200 status",
			feed: aggregatorv1.Feed{
				Spec: aggregatorv1.FeedSpec{
					Url:  "http://example.com/feed",
					Name: "test-feed",
				},
			},
			mockPutResponse: &http.Response{
				StatusCode: http.StatusInternalServerError,
				Body:       io.NopCloser(bytes.NewBufferString(`{}`)),
			},
			mockPutError:  nil,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockHttpClient := controller.NewMockHttpClient(ctrl)

			mockHttpClient.EXPECT().
				Do(gomock.Any()).
				Return(tt.mockPutResponse, tt.mockPutError).
				Times(1)

			reconciler := &FeedReconciler{
				HttpClient: mockHttpClient,
				HttpsLinks: HttpsClientData{
					EndpointForSourceManaging: "http://mock-server/update-feed",
				},
			}

			err := reconciler.updateFeed(tt.feed)

			if (err != nil) != tt.expectedError {
				t.Errorf("updateFeed() error = %v, expectedError %v", err, tt.expectedError)
			}
		})
	}
}

func TestFeedReconcile(t *testing.T) {
	// Create a new scheme and add the necessary schemes
	scheme := runtime.NewScheme()
	_ = clientgoscheme.AddToScheme(scheme)
	_ = aggregatorv1.AddToScheme(scheme)

	// Create the initial Feed object
	initialFeed := &aggregatorv1.Feed{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-feed",
			Namespace: "default",
		},
		Spec: aggregatorv1.FeedSpec{
			Name: "Test Feed",
			Url:  "https://example.com/rss",
		},
		Status: aggregatorv1.FeedStatus{
			Conditions: []aggregatorv1.Condition{},
		},
	}
	client := fake.NewClientBuilder().WithScheme(scheme).WithObjects(initialFeed).Build()

	// Retrieve the initial Feed object
	feed := &aggregatorv1.Feed{}
	err := client.Get(context.Background(), types.NamespacedName{
		Name:      "test-feed",
		Namespace: "default",
	}, feed)
	assert.NoError(t, err, "initial Feed object should be found")

	// Create the mock HTTP client
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockHTTPClient := controller.NewMockHttpClient(ctrl)

	// Set up the mock POST request
	mockHTTPClient.EXPECT().
		Post("https://news-aggregator-service.news-aggregator.svc.cluster.local:443/sources", "application/json", gomock.Any()).
		Return(&http.Response{
			StatusCode: http.StatusCreated,
			Body:       io.NopCloser(bytes.NewBufferString("")),
		}, nil)

	// Create the reconciler with the mock client
	r := &FeedReconciler{
		Client:     client,
		Scheme:     scheme,
		HttpClient: mockHTTPClient,
		Finalizer:  "feed.finalizers.news.teamdev.com",
		HttpsLinks: HttpsClientData{
			ServerUrl:                 "https://news-aggregator-service.news-aggregator.svc.cluster.local:443",
			EndpointForSourceManaging: "/sources",
		},
	}

	// Create the reconcile request
	req := reconcile.Request{
		NamespacedName: types.NamespacedName{
			Name:      "test-feed",
			Namespace: "default",
		},
	}

	// Perform the reconciliation
	res, err := r.Reconcile(context.Background(), req)
	assert.False(t, res.Requeue)

	// Retrieve the Feed object after reconciliation
	err = client.Get(context.Background(), req.NamespacedName, feed)
	assert.NoError(t, err, "Feed object should be found after reconciliation")

	// Verify that the finalizer was added
	if !assert.Contains(t, feed.Finalizers, "feed.finalizers.news.teamdev.com", "Finalizer should be added") {
		t.Logf("Finalizers found: %v", feed.Finalizers)
	}
}
