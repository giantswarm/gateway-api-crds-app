package basic

import (
	"fmt"
	"regexp"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/giantswarm/apptest-framework/v5/pkg/state"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	cr "sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	channelAnnotation       = "gateway.networking.k8s.io/channel"
	bundleVersionAnnotation = "gateway.networking.k8s.io/bundle-version"

	inferenceBundleVersionAnnotation = "inference.networking.k8s.io/bundle-version"

	// Field manager the installer Job applies the CRDs with, see the Job template.
	installerFieldManager = "kubectl-client-side-apply"
)

// gatewayAPICRDs are the Gateway API CRDs the chart installs from the standard channel with
// its default values.
var gatewayAPICRDs = []string{
	"backendtlspolicies.gateway.networking.k8s.io",
	"gatewayclasses.gateway.networking.k8s.io",
	"gateways.gateway.networking.k8s.io",
	"grpcroutes.gateway.networking.k8s.io",
	"httproutes.gateway.networking.k8s.io",
	"listenersets.gateway.networking.k8s.io",
	"referencegrants.gateway.networking.k8s.io",
	"tcproutes.gateway.networking.k8s.io",
	"tlsroutes.gateway.networking.k8s.io",
	"udproutes.gateway.networking.k8s.io",
}

// inferenceCRDs are the Inference Extension CRDs this suite opts into via ./bundle_values.yaml.
var inferenceCRDs = []string{
	"inferencepools.inference.networking.k8s.io",
}

// disabledCRDs are off in the chart defaults and not opted into by this suite, so the
// installer must not create them.
var disabledCRDs = []string{
	"xbackends.gateway.networking.x-k8s.io",
	"xbackendtrafficpolicies.gateway.networking.x-k8s.io",
	"xmeshes.gateway.networking.x-k8s.io",
	"inferenceobjectives.inference.networking.x-k8s.io",
	"inferencepoolimports.inference.networking.x-k8s.io",
	"inferencemodelrewrites.inference.networking.x-k8s.io",
}

// The exact bundle version follows the upstream release the chart vendors, so the tests only
// assert its shape.
var bundleVersionRegex = regexp.MustCompile(`^v\d+\.\d+\.\d+`)

// installedCRDs returns every CRD the chart is expected to install in this suite.
func installedCRDs() []string {
	return append(append([]string{}, gatewayAPICRDs...), inferenceCRDs...)
}

// getCRD fetches a CRD from the workload cluster, failing the test if it isn't there.
func getCRD(name string) *unstructured.Unstructured {
	crd := newCRD(name)
	Expect(wcClient().Get(state.GetContext(), cr.ObjectKey{Name: name}, crd)).To(Succeed())

	return crd
}

// crdEstablishedTests verifies every selected CRD was applied by the installer Job and is
// served by the API server.
func crdEstablishedTests() {
	for _, name := range installedCRDs() {
		By(fmt.Sprintf("checking CRD %s is established", name))
		crd := newCRD(name)
		Eventually(func() error {
			if err := wcClient().Get(state.GetContext(), cr.ObjectKey{Name: name}, crd); err != nil {
				return err
			}
			conditions, _, err := unstructured.NestedSlice(crd.Object, "status", "conditions")
			if err != nil {
				return err
			}
			for _, condition := range conditions {
				condition, ok := condition.(map[string]any)
				if !ok {
					continue
				}
				if condition["type"] == "Established" && condition["status"] == "True" {
					return nil
				}
			}

			return fmt.Errorf("CRD %s is not established", name)
		}).
			WithTimeout(10 * time.Minute).
			WithPolling(10 * time.Second).
			Should(Succeed())

		By(fmt.Sprintf("checking CRD %s was applied by the installer Job", name))
		managers := []string{}
		for _, entry := range crd.GetManagedFields() {
			managers = append(managers, entry.Manager)
		}
		Expect(managers).To(ContainElement(installerFieldManager))

		By(fmt.Sprintf("checking CRD %s has a single served storage version", name))
		versions, _, err := unstructured.NestedSlice(crd.Object, "spec", "versions")
		Expect(err).NotTo(HaveOccurred())
		storageVersions := []string{}
		for _, version := range versions {
			version, ok := version.(map[string]any)
			Expect(ok).To(BeTrue())
			if version["storage"] == true {
				Expect(version["served"]).To(BeTrue(), "storage version %v of %s is not served", version["name"], name)
				storageVersions = append(storageVersions, fmt.Sprintf("%v", version["name"]))
			}
		}
		Expect(storageVersions).To(HaveLen(1))
	}
}

// crdVersionTests verifies all Gateway API CRDs come from the same bundle and the same
// channel, catching an installer that mixes the standard and experimental sets.
func crdVersionTests() {
	bundleVersion := ""
	for _, name := range gatewayAPICRDs {
		By(fmt.Sprintf("checking CRD %s is from the standard channel", name))
		annotations := getCRD(name).GetAnnotations()
		Expect(annotations).To(HaveKeyWithValue(channelAnnotation, "standard"))

		By(fmt.Sprintf("checking CRD %s bundle version", name))
		Expect(annotations).To(HaveKey(bundleVersionAnnotation))
		version := annotations[bundleVersionAnnotation]
		Expect(version).To(MatchRegexp(bundleVersionRegex.String()))

		if bundleVersion == "" {
			bundleVersion = version
		}
		Expect(version).To(Equal(bundleVersion), "CRD %s is from bundle %s while others are from %s", name, version, bundleVersion)
	}

	for _, name := range inferenceCRDs {
		By(fmt.Sprintf("checking CRD %s bundle version", name))
		annotations := getCRD(name).GetAnnotations()
		Expect(annotations).To(HaveKey(inferenceBundleVersionAnnotation))
		Expect(annotations[inferenceBundleVersionAnnotation]).To(MatchRegexp(bundleVersionRegex.String()))
	}
}

// crdNotSelectedTests verifies the installer only applies what the values select.
func crdNotSelectedTests() {
	for _, name := range disabledCRDs {
		By(fmt.Sprintf("checking CRD %s is not installed", name))
		Expect(isNotFound(newCRD(name), cr.ObjectKey{Name: name})()).
			To(BeTrue(), "CRD %s is installed but not selected by the chart values", name)
	}
}

// gatewayAPIBundleVersion returns the bundle version of the installed Gateway API CRDs.
func gatewayAPIBundleVersion() string {
	return getCRD("httproutes.gateway.networking.k8s.io").GetAnnotations()[bundleVersionAnnotation]
}
