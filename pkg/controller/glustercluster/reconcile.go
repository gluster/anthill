package glustercluster

import (
	"context"
	"strconv"

	"k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"

	operatorv1alpha1 "github.com/gluster/anthill/pkg/apis/operator/v1alpha1"

	"github.com/gluster/anthill/pkg/reconciler"
)

var (
	log                                         = logf.Log.WithName("controller_glustercluster")
	allProcedures      reconciler.ProcedureList = []reconciler.Procedure{*ProcedureV1}
	reconcileProcedure *reconciler.Procedure
)

// Reconcile reads that state of the cluster for a GlusterCluster object and makes changes based on the state read
// and what is in the GlusterCluster.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  This example creates
// a Pod as an example
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.

func (r *ReconcileGlusterCluster) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling GlusterCluster")

	// Fetch the GlusterCluster instance
	instance := &operatorv1alpha1.GlusterCluster{}
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	// Get current reconcile version from CR
	version := instance.Status.ReconcileVersion
	if version == 0 {
		// choose the highest compatible version
		reconcileProcedure, _ = allProcedures.NewestCompatible(version)
	} else {
		// If no current version, use highest version to reconcile
		reconcileProcedure, _ = allProcedures.Newest()
	}

	// Execute the reconcile procedure. Not sure how to handle the error
	procedureStatus, _ := reconcileProcedure.Execute(request, r.client, r.scheme)
	// Walk ProcedureStatus.Results and add to the CR status
	for _, result := range procedureStatus.Results {
		instance.Status.ReconcileActions[result.Name] = result.Result
	}

	err = r.client.Status().Update(context.TODO(), instance)
	if err != nil {
		return reconcile.Result{}, err
	}
	// if ProcedureStatus.FullyReconciled
	//   update reconcile version in the CR to match the Procedure version
	//   use a timed reconcile requeue //left this part out. Why requeue?
	if procedureStatus.FullyReconciled {
		instance.Spec.Options["reconcileVersion"] = strconv.Itoa(reconcileProcedure.Version())
		err := r.client.Update(context.TODO(), instance)
		if err != nil {
			if errors.IsNotFound(err) {
				// Request object not found, could have been deleted after reconcile request.
				// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
				// Return and don't requeue
				return reconcile.Result{}, nil
			}
			// Error reading the object - requeue the request.
			return reconcile.Result{}, err
		}
	} else {
		//   requeue immediately
		return reconcile.Result{Requeue: true}, nil
	}

	return reconcile.Result{}, nil
}
