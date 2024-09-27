package v1

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
	"strings"
)

// +kubebuilder:webhook:path=/validate--v1-configmap,mutating=false,failurePolicy=fail,sideEffects=None,groups="",resources=configmaps,verbs=create;update,versions=v1,name=vconfigmap.kb.io,admissionReviewVersions=v1

// ConfigMapValidator is a webhook validator for ConfigMap resources.
type ConfigMapValidator struct {
	client.Client
	ConfigMapName      string
	ConfigMapNamespace string
}

// SetupWebhookWithManager sets up the webhook with the manager.
func (v *ConfigMapValidator) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(&corev1.ConfigMap{}).
		WithValidator(v).
		Complete()
}

// ValidateCreate implements validation logic for ConfigMap creation.
func (v *ConfigMapValidator) ValidateCreate(ctx context.Context, obj runtime.Object) (admission.Warnings, error) {
	configMap, ok := obj.(*corev1.ConfigMap)
	if !ok {
		return nil, fmt.Errorf("not a ConfigMap")
	}

	var errorsList field.ErrorList
	specPath := field.NewPath("data")

	for key, value := range configMap.Data {
		if value == "" {
			errorsList = append(errorsList, field.Required(specPath.Child(key), "feed group defined but no feeds mapped to it"))
		}
	}

	existingFeeds, err := v.getFeedsFromContext(ctx, configMap)

	if err != nil {
		errorsList = append(errorsList, field.InternalError(specPath, fmt.Errorf("failed to retrieve feeds: %v", err)))
	}

	if len(errorsList) > 0 {
		return nil, errorsList.ToAggregate()
	}
	return v.checkFeedsExist(existingFeeds, configMap)
}

// ValidateUpdate implements validation logic for ConfigMap updates.
func (v *ConfigMapValidator) ValidateUpdate(ctx context.Context, oldObj, newObj runtime.Object) (admission.Warnings, error) {
	return v.ValidateCreate(ctx, newObj)
}

// ValidateDelete implements validation logic for ConfigMap deletion.
func (v *ConfigMapValidator) ValidateDelete(ctx context.Context, obj runtime.Object) (admission.Warnings, error) {
	return nil, nil
}

// checkFeedsExist checks whether the feeds specified in the ConfigMap exist in the namespace.
func (v *ConfigMapValidator) checkFeedsExist(existingFeeds map[string]struct{}, configMap *corev1.ConfigMap) (admission.Warnings, error) {

	var errorsList field.ErrorList
	specPath := field.NewPath("data")

	for key, value := range configMap.Data {
		logrus.Info("Feed in the ConfigMap: ", value)
		feeds := strings.Split(value, ",")
		for _, feedName := range feeds {
			feedName = strings.TrimSpace(feedName)
			if _, exists := existingFeeds[feedName]; !exists {
				errorsList = append(errorsList, field.Invalid(specPath.Child(key), feedName, fmt.Sprintf("feed %s does not exist in namespace %s", feedName, v.ConfigMapNamespace)))
			}
		}
	}

	if len(errorsList) > 0 {
		return nil, errorsList.ToAggregate()
	}

	return nil, nil
}

func (v *ConfigMapValidator) getFeedsFromContext(ctx context.Context, configMap *corev1.ConfigMap) (map[string]struct{}, error) {
	logrus.Info("Validate feeds in the ConfigMap: ", configMap.Name)
	logrus.Info("ConfigMap namespace: ", v.ConfigMapNamespace)

	feedList := &FeedList{}
	err := v.List(ctx, feedList, &client.ListOptions{Namespace: v.ConfigMapNamespace})
	if err != nil {
		return nil, fmt.Errorf("failed to list feeds: %v", err)
	}
	logrus.Info("Feeds in the namespace: ", feedList.Items)

	existingFeeds := make(map[string]struct{})
	for _, feed := range feedList.Items {
		existingFeeds[feed.Spec.Name] = struct{}{}
	}
	return existingFeeds, nil
}
