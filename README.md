# coco-play tool

## Introduction

The Confidential Containers (CoCo) is built upon hardware-based Trusted Execution Environment (TEE) technologies that still aren't widespread on developers workstation. This represents a barrier for those who want to develop, test or otherwise are curious about CoCo. Fortunately, we've circumvent that by providing a mechanism to run CoCo on any environment without the use of confidential hardware. This is accomplished with
a sample attester and verifier that bypasses the TEE-based attestation, so that pods created with the `kata-qemu-coco-dev` [runtimeClass](https://kubernetes.io/docs/concepts/containers/runtime-class/) will
run as if they are on TEE. Obviously any pod on this circumstance isn't strictly confidential but it will behave as such, then users are able to "simulate" actions like running a pod from an encrypted signed image. 

The `coco-play` tool is meant to quickly build an environment, called `playground`, on users workstations to allow them run pods with `kata-qemu-coco-dev` runtimeClass easily. 

## What's a coco-play playground

The tool creates a managed environment that consists of a Kubernetes in Docker ([kind](https://kind.sigs.k8s.io/)) cluster in which the CoCo stack as well as a Key Broker Server (KBS) are installed. Users interact with the cluster using regular Kubernetes tools such as [`kubectl`](https://kubernetes.io/docs/reference/kubectl/).

It's only required the following softwares in the system:

* Docker
* Qemu/KVM
* kubectl

## How to use coco-play

`coco-play` has a couple of commands, run `coco-play -h` for a completed and up-to-date list.

You should start by creating a playground. Run `play-create` command as below:

```shell
$ ./coco-play play-create
Creating Kind cluster (coco-play)...
Labeling worker...
node/coco-play-control-plane labeled

Creating CoCo controller...
namespace/confidential-containers-system created
customresourcedefinition.apiextensions.k8s.io/ccruntimes.confidentialcontainers.org created
serviceaccount/cc-operator-controller-manager created
role.rbac.authorization.k8s.io/cc-operator-leader-election-role created
clusterrole.rbac.authorization.k8s.io/cc-operator-manager-role created
clusterrole.rbac.authorization.k8s.io/cc-operator-metrics-reader created
clusterrole.rbac.authorization.k8s.io/cc-operator-proxy-role created
rolebinding.rbac.authorization.k8s.io/cc-operator-leader-election-rolebinding created
clusterrolebinding.rbac.authorization.k8s.io/cc-operator-manager-rolebinding created
clusterrolebinding.rbac.authorization.k8s.io/cc-operator-proxy-rolebinding created
configmap/cc-operator-manager-config created
service/cc-operator-controller-manager-metrics-service created
deployment.apps/cc-operator-controller-manager created

Waiting for deployment "cc-operator-controller-manager" rollout to finish: 0 of 1 updated replicas are available...
deployment "cc-operator-controller-manager" successfully rolled out

Creating CoCo ccruntime CRD...
ccruntime.confidentialcontainers.org/ccruntime-sample created

Install KBS...
Enumerating objects: 5727, done.
Counting objects: 100% (56/56), done.
Compressing objects: 100% (24/24), done.
Total 5727 (delta 35), reused 32 (delta 32), pack-reused 5671 (from 2)
namespace/coco-tenant created
configmap/kbs-config-68g582hg7m created
configmap/policy-config-t6mb856fbh created
secret/kbs-auth-public-key-dh6km9bbdh created
secret/keys-gmf2g75547 created
service/kbs created
deployment.apps/kbs created

```

At the end and if everything went well, it will have a fully running playground. The kubeconfig file is appended to `~/.kube/config` so that you can now use `kubectl` to issue commands targetting the playground cluster. For example, following lists the installed `runtimeClass`:

```shell
$ kubectl --kubeconfig ~/.kube/config --cluster kind-coco-play get runtimeclass
NAME                 HANDLER              AGE
kata                 kata-qemu            4h21m
kata-clh             kata-clh             4h21m
kata-qemu            kata-qemu            4h21m
kata-qemu-coco-dev   kata-qemu-coco-dev   4h21m
kata-qemu-sev        kata-qemu-sev        4h21m
kata-qemu-snp        kata-qemu-snp        4h21m
kata-qemu-tdx        kata-qemu-tdx        4h21m
```

The KBS is installed in the playground. You can use the `kbs-info` command to get information, mainly the service address:

```shell
$ ./coco-play kbs-info
Status: Running
Service address: 172.18.0.2:31207
```

With the address (`172.18.0.2:31207` on example above) in hands, you can try out launching a pod to fetch a resource key from KBS:

```shell
$ cat <<EOF>>coco-demo.yaml
apiVersion: v1
kind: Pod
metadata:
  name: coco-demo
  annotations:
    "io.containerd.cri.runtime-handler": "kata-qemu-coco-dev"
    io.katacontainers.config.hypervisor.kernel_params: " agent.aa_kbc_params=cc_kbc::http://172.18.0.2:31207"
spec:
  runtimeClassName: kata-qemu-coco-dev
  containers:
    - name: busybox
      image: quay.io/prometheus/busybox:latest
      imagePullPolicy: Always
      command:
        - sh
        - -c
        - |
          wget -O- http://127.0.0.1:8006/cdh/resource/reponame/workload_key/key.bin; sleep infinity
  restartPolicy: Never
EOF
$ kubectl apply -f coco-demo.yaml
pod/coco-demo created
$ kubectl wait --for=condition=Ready pod/coco-demo
pod/coco-demo condition met
$ kubectl logs coco-demo
Connecting to 127.0.0.1:8006 (127.0.0.1:8006)
somesecret
writing to stdout
-                    100% |********************************|    11  0:00:00 ETA
written to stdout
```

To add or update a resource key in the KBS, use the `kbs-set-resource` command. For example:

```shell
$ echo "somesecret" > key.txt
$ ./coco-play kbs-set-resource default/tests/key key.txt
```

You may want to delete the playground when you are done. Run the `play-delete` command like below to delete the cluster:

```shell
$ ./coco-play play-delete
Deleting Kind cluster (coco-play)...
```

## How to develop this tool

In order to build this tool you must have `go` installed in your environment, then run:

```shell
go build
```

It's also recommended to have [kind](https://kind.sigs.k8s.io/) installed to use the tools bundled in Kind to allow better debugging of the cluster created if needed.
