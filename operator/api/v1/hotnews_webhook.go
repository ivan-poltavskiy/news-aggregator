package v1

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
	"time"
)

// SetupWebhookWithManager will setup the manager to manage the webhooks
func (r *HotNews) SetupWebhookWithManager(mgr ctrl.Manager) error {
	k8sClient = mgr.GetClient()
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// +kubebuilder:webhook:path=/mutate-aggregator-com-teamdev-v1-hotnews,mutating=true,failurePolicy=fail,sideEffects=None,groups=aggregator.com.teamdev,resources=hotnews,verbs=create;update,versions=v1,name=mhotnews.kb.io,admissionReviewVersions=v1

var _ webhook.Defaulter = &HotNews{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *HotNews) Default() {
	if r.Spec.SummaryConfig.TitlesCount == 0 {
		r.Spec.SummaryConfig.TitlesCount = 10
	}

	if len(r.Spec.FeedsName) == 0 && len(r.Spec.FeedGroups) == 0 {
		feedList := &FeedList{}
		listOpts := client.ListOptions{Namespace: r.Namespace}
		err := k8sClient.List(context.Background(), feedList, &listOpts)
		if err != nil {
			logrus.Errorf("validateFeeds: failed to list feeds: %v", err)
		}
		var feedNameList []string
		for _, feed := range feedList.Items {
			feedNameList = append(feedNameList, feed.Name)
		}
		r.Spec.FeedsName = feedNameList
	}

	logrus.Info("default", "name", r.Name)

}

// +kubebuilder:webhook:path=/validate-aggregator-com-teamdev-v1-hotnews,mutating=false,failurePolicy=fail,sideEffects=None,groups=aggregator.com.teamdev,resources=hotnews,verbs=create;update,versions=v1,name=mhotnews.kb.io,admissionReviewVersions=v1

var _ webhook.Validator = &HotNews{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *HotNews) ValidateCreate() (admission.Warnings, error) {
	logrus.Info("validate create", "name", r.Name)
	return r.validateHotNews()
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *HotNews) ValidateUpdate(old runtime.Object) (admission.Warnings, error) {
	logrus.Info("validate update", "name", r.Name)

	return r.validateHotNews()
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *HotNews) ValidateDelete() (admission.Warnings, error) {
	logrus.Info("validate delete", "name", r.Name)

	return nil, nil
}

func (r *HotNews) validateHotNews() (admission.Warnings, error) {
	var errorsList field.ErrorList
	specPath := field.NewPath("spec")

	err := r.validateDate()

	if err != nil {
		errorsList = append(errorsList, field.Required(specPath.Child("date"), err.Error()))
	}

	if len(r.Spec.Keywords) == 0 {
		errorsList = append(errorsList, field.Required(specPath.Child("keywords"), "keywords is required"))
	}

	err = r.validateFeeds()
	if err != nil {
		errorsList = append(errorsList, field.Required(specPath.Child("FeedsName"), err.Error()))
	}

	logrus.Info("Error list lenght: ", len(errorsList))
	logrus.Info("Errors from error list: ", errorsList.ToAggregate())

	if len(errorsList) > 0 {
		return nil, errorsList.ToAggregate()
	}

	return nil, nil
}

func (r *HotNews) validateKeywords() error {
	if len(r.Spec.Keywords) == 0 {
		return fmt.Errorf("keywords is required")
	}
	return nil
}

func (r *HotNews) validateFeeds() error {
	feedList := &FeedList{}
	listOpts := client.ListOptions{Namespace: r.Namespace}
	err := k8sClient.List(context.Background(), feedList, &listOpts)
	if err != nil {
		return fmt.Errorf("validateFeeds: failed to list feeds: %v", err)
	}

	existingFeeds := make(map[string]bool)
	for _, feed := range feedList.Items {
		existingFeeds[feed.Name] = true
	}

	for _, feedName := range r.Spec.FeedsName {
		if !existingFeeds[feedName] {
			return fmt.Errorf("validateFeeds: feed %s does not exist in namespace %s", feedName, r.Namespace)
		}
	}

	return nil
}

func (r *HotNews) validateDate() error {

	if r.Spec.DateStart != "" || r.Spec.DateEnd != "" {
		if r.Spec.DateStart == "" || r.Spec.DateEnd == "" {
			return fmt.Errorf("both DateStart and DateEnd must be provided if one is specified")
		}

		dateStart, err := time.Parse("2006-01-02", r.Spec.DateStart)
		if err != nil {
			return fmt.Errorf("invalid DateStart format: must be yyyy-mm-dd")
		}
		dateEnd, err := time.Parse("2006-01-02", r.Spec.DateEnd)
		if err != nil {
			return fmt.Errorf("invalid DateEnd format: must be yyyy-mm-dd")
		}

		if !dateStart.Before(dateEnd) {
			return fmt.Errorf("DateStart must be before DateEnd")
		}
	}

	return nil
}
