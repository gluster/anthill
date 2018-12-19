package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GlusterStorageTarget defines a storage target
// used in StorageClass. Each storage class need to be
// backed by a source of storage which is defined by this
// target
type GlusterStorageTarget struct {
	Name        string       `json:"name,omitempty"`
	Addresses   []string     `json:"address"`
	Credentials *Credentials `json:"credentials,omitempty"`
}

// GlusterClusterReplicationDetails defines replication details
// for the cluster if geo-replication is used for volumes. This
// defines the target where volumes from the cluster would get
// geo replicated
type GlusterClusterReplicationDetails struct {
	Credentials *Credentials           `json:"credentials,omitempty"`
	Targets     []GlusterStorageTarget `json:"targets"`
}

// GlusterNodeThreshold defines threshold details for node
type GlusterNodeThreshold struct {
	Nodes          *int               `json:"nodes,omitempty"`
	MinNodes       *int               `json:"minNodes,omitempty"`
	MaxNodes       *int               `json:"maxNodes,omitempty"`
	FreeStorageMin *resource.Quantity `json:"freeStorageMin,omitempty"`
	FreeStorageMax *resource.Quantity `json:"freeStorageMax,omitempty"`
}

// GlusterNodeStorageDetails defines storage class details
type GlusterNodeStorageDetails struct {
	StorageClassName string             `json:"storageClassName,omitempty"`
	Capacity         *resource.Quantity `json:"capacity,omitempty"`
}

// GlusterNodeTemplate defines a gluster node's template
type GlusterNodeTemplate struct {
	Name      string                     `json:"name,omitempty"`
	Zone      string                     `json:"zone,omitempty"`
	Threshold *GlusterNodeThreshold      `json:"threshold,omitempty"`
	Affinity  *corev1.NodeAffinity       `json:"nodeAffinity,omitempty"`
	Storage   *GlusterNodeStorageDetails `json:"storage,omitempty"`
}

// GlusterClusterSpec defines the desired state of GlusterCluster
type GlusterClusterSpec struct {
	Options       map[string]string                 `json:"clusterOptions,omitempty"`
	Drivers       []string                          `json:"drivers"`
	GlusterCA     *Credentials                      `json:"glusterCA,omitempty"`
	Replication   *GlusterClusterReplicationDetails `json:"replication,omitempty"`
	NodeTemplates []GlusterNodeTemplate             `json:"nodeTemplates"`
}

// GlusterClusterStatus defines the observed state of GlusterCluster
type GlusterClusterStatus struct {
	State string `json:"state,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// GlusterCluster is the Schema for the glusterclusters API
// +k8s:openapi-gen=true
type GlusterCluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   GlusterClusterSpec   `json:"spec,omitempty"`
	Status GlusterClusterStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// GlusterClusterList contains a list of GlusterCluster
type GlusterClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []GlusterCluster `json:"items"`
}

func init() {
	SchemeBuilder.Register(&GlusterCluster{}, &GlusterClusterList{})
}
