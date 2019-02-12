package glustercluster

import (
	"github.com/gluster/anthill/pkg/reconciler"
)

// ProcedureV1 is Procedure version 1
var ProcedureV1 = reconciler.NewProcedure(
	0,
	0,
	[]*reconciler.Action{
		etcdClusterCreated,
		glusterFuseProvisionerDeployed,
		glusterFuseAttacherDeployed,
		glusterFuseNodeDeployed,
	},
)
