package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ConditionType represents a condition type for a Feed
type ConditionType string

const (
	// ConditionAdded indicates that the feed has been successfully added
	ConditionAdded ConditionType = "Added"
	// ConditionDeleted indicates that the feed has been successfully deleted
	ConditionDeleted ConditionType = "Deleted"
)

// Condition represents the state of a Feed at a certain point.
type Condition struct {
	// Type of the condition, e.g., Added, Updated, Deleted.
	Type ConditionType `json:"type"`
	// Status of the condition. Could be true or false
	Status bool `json:"status"`
	// If status is False, the reason should be populated
	Reason string `json:"reason,omitempty"`
	// If status is False, the message should be populated
	Message string `json:"message,omitempty"`
	// Last time the condition transitioned from one status to another.
	LastUpdateTime metav1.Time `json:"lastUpdateTime,omitempty"`
}

// FeedStatus defines the observed state of Feed
type FeedStatus struct {
	Conditions []Condition `json:"conditions,omitempty"`
}

type FeedSpec struct {
	Name string `json:"name,omitempty"`
	Url  string `json:"url,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

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

func init() {
	SchemeBuilder.Register(&Feed{}, &FeedList{})
}
