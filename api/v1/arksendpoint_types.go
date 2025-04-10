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
	gatewayv1 "sigs.k8s.io/gateway-api/apis/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ArksEndpointSpec defines the desired state of ArksEndpoint.
type ArksEndpointSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:default=1
	DefaultWeight int32 `json:"defaultWeight"`

	// +kubebuilder:validation:Required
	GatewayRef gatewayv1.ParentReference `json:"gatewayRef"`

	// +optional
	// +kubebuilder:validation:MaxItems=16
	// +kubebuilder:default={{path:{ type: "PathPrefix", value: "/"}}}
	MatchConfigs []gatewayv1.HTTPRouteMatch `json:"matchConfigs,omitempty"`

	// +optional
	// +kubebuilder:validation:MaxItems=16
	RouteConfigs []gatewayv1.HTTPBackendRef `json:"routeConfigs"`
}

// ArksEndpointStatus defines the observed state of ArksEndpoint.
type ArksEndpointStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// +optional
	// +kubebuilder:validation:MaxItems=16
	Routes []gatewayv1.HTTPBackendRef `json:"routes,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:printcolumn:name="Default Weight",type="integer",JSONPath=".spec.defaultWeight"
// ArksEndpoint is the Schema for the arksendpoints API.
type ArksEndpoint struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ArksEndpointSpec   `json:"spec,omitempty"`
	Status ArksEndpointStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ArksEndpointList contains a list of ArksEndpoint.
type ArksEndpointList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ArksEndpoint `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ArksEndpoint{}, &ArksEndpointList{})
}
