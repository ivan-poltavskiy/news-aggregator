package controller

import (
	"bytes"
	aggregatorv1 "com.teamdev/news-aggregator/api/v1"
	controller "com.teamdev/news-aggregator/internal/controller/mock_aggregator"
	"context"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"io"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"
	"net/http"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"testing"
)

func TestReconcileHotNews(t *testing.T) {
	testScheme := runtime.NewScheme()
	_ = scheme.AddToScheme(testScheme)
	_ = aggregatorv1.AddToScheme(testScheme)
	configMap := &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-configmap",
			Namespace: "default",
		},
		Data: map[string]string{
			"group1": "source1,source2",
		},
	}

	fakeClient := fake.NewClientBuilder().
		WithScheme(testScheme).
		WithObjects(
			&aggregatorv1.HotNews{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-hotnews",
					Namespace: "default",
					Finalizers: []string{
						"test-finalizer",
					},
				},
				Spec: aggregatorv1.HotNewsSpec{
					Keywords:   []string{"keyword1"},
					DateStart:  "2024-01-01",
					DateEnd:    "2024-01-31",
					FeedsName:  []string{"feed1"},
					FeedGroups: []string{"group1"},
					SummaryConfig: aggregatorv1.SummaryConfig{
						TitlesCount: 5,
					},
				},
			},
			configMap,
		).Build()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockHttpClient := controller.NewMockHttpClient(ctrl)

	expectedURL := "?endDate=2024-01-31&keywords=keyword1&sources=source1%2Csource2&startDate=2024-01-01"
	mockHttpClient.EXPECT().
		Get(gomock.Eq(expectedURL)).
		Return(&http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewBufferString(`[{"Title": "News1"},{"Title": "News2"},{"Title": "News3"}]`)),
		}, nil).
		Times(1)

	r := &HotNewsReconciler{
		Client:        fakeClient,
		Scheme:        testScheme,
		HttpClient:    mockHttpClient,
		HttpsLinks:    HttpsClientData{},
		Finalizer:     "test-finalizer",
		ConfigMapMame: "test-configmap",
	}

	req := reconcile.Request{
		NamespacedName: client.ObjectKey{
			Name:      "test-hotnews",
			Namespace: "default",
		},
	}
	result, err := r.Reconcile(context.Background(), req)
	assert.Equal(t, reconcile.Result{}, result)

	updatedHotNews := &aggregatorv1.HotNews{}
	err = fakeClient.Get(context.TODO(), req.NamespacedName, updatedHotNews)
	assert.NoError(t, err)
}
