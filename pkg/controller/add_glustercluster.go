package controller

import (
	"github.com/gluster/anthill/pkg/controller/glustercluster"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, glustercluster.Add)
}
