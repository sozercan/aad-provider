# Azure Active Directory (AAD) provider

- Create a new service principal and assign necesssary Graph API permissions to read user profile.

- Add tenant ID, client ID, client secret to the `manifest/secrets.yaml` file.

- Deploy Gatekeeper with external data enabled (`--enable-external-data`)

# Installation

- `kubectl apply -f manifest`

- `kubectl apply -f policy/provider.yaml`
  - Update `proxyURL` if it's not `http://aad-provider.default:8090` (default)

- `kubectl apply -f policy/assignmetadata.yaml`

# Verification

- `kubectl apply -f examples/test.yaml`

- `kubectl get deploy test-deployment -o yaml`
