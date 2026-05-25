package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Appplication is a spec of a simple application
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type Application struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata" protobuf:"bytes,1,opt,name=metadata"`
	Spec              ApplicationSpec   `json:"spec" protobuf:"bytes,2,opt,name=spec"` // Common: shared with Application (different type)
	Status            ApplicationStatus `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

type ApplicationSpec struct {
	Name  string `json:"name" protobuf:"bytes,1,opt,name=name"`
	Image string `json:"image" protobuf:"bytes,1,opt,name=image"`
}

type ApplicationStatus struct {
	Health HealthStatus `json:"health,omitempty" protobuf:"bytes,1,opt,name=health"`
}
