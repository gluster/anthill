package reconcileaction

import (
	"errors"
	"testing"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

var trueAction = ReconcileAction{
	Name: "trueAction",
	action: func(_ reconcile.Request) (corev1.ConditionStatus, error) {
		return corev1.ConditionTrue, nil
	},
}

var falseAction = ReconcileAction{
	Name: "falseAction",
	action: func(_ reconcile.Request) (corev1.ConditionStatus, error) {
		return corev1.ConditionFalse, nil
	},
}

var unknownAction = ReconcileAction{
	Name: "unknownAction",
	action: func(_ reconcile.Request) (corev1.ConditionStatus, error) {
		return corev1.ConditionUnknown, nil
	},
}

var errGeneric = errors.New("an error")
var errorAction = ReconcileAction{
	Name: "errorAction",
	action: func(_ reconcile.Request) (corev1.ConditionStatus, error) {
		return corev1.ConditionUnknown, errGeneric
	},
}

var tfAction = ReconcileAction{
	Name:    "TruePrereqsFalseAction",
	prereqs: []*ReconcileAction{&trueAction, &trueAction},
	action: func(_ reconcile.Request) (corev1.ConditionStatus, error) {
		return corev1.ConditionFalse, nil
	},
}

var ftAction = ReconcileAction{
	Name:    "FalsePrereqsTrueAction",
	prereqs: []*ReconcileAction{&trueAction, &falseAction},
	action: func(_ reconcile.Request) (corev1.ConditionStatus, error) {
		return corev1.ConditionTrue, nil
	},
}

var count int
var countAction = ReconcileAction{
	Name: "CountingAction",
	action: func(_ reconcile.Request) (corev1.ConditionStatus, error) {
		count++
		return corev1.ConditionTrue, nil
	},
}

func TestActionsReturnCorrectValue(t *testing.T) {
	var tests = []struct {
		input    ReconcileAction
		wantCond corev1.ConditionStatus
		wantErr  error
	}{
		{trueAction, corev1.ConditionTrue, nil},
		{falseAction, corev1.ConditionFalse, nil},
		{unknownAction, corev1.ConditionUnknown, nil},
		{errorAction, corev1.ConditionUnknown, errGeneric},
		{tfAction, corev1.ConditionFalse, nil},
		{ftAction, corev1.ConditionUnknown, nil},
	}

	dummy := reconcile.Request{
		NamespacedName: types.NamespacedName{
			Name:      "name",
			Namespace: "namespace",
		},
	}

	for _, test := range tests {
		c, e := test.input.Execute(dummy)
		if c != test.wantCond || e != test.wantErr {
			t.Errorf("%s -- expected: (%v, %v) -- got: (%v, %v)", test.input.Name, test.wantCond, test.wantErr, c, e)
		}
	}
}

func TestActionsCacheValues(t *testing.T) {
	dummy := reconcile.Request{
		NamespacedName: types.NamespacedName{
			Name:      "name",
			Namespace: "namespace",
		},
	}

	countAction.Execute(dummy)
	if count != 1 {
		t.Errorf("execution count should be 1; is %d", count)
	}
	countAction.Execute(dummy)
	if count != 1 {
		t.Errorf("execution was not properly cached")
	}
	countAction.Clear()
	countAction.Execute(dummy)
	if count != 2 {
		t.Errorf("execution count should be 2; cache didn't clear")
	}
}
