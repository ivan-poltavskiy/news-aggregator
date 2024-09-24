package v1

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"net/url"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
	"time"
)

var k8sClient client.Client

func (r *Feed) SetupWebhookWithManager(mgr ctrl.Manager) error {
	k8sClient = mgr.GetClient()
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// +kubebuilder:webhook:path=/validate-aggregator-com-teamdev-v1-feed,mutating=false,failurePolicy=fail,sideEffects=None,groups=aggregator.com.teamdev,resources=feeds,verbs=create;update;delete,versions=v1,name=vfeed.kb.io,admissionReviewVersions=v1

// ValidateCreate validates the input data at the time of Feed's creation
func (r *Feed) ValidateCreate() (admission.Warnings, error) {
	logrus.Info("validate create", "name", r.Name)

	return r.validateFeed()
}

// ValidateUpdate validates the input data at the time of Feed's update
func (r *Feed) ValidateUpdate(old runtime.Object) (admission.Warnings, error) {

	logrus.Info("validate update", "name", r.Name)

	return r.validateFeed()
}

// ValidateDelete validates the input data at the time of Feed's delete
func (r *Feed) ValidateDelete() (admission.Warnings, error) {
	logrus.Info("validate delete", "name", r.Name)
	return nil, nil
}

// validateFeed implements the common validation logic for both create and update operations.
func (r *Feed) validateFeed() (admission.Warnings, error) {
	var errorsList field.ErrorList
	specPath := field.NewPath("spec")

	// Validate name
	if r.Spec.Name == "" {
		errorsList = append(errorsList, field.Required(specPath.Child("name"), "name cannot be empty"))
	} else if len(r.Spec.Name) > 20 {
		errorsList = append(errorsList, field.Invalid(specPath.Child("name"), r.Spec.Name, "name must not exceed 20 characters"))
	}

	// Validate URL
	if r.Spec.Url == "" {
		errorsList = append(errorsList, field.Required(specPath.Child("url"), "URL cannot be empty"))
	} else if !isValidURL(r.Spec.Url) {
		errorsList = append(errorsList, field.Invalid(specPath.Child("url"), r.Spec.Url, "URL must be a valid URL"))
	}

	// Check name uniqueness
	if err := checkNameUniqueness(r); err != nil {
		errorsList = append(errorsList, field.Invalid(specPath.Child("name"), r.Spec.Name, err.Error()))
	}

	if len(errorsList) > 0 {
		return nil, errorsList.ToAggregate()
	}

	return nil, nil
}

// isValidURL checks if the given string is a valid URL.
func isValidURL(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}

// checkNameUniqueness queries the Kubernetes API to ensure that no other Feed with the same name exists in the same namespace.
func checkNameUniqueness(feed *Feed) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	feedList := &FeedList{}
	listOpts := client.ListOptions{Namespace: feed.Namespace}
	err := k8sClient.List(ctx, feedList, &listOpts)
	if err != nil {
		return fmt.Errorf("checkNameUniqueness: failed to list feeds: %v", err)
	}

	for _, existingFeed := range feedList.Items {
		if existingFeed.Spec.Name == feed.Spec.Name && existingFeed.Namespace == feed.Namespace && existingFeed.UID != feed.UID {
			return fmt.Errorf("checkNameUniqueness: a Feed with name '%s' already exists in namespace '%s'", feed.Spec.Name, feed.Namespace)
		}
	}
	return nil
}
