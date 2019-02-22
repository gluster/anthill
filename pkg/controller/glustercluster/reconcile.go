package glustercluster

import (
	"context"

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
	err                error
	procedureStatus    *reconciler.ProcedureStatus
)

// Reconcile reads that state of the cluster for a GlusterCluster object and makes changes based on the state read
// and what is in the GlusterCluster.Spec
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileGlusterCluster) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling GlusterCluster")

	// Fetch the GlusterCluster instance
	instance := &operatorv1alpha1.GlusterCluster{}
	err = r.client.Get(context.TODO(), request.NamespacedName, instance)
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
	// If no current version, use highest version to reconcile
	reconcileProcedure, err = allProcedures.NewestCompatible(version)
	if err != nil {
		log.Error(err, "Failed to get the reconcile version.")
		return reconcile.Result{}, err
	}

	// Execute the reconcile procedure.
	procedureStatus, err = reconcileProcedure.Execute(request, r.client, r.scheme)
	if err != nil {
		log.Error(err, "Failed to execute procedure.")
		return reconcile.Result{}, err
	}
	// Walk ProcedureStatus.Results and add to the CR status
	reconcileActionStatus := make(map[string]reconciler.Result)
	for _, result := range procedureStatus.Results {
		reconcileActionStatus[result.Name] = result.Result
	}
	instance.Status.ReconcileActions = reconcileActionStatus

	if !procedureStatus.FullyReconciled {
		err = r.client.Update(context.TODO(), instance)
		if err != nil {

			return reconcile.Result{}, err
		}
		//requeue immediately
		return reconcile.Result{Requeue: true}, nil
	}
	// if ProcedureStatus.FullyReconciled
	//   update reconcile version in the CR to match the Procedure version
	//   use a timed reconcile requeue //left this part out. Why requeue?
	newVersion := reconcileProcedure.Version()
	instance.Status.ReconcileVersion = &newVersion
	err = r.client.Update(context.TODO(), instance)
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
	// use a timed reconcile requeue //TODO: implement backoff
	return reconcile.Result{RequeueAfter: 3e+10}, nil

}
