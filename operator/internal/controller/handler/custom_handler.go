package handler

import (
	aggregatorv1 "com.teamdev/news-aggregator/api/v1"
	"context"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
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
		requests = append(requests, ctrl.Request{
			NamespacedName: types.NamespacedName{
				Name:      hotNews.Name,
				Namespace: hotNews.Namespace,
			},
		})
	}
	return requests
}
