package predicate

import (
	"github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

// HotNewsPredicate returns a predicate that defines the filtering logic for events
// related to HotNews objects
func HotNewsPredicate() predicate.Predicate {

	return predicate.Funcs{

		CreateFunc: func(e event.CreateEvent) bool {
			return true
		},
		DeleteFunc: func(e event.DeleteEvent) bool {
			return !e.DeleteStateUnknown
		},
		UpdateFunc: func(e event.UpdateEvent) bool {
			return e.ObjectNew.GetGeneration() != e.ObjectOld.GetGeneration()
		},
		GenericFunc: func(e event.GenericEvent) bool {
			return true
		},
	}
}

// ConfigMapNamePredicate check that the config map name is equals to provided name
func ConfigMapNamePredicate(name string) predicate.Predicate {
	logrus.Info("Starting ConfigMapNamePredicate with name of config map: " + name)
	return predicate.NewPredicateFuncs(func(obj client.Object) bool {
		configMap, ok := obj.(*v1.ConfigMap)
		if !ok {
			return false
		}
		return configMap.Name == name
	})
}
