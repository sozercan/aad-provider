# Azure Active Directory (AAD) provider

Azure Active Directory (AAD) provider is used for mutating Kubernetes labels to display name of the AAD user using Microsoft Graph API.

> This repo is meant for testing Gatekeeper external data feature. Do not use for production.

- Make sure you have a Kubernetes user that matches the AAD user you want to query (e.g. `user@example.com`).

- Create a new service principal and assign necessary Microsoft Graph API permissions to read user profile (`profile` and `User.Read.All`).

- Add your tenant ID, client ID, client secret to the `manifest/secret.yaml` file.

- Deploy Gatekeeper with external data enabled (`--enable-external-data`).

# Installation

- `kubectl apply -f manifest`

- `kubectl apply -f policy/provider.yaml`
  - Update `proxyURL` if it's not `http://aad-provider.default:8090` (default)

- `kubectl apply -f policy/assignmetadata.yaml`

# Mutation

- `kubectl apply -f examples/test.yaml`

- `kubectl get deploy test-deployment -o yaml`
  - You should see `owners` label filled with your AAD display name.
    ```
    $ kubectl get cm test-configmap -o yaml
    apiVersion: v1
    kind: ConfigMap
    metadata:
      labels:
        owner: Sertac_Ozercan
      name: test-configmap
      namespace: default
    ```
