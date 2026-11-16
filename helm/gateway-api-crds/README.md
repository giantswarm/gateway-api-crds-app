# gateway-api-crds

App to install Gateway API CRDs to Giantswarm Clusters

**Homepage:** <https://github.com/giantswarm/gateway-api-crds-app>

## Source Code

* <https://github.com/kubernetes-sigs/gateway-api>
* <https://github.com/kubernetes-sigs/gateway-api-inference-extension>

## Requirements

Kubernetes: `>=1.25.0-0`

## Values

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| install.gatewayclasses | string | `"standard"` |  |
| install.gateways | string | `"standard"` |  |
| install.grpcroutes | string | `"standard"` |  |
| install.httproutes | string | `"standard"` |  |
| install.referencegrants | string | `"standard"` |  |
| install.backendtlspolicies | string | `"standard"` |  |
| install.listenersets | string | `"standard"` |  |
| install.tcproutes | string | `"standard"` |  |
| install.tlsroutes | string | `"standard"` |  |
| install.udproutes | string | `"standard"` |  |
| install.xbackends | string | `""` |  |
| install.xbackendtrafficpolicies | string | `""` |  |
| install.xmeshes | string | `""` |  |
| install.inferencepools | string | `""` |  |
| install.inferenceobjectives | string | `""` |  |
| install.inferencepoolimports | string | `""` |  |
| install.inferencemodelrewrites | string | `""` |  |
| install.admissionPolicies | bool | `true` |  |
| crds.image.registry | string | `"gsoci.azurecr.io"` |  |
| crds.image.repository | string | `"giantswarm/gateway-api-crds"` |  |
| crds.image.tag | string | `""` |  |
| ciliumNetworkPolicy.enabled | bool | `true` |  |
