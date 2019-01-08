package v0

import (
	"github.com/gluster/anthill/pkg/reconciler"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

var V0Procedure = reconciler.Procedure{ //not happy with this name
	minVersion: 0,
	version:    0,
	actions:    []*reconciler.Action{&exampleAction},
}

/* action candidates --- (need a naming convention)
//action => prereq => prereq

// I can only think of a procedure with a single action
// since procedureList doesn't preserve order.
nodePool => nodeSVC => nodeCR => nodeCRD
							  => etcdCluster => etcdCRD


*/

//procedure level actions

var ExampleAction = reconciler.Action{
	Name: "etcdExposed",
	action: func(request reconcile.Request, client client.Client, scheme *runtime.Scheme) (Result, error) {
		return Result{Status: corev1.ConditionTrue, Message: "it's true"}, nil
	},
	prereqs: []Action{
		examplePrereqAction,
	},
}

//prereq level actions
var ExamplePrereqAction = reconciler.Action{
	Name: "resourceFound",
	action: func(request reconcile.Request, client client.Client, scheme *runtime.Scheme) (Result, error) {
		return Result{Status: corev1.ConditionTrue, Message: "it's true"}, nil
	},
}
