package basic

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/giantswarm/apptest-framework/v5/pkg/state"

	admissionregistrationv1 "k8s.io/api/admissionregistration/v1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"
	cr "sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	safeUpgradesPolicyName = "safe-upgrades.gateway.networking.k8s.io"

	// Name of the throwaway CRD the enforcement tests apply. It lives in the Gateway API
	// group so the policy matches it, and is removed once the tests are done.
	testCRDName = "e2etests.gateway.networking.k8s.io"

	// Any bundle version before v1.5.0 is rejected by the policy.
	unsafeBundleVersion = "v1.0.0"
)

// admissionPolicyResourceTests verifies the safe-upgrades policy and its binding are
// installed and set to deny.
func admissionPolicyResourceTests() {
	By("checking the ValidatingAdmissionPolicy exists")
	policy := &admissionregistrationv1.ValidatingAdmissionPolicy{}
	Eventually(func() error {
		return wcClient().Get(state.GetContext(), cr.ObjectKey{Name: safeUpgradesPolicyName}, policy)
	}).
		WithTimeout(5 * time.Minute).
		WithPolling(5 * time.Second).
		Should(Succeed())
	Expect(policy.Spec.Validations).ToNot(BeEmpty())

	By("checking the ValidatingAdmissionPolicyBinding denies violations")
	binding := &admissionregistrationv1.ValidatingAdmissionPolicyBinding{}
	Eventually(func() error {
		return wcClient().Get(state.GetContext(), cr.ObjectKey{Name: safeUpgradesPolicyName}, binding)
	}).
		WithTimeout(5 * time.Minute).
		WithPolling(5 * time.Second).
		Should(Succeed())
	Expect(binding.Spec.PolicyName).To(Equal(safeUpgradesPolicyName))
	Expect(binding.Spec.ValidationActions).To(ContainElement(admissionregistrationv1.Deny))
}

// admissionPolicyEnforcementTests verifies the policy is actually evaluated by the API
// server: a Gateway API CRD from an unsupported bundle is rejected while one from the bundle
// the chart ships is accepted.
func admissionPolicyEnforcementTests() {
	// The test CRD must not outlive the suite, not even when the policy fails to reject it.
	DeferCleanup(func() {
		deleteAndWait(&apiextensionsv1.CustomResourceDefinition{
			ObjectMeta: metav1.ObjectMeta{Name: testCRDName},
		})
	})

	By("checking a Gateway API CRD from an unsupported bundle is rejected")
	err := wcClient().Create(state.GetContext(), testCRD(unsafeBundleVersion))
	Expect(err).To(HaveOccurred())
	Expect(err.Error()).To(ContainSubstring("prohibited by default"))

	By("checking a Gateway API CRD from the installed bundle is accepted")
	Expect(wcClient().Create(state.GetContext(), testCRD(gatewayAPIBundleVersion()))).To(Succeed())
}

// testCRD builds a throwaway CRD in the Gateway API group, annotated with the given bundle
// version so the safe-upgrades policy is exercised.
func testCRD(bundleVersion string) *apiextensionsv1.CustomResourceDefinition {
	return &apiextensionsv1.CustomResourceDefinition{
		ObjectMeta: metav1.ObjectMeta{
			Name: testCRDName,
			Annotations: map[string]string{
				// CRDs in the *.k8s.io groups need an approval annotation. This one is
				// never part of the API, so it explicitly bypasses the approval process.
				"api-approved.kubernetes.io": "unapproved.kubernetes.io/e2e-test",
				channelAnnotation:            "standard",
				bundleVersionAnnotation:      bundleVersion,
			},
		},
		Spec: apiextensionsv1.CustomResourceDefinitionSpec{
			Group: "gateway.networking.k8s.io",
			Names: apiextensionsv1.CustomResourceDefinitionNames{
				Kind:     "E2ETest",
				ListKind: "E2ETestList",
				Plural:   "e2etests",
				Singular: "e2etest",
			},
			Scope: apiextensionsv1.NamespaceScoped,
			Versions: []apiextensionsv1.CustomResourceDefinitionVersion{
				{
					Name:    "v1alpha1",
					Served:  true,
					Storage: true,
					Schema: &apiextensionsv1.CustomResourceValidation{
						OpenAPIV3Schema: &apiextensionsv1.JSONSchemaProps{
							Type:                   "object",
							XPreserveUnknownFields: ptr.To(true),
						},
					},
				},
			},
		},
	}
}
