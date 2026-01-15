# Pod Running Control
The `pod-running-control` prevents the start of business containers by running as an init-container within the pod.

This is suitable for scenarios where you want to reserve node resources without immediately starting business containers.

## Examples
1. Apply CRDs and RBAC
```bash
$ kubectl apply -f ./manifests/crds
$ kubectl apply -f ./manifests/rbac.yaml
```

2. **Apply Pod**
```yaml
apiVersion: pod-running-control.io/v1alpha1
kind: PodRunningGate
metadata:
  name: test-pod
spec:
  gates:
  - "test-block"
---
apiVersion: v1
kind: Pod
metadata:
  name: test-pod
spec:
  initContainers:
  - name: running-control
    image: ghcr.io/iceber/pod-running-control:latest
    env:
    - name: POD_RUNNING_GATE_NAMESPACE
      valueFrom:
        fieldRef:
          fieldPath: metadata.namespace
    - name: POD_RUNNING_GATE_NAME
      valueFrom:
        fieldRef:
          fieldPath: metadata.name
  containers:
  - name: nginx
    image: nginx:1.14.2
    ports:
    - containerPort: 80
```
