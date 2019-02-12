package v1alpha1

import (
	"github.com/gluster/anthill/pkg/reconciler"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Credentials defines gluster node secret credentials
type Credentials struct {
	SecretName      string `json:"secretName,omitempty"`
	SecretNamespace string `json:"secreteNamespace,omitempty"`
}

// GlusterNodeExternal defines external details of gluster nodes
type GlusterNodeExternal struct {
	Address     string       `json:"address,omitempty"`
	Credentials *Credentials `json:"credentials,omitempty"`
}

// StorageDevice defines storage details of gluster nodes
type StorageDevice struct {
	Device  string   `json:"device,omitempty"`
	PVCName string   `json:"pvcName,omitempty"`
	Tags    []string `json:"tags"`
}

// GlusterNodeSpec defines the desired state of GlusterNode
type GlusterNodeSpec struct {
	Cluster      string               `json:"cluster,omitempty"`
	Zone         string               `json:"zone,omitempty"`
	DesiredState string               `json:"desiredState,omitempty"`
	ExternalInfo *GlusterNodeExternal `json:"external,omitempty"`
	Storage      []StorageDevice      `json:"storage"`
	Affinity     *corev1.NodeAffinity `json:"nodeAffinity,omitempty"`
}

// GlusterNodeStatus defines the observed state of GlusterNode
type GlusterNodeStatus struct {
	State            string                       `json:"currentState,omitempty"`
	ReconcileVersion *int                         `json:"reconcileVersion,omitempty"`
	ReconcileActions map[string]reconciler.Result `json:"reconcileActions,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// GlusterNode is the Schema for the glusternodes API
// +k8s:openapi-gen=true
type GlusterNode struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   GlusterNodeSpec   `json:"spec,omitempty"`
	Status GlusterNodeStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// GlusterNodeList contains a list of GlusterNode
type GlusterNodeList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []GlusterNode `json:"items"`
}

func init() {
	SchemeBuilder.Register(&GlusterNode{}, &GlusterNodeList{})
}
