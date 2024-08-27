package predicate

import (
	"sigs.k8s.io/controller-runtime/pkg/event"
)

type CustomPredicate struct {
	CreateFunc  func(e event.CreateEvent) bool
	DeleteFunc  func(e event.DeleteEvent) bool
	UpdateFunc  func(e event.UpdateEvent) bool
	GenericFunc func(e event.GenericEvent) bool
}

func (p *CustomPredicate) Create(e event.CreateEvent) bool {
	if p.CreateFunc != nil {
		return p.CreateFunc(e)
	}
	return true
}

func (p *CustomPredicate) Delete(e event.DeleteEvent) bool {
	if p.DeleteFunc != nil {
		return p.DeleteFunc(e)
	}
	return true
}

func (p *CustomPredicate) Update(e event.UpdateEvent) bool {
	if p.UpdateFunc != nil {
		return p.UpdateFunc(e)
	}
	return true
}

func (p *CustomPredicate) Generic(e event.GenericEvent) bool {
	if p.GenericFunc != nil {
		return p.GenericFunc(e)
	}
	return true
}
