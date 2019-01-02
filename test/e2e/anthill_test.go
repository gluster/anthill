package e2e

import (
	goctx "context"
	"testing"
	"time"

	framework "github.com/operator-framework/operator-sdk/pkg/test"
	"github.com/operator-framework/operator-sdk/pkg/test/e2eutil"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	apis "github.com/gluster/anthill/pkg/apis"
	operator "github.com/gluster/anthill/pkg/apis/operator/v1alpha1"
)

var (
	retryInterval        = time.Second * 5
	cleanupRetryInterval = time.Second * 1
	cleanupTimeout       = time.Second * 5
)

func TestAnthill(t *testing.T) {
	glusterClusterList := &operator.GlusterClusterList{
		TypeMeta: metav1.TypeMeta{
			Kind:       "GlusterCluster",
			APIVersion: "operator.gluster.org/v1alpha1",
		},
	}
	err := framework.AddToFrameworkScheme(apis.AddToScheme, glusterClusterList)
	if err != nil {
		t.Fatalf("failed to add custom resource scheme to framework: %v", err)
	}
	// run subtests
	t.Run("anthill-tests", func(t *testing.T) {
		t.Run("SimpleTest", simpleTest)
	})
}

func simpleTest(t *testing.T) {
	t.Parallel()
	ctx := framework.NewTestCtx(t)
	defer ctx.Cleanup()
	err := ctx.InitializeClusterResources(&framework.CleanupOptions{TestContext: ctx, Timeout: cleanupTimeout, RetryInterval: cleanupRetryInterval})
	if err != nil {
		t.Fatalf("failed to initialize cluster resources: %v", err)
	}
	t.Log("Initialized cluster resources")
	namespace, err := ctx.GetNamespace()
	if err != nil {
		t.Fatal(err)
	}
	// get global framework variables
	f := framework.Global

	// wait for anthill operator to be ready
	err = e2eutil.WaitForOperatorDeployment(t, f.KubeClient, namespace, "anthill", 1, retryInterval, time.Minute*5)
	if err != nil {
		t.Fatal(err)
	}

	// create memcached custom resource
	cluster := &operator.GlusterCluster{
		TypeMeta: metav1.TypeMeta{
			Kind:       "GlusterCluster",
			APIVersion: "operator.gluster.org/v1alpha1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "cluster",
			Namespace: namespace,
		},
		Spec: operator.GlusterClusterSpec{
			//Size: 3,
		},
	}
	err = f.Client.Create(goctx.TODO(), cluster, &framework.CleanupOptions{TestContext: ctx, Timeout: time.Second * 5, RetryInterval: time.Second * 1})
	if err != nil {
		t.Fatal(err)
	}

}
