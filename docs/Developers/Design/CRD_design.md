This document describes the set of [Custom Resource Definitions
(CRDs)](https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/custom-resources/)
that will be used to configure a GCS Gluster cluster. The actual implementation
of the resources described here will be phased in during development. The
purpose of this document is to provide the overall structure, ensuring the end
result provides necessary configurability in a user-friendly manner.

# Overview

A single Gluster operator may control one or more individual Gluster clusters.
Each cluster could be either hosted within the Kubernetes cluster as a set of
pods (a converged configuration) or the Gluster nodes could be running on hosts
outside the Kubernetes cluster (referred to as independent mode). The
capabilities of the operator differ significantly between these two modes of
deployment, but the same set of CRDs should be used for both where possible.

A given Gluster cluster is defined by several different Custom Resources (CRs)
that form a hierarchy. At the top level is the "cluster" CR (`GlusterCluster`)
that describes cluster-wide configuration options such as the "name" for the
cluster, the TLS credentials that will be used for securing communication, and
peer clusters for geo-replication.

Incorporated into the cluster definition are a number of node definition
templates. These describe the different configurations of nodes that the
operator can create and how those nodes are spread across failure domains. Only
nodes that use PersistentVolumes for their storage can be created via template.
Other node types must be created manually. This includes both converged nodes
that use local devices (directly accessing the `/dev/...` tree) and independent
nodes that reside on external servers.

Below the cluster definition are node definitions (`GlusterNode`) that track
the state of the individual Gluster pods. Manipulating these objects permits an
administrator to place a node into a "disabled" state for maintenance or to
decommission it entirely (by deleting the node object).

![Hierarchy of Gluster custom resources](crd_hierarchy.dot.png)

# Custom resources

This section describes the fields in each of the custom resources.

## Cluster CR

The cluster CR defines the cluster-level configuration. A commented example is
shown below:

```yaml
apiVersion: "operator.gluster.org/v1alpha1"
kind: GlusterCluster
metadata:
  # Name for the Gluster cluster that will be created by the operator
  name: my-cluster
  # CR is namespaced
  namespace: gcs
spec:
  # Cluster options allows setting "gluster vol set" options that are
  # cluster-wide (i.e. don't take a volname argument).
  clusterOptions:  # (optional)
    "cluster.halo-enabled": "yes"
  # Drivers lists the CSI drivers that should be deployed for use with this
  # cluster
  drivers:
    - gluster-fuse
    - gluster-block
  # Gluster CA to use for generating Gluster TLS keys.
  # Contains Secret w/ CA key & cert
  glusterCA:  # (optional)
    secretName: my-secret
    secretNamespace: my-ns  # default is metadata.namespace
  # Georeplication
  replication:  # (optional)
    # Credentials for using this cluster as a target
    credentials:
      secretName: my-secret
      secretNamespace: my-ns  # default is metadata.namespace
    targets:
      # Each target has a name that can be used in the StorageClass
      - name: foo
        # Addresses of node(s) in the peer cluster
        address:
          - 1.1.1.1
          - my.dns.com
        # Credentials for setting up session (ssh user & key)
        credentials:
          secretName: my-secret
          secretNamespace: my-ns  # default is metadata.namespace
  # Only PV-based nodes are built from templates
  nodeTemplates:  # (optional)
    - name: myTemplate
      # Zone is the "failure domain"
      zone: my-zone  # default is .nodeTemplates.name
      thresholds:
        nodes: 7  # may only be specified if other fields are absent
        minNodes: 3
        maxNodes: 42
        freeStorageMin: 1Ti
        freeStorageMax: 3Ti
      nodeAffinity:  # (optional)
        # https://kubernetes.io/docs/concepts/configuration/assign-pod-node/#affinity-and-anti-affinity
        # Would include "zone"-level affinity
        # Operator will overlay best-effort pod anti-affinity
        ...
      storage:
        storageClassName: my-sc
        capacity: 1Ti
status:
  # TBD operator state
  ...
```

All CRs live within the `operator.gluster.org` group and have version
`v1alpha1`. The cluster `my-cluster`, above would be contained within the `gcs`
namespace. All components of the Gluster cluster would be expected to exist
within this single namespace. The `spec` field provides the main configuration
options.

The `clusterOptions` section are Gluster options (i.e., normally manipulated
via ithe cli `gluster vol set`) that do not take a volume parameter.

The `drivers` list provides the list of CSI drivers that will be deployed by
the operator for use with this Gluster cluster.

The `glusterCA` field holds a reference to a Kubernetes Secret containing the
certificate authority `.key` and `.pem` files from which both client and server
TLS keys can be generated. These will be used to automatically configure data
encryption between the CSI driver and the Gluster bricks.

