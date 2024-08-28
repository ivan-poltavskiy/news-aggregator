package handler_test

import (
	aggregatorv1 "com.teamdev/news-aggregator/api/v1"
	handler2 "com.teamdev/news-aggregator/internal/controller/handler"
	"context"
	"errors"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/client/interceptor"
)

var _ = ginkgo.Describe("HotNewsHandler", func() {
	var (
		ctx        context.Context
		fakeClient client.Client
		handler    *handler2.HotNewsHandler
		hotNews    aggregatorv1.HotNews
		feed       aggregatorv1.Feed
	)

	ginkgo.BeforeEach(func() {
		ctx = context.Background()
		fakeClient = fake.NewClientBuilder().WithScheme(scheme.Scheme).Build()
		handler = &handler2.HotNewsHandler{
			Client: fakeClient,
		}

	})

	ginkgo.It("should return reconcile requests for each matching HotNews", func() {
		hotNews = aggregatorv1.HotNews{
			TypeMeta: metav1.TypeMeta{
				Kind:       "HotNews",
				APIVersion: aggregatorv1.GroupVersion.String(),
			},
			ObjectMeta: metav1.ObjectMeta{
				Namespace: "test-namespace",
				Name:      "test-hotnews",
			},
			Spec: aggregatorv1.HotNewsSpec{
				FeedsName: []string{"test-feed"},
			},
		}
		feed = aggregatorv1.Feed{
			TypeMeta: metav1.TypeMeta{
				Kind:       "Feed",
				APIVersion: aggregatorv1.GroupVersion.String(),
			},
			ObjectMeta: metav1.ObjectMeta{
				Namespace: "test-namespace",
				Name:      "test-feed",
			},
			Spec: aggregatorv1.FeedSpec{
				Name: "test-feed",
			},
		}
		fakeClient.Create(ctx, &feed)
		fakeClient.Create(ctx, &hotNews)

		requests := handler.UpdateHotNews(ctx, &feed)

		gomega.Expect(requests).To(gomega.HaveLen(1))
		gomega.Expect(requests[0].NamespacedName).To(gomega.Equal(types.NamespacedName{
			Namespace: "test-namespace",
			Name:      "test-hotnews",
		}))
	})

	ginkgo.It("should return nil when the Client.List return error", func() {
		fakeClient = fake.NewClientBuilder().WithScheme(scheme.Scheme).WithInterceptorFuncs(
			interceptor.Funcs{List: func(ctx context.Context, client client.WithWatch, list client.ObjectList, opts ...client.ListOption) error {
				return errors.New("test-list-error")
			}}).Build()
		handler.UpdateHotNews(ctx, &feed)

		gomega.Expect(nil)
	})

})
