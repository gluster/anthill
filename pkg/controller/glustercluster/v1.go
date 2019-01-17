package glustercluster

import (
	"github.com/gluster/anthill/pkg/reconciler"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// Procedure versioin 0
var ProcedureV1 = reconciler.NewProcedure( //not happy with this name
	0,
	0,
	[]*reconciler.Action{
		EtcdClusterCreated,
		GlusterFuseProvisionerDeployed,
		GlusterFuseAttachedDeployed,
		GlusterFuseNodeDeployed,
	},
)

/* action candidates --- (need a naming convention)
//action => prereq => prereq

// I can only think of a procedure with a single action
// since procedureList doesn't preserve order.
nodePool => nodeSVC => nodeCR => nodeCRD
							  => etcdCluster => etcdCRD


*/

//procedure level actions
var EtcdClusterCreated = reconciler.NewAction(
	"EtcdClusterCreated",
	[]*reconciler.Action{
		ExamplePrereqAction,
	},
	func(request reconcile.Request, client client.Client, scheme *runtime.Scheme) (reconciler.Result, error) {
		return reconciler.Result{Status: corev1.ConditionTrue, Message: "it's true"}, nil
	},
)

var GlusterFuseProvisionerDeployed = reconciler.NewAction(
	"GlusterFuseProvisionerDeployed",
	[]*reconciler.Action{
		ExamplePrereqAction,
	},
	func(request reconcile.Request, client client.Client, scheme *runtime.Scheme) (reconciler.Result, error) {
		return reconciler.Result{Status: corev1.ConditionTrue, Message: "it's true"}, nil
	},
)

var GlusterFuseAttachedDeployed = reconciler.NewAction(
	"GlusterFuseAttachedDeployed",
	[]*reconciler.Action{
		ExamplePrereqAction,
	},
	func(request reconcile.Request, client client.Client, scheme *runtime.Scheme) (reconciler.Result, error) {
		return reconciler.Result{Status: corev1.ConditionTrue, Message: "it's true"}, nil
	},
)

var GlusterFuseNodeDeployed = reconciler.NewAction(
	"GlusterFuseNodeDeployed",
	[]*reconciler.Action{
		ExamplePrereqAction,
	},
	func(request reconcile.Request, client client.Client, scheme *runtime.Scheme) (reconciler.Result, error) {
		return reconciler.Result{Status: corev1.ConditionTrue, Message: "it's true"}, nil
	},
)

//prereq level actions
var ExamplePrereqAction = reconciler.NewAction(
	"ExamplePrereqAction",
	[]*reconciler.Action{},
	func(request reconcile.Request, client client.Client, scheme *runtime.Scheme) (reconciler.Result, error) {
		return reconciler.Result{Status: corev1.ConditionTrue, Message: "it's true"}, nil
	},
)
