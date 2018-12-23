package reconcileaction

import (
	"errors"
	"testing"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

var trueAction = Action{
	Name: "trueAction",
	action: func(_ reconcile.Request, _ client.Client, _ *runtime.Scheme) (Result, error) {
		return Result{Status: corev1.ConditionTrue, Message: "it's true"}, nil
	},
}

var falseAction = Action{
	Name: "falseAction",
	action: func(_ reconcile.Request, _ client.Client, _ *runtime.Scheme) (Result, error) {
		return Result{Status: corev1.ConditionFalse, Message: "it's false"}, nil
	},
}

var unknownAction = Action{
	Name: "unknownAction",
	action: func(_ reconcile.Request, _ client.Client, _ *runtime.Scheme) (Result, error) {
		return Result{Status: corev1.ConditionUnknown, Message: "who knows?"}, nil
	},
}

var errGeneric = errors.New("an error")
var errorAction = Action{
	Name: "errorAction",
	action: func(_ reconcile.Request, _ client.Client, _ *runtime.Scheme) (Result, error) {
		return Result{Status: corev1.ConditionUnknown, Message: "it was bad"}, errGeneric
	},
}

var tfAction = Action{
	Name:    "TruePrereqsFalseAction",
	prereqs: []*Action{&trueAction, &trueAction},
	action: func(_ reconcile.Request, _ client.Client, _ *runtime.Scheme) (Result, error) {
		return Result{Status: corev1.ConditionFalse, Message: "it's false"}, nil
	},
}

var ftAction = Action{
	Name:    "FalsePrereqsTrueAction",
	prereqs: []*Action{&trueAction, &falseAction},
	action: func(_ reconcile.Request, _ client.Client, _ *runtime.Scheme) (Result, error) {
		return Result{Status: corev1.ConditionTrue, Message: "it's true"}, nil
	},
}

var count int
var countAction = Action{
	Name: "CountingAction",
	action: func(_ reconcile.Request, _ client.Client, _ *runtime.Scheme) (Result, error) {
		count++
		return Result{Status: corev1.ConditionTrue, Message: "it's true"}, nil
	},
}

func TestActionsReturnCorrectValue(t *testing.T) {
	var tests = []struct {
		input    Action
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

	request := reconcile.Request{
		NamespacedName: types.NamespacedName{
			Name:      "name",
			Namespace: "namespace",
		},
	}
	client := fake.NewFakeClient()
	var scheme *runtime.Scheme

	for _, test := range tests {
		r, e := test.input.Execute(request, client, scheme)
		if r.Status != test.wantCond || e != test.wantErr {
			t.Errorf("%s -- expected: (%v, %v) -- got: (%v, %v)", test.input.Name, test.wantCond, test.wantErr, r.Status, e)
		}
	}
}

func TestActionsCacheValues(t *testing.T) {
	request := reconcile.Request{
		NamespacedName: types.NamespacedName{
			Name:      "name",
			Namespace: "namespace",
		},
	}
	client := fake.NewFakeClient()
	var scheme *runtime.Scheme

	countAction.Execute(request, client, scheme)
	if count != 1 {
		t.Errorf("execution count should be 1; is %d", count)
	}
	countAction.Execute(request, client, scheme)
	if count != 1 {
		t.Errorf("execution was not properly cached")
	}
	countAction.Clear()
	countAction.Execute(request, client, scheme)
	if count != 2 {
		t.Errorf("execution count should be 2; cache didn't clear")
	}
}
