package v1

import (
	"context"
	"errors"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/runtime"
	"net/url"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

var feedlog = logf.Log.WithName("feed-resource")

func (r *Feed) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// +kubebuilder:webhook:path=/mutate-aggregator-com-teamdev-v1-feed,mutating=true,failurePolicy=fail,sideEffects=None,groups=aggregator.com.teamdev,resources=feeds,verbs=create;update;delete,versions=v1,name=mfeed.kb.io,admissionReviewVersions=v1

var _ webhook.Defaulter = &Feed{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *Feed) Default() {
	feedlog.Info("default", "name", r.Name)

	// TODO(user): fill in your defaulting logic.
}

// +kubebuilder:webhook:path=/validate-aggregator-com-teamdev-v1-feed,mutating=false,failurePolicy=fail,sideEffects=None,groups=aggregator.com.teamdev,resources=feeds,verbs=create;update;delete,versions=v1,name=vfeed.kb.io,admissionReviewVersions=v1

var _ webhook.Validator = &Feed{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *Feed) ValidateCreate() (admission.Warnings, error) {
	feedlog.Info("validate create", "name", r.Name)
	return r.validateFeed()
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *Feed) ValidateUpdate(old runtime.Object) (admission.Warnings, error) {
	feedlog.Info("validate update", "name", r.Name)
	return r.validateFeed()
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *Feed) ValidateDelete() (admission.Warnings, error) {
	feedlog.Info("validate delete", "name", r.Name)
	return nil, nil
}

// validateFeed implements the common validation logic for both create and update operations.
func (r *Feed) validateFeed() (admission.Warnings, error) {
	// Validate name
	if r.Spec.Name == "" {
		return nil, errors.New("name cannot be empty")
	}
	if len(r.Spec.Name) > 20 {
		return nil, errors.New("name must not exceed 20 characters")
	}

	// Validate link
	if r.Spec.Url == "" {
		return nil, errors.New("link cannot be empty")
	}
	if !isValidURL(r.Spec.Url) {
		return nil, errors.New("link must be a valid URL")
	}
	uniqueInNamespace, err := isUniqueInNamespace(r)
	if err != nil {
		return nil, err
	}
	if !uniqueInNamespace {
		return nil, errors.New("name and link must be unique in the namespace")
	}

	return nil, nil
}

// isUniqueInNamespace checks if the given value is unique in the namespace for a specified field.
func isUniqueInNamespace(feed *Feed) (bool, error) {

	ctx := context.Background()
	config := ctrl.GetConfigOrDie()
	scheme := runtime.NewScheme()
	if err := AddToScheme(scheme); err != nil {
		logrus.Error(err, "failed to add scheme")
		return false, err
	}
	cl, err := client.New(config, client.Options{Scheme: scheme})
	if err != nil {
		logrus.Error(err, "failed to create client")
		return false, err
	}

	var feedList FeedList
	logrus.Info("Listing sources in namespace", "namespace", feed.Namespace)
	if err := cl.List(ctx, &feedList, client.InNamespace(feed.Namespace)); err != nil {
		logrus.Error(err, "failed to list sources")
		return false, err
	}

	logrus.Info("FeedList: ", feedList.Items)

	for _, existingFeed := range feedList.Items {
		if existingFeed.Spec.Name == feed.Spec.Name || existingFeed.Spec.Url == feed.Spec.Url {
			return false, nil
		}
	}

	return true, nil
}

// isValidURL checks if the given string is a valid URL.
func isValidURL(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}
