package reconcileaction

import (
	"errors"
	"fmt"
	"math/rand"
	"sort"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// ProcedureList is used to define all the reconcile Procedures for the operator
type ProcedureList []Procedure

// Newest returns the Procedure with the highest version number
func (pl ProcedureList) Newest() (*Procedure, error) {
	if len(pl) < 1 {
		return nil, fmt.Errorf("empty list of reconcile procedures")
	}

	// Sort descending by version
	sort.Slice(pl, func(i, j int) bool {
		return pl[i].Version() > pl[j].Version()
	})
	p := pl[0]

	return &p, nil
}

// NewestCompatible returns the newest Procedure that is compatible with the
// provided currentVersion
func (pl ProcedureList) NewestCompatible(currentVersion int) (*Procedure, error) {
	if len(pl) < 1 {
		return nil, errors.New("empty list of reconcile procedures")
	}

	// Sort descending by version
	sort.Slice(pl, func(i, j int) bool {
		return pl[i].Version() > pl[j].Version()
	})

	// Walk list to find the first that is compatible w/ the currentVersion
	for _, p := range pl {
		if p.MinVersion() <= currentVersion {
			return &p, nil
		}
	}

	return nil, fmt.Errorf("no procedures are compatible with deployed version %d", currentVersion)
}

// Procedure defines the complete "procedure" (the set of Action) necessary to
// completely reconcile the state.
type Procedure struct {
	minVersion int
	version    int
	actions    []*Action
}

// Version is the reconciler version implemented by this Procedure
func (p *Procedure) Version() int {
	return p.version
}

// MinVersion is the minimum reconciler version this Procedure can be used to
// upgrade from.
func (p *Procedure) MinVersion() int {
	return p.minVersion
}

// ActionResult is the Result of an action, paired with its name
type ActionResult struct {
	Name string
	Result
}

// ProcedureStatus is the result of executing a reconcile Procedure
type ProcedureStatus struct {
	// Results will contain an entry for each action that was attempted
	// while executing the Procedure
	Results []ActionResult
	// FullyReconciled will be true iff the system state was found to be
	// full reconciled according to the Procedure
	FullyReconciled bool
}

// Execute the reconcile Procedure
func (p *Procedure) Execute(request reconcile.Request, client client.Client, scheme *runtime.Scheme) (*ProcedureStatus, error) {
	// All action dependencies MUST be expressed via its prereqs. To enforce
	// that, we intentionally shuffle the list of actions that define a
	// Procedure.
	actions := p.actions
	rand.Shuffle(len(actions), func(i, j int) {
		actions[i], actions[j] = actions[j], actions[i]
	})

	// Clear any cached action state
	for _, step := range actions {
		step.Clear()
	}

	status := ProcedureStatus{
		FullyReconciled: true,
	}

	// Execute the actions
	for _, step := range actions {
		result, err := step.Execute(request, client, scheme)
		if err != nil {
			return nil, err
		}
		// Any component Action not fully reconciled means this
		// Procedure isn't either
		if result.Status != corev1.ConditionTrue {
			status.FullyReconciled = false
		}

		ar := ActionResult{
			Name:   step.Name,
			Result: result,
		}
		status.Results = append(status.Results, ar)
	}

	return &status, nil
}
