package basic

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/giantswarm/apptest-framework/v5/pkg/state"
	"github.com/giantswarm/clustertest/v5/pkg/logger"
	"github.com/giantswarm/clustertest/v5/pkg/wait"

	"k8s.io/apimachinery/pkg/api/errors"
)

// deploymentAppTests verifies the bundle and the CRD app itself are deployed at the version
// under test before any CRD specific assertion runs.
func deploymentAppTests() {
	By("checking the bundle application is created")
	Expect(state.GetBundleApplication()).ToNot(BeNil())
	Expect(state.GetBundleApplication().AppName).ToNot(Equal(state.GetApplication().AppName))

	By("checking the bundle app is deployed")
	Eventually(wait.IsAppDeployed(state.GetContext(), state.GetFramework().MC(), state.GetBundleApplication().InstallName, state.GetBundleApplication().InstallNamespace)).
		WithTimeout(30 * time.Second).
		WithPolling(50 * time.Millisecond).
		Should(BeTrue())

	By("checking the test app is deployed")
	Eventually(func() (bool, error) {
		done, err := wait.IsAppDeployed(state.GetContext(), state.GetFramework().MC(), state.GetApplication().InstallName, state.GetApplication().Organization.GetNamespace())()
		if err != nil {
			if errors.IsNotFound(err) {
				logger.Log("App '%s/%s' doesn't exist yet", state.GetApplication().Organization.GetNamespace(), state.GetApplication().InstallName)
				return false, nil
			}
			return false, err
		}

		return done, nil
	}).
		WithTimeout(15 * time.Minute).
		WithPolling(5 * time.Second).
		Should(BeTrue())

	By("checking the test app is deployed at the correct version")
	Eventually(wait.IsAppVersion(state.GetContext(), state.GetFramework().MC(), state.GetApplication().InstallName, state.GetApplication().Organization.GetNamespace(), state.GetApplication().Version)).
		WithTimeout(5 * time.Minute).
		WithPolling(5 * time.Second).
		Should(BeTrue())
}
