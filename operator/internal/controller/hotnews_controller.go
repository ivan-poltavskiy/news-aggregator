package controller

import (
	aggregatorv1 "com.teamdev/news-aggregator/api/v1"
	additionalsHandlers "com.teamdev/news-aggregator/internal/controller/handler"
	additionalsPredicates "com.teamdev/news-aggregator/internal/controller/predicate"
	"context"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"net/http"
	"net/url"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"strings"
	"time"
)

// HotNewsReconciler reconciles a HotNews object
type HotNewsReconciler struct {
	client.Client
	Scheme        *runtime.Scheme
	HttpClient    HttpClient
	HttpsLinks    HttpsClientData
	Finalizer     string
	ConfigMapName string
}

// +kubebuilder:rbac:groups=aggregator.com.teamdev,resources=hotnews,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=aggregator.com.teamdev,resources=hotnews/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=aggregator.com.teamdev,resources=hotnews/finalizers,verbs=update
// +kubebuilder:rbac:groups="",resources=configmaps,verbs=get;list;watch
// +kubebuilder:rbac:groups=news-aggregator.com.teamdev,resources=feeds,verbs=get;list;watch;create;update;patch;delete

// The Reconcile method brings the current state of the HotNews object to the desired state.
// It checks the existence of the object in the system. Additionally, the method manages
// finalizers and owner references of the HotNews object. When attempting to delete the object,
// it checks for the presence of the owner reference and finalizer.
// The method also retrieves the news from the server based on the parameters defined in HotNewsSpec.
// Regardless of the success of the operation, the method updates the current state of the object in its status.
func (r *HotNewsReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {

	logrus.Info("Starting hot news reconcile")
	var hotNews aggregatorv1.HotNews

	err := r.Client.Get(ctx, req.NamespacedName, &hotNews)
	if err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	if !containsString(hotNews.ObjectMeta.Finalizers, r.Finalizer) {
		logrus.Info("Adding finalizer to: " + hotNews.Name)
		hotNews.ObjectMeta.Finalizers = append(hotNews.ObjectMeta.Finalizers, r.Finalizer)
		if err := r.Client.Update(ctx, &hotNews); err != nil {
			return ctrl.Result{}, err
		}
	}

	if !hotNews.ObjectMeta.DeletionTimestamp.IsZero() {
		if containsString(hotNews.ObjectMeta.Finalizers, r.Finalizer) {
			if cleanupErr := r.CleanupOwnerReferences(ctx, req.Namespace, req.Name); cleanupErr != nil {
				return ctrl.Result{}, cleanupErr
			}
			logrus.Info("Removing finalizer of: " + hotNews.Name)
			hotNews.ObjectMeta.Finalizers = removeString(hotNews.ObjectMeta.Finalizers, r.Finalizer)
			if err := r.Client.Update(ctx, &hotNews); err != nil {
				return ctrl.Result{}, err
			}
		}
		return ctrl.Result{}, nil
	}

	err = r.reconcileHotNews(&hotNews, req.Namespace, ctx)
	if err != nil {
		hotNews.Status.SetCondition(aggregatorv1.Condition{
			Type:    aggregatorv1.ConditionAdded,
			Success: false,
			Message: "Hot News Reconcile Failed",
			Reason:  err.Error(),
		})
		if err := r.Client.Status().Update(ctx, &hotNews); err != nil {
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, err
	}

	hotNews.Status.SetCondition(aggregatorv1.Condition{
		Type:    aggregatorv1.ConditionAdded,
		Success: true,
		Message: "",
		Reason:  "",
	})

	if err := r.UpdateOwnerReferencesForFeeds(ctx, &hotNews); err != nil {
		return ctrl.Result{}, err
	}

	if err := r.Client.Status().Update(ctx, &hotNews); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *HotNewsReconciler) SetupWithManager(mgr ctrl.Manager) error {

	hotNewsHandler := additionalsHandlers.HotNewsHandler{Client: r.Client}

	return ctrl.NewControllerManagedBy(mgr).
		For(&aggregatorv1.HotNews{}, builder.WithPredicates(predicate.GenerationChangedPredicate{})).
		Watches(
			&aggregatorv1.Feed{},
			handler.EnqueueRequestsFromMapFunc(hotNewsHandler.UpdateHotNews),
		).
		Watches(
			&v1.ConfigMap{},
			handler.EnqueueRequestsFromMapFunc(hotNewsHandler.UpdateConfigMap),
			builder.WithPredicates(additionalsPredicates.ConfigMapNamePredicate(r.ConfigMapName)),
		).
		Complete(r)
}

// reconcileHotNews synchronizes the state of the HotNews custom resource
// with the external news sources specified in the ConfigMap.
func (r *HotNewsReconciler) reconcileHotNews(hotNews *aggregatorv1.HotNews, namespace string, ctx context.Context) error {
	var feedGroupConfigMap v1.ConfigMap
	var feedNames []string

	err := r.Get(ctx, client.ObjectKey{Namespace: namespace, Name: r.ConfigMapName}, &feedGroupConfigMap)
	if err != nil && !errors.IsNotFound(err) {
		return err
	}

	if err == nil {
		feedNames = r.getFeedNamesFromConfigMap(hotNews, &feedGroupConfigMap)
		logrus.Info("feeds name from config map length: ", len(feedNames))

		if len(hotNews.Spec.FeedsName) != 0 {
			feedNames = append(feedNames, hotNews.Spec.FeedsName...)
		}
	}

	createdUrl, err := r.createUrl(*hotNews, feedNames)
	if err != nil {
		return err
	}

	logrus.Info("URL= ", createdUrl)

	articles, err := r.fetchNews(createdUrl)
	if err != nil {
		return err
	}

	titlesCount := hotNews.Spec.SummaryConfig.TitlesCount
	if len(articles) < titlesCount {
		titlesCount = len(articles)
	}

	hotNews.Status.ArticlesTitles = getTopTitles(articles, titlesCount)
	logrus.Info("Length of news: ", len(articles))

	hotNews.Status.ArticlesCount = len(articles)
	hotNews.Status.NewsLink = createdUrl

	return nil
}

// getFeedNamesFromConfigMap retrieves the list of feed names from the ConfigMap and removes spaces around feed names.
func (r *HotNewsReconciler) getFeedNamesFromConfigMap(hotNews *aggregatorv1.HotNews, configMap *v1.ConfigMap) []string {
	var feedNames []string

	for _, feedGroup := range hotNews.Spec.FeedGroups {

		if feeds, found := configMap.Data[feedGroup]; found {

			for _, feed := range strings.Split(feeds, ",") {
				feedNames = append(feedNames, strings.TrimSpace(feed))
			}
		}
	}

	return feedNames
}

// createUrl constructs the URL used to fetch news based on the
// configuration provided in the HotNews resource and the related ConfigMap.
func (r *HotNewsReconciler) createUrl(hotNews aggregatorv1.HotNews, feedNames []string) (string, error) {
	baseUrl := r.HttpsLinks.ServerUrl + r.HttpsLinks.EndpointForSourceManaging
	params := url.Values{}

	if len(feedNames) > 0 {
		params.Add("sources", strings.Join(feedNames, ","))
	} else {
		return "", fmt.Errorf("feeds and config maps not present")
	}

	if len(hotNews.Spec.Keywords) > 0 {
		params.Add("keywords", strings.Join(hotNews.Spec.Keywords, ","))
	}

	if hotNews.Spec.DateStart != "" && hotNews.Spec.DateEnd != "" {
		params.Add("startDate", hotNews.Spec.DateStart)
		params.Add("endDate", hotNews.Spec.DateEnd)
	}

	return baseUrl + "?" + params.Encode(), nil
}

// fetchNews sends an HTTP GET request to the specified URL to retrieve a list of news articles.
func (r *HotNewsReconciler) fetchNews(url string) ([]news, error) {
	resp, err := r.HttpClient.Get(url)
	if err != nil {
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logrus.Error("Failed to close response body")
			return
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get articles from news aggregator: %s", resp.Status)
	}

	var articles []news
	if err := json.NewDecoder(resp.Body).Decode(&articles); err != nil {
		return nil, err
	}

	return articles, nil
}

// getTopTitles extracts the titles of the top news based on the specified count.
func getTopTitles(articles []news, count int) []string {
	var titles []string
	for i := 0; i < len(articles) && i < count; i++ { // Adjusted condition to correctly limit the number of titles
		titles = append(titles, articles[i].Title)
	}
	return titles
}

type news struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Link        string    `json:"url"`
	Date        time.Time `json:"publishedAt"`
	SourceName  string
}
