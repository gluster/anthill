package reconcileaction

import (
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// ReconcileAction is an action that reconciles the system state. It has a list
// of prerequisite actions that must be true in order for the action to be
// invoked.
type ReconcileAction struct {
	Name          string
	prereqs       []*ReconcileAction
	action        func(reconcile.Request) (corev1.ConditionStatus, error)
	lastCondition *corev1.ConditionStatus
	lastError     error
}

// Execute or return a previously cached result of the ReconcileAction, checking prereqs first
func (ra *ReconcileAction) Execute(request reconcile.Request) (corev1.ConditionStatus, error) {
	// If we have executed before, return the cached result
	if ra.lastCondition != nil {
		return *ra.lastCondition, ra.lastError
	}

	// Walk through the prereqs; stop and return corev1.ConditionUnknown if a prereq doesn't return corev1.ConditionTrue
	for _, prereq := range ra.prereqs {
		cond, err := prereq.Execute(request)
		if err != nil || cond != corev1.ConditionTrue {
			unknown := corev1.ConditionUnknown
			ra.lastCondition = &unknown
			ra.lastError = nil
			return *ra.lastCondition, ra.lastError
		}
	}

	// Perform the reconcile action
	c, e := ra.action(request)
	ra.lastCondition, ra.lastError = &c, e
	return *ra.lastCondition, ra.lastError
}

// Clear and cached results of this ReconcileAction
func (ra *ReconcileAction) Clear() {
	ra.lastCondition = nil
	ra.lastError = nil
}
