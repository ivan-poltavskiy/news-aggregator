package v1

import (
	"context"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/runtime"
	"net/url"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

var k8sClient client.Client

func (r *Feed) SetupWebhookWithManager(mgr ctrl.Manager) error {
	k8sClient = mgr.GetClient()
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// +kubebuilder:webhook:path=/mutate-aggregator-com-teamdev-v1-feed,mutating=true,failurePolicy=fail,sideEffects=None,groups=aggregator.com.teamdev,resources=feeds,verbs=create;update;delete,versions=v1,name=mfeed.kb.io,admissionReviewVersions=v1

var _ webhook.Defaulter = &Feed{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *Feed) Default() {
	logrus.Info("default", "name", r.Name)
}

// +kubebuilder:webhook:path=/validate-aggregator-com-teamdev-v1-feed,mutating=false,failurePolicy=fail,sideEffects=None,groups=aggregator.com.teamdev,resources=feeds,verbs=create;update;delete,versions=v1,name=vfeed.kb.io,admissionReviewVersions=v1

var _ webhook.Validator = &Feed{}

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

	// Validate name
	if r.Spec.Name == "" {
		return nil, errors.New("ValidateCreate: name cannot be empty")
	}
	if len(r.Spec.Name) > 20 {
		return nil, errors.New("ValidateCreate: name must not exceed 20 characters")
	}

	// Validate link
	if r.Spec.Url == "" {
		return nil, errors.New("ValidateCreate: link cannot be empty")
	}
	if !isValidURL(r.Spec.Url) {
		return nil, errors.New("ValidateCreate: link must be a valid URL")
	}

	err := checkNameUniqueness(r)
	if err != nil {
		return nil, errors.New("ValidateCreate: name is already taken")
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
	feedList := &FeedList{}
	listOpts := client.ListOptions{Namespace: feed.Namespace}
	err := k8sClient.List(context.Background(), feedList, &listOpts)
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
