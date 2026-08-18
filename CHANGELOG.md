# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Changed

- Apply the CRDs from a dedicated installer image via a `pre-install`/`pre-upgrade` hook Job instead of rendering them as Helm templates. The rendered release no longer contains the CRDs, which kept it from fitting in Helm's release Secret.
- CRDs are applied with server-side apply and are no longer Helm-managed, so `helm uninstall` leaves them in place and a CRD rejected by the `safe-upgrades` policy now surfaces as a failed hook Job.
- Reject unknown keys under `install`, which previously were silently ignored.

### Added

- E2E test suite covering CRD installation, installer hook cleanup, the `safe-upgrades` admission policy and the served Gateway API resources.
- `crds.image` values to override the installer image. The tag defaults to the chart version.
- `CiliumNetworkPolicy` allowing the installer Job to reach the Kubernetes API server, without which the Job cannot apply the CRDs on clusters that default to denying egress. Set `ciliumNetworkPolicy.enabled` to `false` to skip it.

## [1.9.0] - 2026-08-17

### Changed

- Upgrade Gateway API CRDs to [v1.6.1](https://github.com/kubernetes-sigs/gateway-api/releases/tag/v1.6.1)
- Upgrade Gateway API Inference Extension CRDs to [v1.5.0](https://github.com/kubernetes-sigs/gateway-api-inference-extension/releases/tag/v1.5.0)
- Set `install.tcproutes` and `install.udproutes` to `standard`, following their graduation to the standard channel. Required by Envoy Gateway 1.9, which reconciles them via `gateway.networking.k8s.io/v1`.
- Install the `safe-upgrades` `ValidatingAdmissionPolicy`, which was previously dropped and left its binding orphaned. Set `install.admissionPolicies` to `false` to skip both.

### Added

- Add `install.xbackends` value for the new experimental XBackend CRD.
- Add `install.inferencepoolimports` and `install.inferencemodelrewrites` values, whose CRDs shipped without a way to enable them.

### Removed

- Remove the alpha `InferencePool` CRD (`inference.networking.x-k8s.io/v1alpha2`), dropped upstream in Inference Extension v1.5.0. `install.inferencepools` now only accepts `standard` or `""`.

## [1.8.1] - 2026-05-28

### Changed

- Update `install.listenersets` value to `standard` (renamed from `xlistenersets`)
- Update `install.tlsroutes` value to `standard`

## [1.8.0] - 2026-05-22

- Upgrade Gateway API Inference Extension CRDs to [v1.4.0](https://github.com/kubernetes-sigs/gateway-api-inference-extension/releases/tag/v1.4.0)
- Upgrade Gateway API CRDs to [v1.5.1](https://github.com/kubernetes-sigs/gateway-api/releases/tag/v1.5.1)
- Add `install.admissionPolicies` value (default `true`) to control installation of `ValidatingAdmissionPolicyBinding` resources

## [1.7.0] - 2026-02-13

### Changed

- Upgrade [Gateway API Inference Extension](https://gateway-api-inference-extension.sigs.k8s.io/) CRDs to [v1.3.0](https://github.com/kubernetes-sigs/gateway-api-inference-extension/releases/tag/v1.3.0)
- Restore Gateway API CRDs from [v1.4.1](https://github.com/kubernetes-sigs/gateway-api/releases/tag/v1.4.1) release.

## [1.6.1] - 2026-01-12

### Changed

- Upgrade Gateway API CRDs to [v1.4.1](https://github.com/kubernetes-sigs/gateway-api/releases/tag/v1.4.1)

## [1.6.0] - 2025-11-04

### Changed

- Upgrade [Gateway API Inference Extension](https://gateway-api-inference-extension.sigs.k8s.io/) CRDs to [v1.1.0](https://github.com/kubernetes-sigs/gateway-api-inference-extension/releases/tag/v1.1.0)

## [1.5.1] - 2025-10-20

### Changed

- Upgrade [Gateway API Inference Extension](https://gateway-api-inference-extension.sigs.k8s.io/) CRDs to [v1.0.2](https://github.com/kubernetes-sigs/gateway-api-inference-extension/releases/tag/v1.0.2)

## [1.5.0] - 2025-10-14

### Changed

- Upgrade Gateway API CRDs to [v1.4.0](https://github.com/kubernetes-sigs/gateway-api/releases/tag/v1.4.0)
  **We strongly suggest reading the upstream release notes.**
- Update `install.backendtlspolicies` value to `standard`
- Add `install.xmeshes` value

## [1.4.0] - 2025-10-14

### Changed

- Upgrade Gateway API CRDs to [v1.3.0](https://github.com/kubernetes-sigs/gateway-api/releases/tag/v1.3.0)
  **We strongly suggest reading the upstream release notes.**
- Remove `install.backendlbpolicies` value
- Add `install.xbackendtrafficpolicies` and `install.xlistenersets` values

## [1.3.0] - 2025-10-14

### Changed

- Update `appVersion` to reflect versions of included Gateway API and Gateway API Inference Extension CRDs

### Added

- Add [Gateway API Inference Extension](https://gateway-api-inference-extension.sigs.k8s.io/) [v1.0.1](https://github.com/kubernetes-sigs/gateway-api-inference-extension/releases/tag/v1.0.1) CRDs

## [1.2.1] - 2025-06-10

### Changed

- Add `helm.sh/resource-policy` annotation to all CRDs to prevent deletion by default

## [1.2.0] - 2025-02-05

### Added

- Allow selecting channel ("standard" or "experimental") for each individual CRD

### Changed

- Upgrade Gateway API CRDs to [v1.2.1](https://github.com/kubernetes-sigs/gateway-api/releases/tag/v1.2.1)

## [1.1.0] - 2024-05-13

### Changed

- Upgrade Gateway API CRDs to [v1.1.0](https://github.com/kubernetes-sigs/gateway-api/releases/tag/v1.1.0)

## [1.0.0] - 2024-01-15

### Added
- Initial release

[Unreleased]: https://github.com/giantswarm/gateway-api-crds-app/compare/v1.9.0...HEAD
[1.9.0]: https://github.com/giantswarm/gateway-api-crds-app/compare/v1.8.1...v1.9.0
[1.8.1]: https://github.com/giantswarm/gateway-api-crds-app/compare/v1.8.0...v1.8.1
[1.8.0]: https://github.com/giantswarm/gateway-api-crds-app/compare/v1.7.0...v1.8.0
[1.7.0]: https://github.com/giantswarm/gateway-api-crds-app/compare/v1.6.1...v1.7.0
[1.6.1]: https://github.com/giantswarm/gateway-api-crds-app/compare/v1.6.0...v1.6.1
[1.6.0]: https://github.com/giantswarm/gateway-api-crds-app/compare/v1.5.1...v1.6.0
[1.5.1]: https://github.com/giantswarm/gateway-api-crds-app/compare/v1.5.0...v1.5.1
[1.5.0]: https://github.com/giantswarm/gateway-api-crds-app/compare/v1.4.0...v1.5.0
[1.4.0]: https://github.com/giantswarm/gateway-api-crds-app/compare/v1.3.0...v1.4.0
[1.3.0]: https://github.com/giantswarm/gateway-api-crds-app/compare/v1.2.1...v1.3.0
[1.2.1]: https://github.com/giantswarm/gateway-api-crds-app/compare/v1.2.0...v1.2.1
[1.2.0]: https://github.com/giantswarm/gateway-api-crds-app/compare/v1.1.0...v1.2.0
[1.1.0]: https://github.com/giantswarm/gateway-api-crds-app/compare/v1.0.0...v1.1.0
[1.0.0]: https://github.com/giantswarm/gateway-api-crds-app/releases/tag/v1.0.0
