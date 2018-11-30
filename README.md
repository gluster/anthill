# Gluster operator for Kubernetes and OpenShift

[![Build Status](https://travis-ci.org/gluster/anthill.svg?branch=master)](https://travis-ci.org/gluster/anthill)
[![Documentation Status](https://readthedocs.org/projects/gluster-anthill/badge/?version=latest)](http://gluster-anthill.readthedocs.io/)
<!-- Badges: TravisCI, CentOS CI, Coveralls, GoDoc, GoReport, ReadTheDocs -->

**Found a bug?** [Let us know.](https://github.com/gluster/operator/issues/new?template=bug_report.md)

**Have a request?** [Tell us about it.](https://github.com/gluster/operator/issues/new?template=feature_request.md)

**Interested in helping out?** Take a look at the [contributing
doc](CONTRIBUTING.md) to find out how.

## Build

The operator is based on the [Operator
SDK](https://github.com/operator-framework/operator-sdk). In order to build the
operator, you first need to install the SDK. [Instructions are
here.](https://github.com/operator-framework/operator-sdk#quick-start)

Once the SDK is installed, Anthill can be built via:

```bash
$ dep ensure

$ operator-sdk build docker.io/gluster/anthill
$ operator-sdk build docker.io/gluster/anthill
INFO[0000] Building Docker image docker.io/gluster/anthill
Sending build context to Docker daemon  152.1MB
Step 1/3 : FROM docker.io/centos:7.5.1804
7.5.1804: Pulling from library/centos
Digest: sha256:adbf1369c8dbeaa2bf9df0b891394e3ca09b5272a2fb54e0fadc6b14dd93fcad
Status: Downloaded newer image for centos:7.5.1804
 ---> 76d6bc25b8a5
Step 2/3 : USER nobody
 ---> Using cache
 ---> dc1dd738f9a4
Step 3/3 : ADD build/_output/bin/anthill /usr/local/bin/anthill
 ---> Using cache
 ---> cffc3a0ef59e
Successfully built cffc3a0ef59e
Successfully tagged gluster/anthill:latest
INFO[0002] Operator build complete.
```

## Installation

Install the CRDs into the cluster:

```bash
$ kubectl apply -f deploy/crds/operator_v1alpha1_glustercluster_crd.yaml
customresourcedefinition.apiextensions.k8s.io "glusterclusters.operator.gluster.org" created
```

Install the service account, role, and rolebinding:

```bash
$ kubectl apply -f deploy/service_account.yaml
serviceaccount "anthill" created

$ kubectl apply -f deploy/role.yaml
role.rbac.authorization.k8s.io "anthill" created
rolebinding.rbac.authorization.k8s.io "anthill" created
```

There are two options for deploying the operator.

1. It can be run normally, inside the cluster. For this, see
   `deploy/operator.yaml` for a skeleton.
1. It can also be run outside the cluster for development purposes. This
   removes the need to push the container to a registry by running the operator
   executable locally. For this:

   ```bash
   $ OPERATOR_NAME=anthill operator-sdk up local --namespace=default
   INFO[0000] Running the operator locally.
   {"level":"info","ts":1542396040.2412076,"logger":"cmd","caller":"manager/main.go:57","msg":"Registering Components."}
   {"level":"info","ts":1542396040.2413611,"logger":"kubebuilder.controller","caller":"controller/controller.go:120","msg":"Starting EventSource","Controller":"glustercluster-controller","Source":"kind source: /, Kind="}
   ...
   ```
