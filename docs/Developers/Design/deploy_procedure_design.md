This document describes the actions that have to be executed successfully for both of anthill's CRDs to be consider reconciled.  They are implemented as `reconciler.Action` objects, and are enumerated in a `reconciler.Procedure` object. Procedure level actions are ones that modify state and should have a corresponding entry populated in `.Status.ReconcileActions map[string]reconciler.Result`. Top level actions are executed in an arbitrary order so they must define any prerequisite actions explicitly. An action may be a top-level action and still defined as a prerequisite and the caching implementation will ensure that it is executed only the necessary amount of times(once?).


## GlusterCluster actions

### etcdClusterReconciled
prereqs:
 -  [etcdCRDRegistered](#etcdCRDRegistered)

### etcdCRDReistered

### csiAttacherReconciled

### csiNodePluginReconciled

### csiProvisionerReconciled
### glusterClusterServicesReconciled
prereqs:
 -  [glusterNodesReconciled](#glusterNodesReconciled)
### glusterNodesReconciled
prereqs:
 -  [etcdClusterReconciled](#etcdClusterReconciled)


## GlusterNode actions

`GlusterNode` CRs can be created manually or by the `GlusterCluster` Controller according to a `template`. The node is associated with the cluster using the `gluster.org/cluster-name` label that they will both share. `GlusterCluster`s that consume local storage via `hostPath` require their `GlusterNode`s to be created manually and have `nodeAffinity` set in a way that it will only be scheduled on that node.

### statefulSetReconciled
prereqs:
 -  [managedDevicesReconciled](#managedDevicesReconciled)
### managedDevicesReconciled
prereqs:
 -  [glusterClusterServicesReconciled](#glusterClusterServicesReconciled)

