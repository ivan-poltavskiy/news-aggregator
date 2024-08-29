package v1

import (
	"context"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// ConfigMapValidator defines a validator for ConfigMaps
type ConfigMapValidator struct {
	client.Client
}

// SetupWebhookWithManager sets up the webhook with the manager
func (v *ConfigMapValidator) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(&corev1.ConfigMap{}).
		WithValidator(v).
		Complete()
}

// ValidateCreate implements validation logic for ConfigMap creation
func (v *ConfigMapValidator) ValidateCreate(ctx context.Context, obj runtime.Object) (admission.Warnings, error) {
	configMap, ok := obj.(*corev1.ConfigMap)
	if !ok {
		return nil, fmt.Errorf("expected ConfigMap, got %T", obj)
	}

	var errorsList field.ErrorList
	specPath := field.NewPath("data")

	for key, value := range configMap.Data {
		if value == "" {
			errorsList = append(errorsList, field.Required(specPath.Child(key), fmt.Sprintf("value for key %s cannot be empty", key)))
		}
	}

	if len(errorsList) > 0 {
		return nil, errorsList.ToAggregate()
	}

	return nil, nil
}

// ValidateUpdate implements validation logic for ConfigMap updates
func (v *ConfigMapValidator) ValidateUpdate(ctx context.Context, oldObj, newObj runtime.Object) (admission.Warnings, error) {
	return v.ValidateCreate(ctx, newObj)
}

// ValidateDelete implements validation logic for ConfigMap deletion
func (v *ConfigMapValidator) ValidateDelete(ctx context.Context, obj runtime.Object) (admission.Warnings, error) {
	// Можно пропустить логику валидации на удаление
	return nil, nil
}
