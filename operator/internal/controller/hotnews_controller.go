package controller

import (
	aggregatorv1 "com.teamdev/news-aggregator/api/v1"
	"context"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"net/http"
	"net/url"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
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
	ConfigMapMame string
}

// +kubebuilder:rbac:groups=aggregator.com.teamdev,resources=hotnews,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=aggregator.com.teamdev,resources=hotnews/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=aggregator.com.teamdev,resources=hotnews/finalizers,verbs=update
// +kubebuilder:rbac:groups="",resources=configmaps,verbs=get;list;watch

func (r *HotNewsReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	var hotNews aggregatorv1.HotNews
	var feedGroupConfigMap v1.ConfigMap

	if err := r.Get(ctx, client.ObjectKey{Namespace: req.Namespace, Name: r.ConfigMapMame}, &feedGroupConfigMap); err != nil {
		if errors.IsNotFound(err) {
			logrus.Print("ConfigMap not found")
			return ctrl.Result{}, err
		}
		logrus.Printf("Error retrieving ConfigMap %s from k8s Cluster: %v", "feed-group-source", err)
		return ctrl.Result{}, err
	}

	err := r.Client.Get(ctx, req.NamespacedName, &hotNews)
	if err != nil {
		if errors.IsNotFound(err) {
			logrus.Info("Reconcile: Feed was not found. Error ignored")
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	err = r.reconcileHotNews(&hotNews, &feedGroupConfigMap)
	if err != nil {
		hotNews.Status.AddCondition(aggregatorv1.Condition{
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
	hotNews.Status.AddCondition(aggregatorv1.Condition{
		Type:    aggregatorv1.ConditionAdded,
		Success: true,
		Message: "",
		Reason:  "",
	})

	if err := r.updateFeedOwnerReferences(ctx, &hotNews); err != nil {
		logrus.Errorf("Failed to update Feed ownerReferences: %v", err)
		return ctrl.Result{}, err
	}

	// Remove owner references from feeds that are not used by HotNews anymore
	if err := r.removeUnusedOwnerReferences(ctx, &hotNews); err != nil {
		logrus.Errorf("Failed to remove unused Feed ownerReferences: %v", err)
		return ctrl.Result{}, err
	}

	if err := r.Client.Status().Update(ctx, &hotNews); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
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
			GenericFunc: func(e event.GenericEvent) bool {
				return true
			},
		}).
		Watches(
			&aggregatorv1.Feed{},
			handler.EnqueueRequestsFromMapFunc(r.updateHotNews),
		).
		Watches(
			&v1.ConfigMap{},
			handler.EnqueueRequestsFromMapFunc(r.updateHotNews),
		).
		Complete(r)
}

// reconcileHotNews synchronizes the state of the HotNews custom resource
// with the external news sources specified in the ConfigMap.
func (r *HotNewsReconciler) reconcileHotNews(hotNews *aggregatorv1.HotNews, configMap *v1.ConfigMap) error {

	createdUrl, err := r.createUrl(*hotNews, configMap)
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

func (r *HotNewsReconciler) updateFeedOwnerReferences(ctx context.Context, hotNews *aggregatorv1.HotNews) error {
	feedList := &aggregatorv1.FeedList{}
	listOpts := client.ListOptions{Namespace: hotNews.Namespace}

	err := r.Client.List(ctx, feedList, &listOpts)
	if err != nil {
		return err
	}

	for _, feed := range feedList.Items {
		if r.isFeedUsedInHotNews(&feed, hotNews) {
			if err := r.addOwnerReference(ctx, &feed, hotNews); err != nil {
				return fmt.Errorf("failed to add ownerReference to Feed %s: %w", feed.Name, err)
			}
		}
	}

	return nil
}

func (r *HotNewsReconciler) removeUnusedOwnerReferences(ctx context.Context, hotNews *aggregatorv1.HotNews) error {
	feedList := &aggregatorv1.FeedList{}
	listOpts := client.ListOptions{Namespace: hotNews.Namespace}

	err := r.Client.List(ctx, feedList, &listOpts)
	if err != nil {
		return err
	}

	for _, feed := range feedList.Items {
		if !r.isFeedUsedInHotNews(&feed, hotNews) {
			if err := r.removeOwnerReference(ctx, &feed, hotNews); err != nil {
				return fmt.Errorf("failed to remove ownerReference from Feed %s: %w", feed.Name, err)
			}
		}
	}

	return nil
}

// updateHotNews is a handler function that is triggered when relevant changes
// occur to resources that the controller watches.
func (r *HotNewsReconciler) updateHotNews(ctx context.Context, obj client.Object) []reconcile.Request {
	var hotNewsList aggregatorv1.HotNewsList

	if err := r.List(ctx, &hotNewsList, client.InNamespace(obj.GetNamespace())); err != nil {
		log.Log.Error(err, "Failed to list HotNews resources")
		return nil
	}

	var requests []ctrl.Request
	for _, hotNews := range hotNewsList.Items {
		var feedGroupConfigMap v1.ConfigMap
		if err := r.Get(ctx, client.ObjectKey{Namespace: hotNews.Namespace, Name: r.ConfigMapMame}, &feedGroupConfigMap); err != nil {
			log.Log.Error(err, "Failed to get ConfigMap")
			continue
		}

		err := r.reconcileHotNews(&hotNews, &feedGroupConfigMap)
		if err != nil {
			log.Log.Error(err, "Failed to reconcile HotNews after Feed/ConfigMap update")
			continue
		}

		requests = append(requests, ctrl.Request{
			NamespacedName: client.ObjectKey{
				Name:      hotNews.Name,
				Namespace: hotNews.Namespace,
			},
		})
	}

	return requests
}

// createUrl constructs the URL used to fetch news based on the
// configuration provided in the HotNews resource and the related ConfigMap.
func (r *HotNewsReconciler) createUrl(hotNews aggregatorv1.HotNews, configMap *v1.ConfigMap) (string, error) {
	baseUrl := r.HttpsLinks.ServerUrl + r.HttpsLinks.EndpointForSourceManaging
	params := url.Values{}

	if len(hotNews.Spec.FeedGroups) > 0 {
		var feedNames []string
		for _, group := range hotNews.Spec.FeedGroups {
			if feeds, found := configMap.Data[group]; found {
				feedList := strings.Split(feeds, ",")
				for _, feedName := range feedList {
					feedName = strings.TrimSpace(feedName)
					var currentFeed aggregatorv1.Feed
					err := r.Client.Get(context.TODO(), client.ObjectKey{Namespace: hotNews.Namespace, Name: feedName}, &currentFeed)
					if err != nil {
						return "", fmt.Errorf("failed to get Feed %s: %w", feedName, err)
					}
					feedNames = append(feedNames, currentFeed.Spec.Name)
				}
			} else {
				logrus.Warnf("Feed group %s not found in ConfigMap", group)
			}
		}
		logrus.Info("Sources from feed groups: ", strings.Join(feedNames, ","))
		hotNews.Spec.FeedsName = feedNames
	}

	if len(hotNews.Spec.FeedsName) > 0 {
		params.Add("sources", strings.Join(hotNews.Spec.FeedsName, ","))
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

type news struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Link        string    `json:"url"`
	Date        time.Time `json:"publishedAt"`
	SourceName  string
}

// fetchNews sends an HTTP GET request to the specified URL to retrieve a list
// of news
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

// getTopTitles extracts the titles of the top news based on the
// specified count.
func getTopTitles(articles []news, count int) []string {
	var titles []string
	for i := 0; i < len(articles) && i <= count; i++ {
		titles = append(titles, articles[i].Title)
	}
	return titles
}

func (r *HotNewsReconciler) isFeedUsedInHotNews(feed *aggregatorv1.Feed, hotNews *aggregatorv1.HotNews) bool {
	for _, feedName := range hotNews.Spec.FeedsName {
		if feedName == feed.Name {
			return true
		}
	}
	return false
}

func (r *HotNewsReconciler) addOwnerReference(ctx context.Context, feed *aggregatorv1.Feed, hotNews *aggregatorv1.HotNews) error {
	ownerRef := metav1.OwnerReference{
		APIVersion: hotNews.APIVersion,
		Kind:       hotNews.Kind,
		Name:       hotNews.Name,
		UID:        hotNews.UID,
	}

	existingOwnerReferences := feed.ObjectMeta.OwnerReferences
	for _, ref := range existingOwnerReferences {
		if ref.UID == hotNews.UID {
			return nil
		}
	}

	feed.ObjectMeta.OwnerReferences = append(feed.ObjectMeta.OwnerReferences, ownerRef)
	err := r.Client.Update(ctx, feed)
	return err
}

func (r *HotNewsReconciler) getFeedsWithOwnerReferences(ctx context.Context, hotNews *aggregatorv1.HotNews) ([]aggregatorv1.Feed, error) {
	var feedList aggregatorv1.FeedList
	listOpts := client.ListOptions{Namespace: hotNews.Namespace}

	err := r.Client.List(ctx, &feedList, &listOpts)
	if err != nil {
		return nil, err
	}

	var feedsWithOwnerReferences []aggregatorv1.Feed
	for _, feed := range feedList.Items {
		for _, ref := range feed.ObjectMeta.OwnerReferences {
			if ref.UID == hotNews.UID {
				feedsWithOwnerReferences = append(feedsWithOwnerReferences, feed)
				break
			}
		}
	}

	return feedsWithOwnerReferences, nil
}

func (r *HotNewsReconciler) removeOwnerReference(ctx context.Context, feed *aggregatorv1.Feed, hotNews *aggregatorv1.HotNews) error {
	var updatedOwnerReferences []metav1.OwnerReference

	for _, ref := range feed.ObjectMeta.OwnerReferences {
		if ref.UID != hotNews.UID {
			updatedOwnerReferences = append(updatedOwnerReferences, ref)
		}
	}

	if len(updatedOwnerReferences) == len(feed.ObjectMeta.OwnerReferences) {
		return nil
	}

	feed.ObjectMeta.OwnerReferences = updatedOwnerReferences
	return r.Client.Update(ctx, feed)
}
