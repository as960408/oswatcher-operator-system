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

// OSStatusSpec defines the desired state of OSStatus.

type OSStatusSpec struct {
	NodeName    string        `json:"nodeName,omitempty"`
	NodeIP      string        `json:"nodeIP,omitempty"`
	CPUUsage    string        `json:"cpuUsage,omitempty"`
	MemUsage    string        `json:"memUsage,omitempty"`
	RootUsage   string        `json:"rootUsage,omitempty"`
	Uptime      string        `json:"uptime,omitempty"`
	TopCPUProcs []ProcessInfo `json:"topCPUProcs,omitempty"`
	TopMemProcs []ProcessInfo `json:"topMemProcs,omitempty"`
	CollectedAt string        `json:"collectedAt,omitempty"`

	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of OSStatus. Edit osstatus_types.go to remove/update
	//Foo string `json:"foo,omitempty"`
}

// OSStatusStatus defines the observed state of OSStatus.

type ProcessInfo struct {
	PID     string `json:"pid,omitempty"`
	User    string `json:"user,omitempty"`
	Command string `json:"command,omitempty"`
	CPU     string `json:"cpu,omitempty"` // CPU 점유율
	Mem     string `json:"mem,omitempty"` // 메모리 점유율
}

// OSStatus is the Schema for the osstatuses API.

// +kubebuilder:printcolumn:name="CPU",type=string,JSONPath=`.spec.cpuUsage`
// +kubebuilder:printcolumn:name="Mem",type=string,JSONPath=`.spec.memUsage`
// +kubebuilder:printcolumn:name="RootUsage",type=string,JSONPath=`.spec.rootUsage`
// +kubebuilder:printcolumn:name="TopCPUProcess",type=string,JSONPath=`.spec.topCPUProcs[0].command`
// +kubebuilder:printcolumn:name="TopMemProcess",type=string,JSONPath=`.spec.topMemProcs[0].command`
// +kubebuilder:printcolumn:name="Rootdisk",type=string,JSONPath=`.spec.rootUsage`
// +kubebuilder:printcolumn:name="UPTIME",type=string,JSONPath=`.spec.uptime`
// +kubebuilder:printcolumn:name="CollectedAt",type=string,JSONPath=`.spec.collectedAt`
// +kubebuilder:object:root=true
type OSStatus struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec OSStatusSpec `json:"spec,omitempty"`
	// Status OSStatusStatus `json:"status,omitempty"`
}

// OSStatusStatus defines the observed state of OSStatus.

// +kubebuilder:object:root=true

// OSStatusList contains a list of OSStatus.
type OSStatusList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []OSStatus `json:"items"`
}

func init() {
	SchemeBuilder.Register(&OSStatus{}, &OSStatusList{})
}
