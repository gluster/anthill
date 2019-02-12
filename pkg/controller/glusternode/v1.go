package glusternode

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
		etcdEndpointValid,
		statefullSetCreated,
	},
)

var etcdEndpointValid = reconciler.NewAction(
	"etcdEndpointValid",
	[]*reconciler.Action{},
	func(request reconcile.Request, client client.Client, scheme *runtime.Scheme) (reconciler.Result, error) {
		return reconciler.Result{Status: corev1.ConditionTrue, Message: "it's true"}, nil
	},
)
var statefullSetCreated = reconciler.NewAction(
	"statefullSetCreated",
	[]*reconciler.Action{etcdEndpointValid},
	func(request reconcile.Request, client client.Client, scheme *runtime.Scheme) (reconciler.Result, error) {
		return reconciler.Result{Status: corev1.ConditionTrue, Message: "it's true"}, nil
	},
)
