package handler

import (
	aggregatorv1 "com.teamdev/news-aggregator/api/v1"
	"context"
	"github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"strings"
)

type HotNewsHandler struct {
	Client client.Client
}

// UpdateHotNews is a handler function that is triggered when relevant changes
// occur to resources that the controller watches.
func (h *HotNewsHandler) UpdateHotNews(ctx context.Context, obj client.Object) []reconcile.Request {
	var hotNewsList aggregatorv1.HotNewsList
	// List only the HotNews resources in the same namespace as the changed object
	if err := h.Client.List(ctx, &hotNewsList, client.InNamespace(obj.GetNamespace())); err != nil {
		log.Log.Error(err, "Failed to list HotNews resources")
		return nil
	}

	var requests []ctrl.Request
	for _, hotNews := range hotNewsList.Items {
		if isOwner(hotNews, obj) {
			logrus.Info("Creating request for hot news " + hotNews.Name + " when child feed " + obj.GetName() + " updated")
			requests = append(requests, ctrl.Request{
				NamespacedName: types.NamespacedName{
					Name:      hotNews.Name,
					Namespace: hotNews.Namespace,
				},
			})
		}
	}
	return requests
}

// UpdateConfigMap processes changes to ConfigMap objects and generates reconcile requests for HotNews resources.
func (h *HotNewsHandler) UpdateConfigMap(ctx context.Context, obj client.Object) []ctrl.Request {

	configMap, ok := obj.(*v1.ConfigMap)
	if !ok {
		logrus.Info("Object isn't ConfigMap: ", obj)

		return nil
	}

	namespace := configMap.Namespace
	hotNewsList := &aggregatorv1.HotNewsList{}
	err := h.Client.List(ctx, hotNewsList, client.InNamespace(namespace))
	if err != nil {
		logrus.Info("Error listing HotNews in namespace ", namespace, err)
		return nil
	}

	var requests []ctrl.Request
	for _, hotNews := range hotNewsList.Items {
		feedGroups := getFeedsNamesFromFeedGroups(*configMap, hotNews)
		if len(feedGroups) > 0 {
			requests = append(requests, ctrl.Request{
				NamespacedName: client.ObjectKey{
					Namespace: hotNews.Namespace,
					Name:      hotNews.Name,
				},
			})
		}
	}

	logrus.Info("Completed UpdateConfigMap")

	return requests
}

func getFeedsNamesFromFeedGroups(configMap v1.ConfigMap, hotNews aggregatorv1.HotNews) []string {
	var feedsNames []string

	for _, feedGroup := range hotNews.Spec.FeedGroups {
		if value, ok := configMap.Data[feedGroup]; ok {
			feedsNames = append(feedsNames, strings.Split(value, ",")...)
		}
	}
	return feedsNames
}

// check that provided object is the child of provided hot news
func isOwner(hotNews aggregatorv1.HotNews, obj client.Object) bool {
	for _, ownerRef := range obj.GetOwnerReferences() {
		if ownerRef.UID == hotNews.UID {
			return true
		}
	}
	return false
}
