/*
Copyright 2022 Kong Inc.

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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func init() {
	SchemeBuilder.Register(&DataPlane{}, &DataPlaneList{})
}

//+genclient
//+k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:shortName=kdp,categories=kong;all
//+kubebuilder:printcolumn:name="Ready",description="The Resource is ready",type=string,JSONPath=`.status.conditions[?(@.type=='Ready')].status`
//+kubebuilder:printcolumn:name="Provisioned",description="The Resource is provisioned",type=string,JSONPath=`.status.conditions[?(@.type=='Provisioned')].status`

// DataPlane is the Schema for the dataplanes API
type DataPlane struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DataPlaneSpec   `json:"spec,omitempty"`
	Status DataPlaneStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// DataPlaneList contains a list of DataPlane
type DataPlaneList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []DataPlane `json:"items"`
}

// DataPlaneSpec defines the desired state of DataPlane
type DataPlaneSpec struct {
	DataPlaneDeploymentOptions `json:",inline"`
}

// DataPlaneDeploymentOptions defines the information specifically needed to
// deploy the DataPlane.
type DataPlaneDeploymentOptions struct {
	DeploymentOptions `json:",inline"`
}

// DataPlaneStatus defines the observed state of DataPlane
type DataPlaneStatus struct {
	// Conditions describe the status of the DataPlane.
	// +listType=map
	// +listMapKey=type
	// +kubebuilder:validation:MinItems=1
	// +kubebuilder:validation:MaxItems=8
	Conditions []metav1.Condition `json:"conditions,omitempty"`

	// Service indicates the Service that exposes the DataPlane's configured routes
	Service string `json:"service,omitempty"`

	// Addresses lists the addresses that have actually been bound to the DataPlane.
	//
	// +optional
	Addresses []Address `json:"addresses,omitempty"`
}

// Address describes an address which can be either an IP address or a hostname.
type Address struct {
	// Type of the address.
	//
	// +optional
	// +kubebuilder:default=IPAddress
	Type *AddressType `json:"type,omitempty"`

	// Value of the address. The validity of the values will depend
	// on the type and support by the controller.
	//
	// Examples: `1.2.3.4`, `128::1`, `my-ip-address`.
	//
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=253
	Value string `json:"value"`
}

// AddressType defines how a network address is represented as a text string.
//
// +kubebuilder:validation:Pattern=`^IPAddress|Hostname$`
type AddressType string

const (
	// A textual representation of a numeric IP address. IPv4
	// addresses must be in dotted-decimal form. IPv6 addresses
	// must be in a standard IPv6 text representation
	// (see [RFC 5952](https://tools.ietf.org/html/rfc5952)).
	//
	// This type is intended for specific addresses. Address ranges are not
	// supported (e.g. you can not use a CIDR range like 127.0.0.0/24 as an
	// IPAddress).
	IPAddressType AddressType = "IPAddress"

	// A Hostname represents a DNS based ingress point. This is similar to the
	// corresponding hostname field in Kubernetes load balancer status. For
	// example, this concept may be used for cloud load balancers where a DNS
	// name is used to expose a load balancer.
	HostnameAddressType AddressType = "Hostname"
)

// GetConditions retrieves the DataPlane Status Conditions
func (d *DataPlane) GetConditions() []metav1.Condition {
	return d.Status.Conditions
}

// SetConditions sets the DataPlane Status Conditions
func (d *DataPlane) SetConditions(conditions []metav1.Condition) {
	d.Status.Conditions = conditions
}
