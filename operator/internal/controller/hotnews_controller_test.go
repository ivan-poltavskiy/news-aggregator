package controller

import (
	"bytes"
	aggregatorv1 "com.teamdev/news-aggregator/api/v1"
	controller "com.teamdev/news-aggregator/internal/controller/mock_aggregator"
	"context"
	"encoding/json"
	"errors"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"io"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"net/http"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/client/interceptor"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"time"
)

var _ = Describe("Negative reconcile tests for HotNewsReconciler", func() {

	var (
		configMapName = "feed-group-source"
		reconciler    HotNewsReconciler
		httpClient    *controller.MockHttpClient
		fakeClient    client.Client
	)

	BeforeEach(func() {
		t := GinkgoT()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		httpClient = controller.NewMockHttpClient(ctrl)
		fakeClient = fake.NewClientBuilder().WithScheme(scheme.Scheme).WithStatusSubresource(&aggregatorv1.HotNews{}).Build()

		reconciler = HotNewsReconciler{
			Client:     fakeClient,
			Scheme:     scheme.Scheme,
			HttpClient: httpClient,
			HttpsLinks: HttpsClientData{
				ServerUrl:                 "serverUrl",
				EndpointForSourceManaging: "endpointForGetNews",
			},
			Finalizer:     "feed.finalizers.news.teamdev.com",
			ConfigMapName: configMapName,
		}
	})

	AfterEach(func() {})

	ctx := context.Background()

	Describe("Error scenarios", func() {

		It("ConfigMap is not provided", func() {
			namespacedName := types.NamespacedName{Namespace: "default", Name: "hotnews"}
			_, err := reconciler.Reconcile(ctx, reconcile.Request{NamespacedName: namespacedName})
			gomega.Expect(err)
		})

		It("Hot News is not provided", func() {
			configMap := v1.ConfigMap{
				TypeMeta: metav1.TypeMeta{},
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "default",
					Name:      configMapName,
				},
				Immutable:  nil,
				Data:       nil,
				BinaryData: nil,
			}
			fakeClient.Create(ctx, &configMap)
			namespacedName := types.NamespacedName{Namespace: "default", Name: "hotnews"}
			_, err := reconciler.Reconcile(ctx, reconcile.Request{NamespacedName: namespacedName})
			gomega.Expect(err)
		})

		It("Https server returns error when trying to fetch news", func() {
			feed := aggregatorv1.Feed{
				TypeMeta: metav1.TypeMeta{},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test",
					Namespace: "default",
				},
				Spec: aggregatorv1.FeedSpec{
					Name: "test",
					Url:  "test.com",
				},
				Status: aggregatorv1.FeedStatus{},
			}
			configMap := v1.ConfigMap{
				TypeMeta: metav1.TypeMeta{
					Kind:       "ConfigMap",
					APIVersion: "v1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "default",
					Name:      configMapName,
				},
				Immutable:  nil,
				Data:       nil,
				BinaryData: nil,
			}
			hotNews := aggregatorv1.HotNews{
				TypeMeta: metav1.TypeMeta{
					Kind:       "HotNews",
					APIVersion: "v1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "hotnews",
					Namespace: "default",
				},
				Spec: aggregatorv1.HotNewsSpec{
					Keywords:      []string{"test"},
					DateStart:     "",
					DateEnd:       "",
					FeedsName:     []string{"test"},
					FeedGroups:    nil,
					SummaryConfig: aggregatorv1.SummaryConfig{},
				},
				Status: aggregatorv1.HotNewsStatus{},
			}
			fakeClient.Create(ctx, &configMap)
			fakeClient.Create(ctx, &feed)
			fakeClient.Create(ctx, &hotNews)

			httpClient.EXPECT().Get("serverUrlendpointForGetNews?keywords=test&sources=test").Return(nil, errors.New("TestErr"))

			namespacedName := types.NamespacedName{Namespace: "default", Name: "hotnews"}
			_, err := reconciler.Reconcile(ctx, reconcile.Request{NamespacedName: namespacedName})
			gomega.Expect(err)
			conditions := hotNews.Status.GetCurrentCondition()
			gomega.Expect(!conditions.Success)
		})

		It("Https server returns non-200 status code when trying to fetch news", func() {
			feed := aggregatorv1.Feed{
				TypeMeta: metav1.TypeMeta{},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test",
					Namespace: "default",
				},
				Spec: aggregatorv1.FeedSpec{
					Name: "test",
					Url:  "test.com",
				},
				Status: aggregatorv1.FeedStatus{},
			}
			configMap := v1.ConfigMap{
				TypeMeta: metav1.TypeMeta{
					Kind:       "ConfigMap",
					APIVersion: "v1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "default",
					Name:      configMapName,
				},
				Immutable:  nil,
				Data:       nil,
				BinaryData: nil,
			}
			hotNews := aggregatorv1.HotNews{
				TypeMeta: metav1.TypeMeta{
					Kind:       "HotNews",
					APIVersion: "v1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "hotnews",
					Namespace: "default",
				},
				Spec: aggregatorv1.HotNewsSpec{
					Keywords:      []string{"test"},
					DateStart:     "",
					DateEnd:       "",
					FeedsName:     []string{"test"},
					FeedGroups:    nil,
					SummaryConfig: aggregatorv1.SummaryConfig{},
				},
				Status: aggregatorv1.HotNewsStatus{},
			}
			fakeClient.Create(ctx, &configMap)
			fakeClient.Create(ctx, &feed)
			fakeClient.Create(ctx, &hotNews)

			httpClient.EXPECT().Get("serverUrlendpointForGetNews?keywords=test&sources=test").Return(&http.Response{StatusCode: http.StatusBadGateway, Body: http.NoBody}, nil)

			namespacedName := types.NamespacedName{Namespace: "default", Name: "hotnews"}
			_, err := reconciler.Reconcile(ctx, reconcile.Request{NamespacedName: namespacedName})
			gomega.Expect(err)
			conditions := hotNews.Status.GetCurrentCondition()
			gomega.Expect(!conditions.Success)
		})

		It("Provided feeds are not present in the ConfigMap", func() {
			configMap := v1.ConfigMap{
				TypeMeta: metav1.TypeMeta{
					Kind:       "ConfigMap",
					APIVersion: "v1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "default",
					Name:      configMapName,
				},
				Immutable:  nil,
				Data:       nil,
				BinaryData: nil,
			}
			hotNews := aggregatorv1.HotNews{
				TypeMeta: metav1.TypeMeta{
					Kind:       "HotNews",
					APIVersion: "v1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "hotnews",
					Namespace: "default",
				},
				Spec: aggregatorv1.HotNewsSpec{
					Keywords:      []string{"test"},
					DateStart:     "",
					DateEnd:       "",
					FeedsName:     nil,
					FeedGroups:    []string{"test"},
					SummaryConfig: aggregatorv1.SummaryConfig{},
				},
				Status: aggregatorv1.HotNewsStatus{},
			}
			fakeClient.Create(ctx, &configMap)
			fakeClient.Create(ctx, &hotNews)

			namespacedName := types.NamespacedName{Namespace: "default", Name: "hotnews"}
			_, err := reconciler.Reconcile(ctx, reconcile.Request{NamespacedName: namespacedName})
			gomega.Expect(err)
			conditions := hotNews.Status.GetCurrentCondition()
			gomega.Expect(!conditions.Success)
		})
	})

	Describe("Positive scenario", func() {

		It("Successfully fetches news from the server", func() {
			feed := aggregatorv1.Feed{
				TypeMeta: metav1.TypeMeta{
					Kind:       "Feed",
					APIVersion: "v1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test",
					Namespace: "default",
				},
				Spec: aggregatorv1.FeedSpec{
					Name: "test",
					Url:  "test.com",
				},
				Status: aggregatorv1.FeedStatus{},
			}

			configMap := v1.ConfigMap{
				TypeMeta: metav1.TypeMeta{
					Kind:       "ConfigMap",
					APIVersion: "v1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "default",
					Name:      configMapName,
				},
				Immutable:  nil,
				Data:       map[string]string{"test": "test"},
				BinaryData: nil,
			}

			hotNews := aggregatorv1.HotNews{
				TypeMeta: metav1.TypeMeta{
					Kind:       "HotNews",
					APIVersion: "v1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "hotnews",
					Namespace: "default",
				},
				Spec: aggregatorv1.HotNewsSpec{
					Keywords:      []string{"test"},
					DateStart:     "",
					DateEnd:       "",
					FeedsName:     []string{"test"},
					FeedGroups:    nil,
					SummaryConfig: aggregatorv1.SummaryConfig{},
				},
				Status: aggregatorv1.HotNewsStatus{},
			}

			fakeClient.Create(ctx, &configMap)
			fakeClient.Create(ctx, &feed)
			fakeClient.Create(ctx, &hotNews)

			returnedNews := []news{{
				Title:       "TestTile",
				Description: "test",
				Link:        "test",
				Date:        time.Time{},
				SourceName:  "test",
			}}
			jsonData, err := json.Marshal(returnedNews)
			readCloser := io.NopCloser(bytes.NewReader(jsonData))

			httpClient.EXPECT().Get("serverUrlendpointForGetNews?keywords=test&sources=test").Return(&http.Response{
				StatusCode: http.StatusOK,
				Body:       readCloser,
			}, nil)

			namespacedName := types.NamespacedName{Namespace: "default", Name: "hotnews"}
			_, err = reconciler.Reconcile(ctx, reconcile.Request{NamespacedName: namespacedName})

			gomega.Expect(err).To(gomega.BeNil())
			gomega.Expect(!hotNews.Status.GetCurrentCondition().Success)
		})
	})
})

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
		gomega.Expect(err).To(gomega.BeNil(), "Expected no error when Feed is not found")
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
		gomega.Expect(err)
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
		gomega.Expect(err).To(gomega.HaveOccurred(), "Expected error when updating feed finalizer fails")
		gomega.Expect(feed.Status.GetCurrentCondition().Success == false)
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

		gomega.Expect(err)
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

		gomega.Expect(err)
		gomega.Expect(feed.Status.GetCurrentCondition().Success == false)
	})

})
