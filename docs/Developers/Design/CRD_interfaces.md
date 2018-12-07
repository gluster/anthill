# Overview

This document describes the functional interfaces for the
[GlusterCluster,
GlusterNodeTemplate](../pkg/apis/operator/v1alpha1/glustercluster_types.go),
and [GlusterNode](../pkg/apis/operator/v1alpha1/glusternode_types.go) structs.
These interfaces will be used by the Operator to implement any maintenance
tasks deemed necessary by a comparison between a resource's desired
configuration and it's current working state.

## Operator Workflow

[[DIAGRAM]]

# GlusterCluster

A GlusterCluster represents a single Gluster cluster being monitored by the
Anthill Operator. Functions of this struct represent actions that apply across
all nodes or volumes in a cluster.

## Interface

```go
func (*gc GlusterCluster) GetNodes() (*GlusterNodeList, error)
func (*gc GlusterCluster) Reconcile() error
func (*gc GlusterCluster) ValidateSpec() error
```

### GetNodes

Returns a `GlusterNodeList` of all nodes in the cluster.

### Reconcile

Reconciles the state of the cluster.

1. Ensure all appropriate `drivers` are deployed in the same namespace
1. Reconcile any GlusterNodeTemplates

### ValidateSpec

Verifies that the spec provided defines a valid configuration. A spec should
pass the following checks:

* `drivers`: verify valid drivers
* `glusterCA`: if specified, verify secret exists
* `replication`: TODO

# GlusterNodeTemplate

A GlusterNodeTemplate represents a group of identical nodes that are part of the same Gluster cluster.

## Interface

```go
func (*gnt GlusterNodeTemplate) GetFreeSpace() (*resource.Quantity, error)
func (*gnt GlusterNodeTemplate) GetNodes() (*GlusterNodeList, error)
func (*gnt GlusterNodeTemplate) Reconcile() error
func (*gnt GlusterNodeTemplate) ValidateSpec() error
```

### GetFreeSpace

Returns a `resource.Quantity` of all available storage capacity in the node
group. If any node cannot report it's available capacity, return an error.

### GetNodes

Returns a `GlusterNodeList` of all nodes that match the template name.

### Reconcile

Reconciles the state of the node template.

1. Ensure there are at least `nodes` or `minNodes` number of nodes created. If
   not, create new GlusterNodes to meet this minimum.
1. If `nodes` is not specified, check available storage capacity in the node
   group
    * If neither `minNodes` nor `freeStorageMin` are specified, create new
      GlusterNodes until either `maxNodes` or `freeStorageMax` is met
    * If free space is less than `freeStorageMin`, create new GlusterNodes to
      meet that threshold up to `maxNodes` (if specified)
    * If free space is greater than `freeStorageMax`, delete GlusterNodes to
      meet that threshold down to `minNodes` (if specified)

### ValidateSpec

Verifies that the spec provided defines a valid configuration. A spec should
pass the following checks:

* `threshold`:
  * If `nodes` is specified, all other fields should be empty
  * If any other fields are specified, `nodes` should be empty
  * If `minNodes` and `maxNodes` are specified, `minNodes` < `maxNodes`
  * If `freeStorageMin` and `freeStorageMax` are specified,
    `freeStorageMin` < `freeStorageMax`
* `nodeAffinity`: verify it's a valid `NodeAffinity`
* `storage`: verify StorageClass exists and `capacity` is specified

# GlusterNode

A GlusterNode represents a node in a GlusterCluster. Functions of this struct
represent actions that apply only to a specific node.

## Interface

```go
func (*gn GlusterNode) Disable() error
func (*gn GlusterNode) Enable() error
func (*gn GlusterNode) GetFreeSpace() (*resource.Quantity, error)
func (*gn GlusterNode) Reconcile() error
func (*gn GlusterNode) ReconcileStorageDevices() error
func (*gn GlusterNode) ValidateSpec() error
```

### Disable

Marks a node not ready for modifications or storage allocation

### Enable

Marks a node ready for modifications and storage allocation

### GetFreeSpace

Returns a `resource.Quantity` of all available storage capacity on the node.

### Reconcile

Reconciles the state of the node.

If `external` is specified, ensure connectivity to the node.

If `external` is not specified:

1. Ensure the correct `desiredState` is applied
1. Ensure any required PVCs are available in the same namespace
1. Ensure a pod is running and ready for the node in the same namespace
1. Ensure the node has joined the cluster

### ReconcileStorageDevices

For each known storage device:

1. Ensure device is available
1. Ensure device is part of the storage pool

### ValidateSpec

Verifies that the spec provided defines a valid configuration. A spec should
pass the following checks:

* `cluster`: verify GlusterCluster exists in the same namespace
* `desiredState`: valid state, enabled or disabled
* Only one of `external` or `storage` is specified
* `external`: credentials Secret exists
* `storage`: only one of `device` or `pvcName` exists
