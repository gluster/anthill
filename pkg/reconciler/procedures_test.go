package reconciler

import (
	"testing"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

var v9 = Procedure{
	version:    9,
	minVersion: 6,
	actions:    []*Action{&trueAction, &errorAction, &trueAction},
}

var v8 = Procedure{
	version:    8,
	minVersion: 7,
	actions:    []*Action{&trueAction, &falseAction, &trueAction},
}

var v7 = Procedure{
	version:    7,
	minVersion: 2,
	actions:    []*Action{&trueAction, &countAction},
}

func TestProcedureClearsCached(t *testing.T) {
	request := reconcile.Request{
		NamespacedName: types.NamespacedName{
			Name:      "name",
			Namespace: "namespace",
		},
	}
	client := fake.NewFakeClient()
	var scheme *runtime.Scheme

	count = 0
	ps, err := v7.Execute(request, client, scheme)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if !ps.FullyReconciled {
		t.Errorf("procedure should have been fully reconciled")
	}
	if count != 1 {
		t.Errorf("countAction should have been called once; actual: %d", count)
	}
	_, _ = v7.Execute(request, client, scheme)
	if count != 2 {
		t.Errorf("countAction should have been called twice; actual: %d", count)
	}
}

func TestProcedureReturnsActionError(t *testing.T) {
	request := reconcile.Request{
		NamespacedName: types.NamespacedName{
			Name:      "name",
			Namespace: "namespace",
		},
	}
	client := fake.NewFakeClient()
	var scheme *runtime.Scheme

	_, err := v9.Execute(request, client, scheme)
	if err == nil {
		t.Errorf("Execute should have returned an error")
	}
}

func TestProcedureNotFullyReconciledIfFalseAction(t *testing.T) {
	request := reconcile.Request{
		NamespacedName: types.NamespacedName{
			Name:      "name",
			Namespace: "namespace",
		},
	}
	client := fake.NewFakeClient()
	var scheme *runtime.Scheme

	ps, err := v8.Execute(request, client, scheme)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if ps.FullyReconciled {
		t.Errorf("procedure should not have been fully reconciled")
	}
}

func TestPrereqsHaveCacheCleared(t *testing.T) {
	request := reconcile.Request{
		NamespacedName: types.NamespacedName{
			Name:      "name",
			Namespace: "namespace",
		},
	}
	client := fake.NewFakeClient()
	var scheme *runtime.Scheme

	// a is an action w/ a counting prereq, so we can tell how often the
	// prereqs get called.
	a := Action{
		prereqs: []*Action{&countAction},
		action:  trueAction.action,
	}
	p := Procedure{
		version:    8,
		minVersion: 7,
		actions:    []*Action{&a},
	}
	count = 0
	_, err := p.Execute(request, client, scheme)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if count != 1 {
		t.Errorf("countAction prereq should have been called once; actual: %d", count)
	}
	// execute again. The cache should be cleared, causing count to
	// increase.
	_, err = p.Execute(request, client, scheme)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if count != 2 {
		t.Errorf("countAction prereq should have been called a 2nd time; actual: %d", count)
	}

}

var pl = ProcedureList{v8, v9, v7} // out of order to check sorting

func TestEmptyListReturnsError(t *testing.T) {
	var empty ProcedureList

	if _, err := empty.Newest(); err == nil {
		t.Error("Newest() should have returned an error")
	}
	version := new(int)
	*version = 5
	if _, err := empty.NewestCompatible(version); err == nil {
		t.Error("NewestCompatible() should have returned an error")
	}
}

func TestNewestReturnsHighestVersion(t *testing.T) {
	expected := 9
	p, err := pl.Newest()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if p.Version() != expected {
		t.Errorf("expected version %d, got version %d", expected, p.Version())
	}
}

func TestNewestCompatible(t *testing.T) {
	expected := 7
	version := new(int)
	*version = 4
	p, err := pl.NewestCompatible(version)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if p.Version() != expected {
		t.Errorf("expected version %d, got version %d", expected, p.Version())
	}

	// All versions are compatible w/ v7, so we should get back the highest
	// (9)
	expected = 9
	*version = 7
	p, err = pl.NewestCompatible(version)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if p.Version() != expected {
		t.Errorf("expected version %d, got version %d", expected, p.Version())
	}

	// Nothing is compatible w/ v1
	*version = 1
	if _, err = pl.NewestCompatible(version); err == nil {
		t.Error("NewestCompatible() should have returned an error")
	}

	// nil version should return the Newest
	version = nil
	expected = 9
	p, err = pl.NewestCompatible(version)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if p.Version() != expected {
		t.Errorf("expected version %d, got version %d", expected, p.Version())
	}

}
