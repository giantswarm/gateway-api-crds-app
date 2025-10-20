# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [1.5.1] - 2025-10-20

### Changed

- Upgrade [Gateway API Inference Extension](https://gateway-api-inference-extension.sigs.k8s.io/) CRDs to [v1.0.2](https://github.com/kubernetes-sigs/gateway-api-inference-extension/releases/tag/v1.0.1)

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

[Unreleased]: https://github.com/giantswarm/gateway-api-crds-app/compare/v1.5.1...HEAD
[1.5.1]: https://github.com/giantswarm/gateway-api-crds-app/compare/v1.5.0...v1.5.1
[1.5.0]: https://github.com/giantswarm/gateway-api-crds-app/compare/v1.4.0...v1.5.0
[1.4.0]: https://github.com/giantswarm/gateway-api-crds-app/compare/v1.3.0...v1.4.0
[1.3.0]: https://github.com/giantswarm/gateway-api-crds-app/compare/v1.2.1...v1.3.0
[1.2.1]: https://github.com/giantswarm/gateway-api-crds-app/compare/v1.2.0...v1.2.1
[1.2.0]: https://github.com/giantswarm/gateway-api-crds-app/compare/v1.1.0...v1.2.0
[1.1.0]: https://github.com/giantswarm/gateway-api-crds-app/compare/v1.0.0...v1.1.0
[1.0.0]: https://github.com/giantswarm/gateway-api-crds-app/releases/tag/v1.0.0
