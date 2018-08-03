The operator is responsible for incorporating the domain-specific expertise
necessary to perform upgrades of the Gluster cluster. An administrator, either
directly or through an upgrade service, will choose the desired version of
Gluster (GCS) that should be running on the cluster, and the operator must
properly and autonomously upgrade all components to match that request.

# Problem overview

A version of GCS is comprised of a specific set of container images plus the
associated object manifests required to deploy them. In order to upgrade, the
container images must be changed to those in a subsequent release. These
updates, however, must be mindful of version compatibility between the
components. Particularly when upgrading across a series of releases, it may be
necessary to take intermediate upgrade steps to ensure client (CSI driver)
versions remain compatible with server (Gluster container) versions. Such
dependencies may also arise from other components such as dashboards and
alerting.

## Interaction with Operator Lifecycle Manager

[OLM](https://github.com/operator-framework/operator-lifecycle-manager) is
designed to manage the life cycle of operators (i.e., ensure dependencies,
install, upgrade, etc.). It accomplishes this by maintaining a set of
ClusterServiceVersion objects that describe the available versions of an
operator. These objects are then assembled to create an upgrade path for a
given operator.

OLM attempts to install the latest version of an operator by walking the CSV
versions. It does this by replacing the running operator with the next version,
waiting for it to become ready, then repeating until the operator is
up-to-date. This rapid-succession replacement does not have a mechanism to
pause on a version while underlying services are updated. The implication of
this is:

- The operator is required to support upgrading its internal state from version
  *n-1* to version *n* only. (OLM walks through each version.)
- Every version of the operator must be capable of upgrading the Gluster
  cluster from an arbitrarily old version to the current version. (OLM will not
  wait for lower-level upgrades.)
  > One option considered was to force OLM to pause upgrading the operator by
  > not returning ready until the GCS deployment is fully upgraded. There were
  > several concerns with this approach, including:
  >
  > 1. The timeout for the readiness probe would need to be set arbitrarily
  >    high to permit time for upgrades to complete.
  > 1. While the operator is "unready" and Services would not direct traffic
  >    to it, potentially interfering with reporting metrics.
  > 1. In non-OLM cases, admins would be responsible for installing the correct
  >    sequence of operator versions to properly upgrade their cluster.

Because of the above issues, we must have a method for representing upgrade and
deployment actions that can be maintained indefinitely as a part of the
operator's logic.

# Proposed solution

In each version of the operator, the "reconcile" actions determine what
ultimately gets deployed and how it is configured in the cluster. This makes the
main reconcile operation a good place to apply versioning.

Each versioned reconcile operation consists of a tuple: `minVersion,
currentVersion, reconcile()`. These fields are:

- `minVersion`: The minimum deployed configuration version that this
  `reconcile()` can properly upgrade from.
- `currentVersion`: The configuration version that this `reconcile()` will
  apply.
- `reconcile()`: The entry point for reconciling, (i.e., an implementation of
  `func (r *ReconcileGlusterCluster) Reconcile(request reconcile.Request)
  (reconcile.Result, error)` or `func (r *ReconcileGlusterNode)
  Reconcile(request reconcile.Request) (reconcile.Result, error)`)

At the start of a reconcile cycle, the current deployed version is queried from
the status field of the `GlusterCluster` and used to search for a compatible
reconcile function. A compatible function is located by traversing the list of
tuples and looking for the reconcile function with the highest `currentVersion`
whose `minVersion` is less than or equal to the configuration version that is
currently deployed.

**Note:** *A new deployment will not have a `status.currentVersion` until it has
*been fully reconciled once. In its absence, we know this is a new deployment
*and it should use the most recent reconcile function (or
*GlusterCluster.status.currentVersion in the case of the GlusterNode
*reconciler).*

When a reconcile loop is completed that finds there are no changes necessary
(i.e., the deployment is fully reconciled), it updates `status.currentVersion`
to match the chosen reconciler's `currentVersion`. This ensures a given
reconciler is allowed to converge the state to a known configuration before the
next higher versioned reconciler is invoked. Care must be taken to ensure the
`GlusterCluster` `status.currentVersion` is only advanced once the `GlusterNode`
objects have also converged. To accomplish this, `status.currentVersion` should
be applied by the `GlusterNode` reconciler once it has converged, but the choice
of which reconcile function to use should be based on the
`GlusterCluster.status.currentVersion` field. This allows the operator to mark
the lower-level `GlusterNode` as "fully reconciled" but to avoid advancing its
configuration ahead of the rest of the cluster. The `GlusterCluster` reconciler
will only advance its version when the cluster-level items have converged and
all of the child `GlusterNode` objects have as well.

## Operator evolution

With the above process, new versions of the operator can be released by either:

1. Modifying the most recent reconcile function(s)
1. Adding a new set of reconcile function(s) with a higher version and leaving
   the older ones intact.

The choice of the proper approach is to consider operator upgrade. Consider a
scenario in which the current reconcile function is v3 and it can upgrade from
v2. The v3 reconciler could be modified directly (and changed to v4) if it can
directly upgrade from either v2 or v3. The result would be a v4 reconciler with
a `minVersion` of v2, and there would me no reconciler of version v3 remaining.
A good example of this scenario would be changing the memory request for a GD2
pod. This would result in a change to the reconciler (to apply the new request
to the StatefulSet), but it would not affect any upgrade compatibility.

A situation requiring a new reconciler function would be one in which the
container version of GD2 advanced beyond the test matrix for the client version
that is deployed in v2. In this case, it would be unsafe to jump directly from
v2 to v4 by editing the reconciler in place. In this case, v3 must be left
unchanged to upgrade all components prior to the move to the new GD2 version
deployed by reconciler v4.
