This document describes the actions that have to be executed successfully for
both of anthill's CRDs to be considered reconciled. They are implemented as
`reconciler.Action` objects, and are enumerated in a `reconciler.Procedure`
object. Procedure level actions are ones that modify state and should have a
corresponding entry populated in
`.Status.ReconcileActions map[string]reconciler.Result`.
Top level actions are executed in an arbitrary order so they must define any
prerequisite actions explicitly.
An action may be a top-level action and still defined as a prerequisite and the
caching implementation will ensure that it is executed a maximum of once per
Procedure execution.

# GlusterCluster actions

## etcdClusterReconciled

prereqs:

- [etcdCRDRegistered](#etcdCRDRegistered)

## etcdCRDRegistered

## csiAttacherReconciled

## csiNodePluginReconciled

## csiProvisionerReconciled

## glusterClusterServicesReconciled

prereqs:

- [glusterNodesReconciled](#glusterNodesReconciled)

## glusterNodesReconciled

prereqs:

- [etcdClusterReconciled](#etcdClusterReconciled)

## managedDevicesReconciled

- [glusterClusterServicesReconciled](#glusterClusterServicesReconciled)

# GlusterNode actions

`GlusterNode` CRs can be created manually or by the `GlusterCluster` Controller
according to a `template`.`GlusterCluster`s that consume local storage via
`hostPath` require their `GlusterNode`s to be created manually and have
`nodeAffinity` set in a way that it will only be scheduled on that node.

## statefulSetReconciled
