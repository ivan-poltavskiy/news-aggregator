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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// HotNewsSpec contains the specification's fields of the HotNews CRD
type HotNewsSpec struct {
	Keywords      []string      `json:"keywords,omitempty"`
	DateStart     string        `json:"dateStart,omitempty"`
	DateEnd       string        `json:"dateEnd,omitempty"`
	FeedsName     []string      `json:"feedsName,omitempty"`
	FeedGroups    []string      `json:"feedGroups,omitempty"`
	SummaryConfig SummaryConfig `json:"summaryConfig,omitempty"`
}

// SummaryConfig contains the summary configuration of the HotNews CRD
type SummaryConfig struct {

	// TitlesCount contains the quantity of the titles which will be stored in CRD
	TitlesCount int `json:"titlesCount,omitempty"`
}

// HotNewsStatus defines the observed state of HotNews
type HotNewsStatus struct {
	ArticlesCount  int      `json:"articlesCount,omitempty"`
	NewsLink       string   `json:"newsLink,omitempty"`
	ArticlesTitles []string `json:"articlesTitles,omitempty"`
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

func init() {
	SchemeBuilder.Register(&HotNews{}, &HotNewsList{})
}
