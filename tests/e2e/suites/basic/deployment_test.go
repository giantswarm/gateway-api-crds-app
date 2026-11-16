package basic

import (
	"fmt"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/giantswarm/apptest-framework/v5/pkg/client"
	"github.com/giantswarm/apptest-framework/v5/pkg/state"
)

const bundleName = "gateway-api-bundle"

// deploymentAppTests verifies the bundle HelmRelease and the HelmRelease it renders for the CRD
// app are ready, the latter at the version under test, before any CRD specific assertion runs.
func deploymentAppTests() {
	ctx := state.GetContext()
	clusterName := state.GetCluster().Name
	namespace := state.GetCluster().Organization.GetNamespace()
	bundleHR := fmt.Sprintf("%s-%s", clusterName, bundleName)
	appHR := fmt.Sprintf("%s-%s", clusterName, state.GetApplication().AppName)

	By("checking the bundle HelmRelease is ready")
	Eventually(func() (bool, error) {
		return client.IsHelmReleaseReady(ctx, bundleHR, namespace)
	}).
		WithTimeout(5*time.Minute).
		WithPolling(5*time.Second).
		Should(BeTrue(), client.FailureDiagnostics("bundle HelmRelease '%s/%s' is not ready", namespace, bundleHR))

	By("checking the app HelmRelease is ready")
	Eventually(func() (bool, error) {
		return client.IsHelmReleaseReady(ctx, appHR, namespace)
	}).
		WithTimeout(15*time.Minute).
		WithPolling(5*time.Second).
		Should(BeTrue(), client.FailureDiagnostics("HelmRelease '%s/%s' is not ready", namespace, appHR))

	By("checking the app HelmRelease is at the version under test")
	Eventually(func() (bool, error) {
		return client.IsHelmReleaseVersion(ctx, appHR, namespace, state.GetApplication().Version)
	}).
		WithTimeout(5*time.Minute).
		WithPolling(5*time.Second).
		Should(BeTrue(), client.FailureDiagnostics("HelmRelease '%s/%s' is not at version '%s'", namespace, appHR, state.GetApplication().Version))
}
