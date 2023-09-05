/*
Copyright The Kubernetes Authors.

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

package v1beta2

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// AWSNetworkSpec defines the desired state of AWSNetwork
type AWSNetworkSpec struct {
	// TODO
	// Foo is an example field of AWSNetwork. Edit awsnetwork_types.go to remove/update
	Foo string `json:"foo,omitempty"`
}

// AWSNetworkStatus defines the observed state of AWSNetwork
type AWSNetworkStatus struct {
	// TODO
}

//+kubebuilder:object:root=true
//+kubebuilder:resource:path=awsnetworks,scope=Namespaced,categories=cluster-api,shortName=awsn
//+kubebuilder:subresource:status

// AWSNetwork is the Schema for the awsnetworks API
type AWSNetwork struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AWSNetworkSpec   `json:"spec,omitempty"`
	Status AWSNetworkStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// AWSNetworkList contains a list of AWSNetwork
type AWSNetworkList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AWSNetwork `json:"items"`
}

func init() {
	SchemeBuilder.Register(&AWSNetwork{}, &AWSNetworkList{})
}
