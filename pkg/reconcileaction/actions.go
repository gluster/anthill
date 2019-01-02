package reconcileaction

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// Result is the result of a Action
type Result struct {
	// Status describes the outcome of the action. If True, the action found
	// the system state to already be fully reconciled, thus requiring no
	// action.
	Status corev1.ConditionStatus
	// Message is a short human-readable explanation of the result
	Message string
}

// Action is an action that reconciles the system state. It has a list
// of prerequisite actions that must be true in order for the action to be
// invoked.
type Action struct {
	// Name is a name for the Action, to be used in the CR status and log
	// messages
	Name string
	// prereqs are the (ordered) list of prerequisites that must be true
	// prior to attempting the reconcile action
	prereqs []*Action
	// action attempts to perform the actual reconcile
	action func(reconcile.Request, client.Client, *runtime.Scheme) (Result, error)
	// lastResult holds the result of the last execution of action() or nil
	lastResult *Result
	// lastError holds the error of the last execution of action() or nil
	lastError error
}

// Execute or return a previously cached result of the Action, checking prereqs first
func (ra *Action) Execute(request reconcile.Request, client client.Client, scheme *runtime.Scheme) (Result, error) {
	// If we have executed before, return the cached result
	if ra.lastResult != nil {
		return *ra.lastResult, ra.lastError
	}

	// Walk through the prereqs; stop and return corev1.ConditionUnknown if a prereq doesn't return corev1.ConditionTrue
	for _, prereq := range ra.prereqs {
		result, err := prereq.Execute(request, client, scheme)
		if err != nil || result.Status != corev1.ConditionTrue {
			ra.lastResult = &Result{
				Status:  corev1.ConditionUnknown,
				Message: fmt.Sprintf("prequisite %s not met", prereq.Name),
			}
			ra.lastError = nil
			return *ra.lastResult, ra.lastError
		}
	}

	// Perform the reconcile action
	result, err := ra.action(request, client, scheme)
	ra.lastResult, ra.lastError = &result, err
	return *ra.lastResult, ra.lastError
}

// Clear and cached results of this Action
func (ra *Action) Clear() {
	ra.lastResult = nil
	ra.lastError = nil
}
