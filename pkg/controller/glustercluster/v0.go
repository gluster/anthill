package glustercluster

var v0 = Procedure{
	version:    0,
	minVersion: 0,
	actions:    []*Action{&etcdExposed,&glusterNodesUp}
}

//procedure level actions

var etcdExposed = Action{
	Name: "etcdExposed",
	action: func(request reconcile.Request, client client.Client, scheme *runtime.Scheme) (Result, error) {
		return Result{Status: corev1.ConditionTrue, Message: "it's true"}, nil
	},
	prereqs: []Action{ 
		etcdOperatorRunning,
		etcdCRAccurate,
		etcdSvcUp
	}
}
var glusterNodesUp = Action{
	Name: "glusterNodesUp",
	action: func(request reconcile.Request, client client.Client, scheme *runtime.Scheme) (Result, error) {
		return Result{Status: corev1.ConditionTrue, Message: "it's true"}, nil
	},
	prereqs: []Action{ 
		nodeResourcesUp,
		nodeEnpointsUp,
		clusterServiceUp,
	}
}

//prereqs
var resourceFound = Action{
	Name: "resourceFound",
	action: func(request reconcile.Request, client client.Client, scheme *runtime.Scheme) (Result, error) {
			return Result{Status: corev1.ConditionTrue, Message: "it's true"}, nil
	}
}
var clusterServiceUp = Action{
	Name: "etcdExposed",
	action: func(request reconcile.Request, client client.Client, scheme *runtime.Scheme) (Result, error) {
		return Result{Status: corev1.ConditionTrue, Message: "it's true"}, nil
	},
}
var etcdExposed = Action{
	Name: "etcdExposed",
	action: func(request reconcile.Request, client client.Client, scheme *runtime.Scheme) (Result, error) {
		return Result{Status: corev1.ConditionTrue, Message: "it's true"}, nil
	},
}