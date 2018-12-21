This document covers writing tests for the operator.

# Unit tests

We need to figure out something here...

# E2E testing

The [operator-sdk](https://github.com/operator-framework/operator-sdk) supports
running end-to-end tests of the operator. The framework provides an overview of
how to write tests, and the tests for Anthill can be found in the `tests/e2e`
directory.

Just like running the operator, the tests can be executed by running the
operator either in the cluster or by running it locally. For testing out
changes, running locally is generally the easiest approach.

Assuming you have a `kubeconfig` file to access your cluster:

Start by creating a namespace for testing:

```
$ kubectl --kubeconfig=kubeconfig create ns anthill-e2e
namespace "anthill-e2e" created
```

Then run the e2e tests from the opt level of the repo. It is necessary to
provide the target namespace and the kubeconfig file.

```
$ operator-sdk test local ./test/e2e --up-local --namespace anthill-e2e --kubeconfig=kubeconfig
INFO[0000] Testing operator locally.
ok      github.com/gluster/anthill/test/e2e     0.129s
INFO[0001] Local operator test successfully completed.
```

The above command runs the operator executable on the local machine instead of
as a deployment in the cluster. This tends to provide quicker development
iteration than pushing a build container then pulling and running it in the
cluster.
