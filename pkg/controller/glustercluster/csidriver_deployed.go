package glustercluster

import (
	"fmt"

	"github.com/gluster/anthill/pkg/reconciler"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

var glusterFuseProvisionerDeployed = reconciler.NewAction(
	"glusterFuseProvisionerDeployed",
	[]*reconciler.Action{},
	func(request reconcile.Request, client client.Client, scheme *runtime.Scheme) (reconciler.Result, error) {
		clustername := "clustername"
		clusternamespace := "clusternamespace"
		servicename := fmt.Sprintf("%v-csi-provisioner", clustername)
		statefulSet := appsv1.StatefulSet{
			ObjectMeta: metav1.ObjectMeta{
				Name:      fmt.Sprintf("csi-provisioner-%v", clustername),
				Namespace: fmt.Sprintf("csi-provisioner-%v", clusternamespace),
				Labels: map[string]string{
					"app.kubernetes.io/part-of":   fmt.Sprintf("glustercluster/%v", clustername),
					"app.kubernetes.io/component": "csi-driver",
					"app.kubernetes.io/name":      "csi-provisioner",
				},
			},
			Spec: appsv1.StatefulSetSpec{
				ServiceName: servicename,
			},
		}
		fmt.Printf(statefulSet.Name)
		return reconciler.Result{Status: corev1.ConditionTrue, Message: "it's true"}, nil
	},
)

var glusterFuseAttacherDeployed = reconciler.NewAction(
	"glusterFuseAttachedDeployed",
	[]*reconciler.Action{},
	func(request reconcile.Request, client client.Client, scheme *runtime.Scheme) (reconciler.Result, error) {

		return reconciler.Result{Status: corev1.ConditionTrue, Message: "it's true"}, nil
	},
)

var glusterFuseNodeDeployed = reconciler.NewAction(
	"glusterFuseNodeDeployed",
	[]*reconciler.Action{},
	func(request reconcile.Request, client client.Client, scheme *runtime.Scheme) (reconciler.Result, error) {
		return reconciler.Result{Status: corev1.ConditionTrue, Message: "it's true"}, nil
	},
)

var serviceAccountDeployed = reconciler.NewAction(
	"serviceAccountDeployed",
	[]*reconciler.Action{},
	func(request reconcile.Request, client client.Client, scheme *runtime.Scheme) (reconciler.Result, error) {

		return reconciler.Result{Status: corev1.ConditionTrue, Message: "it's true"}, nil
	},
)
