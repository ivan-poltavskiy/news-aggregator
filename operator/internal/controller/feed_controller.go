package controller

import (
	"bytes"
	aggregatorv1 "com.teamdev/news-aggregator/api/v1"
	"context"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"net/http"
	"net/url"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

// HttpClient defines methods that an HTTP client should implement
//
//go:generate mockgen -source=feed_controller.go -destination=mock_aggregator/mock_http_client.go -package=controller  news-aggregator/operator/internal/controller HttpClient
type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
	Post(url, contentType string, body io.Reader) (*http.Response, error)
	Get(url string) (resp *http.Response, err error)
}

// FeedReconciler reconciles a Feed object
type FeedReconciler struct {
	Client     client.Client
	Scheme     *runtime.Scheme
	HttpClient HttpClient
	Finalizer  string
	HttpsLinks HttpsClientData
}

// HttpsClientData contains information for connecting and working with the http client of news aggregator
type HttpsClientData struct {
	ServerUrl                 string
	EndpointForSourceManaging string
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
			logrus.Info("Reconcile: Feed was not found. Error ignored")
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	if !containsFinalizer(feed.ObjectMeta.Finalizers, r.Finalizer) {
		feed.ObjectMeta.Finalizers = append(feed.ObjectMeta.Finalizers, r.Finalizer)
		if err := r.Client.Update(ctx, &feed); err != nil {
			return ctrl.Result{}, err
		}
	}

	if !feed.ObjectMeta.DeletionTimestamp.IsZero() && containsFinalizer(feed.ObjectMeta.Finalizers, r.Finalizer) {

		if err := r.deleteFeed(&feed.Spec.Name); err != nil {

			feed.Status.SetCondition(aggregatorv1.Condition{
				Type:    aggregatorv1.ConditionDeleted,
				Success: false,
				Message: "Reconcile: Failed to delete feed",
				Reason:  err.Error(),
			})

			if err := r.Client.Status().Update(ctx, &feed); err != nil {
				return ctrl.Result{}, err
			}
			return ctrl.Result{}, err
		}

		feed.ObjectMeta.Finalizers = removeFinalizer(feed.ObjectMeta.Finalizers, r.Finalizer)
		if err := r.Client.Update(ctx, &feed); err != nil {
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, nil
	}

	if feed.Status.GetCurrentCondition().Type != aggregatorv1.ConditionAdded {
		if err := r.addFeed(feed); err != nil {

			feed.Status.SetCondition(aggregatorv1.Condition{
				Type:    aggregatorv1.ConditionAdded,
				Success: false,
				Message: "Reconcile: Failed to add feed",
				Reason:  err.Error(),
			})

			if err := r.Client.Status().Update(ctx, &feed); err != nil {
				return ctrl.Result{}, err
			}
			return ctrl.Result{}, err
		}
	} else {
		if err := r.updateFeed(feed); err != nil {
			feed.Status.SetCondition(aggregatorv1.Condition{
				Type:    aggregatorv1.ConditionAdded,
				Success: false,
				Message: "Reconcile: Failed to update feed",
				Reason:  err.Error(),
			})
			if err := r.Client.Status().Update(ctx, &feed); err != nil {
				return ctrl.Result{}, err
			}
			return ctrl.Result{}, err
		}
	}

	aggregatorv1.AddPositiveCondition(&feed)

	if err := r.Client.Status().Update(ctx, &feed); err != nil {
		return ctrl.Result{}, err
	}

	logrus.Info("UpdateCondition: ", feed.Status.Conditions)

	logrus.Info("Success updated. Feed NewName and Feed Link: ", feed.Spec.Name, feed.Spec.Url)

	return ctrl.Result{}, nil
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

// addFeed call the news aggregator server for adding source to the storage
func (r *FeedReconciler) addFeed(feed aggregatorv1.Feed) error {
	feedCreateRequest := feedCreateRequest{
		Url:  feed.Spec.Url,
		Name: feed.Spec.Name,
	}

	reqBody, err := json.Marshal(feedCreateRequest)
	if err != nil {
		logrus.Error("Failed to marshal source request: ", err)
		return err
	}
	path, err := url.JoinPath(r.HttpsLinks.ServerUrl, r.HttpsLinks.EndpointForSourceManaging)
	if err != nil {
		return err
	}
	resp, err := r.HttpClient.Post(path, "application/json", bytes.NewBuffer(reqBody))
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
		return fmt.Errorf("failed to create source, status code: %d", resp.StatusCode)
	}
	return nil
}

// deleteFeed call the news aggregator server for delete source from the storage
func (r *FeedReconciler) deleteFeed(feedName *string) error {
	deleteRequest := feedDeleteRequest{
		Name: *feedName,
	}

	reqBody, err := json.Marshal(deleteRequest)
	if err != nil {
		logrus.Error("Failed to marshal delete request: ", err)
		return err
	}

	logrus.Infof("Feed for delete name: %s", *feedName)
	path, err := url.JoinPath(r.HttpsLinks.ServerUrl, r.HttpsLinks.EndpointForSourceManaging)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("DELETE", path, bytes.NewBuffer(reqBody))
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
		return fmt.Errorf("failed to delete source, status code: %d", resp.StatusCode)
	}

	logrus.Info("Feed removes successfully.")
	return nil
}

func (r *FeedReconciler) updateFeed(feed aggregatorv1.Feed) error {
	feedUpdateRequest := feedUpdateRequest{
		NewName: feed.Spec.Name,
		OldName: feed.Status.GetCurrentCondition().LastUpdatedName,
		Url:     feed.Spec.Url,
	}

	reqBody, err := json.Marshal(feedUpdateRequest)
	if err != nil {
		logrus.Error("Failed to marshal source request: ", err)
		return err
	}
	path, err := url.JoinPath(r.HttpsLinks.ServerUrl, r.HttpsLinks.EndpointForSourceManaging)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("PUT", path, bytes.NewBuffer(reqBody))
	if err != nil {
		logrus.Error("Failed to create PUT request: ", err)
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := r.HttpClient.Do(req)
	if err != nil {
		logrus.Error("Failed to make PUT request: ", err)
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
		logrus.Error("Failed to update source, status code: ", resp.StatusCode, " response: ", string(body))
		return fmt.Errorf("failed to update source, status code: %d", resp.StatusCode)
	}
	return nil
}

// feedCreateRequest contains the URL of the feed to save it
type feedCreateRequest struct {
	Url  string `json:"url"`
	Name string `json:"name"`
}

// feedUpdateRequest contains the data of the feed to update it
type feedUpdateRequest struct {
	OldName string `json:"old_name"`
	NewName string `json:"new_name"`
	Url     string `json:"url"`
}

// feedDeleteRequest contains the name of the feed to delete it
type feedDeleteRequest struct {
	Name string `json:"name"`
}

// containsFinalizer checks if a string exists in a slice.
func containsFinalizer(slice []string, str string) bool {
	for _, item := range slice {
		if item == str {
			return true
		}
	}
	return false
}

// removeFinalizer removes a string from a slice.
func removeFinalizer(slice []string, str string) []string {
	var newSlice []string
	for _, item := range slice {
		if item == str {
			continue
		}
		newSlice = append(newSlice, item)
	}
	return newSlice
}
