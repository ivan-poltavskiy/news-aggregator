package predicate

import (
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"testing"
)

var _ = ginkgo.Describe("CustomPredicate", func() {
	var pod *corev1.Pod

	ginkgo.BeforeEach(func() {
		pod = &corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{Namespace: "biz", Name: "baz"},
		}
	})

	ginkgo.Describe("Funcs", func() {
		var instance *CustomPredicate

		ginkgo.BeforeEach(func() {
			instance = &CustomPredicate{
				CreateFunc: func(e event.CreateEvent) bool {
					defer ginkgo.GinkgoRecover()
					ginkgo.Fail("Did not expect CreateFunc to be called.")
					return false
				},
				DeleteFunc: func(e event.DeleteEvent) bool {
					defer ginkgo.GinkgoRecover()
					ginkgo.Fail("Did not expect DeleteFunc to be called.")
					return false
				},
				UpdateFunc: func(e event.UpdateEvent) bool {
					defer ginkgo.GinkgoRecover()
					ginkgo.Fail("Did not expect UpdateFunc to be called.")
					return false
				},
				GenericFunc: func(e event.GenericEvent) bool {
					defer ginkgo.GinkgoRecover()
					ginkgo.Fail("Did not expect GenericFunc to be called.")
					return false
				},
			}
		})

		ginkgo.It("should call Create", func() {
			instance.CreateFunc = func(e event.CreateEvent) bool {
				defer ginkgo.GinkgoRecover()
				gomega.Expect(e.Object).To(gomega.Equal(pod))
				return true
			}
			evt := event.CreateEvent{
				Object: pod,
			}
			gomega.Expect(instance.Create(evt)).To(gomega.BeTrue())

			instance.CreateFunc = nil
			gomega.Expect(instance.Create(evt)).To(gomega.BeTrue())
		})

		ginkgo.It("should call Update", func() {
			newPod := pod.DeepCopy()
			newPod.Name = "baz2"
			newPod.Namespace = "biz2"

			instance.UpdateFunc = func(e event.UpdateEvent) bool {
				defer ginkgo.GinkgoRecover()
				gomega.Expect(e.ObjectOld).To(gomega.Equal(pod))
				gomega.Expect(e.ObjectNew).To(gomega.Equal(newPod))
				return true
			}
			evt := event.UpdateEvent{
				ObjectOld: pod,
				ObjectNew: newPod,
			}
			gomega.Expect(instance.Update(evt)).To(gomega.BeTrue())

			instance.UpdateFunc = nil
			gomega.Expect(instance.Update(evt)).To(gomega.BeTrue())
		})

		ginkgo.It("should call Delete", func() {
			instance.DeleteFunc = func(e event.DeleteEvent) bool {
				defer ginkgo.GinkgoRecover()
				gomega.Expect(e.Object).To(gomega.Equal(pod))
				return true
			}
			evt := event.DeleteEvent{
				Object: pod,
			}
			gomega.Expect(instance.Delete(evt)).To(gomega.BeTrue())

			instance.DeleteFunc = nil
			gomega.Expect(instance.Delete(evt)).To(gomega.BeTrue())
		})

		ginkgo.It("should call Generic", func() {
			instance.GenericFunc = func(e event.GenericEvent) bool {
				defer ginkgo.GinkgoRecover()
				gomega.Expect(e.Object).To(gomega.Equal(pod))
				return true
			}
			evt := event.GenericEvent{
				Object: pod,
			}
			gomega.Expect(instance.Generic(evt)).To(gomega.BeTrue())

			instance.GenericFunc = nil
			gomega.Expect(instance.Generic(evt)).To(gomega.BeTrue())
		})
	})
})

func TestPredicates(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "CustomPredicate Suite")
}
