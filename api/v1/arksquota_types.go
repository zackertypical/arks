/*
Copyright 2025.

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

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

type QuotaType string

const (
	QuotaTypePrompt   QuotaType = "prompt"
	QuotaTypeResponse QuotaType = "response"
	QuotaTypeTotal    QuotaType = "total"
	// TODO: support more types
)

// QuotaItem defines a single quota configuration
type QuotaItem struct {
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Enum=prompt;response;total
	Type string `json:"type"`

	// Value of the quota
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Minimum=0
	Value int64 `json:"value"`
}

// ArksQuotaSpec defines the desired state of ArksQuota
type ArksQuotaSpec struct {
	// Quotas is a list of quota configurations
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinItems=1
	Quotas []QuotaItem `json:"quotas"`
}

// QuotaStatus represents the current usage of a quota
type QuotaStatus struct {
	// Name of the quota type
	Type string `json:"type"`

	// Current usage value
	Used int64 `json:"used"`

	// Last update time of the status
	// +optional
	LastUpdateTime *metav1.Time `json:"lastUpdateTime,omitempty"`
}

// ArksQuotaStatus defines the observed state of ArksQuota
type ArksQuotaStatus struct {
	// List of quota usage status
	// +optional
	QuotaStatus []QuotaStatus `json:"quotaStatus,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:resource:shortName=aq

// ArksQuota is the Schema for the arksquotas API
type ArksQuota struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ArksQuotaSpec   `json:"spec,omitempty"`
	Status ArksQuotaStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ArksQuotaList contains a list of ArksQuota
type ArksQuotaList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ArksQuota `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ArksQuota{}, &ArksQuotaList{})
}
