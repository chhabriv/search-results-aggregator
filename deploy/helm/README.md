## Deployment

### Pre-requisites

- Helm v3.10.3
- Kubernetes v1.25.2

Create a namespace:

```shell
kubectl create namespace assignment132
```

Verify the namespace has been created:

```shell
kubectl get namespaces
```

Set the created namespace as default:

```shell
kubectl config set-context --current --namespace=assignment132
```

To pull images during deployment from a private registry:

```shell
kubectl -n assignmen132 create secret docker-registry regcred --docker-server ghcr.io --docker-username REGISTRY_USERNAME --docker-password REGISTRY_PASSWORD_OR_TOKEN  --docker-email EMAIL_ADDRESS
```

Verify the registry credentials have been created:

```shell
kubectl get secrets
```

### Deploy the service using helm 3

To deply the service:

```shell
helm upgrade --install --namespace=assignment132 search-results-aggregator search-results-aggregator/
```

Verify the deployments have been created:

```shell
kubectl get deployments
```

Verify the service has been created:

```shell
kubectl get svc
```

Verify the pods have been created:

```shell
kubectl get po
```

Verify the hpa has been created:

```shell
kubectl get hpa
```

To view pod logs:
```shell
kubectl logs POD_NAME
```

example:
```shell
kubectl logs search-results-aggregator-bd67957f-48ws8
```

### Connect to the deployed service

Requests can be made to the deployed service using port-forwarding.
To enable port forwarding from your machine to the service:

```shell
kubectl port-forward --namespace assignment132 $(kubectl get pod --namespace assignment132 --selector="app=search-results-aggregator" --output jsonpath='{.items[0].metadata.name}') 8080:8080
```

Alternatively, if the pod name is known:

```shell
kubectl port-forward --namespace assignment132 search-results-aggregator-bd67957f-48ws8 8080:8080
```

### Delete a deployed service

To delete a helm release:

```shell
helm delete search-results-aggregator
```
