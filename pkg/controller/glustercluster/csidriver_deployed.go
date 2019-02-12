package glustercluster

import (
	"github.com/gluster/anthill/pkg/reconciler"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

var glusterFuseProvisionerDeployed = reconciler.NewAction(
	"glusterFuseProvisionerDeployed",
	[]*reconciler.Action{},
	func(request reconcile.Request, client client.Client, scheme *runtime.Scheme) (reconciler.Result, error) {
		return reconciler.Result{Status: corev1.ConditionTrue, Message: "it's true"}, nil
	},
)

var glusterFuseAttacherDeployed = reconciler.NewAction(
	"glusterFuseAttachedDeployed",
	[]*reconciler.Action{},
	func(request reconcile.Request, client client.Client, scheme *runtime.Scheme) (reconciler.Result, error) {
		return reconciler.Result{Status: corev1.ConditionTrue, Message: "it's true"}, nil
	},
)

var glusterFuseNodeDeployed = reconciler.NewAction(
	"glusterFuseNodeDeployed",
	[]*reconciler.Action{},
	func(request reconcile.Request, client client.Client, scheme *runtime.Scheme) (reconciler.Result, error) {
		return reconciler.Result{Status: corev1.ConditionTrue, Message: "it's true"}, nil
	},
)

var serviceAccountDeployed = reconciler.NewAction(
	"serviceAccountDeployed",
	[]*reconciler.Action{},
	func(request reconcile.Request, client client.Client, scheme *runtime.Scheme) (reconciler.Result, error) {

		return reconciler.Result{Status: corev1.ConditionTrue, Message: "it's true"}, nil
	},
)
