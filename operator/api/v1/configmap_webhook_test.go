package v1

import (
	"context"
	"fmt"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/client/interceptor"
)

var _ = Describe("Tests for ConfigMap Webhook", func() {
	var (
		ctx        context.Context
		configMap  *v1.ConfigMap
		webhook    *ConfigMapValidator
		fakeClient client.Client
	)
	BeforeEach(func() {
		ctx = context.Background()
		_ = AddToScheme(scheme.Scheme)

		testFeed1 := &Feed{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "testFeed1",
				Namespace: "operator-system",
			},
			Spec: FeedSpec{
				Name: "testFeed1",
			},
		}
		testFeed2 := &Feed{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "testFeed2",
				Namespace: "operator-system",
			},
			Spec: FeedSpec{
				Name: "testFeed2",
			},
		}

		fakeClient = fake.NewClientBuilder().WithScheme(scheme.Scheme).Build()

		webhook = &ConfigMapValidator{
			Client:             fakeClient,
			ConfigMapNamespace: "operator-system",
			ConfigMapName:      "configmap",
		}

		configMap = &v1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "feed-groups",
				Namespace: "operator-system",
			},
			Data: map[string]string{
				"testCategory1": "testFeed1",
				"testCategory2": "testFeed2",
				"testCategory3": "testFeed1, testFeed2",
			},
		}
		Expect(fakeClient.Create(ctx, testFeed1)).Should(Succeed())
		Expect(fakeClient.Create(ctx, testFeed2)).Should(Succeed())
	})

	Context("When webhook validate creating of updating config map ", func() {
		It("should return error when not config map provided", func() {
			configMap.Data[""] = ""
			_, err := webhook.ValidateCreate(ctx, &v1.Pod{})
			Expect(err)
		})
		It("should return error when data in the config map is empty", func() {
			configMap.Data[""] = ""
			_, err := webhook.ValidateCreate(ctx, configMap)
			Expect(err)
		})
	})
	Context("When webhook validates feeds in the config map", func() {
		It("should return error when invalid feeds are provided", func() {
			// Adding a fake feed not present in the namespace
			configMap.Data["testCategory1"] = "fakeFeed"

			// Get the list of existing feeds from context
			existingFeeds, err := webhook.getFeedsFromContext(ctx, configMap)
			Expect(err).ToNot(HaveOccurred())

			// Check if feeds exist
			_, err = webhook.checkFeedsExist(existingFeeds, configMap)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("feed fakeFeed does not exist in namespace"))
		})

		It("should return an error if listing feeds fails", func() {
			fakeClient = fake.NewClientBuilder().WithScheme(scheme.Scheme).WithInterceptorFuncs(interceptor.Funcs{List: func(ctx context.Context, client client.WithWatch, list client.ObjectList, opts ...client.ListOption) error {
				return fmt.Errorf("fake list error")
			}}).Build()
			webhook = &ConfigMapValidator{
				Client:             fakeClient,
				ConfigMapNamespace: "non-existent-namespace",
				ConfigMapName:      "configmap",
			}

			_, err := webhook.getFeedsFromContext(ctx, configMap)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("failed to list feeds"))
		})

		It("should succeed when valid feeds are provided", func() {
			// Get the list of existing feeds from context
			existingFeeds, err := webhook.getFeedsFromContext(ctx, configMap)
			Expect(err).ToNot(HaveOccurred())

			// Check if feeds exist
			_, err = webhook.checkFeedsExist(existingFeeds, configMap)
			Expect(err).ToNot(HaveOccurred())
		})
	})
})
