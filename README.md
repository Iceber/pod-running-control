# Pod Running Control
The `pod-running-control` prevents the start of business containers by running as an init-container within the pod.

This is suitable for scenarios where you want to reserve node resources without immediately starting business containers.

## Highlights
* Can be used with **any versions** of Kubernetes.
* Extremely lightweight and **non-intrusive** to business containers.
* Does not depend on CRD, supports [CEL](https://kubernetes.io/docs/reference/using-api/cel/) and monitors any resource type.

## Quick Start
1. **Apply RBAC** for `running-control` init container to monitor resources.
```bash
$ kubectl apply -f ./examples/rbac.yaml
```

2. **Apply Pod**
```yaml
apiVersion: v1
kind: Pod
metadata:
  name: test-pod
  annotations:
    pod-running-control.io/break: 'true'
spec:
  initContainers:
  - name: running-control
    image: ghcr.io/iceber/pod-running-control:latest
    args:
      - "--gate-gvr=pods.v1."
      - "--gate-namespace=default"
      - "--gate-name=test-pod"
      - |
        --gate-expression=!has(object.metadata.annotations) || !('pod-running-control.io/break' in object.metadata.annotations) || object.metadata.annotations['pod-running-control.io/break'] != 'true'
  containers:
  - name: nginx
    image: nginx:1.14.2
    ports:
    - containerPort: 80
```
The `nginx` container is not running now, the pod is blocked at the init-container stage.

3. It's time to end the blocking and get the `nginx` container running
```bash
$ kubectl annotate pod test-pod pod-running-control.io/break='false'
```


## Roadmap
This is a project with a very clear objective.

While the core logic of the code is not complex, we still need to address and optimize engineering and adaptation efforts.

* Provide additional engineering examples, such as:
    * Mounting a dedicated service account for the running-control init-container
    * Completing pre-loading of business containers during the init-container phase
    * Implementing pod group running control
* Utilize webhooks to detect annotations and automate the injection of running-control init-containers
* Based on feedback, support for more complex CELs
* **Open to further ideas**
