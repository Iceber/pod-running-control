/*
Copyright 2026.

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

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true

// PodRunningGate is the Schema for the podrunninggates API
type PodRunningGate struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is a standard object metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitzero"`

	// spec defines the desired state of PodRunningGate
	// +required
	Spec PodRunningGateSpec `json:"spec"`
}

// PodRunningGateSpec defines the desired state of PodRunningGate
type PodRunningGateSpec struct {
	// +optional
	Gates []string `json:"gates,omitempty"`
}

// PodRunningGateStatus defines the observed state of PodRunningGate.
type PodRunningGateStatus struct {
}

// +kubebuilder:object:root=true

// PodRunningGateList contains a list of PodRunningGate
type PodRunningGateList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitzero"`
	Items           []PodRunningGate `json:"items"`
}
