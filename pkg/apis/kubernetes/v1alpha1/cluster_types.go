// Copyright 2019 (c) Microsoft and contributors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

var localSchemeBuilder = &SchemeBuilder

const (
	Provisioned  ClusterPhase = "Provisioned"
	Provisioning ClusterPhase = "Provisioning"
	Deleting     ClusterPhase = "Deleting"
)

type ClusterPhase string

// ClusterSpec defines the desired state of Cluster
type ClusterSpec struct {
	Cloud        Cloud        `json:"cloud,omitempty"`
	ControlPlane ControlPlane `json:"controlplane,omitempty"`
	Network      Network      `json:"network,omitempty"`
}

// ClusterStatus defines the observed state of Cluster
type ClusterStatus struct {
	Phase   ClusterPhase `json:"phase,omitempty"`
	Message string       `json:"message,omitempty"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Cluster is the Schema for the clusters API
// +k8s:openapi-gen=true
type Cluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ClusterSpec   `json:"spec,omitempty"`
	Status ClusterStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ClusterList contains a list of Cluster
type ClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Cluster `json:"items"`
}

// Cloud is the Schema for cluster cloud configuration
type Cloud struct {
	SubscriptionID string `json:"subscriptionID,omitempty"`
	ResourceGroup  string `json:"resourceGroup,omitempty"`
	Location       string `json:"location,omitempty"`
}

// ControlPlane is the Schema for a cluster's control plane
type ControlPlane struct {
	Fqdn string `json:"fqdn,omitempty"`
}

// Network is the Schema for cluster networking
type Network struct {
	CIDR       string `json:"cidr,omitempty"`
	SubnetCIDR string `json:"subnetCIDR,omitempty"`
}

func init() {
	SchemeBuilder.Register(addDefaultingFuncs, addConversionFuncs)
}

func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(SchemeGroupVersion, &Cluster{}, &ClusterList{})

	// Add common types
	scheme.AddKnownTypes(SchemeGroupVersion, &metav1.Status{})

	// Add the watch version that applies
	metav1.AddToGroupVersion(scheme, SchemeGroupVersion)

	return nil
}
