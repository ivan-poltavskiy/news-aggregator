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
	aggregatorv1 "com.teamdev/news-aggregator/api/v1"
	"context"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"net/http"
	"net/url"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"strings"
	"time"
)

// HotNewsReconciler reconciles a HotNews object
type HotNewsReconciler struct {
	client.Client
	Scheme     *runtime.Scheme
	HttpClient HttpClient
	HttpsLinks HttpsClientData
	Finalizer  string
}

// +kubebuilder:rbac:groups=aggregator.com.teamdev,resources=hotnews,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=aggregator.com.teamdev,resources=hotnews/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=aggregator.com.teamdev,resources=hotnews/finalizers,verbs=update

func (r *HotNewsReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {

	var hotNews aggregatorv1.HotNews

	err := r.Client.Get(ctx, req.NamespacedName, &hotNews)
	if err != nil {
		if errors.IsNotFound(err) {
			logrus.Info("Reconcile: Feed was not found. Error ignored")
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	if !containsString(hotNews.ObjectMeta.Finalizers, r.Finalizer) {
		hotNews.ObjectMeta.Finalizers = append(hotNews.ObjectMeta.Finalizers, r.Finalizer)
		if err := r.Client.Update(ctx, &hotNews); err != nil {
			return ctrl.Result{}, err
		}
	}

	if !hotNews.ObjectMeta.DeletionTimestamp.IsZero() {
		if containsString(hotNews.ObjectMeta.Finalizers, r.Finalizer) {
			hotNews.ObjectMeta.Finalizers = removeString(hotNews.ObjectMeta.Finalizers, r.Finalizer)
			if err := r.Client.Update(ctx, &hotNews); err != nil {
				return ctrl.Result{}, err
			}
		}
		return ctrl.Result{}, nil
	}

	url := r.createUrl(hotNews)
	logrus.Info("URL= ", url)

	articles, err := r.fetchArticles(url)

	if err != nil {
		return ctrl.Result{}, err
	}
	hotNews.Status.ArticlesCount = len(articles)
	hotNews.Status.NewsLink = url

	if err := r.Client.Status().Update(ctx, &hotNews); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *HotNewsReconciler) createUrl(hotNews aggregatorv1.HotNews) string {
	baseUrl := r.HttpsLinks.ServerUrl + r.HttpsLinks.EndpointForSourceManaging
	params := url.Values{}

	if len(hotNews.Spec.FeedsName) > 0 {
		params.Add("sources", strings.Join(hotNews.Spec.FeedsName, ","))
	}

	if len(hotNews.Spec.Keywords) > 0 {
		params.Add("keywords=", strings.Join(hotNews.Spec.Keywords, ","))
	}

	if hotNews.Spec.DateStart != "" && hotNews.Spec.DateEnd != "" {
		params.Add("startDate", hotNews.Spec.DateStart)
		params.Add("endDate", hotNews.Spec.DateEnd)
	}

	return baseUrl + "?" + params.Encode()
}

type News struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Link        string    `json:"url"`
	Date        time.Time `json:"publishedAt"`
	SourceName  string
}

func (r *HotNewsReconciler) fetchArticles(url string) ([]News, error) {
	resp, err := r.HttpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get articles from news aggregator: %s", resp.Status)
	}

	var articles []News
	if err := json.NewDecoder(resp.Body).Decode(&articles); err != nil {
		return nil, err
	}

	return articles, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *HotNewsReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&aggregatorv1.HotNews{}).
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
