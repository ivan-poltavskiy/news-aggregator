package controller

import (
	"bytes"
	aggregatorv1 "com.teamdev/news-aggregator/api/v1"
	"context"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"io"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"net/http"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

// FeedReconciler reconciles a Feed object
type FeedReconciler struct {
	Client     client.Client
	Scheme     *runtime.Scheme
	HttpClient http.Client
	Finalizer  string
}

// FeedCreateRequest contains the URL of the feed to save it
type FeedCreateRequest struct {
	Url string `json:"url"`
}

// DeleteRequest contains the name of the feed to delete it
type DeleteRequest struct {
	Name string `json:"name"`
}

// +kubebuilder:rbac:groups=aggregator.com.teamdev,resources=feeds,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=aggregator.com.teamdev,resources=feeds/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=aggregator.com.teamdev,resources=feeds/finalizers,verbs=update

// Reconcile attempts to bring the state for the Feed CRD from the desired state to the current state.
func (r *FeedReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	var feed aggregatorv1.Feed

	err := r.Client.Get(ctx, req.NamespacedName, &feed)
	if err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	if !containsString(feed.ObjectMeta.Finalizers, r.Finalizer) {
		feed.ObjectMeta.Finalizers = append(feed.ObjectMeta.Finalizers, r.Finalizer)
		if err := r.Client.Update(ctx, &feed); err != nil {
			return ctrl.Result{}, err
		}
	}

	if !feed.ObjectMeta.DeletionTimestamp.IsZero() {
		if containsString(feed.ObjectMeta.Finalizers, r.Finalizer) {
			if err := r.deleteFeed(&feed); err != nil {
				updateCondition(&feed, aggregatorv1.ConditionDeleted, false, "Failed to delete feed", err.Error())
				if err := r.Client.Status().Update(ctx, &feed); err != nil {
					return ctrl.Result{}, err
				}
				return ctrl.Result{}, err
			}

			feed.ObjectMeta.Finalizers = removeString(feed.ObjectMeta.Finalizers, r.Finalizer)
			if err := r.Client.Update(ctx, &feed); err != nil {
				return ctrl.Result{}, err
			}
		}
		return ctrl.Result{}, nil
	}

	if err := r.addFeed(feed); err != nil {
		updateCondition(&feed, aggregatorv1.ConditionAdded, false, "Failed to add feed", err.Error())
		if err := r.Client.Status().Update(ctx, &feed); err != nil {
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, err
	}

	updateCondition(&feed, aggregatorv1.ConditionAdded, true, "", "")
	if err := r.Client.Status().Update(ctx, &feed); err != nil {
		return ctrl.Result{}, err
	}

	logrus.Info("Status updated.")

	return ctrl.Result{}, nil
}

func (r *FeedReconciler) addFeed(feed aggregatorv1.Feed) error {
	feedCreateRequest := FeedCreateRequest{
		Url: feed.Spec.Url,
	}

	reqBody, err := json.Marshal(feedCreateRequest)
	if err != nil {
		logrus.Error("Failed to marshal source request: ", err)
		return err
	}

	resp, err := r.HttpClient.Post("https://news-aggregator-service.news-aggregator.svc.cluster.local:443/sources", "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		logrus.Error("Failed to make POST request: ", err)
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		logrus.Error("Failed to create source, status code: ", resp.StatusCode, " response: ", string(body))
		return err
	}
	return nil
}

func (r *FeedReconciler) deleteFeed(feed *aggregatorv1.Feed) error {
	deleteRequest := DeleteRequest{
		Name: feed.Spec.Name,
	}

	reqBody, err := json.Marshal(deleteRequest)
	if err != nil {
		logrus.Error("Failed to marshal delete request: ", err)
		return err
	}

	logrus.Infof("Feed for delete name: %s", feed.Spec.Name)

	req, err := http.NewRequest("DELETE", "https://news-aggregator-service.news-aggregator.svc.cluster.local:443/sources", bytes.NewBuffer(reqBody))
	if err != nil {
		logrus.Error("Failed to create DELETE request: ", err)
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := r.HttpClient.Do(req)
	if err != nil {
		logrus.Error("Failed to make DELETE request: ", err)
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		logrus.Error("Failed to delete source, status code: ", resp.StatusCode, " response: ", string(body))
		return err
	}

	logrus.Info("Feed finalized successfully.")
	return nil
}

func (r *FeedReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&aggregatorv1.Feed{}).
		WithEventFilter(predicate.Funcs{
			CreateFunc: func(e event.CreateEvent) bool {
				return true
			},
			DeleteFunc: func(e event.DeleteEvent) bool {
				return !e.DeleteStateUnknown
			},
			UpdateFunc: func(e event.UpdateEvent) bool {
				return e.ObjectNew.GetGeneration() != e.ObjectOld.GetGeneration()
			},
		}).
		Complete(r)
}

// containsString checks if a string exists in a slice.
func containsString(slice []string, str string) bool {
	for _, item := range slice {
		if item == str {
			return true
		}
	}
	return false
}

// removeString removes a string from a slice.
func removeString(slice []string, str string) []string {
	var newSlice []string
	for _, item := range slice {
		if item == str {
			continue
		}
		newSlice = append(newSlice, item)
	}
	return newSlice
}

// updateCondition updates or adds a condition in the feed's status
func updateCondition(feed *aggregatorv1.Feed, conditionType aggregatorv1.ConditionType, statusBool bool, reason, message string) {
	newCondition := aggregatorv1.Condition{
		Type:           conditionType,
		Status:         statusBool,
		Reason:         reason,
		Message:        message,
		LastUpdateTime: metav1.Now(),
	}
	for i, condition := range feed.Status.Conditions {
		if condition.Type == conditionType {
			feed.Status.Conditions[i] = newCondition
			return
		}
	}

	feed.Status.Conditions = append(feed.Status.Conditions, newCondition)
}
