package controller

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"k8s.io/apimachinery/pkg/api/errors"
	"net/http"
	"os"
	"time"

	aggregatorv1 "com.teamdev/news-aggregator/api/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func getHttpClient() (*http.Client, error) {
	certPool := x509.NewCertPool()
	certs, err := os.ReadFile("certificates/server.crt")
	if err != nil {
		return nil, err
	}

	certPool.AppendCertsFromPEM(certs)

	keyPEMBlock, err := os.ReadFile("certificates/key.unencrypted.pem")
	if err != nil {
		return nil, err
	}

	cert, err := tls.X509KeyPair(certs, keyPEMBlock)
	if err != nil {
		return nil, err
	}

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			RootCAs:      certPool,
			Certificates: []tls.Certificate{cert},
		},
	}

	return &http.Client{Transport: transport}, nil
}

const finalizer = "feed.finalizers.news.teamdev.com"

// FeedReconciler reconciles a Feed object
type FeedReconciler struct {
	Client client.Client
	Scheme *runtime.Scheme
}

// FeedCreateRequest contains the url of the feed to save it
type FeedCreateRequest struct {
	Url string `json:"url"`
}

// +kubebuilder:rbac:groups=aggregator.com.teamdev,resources=feeds,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=aggregator.com.teamdev,resources=feeds/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=aggregator.com.teamdev,resources=feeds/finalizers,verbs=update

// Reconcile attempts to bring the state for the CRD Feed from the desired state to the current state.
// It receives data from the incoming request to further add the feed and update the status.
func (r *FeedReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	var feed aggregatorv1.Feed

	err := r.Client.Get(ctx, req.NamespacedName, &feed)
	if err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	// Ensure the feed is not nil and has required fields
	if feed.Spec.Name == "" {
		logrus.Error("Feed resource is missing required fields: name or url")
		return ctrl.Result{}, fmt.Errorf("feed resource is missing required fields: name or url")
	}
	if feed.Spec.Url == "" {
		logrus.Error("Feed resource is missing required NAME")
		return ctrl.Result{}, fmt.Errorf("feed resource is missing URL")
	}

	if !containsString(feed.ObjectMeta.Finalizers, finalizer) {
		feed.ObjectMeta.Finalizers = append(feed.ObjectMeta.Finalizers, finalizer)
		if err := r.Client.Update(ctx, &feed); err != nil {
			return ctrl.Result{}, err
		}
	}

	// Check if the object is being deleted
	if !feed.ObjectMeta.DeletionTimestamp.IsZero() {
		// The object is being deleted
		if containsString(feed.ObjectMeta.Finalizers, finalizer) {
			// Run finalization logic for feed finalizer
			if err := r.deleteFeed(&feed); err != nil {
				return ctrl.Result{}, err
			}

			// Remove finalizer from list and update it
			feed.ObjectMeta.Finalizers = removeString(feed.ObjectMeta.Finalizers, finalizer)
			if err := r.Client.Update(ctx, &feed); err != nil {
				return ctrl.Result{}, err
			}
		}
		return ctrl.Result{}, nil
	}

	// Normal reconciliation logic
	logrus.Info("Feed name: " + feed.Spec.Name + ". Feed link: " + feed.Spec.Url)

	feedCreateRequest := FeedCreateRequest{
		Url: feed.Spec.Url,
	}
	logrus.Info(feedCreateRequest.Url)

	reqBody, err := json.Marshal(feedCreateRequest)
	if err != nil {
		logrus.Error("Failed to marshal source request: ", err)
		return ctrl.Result{}, err
	}
	httpClient := &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	resp, err := httpClient.Post("https://news-aggregator-service.news-aggregator.svc.cluster.local:443/sources", "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		logrus.Error("Failed to make POST request: ", err)
		return ctrl.Result{}, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		logrus.Error("Failed to create source, status code: ", resp.StatusCode)
		return ctrl.Result{}, err
	}

	feed.Status.Status = "Source is added"
	err = r.Client.Status().Update(ctx, &feed)
	if err != nil {
		return ctrl.Result{}, err
	}

	logrus.Info("Status updated.")

	return ctrl.Result{}, nil
}

func (r *FeedReconciler) deleteFeed(feed *aggregatorv1.Feed) error {
	httpClient := &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	type DeleteRequest struct {
		Name string `json:"name"`
	}

	deleteRequest := DeleteRequest{
		Name: feed.Spec.Name,
	}

	reqBody, err := json.Marshal(deleteRequest)
	if err != nil {
		logrus.Error("Failed to marshal delete request: ", err)
		return err
	}

	logrus.Infof("Feed for delete name: %s", feed.Spec.Name)

	// Создайте запрос DELETE с телом
	req, err := http.NewRequest("DELETE", "https://news-aggregator-service.news-aggregator.svc.cluster.local:443/sources", bytes.NewBuffer(reqBody))
	if err != nil {
		logrus.Error("Failed to create DELETE request: ", err)
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := httpClient.Do(req)
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
		logrus.Error("Failed to delete source, status code: ", resp.StatusCode)
		return err
	}

	logrus.Info("Feed finalized successfully.")
	return nil
}

func (r *FeedReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&aggregatorv1.Feed{}).
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
	newSlice := []string{}
	for _, item := range slice {
		if item == str {
			continue
		}
		newSlice = append(newSlice, item)
	}
	return newSlice
}
