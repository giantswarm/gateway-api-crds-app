package basic

import (
	"fmt"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	cr "sigs.k8s.io/controller-runtime/pkg/client"
)

// installerName is the name shared by all resources of the installer hook.
const installerName = "gateway-api-crds-installer"

// installerCleanupTests verifies the installer hook leaves nothing behind. All its resources
// carry the "hook-succeeded" delete policy, so a leftover means the Job never completed or
// Helm failed to clean up, both of which would break the next upgrade.
func installerCleanupTests() {
	ciliumNetworkPolicy := &unstructured.Unstructured{}
	ciliumNetworkPolicy.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   "cilium.io",
		Version: "v2",
		Kind:    "CiliumNetworkPolicy",
	})

	namespacedKey := cr.ObjectKey{Name: installerName, Namespace: installNamespace}
	clusterKey := cr.ObjectKey{Name: installerName}

	resources := []struct {
		description string
		object      cr.Object
		key         cr.ObjectKey
	}{
		{"Job", &batchv1.Job{}, namespacedKey},
		{"ServiceAccount", &corev1.ServiceAccount{}, namespacedKey},
		{"CiliumNetworkPolicy", ciliumNetworkPolicy, namespacedKey},
		{"ClusterRole", &rbacv1.ClusterRole{}, clusterKey},
		{"ClusterRoleBinding", &rbacv1.ClusterRoleBinding{}, clusterKey},
	}

	for _, resource := range resources {
		By(fmt.Sprintf("checking installer %s is removed", resource.description))
		Eventually(isNotFound(resource.object, resource.key)).
			WithTimeout(5 * time.Minute).
			WithPolling(10 * time.Second).
			Should(BeTrue())
	}
}
