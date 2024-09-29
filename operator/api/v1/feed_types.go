package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ConditionType represents the type of condition in the Feed lifecycle
type ConditionType string

const (
	// ConditionAdded indicates that the feed has been successfully added
	ConditionAdded ConditionType = "Added"

	// ConditionDeleted indicates that the feed has been successfully deleted
	ConditionDeleted ConditionType = "Deleted"
)

// Condition describes the states of a feed during its life cycle in the system
type Condition struct {
	// Type of the condition, e.g., Added, Deleted.
	Type ConditionType `json:"type"`
	// Success of the condition. Could be true or false
	Success bool `json:"status"`
	// If Success is False, the reason should be populated
	Reason string `json:"reason,omitempty"`
	// If Success is False, the message should be populated
	Message         string `json:"message,omitempty"`
	LastUpdatedName string `json:"lastUpdatedName,omitempty"`
	// Last time the condition transitioned from one status to another.
	LastUpdateTime metav1.Time `json:"lastUpdateTime,omitempty"`
}

// FeedStatus describes the status of a feed during its full life cycle in the system
type FeedStatus struct {
	Conditions []Condition `json:"conditions,omitempty"`
}

// FeedSpec contains the specification's fields of the Feed
type FeedSpec struct {
	// Describe name of the news feed
	Name string `json:"name,omitempty"`
	// URL for fetching news articles for feed
	Url string `json:"url,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// Feed describe the information of the news source for news aggregator in the K8S cluster
type Feed struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   FeedSpec   `json:"spec,omitempty"`
	Status FeedStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// FeedList contains a list of Feed
type FeedList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Feed `json:"items"`
}

// SetCondition adds new condition to the Feed's status
func (f *FeedStatus) SetCondition(condition Condition) {

	for i, currentCondition := range f.Conditions {
		if currentCondition.Type == condition.Type {
			f.Conditions[i] = condition
			return
		}
	}
	f.Conditions = append(f.Conditions, condition)
}

// GetCurrentCondition returns the current condition of the Feed
func (f *FeedStatus) GetCurrentCondition() Condition {
	if len(f.Conditions) == 0 {
		return Condition{}
	}
	return f.Conditions[len(f.Conditions)-1]
}

// AddPositiveCondition Set the success status to the condition of the feed
func AddPositiveCondition(feed *Feed) {

	feed.Status.SetCondition(Condition{
		Type:            ConditionAdded,
		Success:         true,
		LastUpdatedName: feed.Spec.Name,
		Message:         "",
		Reason:          "",
		LastUpdateTime:  metav1.Now(),
	})
}

func init() {
	SchemeBuilder.Register(&Feed{}, &FeedList{})
}
