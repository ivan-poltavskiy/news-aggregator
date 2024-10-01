package predicate

import (
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

var _ = ginkgo.Describe("HotNewsPredicate", func() {
	var pod *v1.Pod

	ginkgo.BeforeEach(func() {
		pod = &v1.Pod{
			ObjectMeta: metav1.ObjectMeta{Namespace: "test-namespace", Name: "test-name"},
		}
	})

	ginkgo.Describe("Check correct call to the funcs of predicate", func() {
		var instance predicate.Predicate

		ginkgo.BeforeEach(func() {
			instance = HotNewsPredicate()
		})

		ginkgo.It("should call Create", func() {
			evt := event.CreateEvent{
				Object: pod,
			}
			gomega.Expect(instance.Create(evt)).To(gomega.BeTrue())
		})

		ginkgo.It("should call Update when generation changes", func() {
			newPod := pod.DeepCopy()
			newPod.Generation = 2

			evt := event.UpdateEvent{
				ObjectOld: pod,
				ObjectNew: newPod,
			}
			gomega.Expect(instance.Update(evt)).To(gomega.BeTrue())
		})

		ginkgo.It("should not call Update when generation is the same", func() {
			newPod := pod.DeepCopy()

			evt := event.UpdateEvent{
				ObjectOld: pod,
				ObjectNew: newPod,
			}
			gomega.Expect(instance.Update(evt)).To(gomega.BeFalse())
		})

		ginkgo.It("should call Delete when DeleteStateUnknown is false", func() {
			evt := event.DeleteEvent{
				Object:             pod,
				DeleteStateUnknown: false,
			}
			gomega.Expect(instance.Delete(evt)).To(gomega.BeTrue())
		})

		ginkgo.It("should not call Delete when DeleteStateUnknown is true", func() {
			evt := event.DeleteEvent{
				Object:             pod,
				DeleteStateUnknown: true,
			}
			gomega.Expect(instance.Delete(evt)).To(gomega.BeFalse())
		})

		ginkgo.It("should call Generic", func() {
			evt := event.GenericEvent{
				Object: pod,
			}
			gomega.Expect(instance.Generic(evt)).To(gomega.BeTrue())
		})
	})
})

var _ = ginkgo.Describe("ConfigMapNamePredicate", func() {
	var configMap *v1.ConfigMap

	ginkgo.BeforeEach(func() {
		configMap = &v1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{Namespace: "test-namespace", Name: "test-configmap"},
		}
	})

	ginkgo.It("should return true for matching ConfigMap name", func() {
		instance := ConfigMapNamePredicate("test-configmap")
		gomega.Expect(instance.Create(event.CreateEvent{
			Object: configMap,
		})).To(gomega.BeTrue())
	})

	ginkgo.It("should return false for non-matching ConfigMap name", func() {
		instance := ConfigMapNamePredicate("different-configmap")
		gomega.Expect(instance.Create(event.CreateEvent{
			Object: configMap,
		})).To(gomega.BeFalse())
	})

	ginkgo.It("should return false if object is not a ConfigMap", func() {
		pod := &v1.Pod{
			ObjectMeta: metav1.ObjectMeta{Namespace: "test-namespace", Name: "test-pod"},
		}
		instance := ConfigMapNamePredicate("test-configmap")
		gomega.Expect(instance.Create(event.CreateEvent{
			Object: pod,
		})).To(gomega.BeFalse())
	})
})