The `replication` set of parameters define the geo-replication configuration
for this cluster, optionally as both a source and target. If this cluster is to
be used as a target for replication, the `replication.credentials` field must
be supplied. This is a reference to a Secret that contains the inbound ssh user
& key information. If this cluster is used as a source, the replication targets
are presented as a list in `replication.targets`, providing a `name` for each
remote cluster, the address(es) via `address`, and the ssh credentials via the
`credentials` field.

The `nodeTemplates` list provides a set of templates that the operator can use
to automatically scale the Gluster cluster as required and to automatically
replace failed storage nodes.

Within this template, there is a `zone` tag that allows the nodes cretaed from
this template to be assigned to a specific failure domain. The default is to
have the zone name equal to the template name. These zones can then be used to
direct storage placement by referencing them in the StorageClass. Unless
otherwise specified, volumes will be created with bricks from different zones.

The `thresholds` block places limits on the amount that the operator can scale
each template up or down. Additionally, it provides thresholds to determine
when scaling should be invoked. The template can have a fixed (constant) number
of nodes by setting `nodes` to the desired value. The operator can also
dynamically size the template if, instead of setting `nodes`, the `minNodes`,
`maxNodes`, `freeStorageMin`, and `freeStorageMax` are configured. In this
case, the number of stroage nodes always remains between the min and max, and
scaling within that range is triggered based on the amount of free storage (not
assigned to a brick) exists across the nodes in that template.

Each template is likely to have a `nodeAffinity` entry to guide the placement
of the Gluster pods to a single failure domain within the cluster.

The `storage` block defines how the backing storage for the templated nodes are
created. This includes the name of a StorageClass that can be used to allocate
block-mode PVs, and the capacity that should be requested from this class.

## Node CR

The Node CR defines a single Gluster server that is a part of the cluster.
These node objects can either be created automatically by the operator from a
template or they can be manually created.

```yaml
apiVersion: "operator.gluster.org/v1alpha1"
kind: GlusterNode
metadata:
  # Name for this node
  name: az1-001
  # CRD is namespaced
  namespace: gcs
  annotations:
    # Applied by operator when it creates/manages this object from a template.
    # When this is present, contents will be dynamically adjusted accd to the
    # template in the cluster CR.
    # When this annotation is present, the admin may only modify
    # .spec.desiredState or delete the CR. Any other change will be
    # overwritten.
    anthill.gluster.org/template: template-name
spec:
  # Nodes belong to a cluster
  cluster: my-cluster
  # Nodes belong to a zone
  zone: az1
  # Admin (or operator) sets desired state for the node.
  desiredState: enabled  # (enabled | disabled)
  # Only 1 of external | storage
  external:
    address: my.host.com
    credentials:
      secretName: my-secret
      secretNamespace: my-ns  # default is metadata.namespace
  storage:
    # Only 1 of device | pvcName
    # Device names must to be stable on the host
    - device: /dev/sd[b-d]
      pvcName: my-pvc
      tags: [tag1, tag2]
  nodeAffinity:
    # https://kubernetes.io/docs/concepts/configuration/assign-pod-node/#affinity-and-anti-affinity
    # For admin created GlusterNodes, this needs to specify a node selector
    # that matches exactly one node. For template-based GNs, this will inherit
    # from the template.
    ...
status:
  # TBD operator state
  # Possible states: (enabled | deleting | disabled)
  currentState: enabled
```

There will be one Node object for each Gluster node in the cluster. By
manipulating this object, the administrator can perform a number of maintenance
actions:

- By deleting a given Node object, the administrator signals the operator that
  the corresponding Gluster node should be de-commissioned and removed from the
  cluster. In the case of a converged deployment, the resources of the
  corresponding pod would also be freed.
- By changing the `.spec.desiredState` of the Node, the administrator can notify
  the operator (and by extension, other Gluster management layers) that the
  particular node should be considered "in maintenance" (`disabled`)such that it
  could be down for an extended time and should not be used as the source or
  target of migration, nor should it be used for new data allocation.

The annotation `anthill.gluster.org/template`, when present, indicates that
this node object was created by the operator from the named template field. As
such, the operator will keep the fields of this object in-sync with the
template definition. When this annotation is present, the administrator should
not modify any field other than `.spec.desiredState`. However, the
administrator may still delete the object to signal that it should be
decommissioned. Manually created `GlusterNode` objects should not have this
annotation, and the administrator is free to modify all `.spec.*` fields.

Within the `.spec`, the `cluster` field contains the name of the Gluster
cluster to which this node belongs, and the `zone` field denotes the failure
domain zone name for this node. The `desiredState` field denotes whether this
node should be considered disabled (not used for new allocation and may be
unavailable as well).

Only one of `external` or `storage` may be present. If `external` exists, this
object represents a Gluster server that is running external to the Kubernetes
cluster. The fields in this section, provide the connection information
(`address`) to access this node. For cases where Heketi and glusterd are in
use, the `credentials` field can be used to provide authentication information
so that Heketi can ssh to the node to perform management operatons. This field
is not necessary when running glusterd2.

