package glustercluster

import (
	"github.com/gluster/anthill/pkg/reconciler"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// ProcedureV1 is Procedure version 1
var ProcedureV1 = reconciler.NewProcedure(
	0,
	0,
	[]*reconciler.Action{
		EtcdClusterCreated,
		GlusterFuseProvisionerDeployed,
		GlusterFuseAttacherDeployed,
		GlusterFuseNodeDeployed,
	},
)

//Move the definitions below out to their own files when implementing them.

//GlusterFuseProvisionerDeployed deploys the GlusterFuseProvisioner
var GlusterFuseProvisionerDeployed = reconciler.NewAction(
	"GlusterFuseProvisionerDeployed",
	[]*reconciler.Action{},
	func(request reconcile.Request, client client.Client, scheme *runtime.Scheme) (reconciler.Result, error) {
		return reconciler.Result{Status: corev1.ConditionTrue, Message: "it's true"}, nil
	},
)

//GlusterFuseAttacherDeployed deploys the GlusterFuseAttacher
var GlusterFuseAttacherDeployed = reconciler.NewAction(
	"GlusterFuseAttachedDeployed",
	[]*reconciler.Action{},
	func(request reconcile.Request, client client.Client, scheme *runtime.Scheme) (reconciler.Result, error) {
		return reconciler.Result{Status: corev1.ConditionTrue, Message: "it's true"}, nil
	},
)

//GlusterFuseNodeDeployed deployes the GlusterFuseNode
var GlusterFuseNodeDeployed = reconciler.NewAction(
	"GlusterFuseNodeDeployed",
	[]*reconciler.Action{},
	func(request reconcile.Request, client client.Client, scheme *runtime.Scheme) (reconciler.Result, error) {
		return reconciler.Result{Status: corev1.ConditionTrue, Message: "it's true"}, nil
	},
)
