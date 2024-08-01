/*
Copyright 2024.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	aggregatorv1 "com.teamdev/news-aggregator/api/v1"
)

// FeedReconciler reconciles a Feed object
type FeedReconciler struct {
	Client client.Client
	Scheme *runtime.Scheme
}
type SourceRequest struct {
	Link string `json:"link"`
}

// +kubebuilder:rbac:groups=aggregator.com.teamdev,resources=feeds,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=aggregator.com.teamdev,resources=feeds/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=aggregator.com.teamdev,resources=feeds/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Feed object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.18.4/pkg/reconcile
func (r *FeedReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {

	_ = log.FromContext(ctx)
	logrus.Info("Hello from reconcile")

	var feed aggregatorv1.Feed

	err := r.Client.Get(ctx, req.NamespacedName, &feed)
	if err != nil {
		return ctrl.Result{}, err
	}
	logrus.Info("Feed name: " + feed.Spec.Name + ". Feed link: " + feed.Spec.Link)

	sourceReq := SourceRequest{
		Link: feed.Spec.Link,
	}

	reqBody, err := json.Marshal(sourceReq)
	if err != nil {
		logrus.Error("Failed to marshal source request: ", err)
		return ctrl.Result{}, err
	}

	// Create a custom HTTP client with TLS verification disabled
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	resp, err := client.Post("https://news-aggregator-service.news-aggregator.svc.cluster.local/sources", "application/json", bytes.NewBuffer(reqBody))
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

func (r *FeedReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&aggregatorv1.Feed{}).
		Complete(r)
}