Converged nodes (running as pods) will have a `storage` section. This section
provides a list of either devices or PersistentVolumeClaims that will be used
by the node for creating bricks. Care must be taken when providing device names
directly that the devices correcponding to the provided names remain constant
at all times.

The `nodeAffitity` section provides the ability to limit the cluster nodes to
which this Gluster server can be assigned. In the case of specifying devices
directly, this should be used to limit the Gluster node to a single Kubernetes
node. When using a PVC, the node affinity should be such that the PVC is
accessible from the nodes that match the affinity, and the affinity should be
further restricted to comply with the desired failure domain zone tag.

# Examples

Below are some example Gluster configurations using the custom resources
defined above.

## AWS cluster, single AZ

This provides a very simple, single availability zone deployment with most
options remaining as default. The Gluster pods can be placed arbitrarity within
the cluster, and the number of nodes can be scaled as required to meet capacity
demands.

```yaml
apiVersion: "operator.gluster.org/v1alpha1"
kind: GlusterCluster
metadata:
  name: my-cluster
  namespace: gluster
spec:
  drivers:
    - gluster-fuse
  glusterCA:
    secretName: ca-secret
  nodeTemplates:
    - name: default
      thresholds:
        minNodes: 3
        maxNodes: 99
        maxStorage: 100Ti
        freeStorageMin: 500Gi
        freeStorageMax: 2Ti
      storage:
        storageClassName: ebs
        size: 1Ti
```

## AWS cluster, multi AZ

Building upon the previous single-AZ deployment is the following configuration
that uses three different AZs for Gluster pods. Here, each Zone definition
provides a unique `storageClassName` to ensure the pod's backing storage is
allocated from the correct EBS AZ, and it provides `nodeAffinity` such that the
Gluster pod will be placed in on a node that is compatible with the chosen EBS
AZ.

The zone names used here (`az1a`, `az1b`, and `az1c`) can be referenced in the
CSI driver's "data zones" list to control placement.

```yaml
apiVersion: "operator.gluster.org/v1alpha1"
kind: GlusterCluster
metadata:
  name: my-cluster
  namespace: gluster
spec:
  drivers:
    - gluster-fuse
  glusterCA:
    secretName: ca-secret
  nodeTemplates:
    - name: az1a
      thresholds:
        minNodes: 3
        maxNodes: 99
        maxStorage: 100Ti
        freeStorageMin: 500Gi
        freeStorageMax: 2Ti
      nodeAffinity:
        requiredDuringSchedulingIgnoredDuringExecution:
          nodeSelectorTerms:
            - matchExpressions:
              - key: failure-domain.beta.kubernetes.io/zone
                operator: In
                values:
                  - us-east-1a
      storage:
        storageClassName: ebs-1a
        size: 1Ti
    - name: az1b
      thresholds:
        minNodes: 3
        maxNodes: 99
        maxStorage: 100Ti
        freeStorageMin: 500Gi
        freeStorageMax: 2Ti
      nodeAffinity:
        requiredDuringSchedulingIgnoredDuringExecution:
          nodeSelectorTerms:
            - matchExpressions:
              - key: failure-domain.beta.kubernetes.io/zone
                operator: In
                values:
                  - us-east-1b
      storage:
        storageClassName: ebs-1b
        size: 1Ti
    - name: az1c
      thresholds:
        minNodes: 3
        maxNodes: 99
        maxStorage: 100Ti
        freeStorageMin: 500Gi
        freeStorageMax: 2Ti
      nodeAffinity:
        requiredDuringSchedulingIgnoredDuringExecution:
          nodeSelectorTerms:
            - matchExpressions:
              - key: failure-domain.beta.kubernetes.io/zone
                operator: In
                values:
                  - us-east-1c
      storage:
        storageClassName: ebs-1c
        size: 1Ti
```

## Bare-metal or virtualized on-prem

With an on-prem installation, it is likely that raw storage will be exposed to
nodes as direct-attached storage. This would be either as physical disks (for
bare metal) or as a statically mapped device (VMWare) or LUN. In these cases,
local block-mode PVs would be used for the storage backing the Gluster pods,
leading to template definitions very similar to above:

```yaml
apiVersion: "operator.gluster.org/v1alpha1"
kind: GlusterCluster
metadata:
  name: my-cluster
  namespace: gluster
spec:
  drivers:
    - gluster-fuse
  glusterCA:
    secretName: ca-secret
  nodeTemplates:
    - name: default
      thresholds:
        minNodes: 3
        maxNodes: 99
        maxStorage: 100Ti
        freeStorageMin: 500Gi
        freeStorageMax: 2Ti
      storage:
        storageClassName: local-pv
        size: 1Ti
```
