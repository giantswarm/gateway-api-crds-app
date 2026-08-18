package basic

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/giantswarm/apptest-framework/v5/pkg/state"
	"github.com/giantswarm/apptest-framework/v5/pkg/suite"
	"github.com/giantswarm/clustertest/v5/pkg/client"

	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/kubernetes/scheme"
	cr "sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	isUpgrade = false

	// Namespace the bundle installs the CRD app into, and therefore the namespace the
	// installer hook resources are created in.
	installNamespace = "kube-system"
)

// The clustertest clients are built on the client-go scheme, which doesn't know about
// CustomResourceDefinitions. Register them so the tests can work with typed objects.
func init() {
	utilruntime.Must(apiextensionsv1.AddToScheme(scheme.Scheme))
}

func TestBasic(t *testing.T) {
	suite.New().
		InAppBundle("gateway-api-bundle").
		WithInstallNamespace(installNamespace).
		WithIsUpgrade(isUpgrade).
		WithValuesFile("./values.yaml").
		WithBundleValuesFile("./bundle_values.yaml").
		Tests(func() {
			It("should have the app correctly deployed", func() {
				deploymentAppTests()
			})
			It("should have the selected CRDs installed", func() {
				crdEstablishedTests()
				crdVersionTests()
				crdNotSelectedTests()
			})
			It("should have removed the installer hook resources", func() {
				installerCleanupTests()
			})
			It("should have the safe-upgrades admission policy enforced", func() {
				admissionPolicyResourceTests()
				admissionPolicyEnforcementTests()
			})
			It("should serve the Gateway API resources", func() {
				gatewayAPIResourceTests()
				gatewayAPIValidationTests()
			})
		}).
		Run(t, "Gateway API CRDs Test")
}

// wcClient returns a client for the workload cluster the app under test is installed into.
func wcClient() *client.Client {
	wcClient, err := state.GetFramework().WC(state.GetCluster().Name)
	Expect(err).NotTo(HaveOccurred())
	Expect(wcClient).NotTo(BeNil())

	return wcClient
}

// isNotFound gets the given object and reports whether the API server doesn't know it.
// A missing kind counts as not found as well, the resource cannot exist without it.
func isNotFound(obj cr.Object, key cr.ObjectKey) func() (bool, error) {
	return func() (bool, error) {
		err := wcClient().Get(state.GetContext(), key, obj)
		if err == nil {
			return false, nil
		}
		if errors.IsNotFound(err) || meta.IsNoMatchError(err) {
			return true, nil
		}

		return false, err
	}
}
