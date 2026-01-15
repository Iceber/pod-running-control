# Pod Running Control
The `pod-running-control` prevents the start of business containers by running as an init-container within the pod.

This is suitable for scenarios where you want to reserve node resources without immediately starting business containers.

## Highlights
* Can be used with **any versions** of Kubernetes.
* Extremely lightweight and **non-intrusive** to business containers

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

## Roadmap
This is a project with a very clear objective.

While the core logic of the code is not complex, we still need to address and optimize engineering and adaptation efforts.

* Compatibility with CEL, support configuring arbitrary resource types and conditional fields for run gates, so pod running control will not depend on fixed CRDs.
* Provide additional engineering examples, such as:
    * Mounting a dedicated service account for the running-control init-container
    * Completing pre-loading of business containers during the init-container phase
    * Implementing pod group running control
* Utilize webhooks to detect annotations and automate the injection of running-control init-containers
* **Open to further ideas**
