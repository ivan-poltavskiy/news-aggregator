package controller

import (
	aggregatorv1 "com.teamdev/news-aggregator/api/v1"
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// UpdateOwnerReferencesForFeeds manages owner references for feeds based on their usage in HotNews.
func (r *HotNewsReconciler) UpdateOwnerReferencesForFeeds(ctx context.Context, hotNews *aggregatorv1.HotNews) error {
	feedList := &aggregatorv1.FeedList{}
	listOpts := client.ListOptions{Namespace: hotNews.Namespace}

	err := r.Client.List(ctx, feedList, &listOpts)
	if err != nil {
		return err
	}

	for _, feed := range feedList.Items {
		logrus.Info("UpdateOwnerReferencesForFeeds: FEED NAME: ", &feed.Spec.Name)
		logrus.Info("UpdateOwnerReferencesForFeeds: hotNews.Spec.FeedsName: ", hotNews.Spec.FeedsName)

		if r.isFeedUsedInHotNews(&feed, hotNews.Spec.FeedsName) {
			if err := r.AddOwnerReference(ctx, &feed, hotNews); err != nil {
				return fmt.Errorf("failed to add ownerReference to Feed %s: %w", feed.Name, err)
			}
		} else {
			if err := r.RemoveOwnerReference(ctx, &feed, hotNews); err != nil {
				return fmt.Errorf("failed to remove ownerReference from Feed %s: %w", feed.Name, err)
			}
		}
	}

	return nil
}

// isFeedUsedInHotNews checks if a feed is used in the current HotNews resource.
func (r *HotNewsReconciler) isFeedUsedInHotNews(feed *aggregatorv1.Feed, feedNames []string) bool {
	for _, feedName := range feedNames {
		if feedName == feed.Spec.Name {
			return true
		}
	}
	return false
}

// AddOwnerReference adds an OwnerReference to a feed.
func (r *HotNewsReconciler) AddOwnerReference(ctx context.Context, feed *aggregatorv1.Feed, hotNews *aggregatorv1.HotNews) error {
	ownerRef := metav1.OwnerReference{
		APIVersion:         hotNews.APIVersion,
		Kind:               hotNews.Kind,
		Name:               hotNews.Name,
		UID:                hotNews.UID,
		BlockOwnerDeletion: pointer.BoolPtr(false),
	}

	existingOwnerReferences := feed.ObjectMeta.OwnerReferences
	for _, ref := range existingOwnerReferences {
		if ref.UID == hotNews.UID {
			return nil
		}
	}

	feed.ObjectMeta.OwnerReferences = append(feed.ObjectMeta.OwnerReferences, ownerRef)
	return r.Client.Update(ctx, feed)
}

// RemoveOwnerReference removes an OwnerReference from a feed.
func (r *HotNewsReconciler) RemoveOwnerReference(ctx context.Context, feed *aggregatorv1.Feed, hotNews *aggregatorv1.HotNews) error {
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

// CleanupOwnerReferences removes OwnerReferences to the deleted HotNews from all related Feeds
func (r *HotNewsReconciler) CleanupOwnerReferences(ctx context.Context, namespace, hotNewsName string) error {
	var feedList aggregatorv1.FeedList
	listOpts := client.ListOptions{Namespace: namespace}

	if err := r.Client.List(ctx, &feedList, &listOpts); err != nil {
		return fmt.Errorf("failed to list Feeds: %w", err)
	}

	for _, feed := range feedList.Items {
		var updatedOwnerReferences []metav1.OwnerReference
		ownerReferenceRemoved := false

		for _, ref := range feed.OwnerReferences {
			if ref.Name == hotNewsName && ref.Kind == "HotNews" {
				ownerReferenceRemoved = true
				continue
			}
			updatedOwnerReferences = append(updatedOwnerReferences, ref)
		}

		if ownerReferenceRemoved {
			feed.OwnerReferences = updatedOwnerReferences
			if err := r.Client.Update(ctx, &feed); err != nil {
				logrus.Errorf("Failed to remove OwnerReference from Feed %s: %v", feed.Name, err)
				return fmt.Errorf("failed to update Feed %s: %w", feed.Name, err)
			}
			logrus.Infof("Successfully removed OwnerReference from Feed %s", feed.Name)
		}
	}

	return nil
}
