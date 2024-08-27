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
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"time"
)

var _ = Describe("Negative reconcile tests", func() {
	var configMapName = "feed-group-source"
	var reconciler HotNewsReconciler
	var httpClient *controller.MockHttpClient
	var fakeClient client.Client
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
			ConfigMapMame: configMapName,
		}
	})

	AfterEach(func() {})

	ctx := context.Background()
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

	It("Https server return error when it try to fetch news", func() {
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

	It("Https server return not 200 status code when it try to fetch news", func() {
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

	It("Provided feeds is not present in the ConfigMap", func() {
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

	It("Provided in the ConfigMap feeds is not present in the cluster ", func() {

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
			SourceName:  "test"},
		}
		jsonData, err := json.Marshal(returnedNews)
		readCloser := io.NopCloser(bytes.NewReader(jsonData))

		httpClient.EXPECT().Get("serverUrlendpointForGetNews?keywords=test&sources=test").Return(&http.Response{
			StatusCode: http.StatusOK,
			Body:       readCloser,
		}, nil)

		namespacedName := types.NamespacedName{Namespace: "default", Name: "hotnews"}
		_, err = reconciler.Reconcile(ctx, reconcile.Request{NamespacedName: namespacedName})

		gomega.Expect(err).To(gomega.BeNil())
	})

})
