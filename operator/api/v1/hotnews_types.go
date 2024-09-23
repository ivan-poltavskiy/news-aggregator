package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// HotNewsSpec contains the specification's fields of the HotNews CRD
type HotNewsSpec struct {
	// Keywords specifies the list of keywords used to filter news articles
	Keywords []string `json:"keywords,omitempty"`

	// DateStart defines the start date for collecting news, in the format YYYY-MM-DD
	DateStart string `json:"dateStart,omitempty"`

	// DateEnd defines the end date for collecting news, in the format YYYY-MM-DD
	DateEnd string `json:"dateEnd,omitempty"`

	// FeedsName lists the names of the news sources to collect articles from
	FeedsName []string `json:"feedsName,omitempty"`

	// FeedGroups specifies the groups of news sources for aggregation
	FeedGroups []string `json:"feedGroups,omitempty"`

	// SummaryConfig contains configuration options for summarizing the news
	SummaryConfig SummaryConfig `json:"summaryConfig,omitempty"`
}

// SummaryConfig contains the summary configuration of the HotNews CRD
type SummaryConfig struct {

	// TitlesCount contains the quantity of the titles which will be stored in CRD
	TitlesCount int `json:"titlesCount,omitempty"`
}

// HotNewsStatus defines the observed state of HotNews
type HotNewsStatus struct {
	// ArticlesCount is the number of news articles collected
	ArticlesCount int `json:"articlesCount,omitempty"`

	// NewsLink provides a URL to access the aggregated news
	NewsLink string `json:"newsLink,omitempty"`

	// ArticlesTitles lists the titles of the collected news articles
	ArticlesTitles []string `json:"articlesTitles,omitempty"`

	// Conditions represent the current state or conditions of the HotNews resource
	Conditions []Condition `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// HotNews is the Schema for the hotnews API
type HotNews struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   HotNewsSpec   `json:"spec,omitempty"`
	Status HotNewsStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// HotNewsList contains a list of HotNews
type HotNewsList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []HotNews `json:"items"`
}

// AddCondition adds new condition to the HotNews's status
func (f *HotNewsStatus) AddCondition(condition Condition) {
	newCondition := Condition{
		Type:            condition.Type,
		Success:         condition.Success,
		Reason:          condition.Reason,
		Message:         condition.Message,
		LastUpdatedName: condition.LastUpdatedName,
		LastUpdateTime:  metav1.Now(),
	}

	f.Conditions = append(f.Conditions, newCondition)
}

// GetCurrentCondition returns the current condition of the HotNews
func (f *HotNewsStatus) GetCurrentCondition() Condition {
	if len(f.Conditions) == 0 {
		return Condition{}
	}
	return f.Conditions[len(f.Conditions)-1]
}

func init() {
	SchemeBuilder.Register(&HotNews{}, &HotNewsList{})
}
