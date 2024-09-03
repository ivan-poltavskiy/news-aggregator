package v1

import (
	"context"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
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
				Namespace: "default",
			},
		}
		testFeed2 := &Feed{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "testFeed2",
				Namespace: "default",
			},
		}

		fakeClient = fake.NewClientBuilder().WithScheme(scheme.Scheme).Build()

		webhook = &ConfigMapValidator{
			Client:             fakeClient,
			ConfigMapNamespace: "default",
			ConfigMapName:      "configmap",
		}

		configMap = &v1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "feed-groups",
				Namespace: "default",
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
	Context("When webhook validate that feeds in the config map is provided ", func() {

		It("should return error when invalid feeds provided in the config map", func() {
			configMap.Data["test"] = "fakeFeed"
			_, err := webhook.checkFeedsExist(ctx, configMap)
			Expect(err)
		})

		It("should work correctly when provided feeds are correct", func() {
			_, err := webhook.checkFeedsExist(ctx, configMap)
			Expect(err).ToNot(HaveOccurred())
		})

	})
})
