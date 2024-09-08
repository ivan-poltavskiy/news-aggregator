package controller

import (
	"bytes"
	aggregatorv1 "com.teamdev/news-aggregator/api/v1"
	controller "com.teamdev/news-aggregator/internal/controller/mock_aggregator"
	"context"
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"io"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"net/http"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/client/interceptor"
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

var _ = Describe("Negative FeedReconcile tests", func() {
	var reconciler FeedReconciler
	var httpClient *controller.MockHttpClient
	var fakeClient client.Client

	BeforeEach(func() {
		t := GinkgoT()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		httpClient = controller.NewMockHttpClient(ctrl)
		fakeClient = fake.NewClientBuilder().WithScheme(scheme.Scheme).WithStatusSubresource(&aggregatorv1.Feed{}).Build()

		reconciler = FeedReconciler{
			Client: fakeClient,
			Scheme: scheme.Scheme,
			HttpsLinks: HttpsClientData{
				ServerUrl:                 "http://newsaggregator.com/manage",
				EndpointForSourceManaging: "/source",
			},
			HttpClient: httpClient,
			Finalizer:  "feed.finalizers.news.teamdev.com",
		}
	})

	AfterEach(func() {})

	ctx := context.Background()

	It("Feed resource is not found", func() {
		namespacedName := types.NamespacedName{Namespace: "default", Name: "non-existent-feed"}
		_, err := reconciler.Reconcile(ctx, ctrl.Request{NamespacedName: namespacedName})
		Expect(err).To(BeNil(), "Expected no error when Feed is not found")
	})

	It("Cannot update the status", func() {
		feed := aggregatorv1.Feed{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-feed",
				Namespace: "default",
			},
			Spec: aggregatorv1.FeedSpec{
				Name: "test-feed",
				Url:  "http://example.com",
			},
			Status: aggregatorv1.FeedStatus{},
		}
		fakeClient.Create(ctx, &feed)
		reconciler.Client = fake.NewClientBuilder().WithScheme(scheme.Scheme).WithInterceptorFuncs(
			interceptor.Funcs{Update: func(ctx context.Context, client client.WithWatch, obj client.Object, opts ...client.UpdateOption) error {
				return errors.New("new error")
			}}).WithStatusSubresource(&aggregatorv1.Feed{}).Build()

		namespacedName := types.NamespacedName{Namespace: "default", Name: "non-existent-feed"}
		_, err := reconciler.Reconcile(ctx, ctrl.Request{NamespacedName: namespacedName})
		Expect(err)
	})

	It("Https server returned the error from POST request", func() {
		feed := aggregatorv1.Feed{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-feed",
				Namespace: "default",
			},
			Spec: aggregatorv1.FeedSpec{
				Name: "test-feed",
				Url:  "http://example.com",
			},
			Status: aggregatorv1.FeedStatus{},
		}
		fakeClient.Create(ctx, &feed)

		reconciler.Client = fakeClient

		expectedURL := "http://newsaggregator.com/manage/source"
		httpClient.EXPECT().Post(expectedURL, "application/json", gomock.Any()).Return(&http.Response{StatusCode: http.StatusBadRequest}, errors.New("new error"))

		namespacedName := types.NamespacedName{Namespace: "default", Name: "test-feed"}
		_, err := reconciler.Reconcile(ctx, ctrl.Request{NamespacedName: namespacedName})
		Expect(err).To(HaveOccurred(), "Expected error when updating feed finalizer fails")
		Expect(feed.Status.GetCurrentCondition().Success == false)
	})

	It("Failed to create the Feed resource in the system", func() {
		reconciler.Client = fake.NewClientBuilder().WithScheme(scheme.Scheme).WithInterceptorFuncs(
			interceptor.Funcs{Create: func(ctx context.Context, client client.WithWatch, obj client.Object, opts ...client.CreateOption) error {
				return errors.New("failed to create resource")
			}}).Build()

		feed := aggregatorv1.Feed{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "failed-create-feed",
				Namespace: "default",
			},
			Spec: aggregatorv1.FeedSpec{
				Name: "failed-create-feed",
				Url:  "http://example.com",
			},
		}
		fakeClient.Create(ctx, &feed)
		namespacedName := types.NamespacedName{Namespace: "default", Name: "failed-create-feed"}
		_, err := reconciler.Reconcile(ctx, ctrl.Request{NamespacedName: namespacedName})

		Expect(err)
	})

	It("Https server returns error on PUT request", func() {
		feed := aggregatorv1.Feed{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-feed-to-update",
				Namespace: "default",
			},
			Spec: aggregatorv1.FeedSpec{
				Name: "test-feed-to-update",
				Url:  "http://example.com",
			},
			Status: aggregatorv1.FeedStatus{
				Conditions: []aggregatorv1.Condition{
					{Type: aggregatorv1.ConditionAdded},
				},
			},
		}
		fakeClient.Create(ctx, &feed)

		reconciler.Client = fakeClient

		expectedURL := "http://newsaggregator.com/manage/source"
		httpClient.EXPECT().
			Do(gomock.Any()).
			DoAndReturn(func(req *http.Request) (*http.Response, error) {
				if req.URL.String() == expectedURL && req.Method == http.MethodPut {
					return &http.Response{StatusCode: http.StatusInternalServerError}, errors.New("PUT error")
				}
				return nil, errors.New("unexpected request")
			})

		namespacedName := types.NamespacedName{Namespace: "default", Name: "test-feed-to-update"}
		_, err := reconciler.Reconcile(ctx, ctrl.Request{NamespacedName: namespacedName})

		Expect(err)
		Expect(feed.Status.GetCurrentCondition().Success == false)
	})

})
