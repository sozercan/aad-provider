# AAD provider

- Deploy Gatekeeper with external data enabled

- `kubectl apply -f manifest`

- `kubectl apply -f policy/provider.yaml`
  - Update `proxyURL` if it's not `http://aad-provider.default:8090`

- `kubectl apply -f policy/assignmetadata.yaml`

- `kubectl apply -f examples/test.yaml`

- `kubectl get deploy test-deployment -o yaml`
