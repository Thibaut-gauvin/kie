# KIE Kubernetes Image Exporter

## About
`kubernetes-image-exporter`, e.g. `kie` exports Prometheus metrics about image usage in your cluster.

## Install

### Helm
```sh
git clone git@github.com:Thibaut-gauvin/kie.git
cd kie
helm upgrade -i kie -n default ./charts/kubernetes-image-exporter
helm test kie -n default
```

## Usages

```sh
kie -h
Available Commands:
  help        Help about any command
  serve       Start kie server.
  version     print current kie version

Flags:
  -h, --help               help for kie
  -l, --log-level string   Log level. Can be any standard log-level ("info", "debug", etc...) (default "info")
```

```sh
kie serve -h
Start kie server.

Usage:
  kie serve [flags]

Flags:
  -h, --help                   help for serve
  -k, --kubeconfig string      path to kubeconfig used to authenticate with cluster when your running kie locally. If not provided, use service-account from pod.
  -p, --listen-port string     the listening port that kie server will use. (default "9145")
  -i, --refresh-interval int   metrics values refresh interval in seconds. (default 30)

Global Flags:
  -l, --log-level string   Log level. Can be any standard log-level ("info", "debug", etc...) (default "info")
```


## Exported metrics


- `kie_cluster_images_running`  
   labels: `{"name", "tag", "digest", "pod"}`  
   type: `gauge`
