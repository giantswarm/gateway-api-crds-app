package basic

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/giantswarm/apptest-framework/v5/pkg/state"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"
	cr "sigs.k8s.io/controller-runtime/pkg/client"
	gatewayv1 "sigs.k8s.io/gateway-api/apis/v1"
)

const (
	testResourceNamespace = "default"
	testGatewayName       = "e2e-crds-gateway"
	testHTTPRouteName     = "e2e-crds-httproute"
)

// gatewayAPIResourceTests verifies the installed CRDs actually serve the Gateway API types.
// Nothing reconciles them in this suite, the point is that the API server accepts and stores
// them.
func gatewayAPIResourceTests() {
	By("creating a Gateway")
	gateway := &gatewayv1.Gateway{
		ObjectMeta: metav1.ObjectMeta{
			Name:      testGatewayName,
			Namespace: testResourceNamespace,
		},
		Spec: gatewayv1.GatewaySpec{
			GatewayClassName: "e2e-crds-test",
			Listeners: []gatewayv1.Listener{
				{
					Name:     "http",
					Port:     80,
					Protocol: gatewayv1.HTTPProtocolType,
				},
			},
		},
	}
	Expect(wcClient().Create(state.GetContext(), gateway)).To(Succeed())
	DeferCleanup(func() {
		deleteAndWait(gateway)
	})

	By("creating an HTTPRoute attached to the Gateway")
	route := &gatewayv1.HTTPRoute{
		ObjectMeta: metav1.ObjectMeta{
			Name:      testHTTPRouteName,
			Namespace: testResourceNamespace,
		},
		Spec: gatewayv1.HTTPRouteSpec{
			CommonRouteSpec: gatewayv1.CommonRouteSpec{
				ParentRefs: []gatewayv1.ParentReference{
					{Name: gatewayv1.ObjectName(testGatewayName)},
				},
			},
			Rules: []gatewayv1.HTTPRouteRule{
				{
					Matches: []gatewayv1.HTTPRouteMatch{
						{
							Path: &gatewayv1.HTTPPathMatch{
								Type:  ptr.To(gatewayv1.PathMatchPathPrefix),
								Value: ptr.To("/"),
							},
						},
					},
				},
			},
		},
	}
	Expect(wcClient().Create(state.GetContext(), route)).To(Succeed())
	DeferCleanup(func() {
		deleteAndWait(route)
	})

	By("checking the resources are stored and served back")
	storedGateway := &gatewayv1.Gateway{}
	Expect(wcClient().Get(state.GetContext(), cr.ObjectKeyFromObject(gateway), storedGateway)).To(Succeed())
	Expect(storedGateway.Spec.Listeners).To(HaveLen(1))
	Expect(storedGateway.Spec.Listeners[0].Protocol).To(Equal(gatewayv1.HTTPProtocolType))

	storedRoute := &gatewayv1.HTTPRoute{}
	Expect(wcClient().Get(state.GetContext(), cr.ObjectKeyFromObject(route), storedRoute)).To(Succeed())
	Expect(storedRoute.Spec.ParentRefs).To(HaveLen(1))
}

// gatewayAPIValidationTests verifies the CRDs carry their CEL validation rules, so a
// truncated or partially applied schema is caught.
func gatewayAPIValidationTests() {
	By("checking an HTTPRoute with a relative path match is rejected")
	route := &gatewayv1.HTTPRoute{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "e2e-crds-invalid-httproute",
			Namespace: testResourceNamespace,
		},
		Spec: gatewayv1.HTTPRouteSpec{
			Rules: []gatewayv1.HTTPRouteRule{
				{
					Matches: []gatewayv1.HTTPRouteMatch{
						{
							Path: &gatewayv1.HTTPPathMatch{
								Type:  ptr.To(gatewayv1.PathMatchPathPrefix),
								Value: ptr.To("no-leading-slash"),
							},
						},
					},
				},
			},
		},
	}
	err := wcClient().Create(state.GetContext(), route)
	Expect(err).To(HaveOccurred())
	Expect(err.Error()).To(ContainSubstring("absolute path"))
}

// deleteAndWait removes an object and waits for it to be gone.
func deleteAndWait(obj cr.Object) {
	Expect(cr.IgnoreNotFound(wcClient().Delete(state.GetContext(), obj))).To(Succeed())
	Eventually(isNotFound(obj, cr.ObjectKeyFromObject(obj))).
		WithTimeout(2 * time.Minute).
		WithPolling(5 * time.Second).
		Should(BeTrue())
}
