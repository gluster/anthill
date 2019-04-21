This document provides a design for how to compose the reconcile operations
necessary to configure and maintain the Gluster cluster.

# Background

The actions of the operator are guided by reconciling a desired state
(communicated via GlusterCluster and GlusterNode CRs) with the objects deployed
in the kubernetes cluster. Many of the actions necessary to perform this
reconciliation are multi-step operations that have a set of restrictions on when
they can be performed. For example, the CSI driver can only be upgraded on a
node when there are no pods using Gluster storage on that node. Likewise, a GD2
pod can only be upgraded if the volumes that it hosts are fully healed.

This set of restrictions on when reconcile actions can be performed can easily
lead to very complicated reconciliation logic and repeated testing of somewhat
expensive conditions (Are volumes healed?).

# Proposal

To better express the structure of the reconcile actions, they can be expressed
as:

- A list of prerequisites that must be met in order to perform the operation.
- The function that performs the reconcile operation.

An example of the above could be a hypothetical `NodeTaggedWithZone` action that
ensures the GD2 node is tagged with the proper topology zone that it is a member
of. This would have prequisites of: `NodeUp` (The GD2 pod is up) and
`NodeInCluster` (The GD2 pod is a member of the Gluster cluster). Only when
these two reconcile completely can we reconcile the node tag.

The full set of actions required to reconcile the desired state would be the
list of prereq+operaton tuples.

## Prerequisites are reconcile actions

In the above discussion, prerequisites were presented as being distinct from
reconcile actions (i.e., they represent true/false conditions as opposed to the
actions that "do something"). However, the prereqs and the actions can actually
be one and the same (an action can act as a prereq).

An action can viewed as a function that returns a tuple: `{Condition,
Error}`. The meaning of each:

- `Condition`: This is a kubernetes `Condition` type, having the possible values
  of `True`, `False`, and `Unknown`.
- `Error`: Non-nil if an error occurred.

The actions must return `True` only if:

- The objective is met (e.g., `NodeInCluster` returns `True` iff the node is
  actually a member of the cluster).
- No action was taken during execution (i.e., the node was determined to
  *already be a member of the cluster* when the action was invoked).

The actions should return `False` if either the condition is untrue or the
function made a change to the state during its execution.

The result would be `Unknown` if the result could not be determined (i.e., one
of its prereqs were not satisfied, so it could not be executed).

Following along the `NodeInCluster` example, when used as a prereq, we only get
a positive result if we are sure the node is currently a member of the cluster,
so that result can serve to gate operations that require the node to be a
member. Likewise, when it is invoked, if if finds that the node is not a member,
it can (potentially) adjust the system state to make the node a member (thus
acting as an action).

# Potential prereqs/actions

Below is a list of potential actions/prereqs. The items in this list are likely
neither necessary nor sufficient.

## Cluster-level

- HasSufficientCapacity - The cluster has sufficient capacity
- NoExcessCapacity - The cluster does not have too much capacity
- HasSufficientNodes - The cluster has enough nodes
- NoExcessNodes - The cluster does not have too many nodes
- ClusterHasFinalizer - The GlusterCluster object has a finalizer to prevent
  deletion
- AllVolumesHealed - All volumes in the cluster are fully healed
- HasValidDisruptionBudget - The PodDisruptionBudget is configured correctly
- HasEtcdCluster - The cluster has a functioning etcd cluster to use

## Node-level

- DevicesRegistered - The devices attached to the gluster pod are properly
  registered w/ GD2 for use creating bricks
- NodeIsUp - The GD2 node is up
- NodeInCluster - The GD2 node is a member of the cluster
- HasFinalizer - The GlusterNode object has a finalizer to prevent deletion
- NodeVolumesHealed - All volumes with a brick on this node are healed
- NodeIsAbandoned - This node has been marked abandoned by GD2 (see discussion
  of the state machine)
- NodeUsesPVC - This node gets its storage via a block-mode PVC
- NodeUsesDevice - This node gets its storage via a device in /dev
- HasDeployment - This node object has a valid Deployment/StatefulSet object

# Reporting

By providing a standard structure to the reconcile actions, the current state of
the system can be exposed via the `status:` field of the CRs. For example,
`status.conditions` for the `GlusterCluster` CR could be an array of the
condition results in the same format as other kube objects.

For illustration of how conditions are typically reported, the following were
copied from a kubernetes node:

```yaml
conditions:
  - lastHeartbeatTime: 2018-12-19T00:01:42Z
    lastTransitionTime: 2018-11-27T05:21:19Z
    message: kernel has no deadlock
    reason: KernelHasNoDeadlock
    status: "False"
    type: KernelDeadlock
  - lastHeartbeatTime: 2018-12-19T00:01:51Z
    lastTransitionTime: 2018-11-27T05:20:39Z
    message: kubelet has sufficient disk space available
    reason: KubeletHasSufficientDisk
    status: "False"
    type: OutOfDisk
  - lastHeartbeatTime: 2018-12-19T00:01:51Z
    lastTransitionTime: 2018-11-27T05:20:39Z
    message: kubelet has sufficient memory available
    reason: KubeletHasSufficientMemory
    status: "False"
    type: MemoryPressure
  - lastHeartbeatTime: 2018-12-19T00:01:51Z
    lastTransitionTime: 2018-11-27T05:20:39Z
    message: kubelet has no disk pressure
    reason: KubeletHasNoDiskPressure
    status: "False"
    type: DiskPressure
  - lastHeartbeatTime: 2018-12-19T00:01:51Z
    lastTransitionTime: 2018-12-12T17:20:30Z
    message: kubelet is posting ready status
    reason: KubeletReady
    status: "True"
    type: Ready
  - lastHeartbeatTime: 2018-12-19T00:01:51Z
    lastTransitionTime: 2018-11-02T18:45:34Z
    message: kubelet has sufficient PID available
    reason: KubeletHasSufficientPID
    status: "False"
    type: PIDPressure
```

The idea is to report the status of each reconcile action via one of these
condition entries. This will provide some explanation of how the operator is
doing at managing the cluster. The `type` and `status` fields come directly from
the reconcile action, but the origin of `reason` and `message` are still TBD.

# Caching

The actions are likely to be invoked multiple times during each reconcile
iteration-- Once to reconcile (as an action) and potentially several times when
used as a prereq for other actions. However, it is not desirable to run the
actions multiple times per iteration as the system state is unlikely to have
converged in such a short amount of time.

To avoid invoking actions multiple times during each cycle, the result
(condition, error) should be cached for the duration of the cycle. This ensures
each action is called once, when first referenced, and that result can then be
used for the rest of the iteration. This should significantly help the
performance of expensive operations such as checking pending heals.
