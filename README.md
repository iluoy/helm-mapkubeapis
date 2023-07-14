## Usage

### Map Helm deprecated or removed Kubernetes APIs

Map release deprecated or removed Kubernetes APIs in-place:

```console
$ helm mapkubeapis -h
Map release deprecated or removed Kubernetes APIs in-place

Usage:
  mapkubeapis [flags]

Flags:
  -A, --all-namespaces                       map kube api of all releases across all namespaces
      --dry-run                              simulate a command
      --except-namespaces strings            except multiple namespaces, for example: --except-namespaces NS1 NS2
      --except-releases-namespaces strings   except multiple releases namespaces, for example: --except-releases-namespaces Release1.NS1 Release2.NS2
  -h, --help                                 help for mapkubeapis
      --kube-context string                  name of the kubeconfig context to use
      --kubeconfig string                    path to the kubeconfig file
      --mapfile string                       path to the API mapping file (default "config/Map.yaml")
      --namespaces strings                   multiple namespaces, for example: --namespaces NS1 NS2
      --releases-namespaces strings          multiple releases, for example: --releases-namespaces Release1.NS1 Release2.NS2

```

Example:

```console
$ helm mapkubeapis mapkubeapis --namespaces "ns1,ns2" --releases-namespaces release3.ns3 --except-namespaces ns4 --except-releases-namespaces release5.ns2 --except-releases-namespaces release3.ns3 -A

```

### ToDo
improve the performance
