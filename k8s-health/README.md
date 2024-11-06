# Kubernetes DaemonSet health check service

Typically, you would like to postpone application startup until Node is healthy.
There are two types of checks:
* DaemonSets on Node are ready
* Node CPU is not overloaded 

This util constantly performs health checks and respond with `200 OK` if they passed.

Just add this service as a dependent endpoint to the Lock service, and lock won't be acquired until Node is healthy.  

## How it works

##### 1. Starts and listens for HTTP requests

Responds with `412 Precondition Failed` until healthcheck succeeds.
Binding `host` and `port` are configurable.

##### 2. Uses Kubernetes API to get a list of DaemonSets and Node load

`NODE_NAME` environment variable must be set.
No additional authentication required to access Kubernetes API if running as a cluster resource.

##### 3. When all DaemonSet Pods are up and running and Node CPU utilisation below defined threshold, return success

Responds with `200 OK`. Constantly repeats health checks.

## Which DaemonSets to check

You may need to check only certain DaemonSets and ignore the others:

* Add `PSL_HC_DAEMONSET_HOST_NETWORK=true` env var to check only ones with binding to the host network, i.e. with `spec.template.spec.hostNetwork: true`
* Use `PSL_HC_DAEMONSET_NAMESPACE=XXX` to check only ones in a specific namespace. All namespaces by default.
* List excluded labels with `PSL_HC_DAEMONSET_EXCLUDE_LABELS` env var. DaemonSets having **at least one** matching label will be excluded.  
* **OR**, list included labels with `PSL_HC_DAEMONSET_INCLUDE_LABELS` flag. DaemonSets having **all** matching labels will be included, rest excluded.
  You can't specify both `PSL_HC_DAEMONSET_EXCLUDE_LABELS` and `PSL_HC_DAEMONSET_INCLUDE_LABELS` flags, choose one.  

## In Cluster / Out Of Cluster configuration

[`kubernetes-go-client`](https://github.com/kubernetes/client-go) is used under the hood.

By default, `kubernetes-go-client` uses *in-cluster* configuration.
It will try to create configuration basing on cluster info from the running pod.
If you set `PSL_K8S_API_URL` parameter, the *out-of-cluster* configuration will be used,
assuming you are running k8s proxy with `kubectl`.

## Configuration

You may specify environment variables to override defaults:

| Option                            | Default | Required | Description                                                                                           |
|-----------------------------------|---------|----------|-------------------------------------------------------------------------------------------------------|
| `PSL_BIND_HOST`                   | 0.0.0.0 |          | Address to bind                                                                                       |
| `PSL_BIND_PORT`                   | 8080    |          | Port to bind                                                                                          |
| `PSL_NODE_NAME`                   |         | +        | K8s node name which the current app instance runs on; to indicate which node health should be checked |
| `PSL_K8S_API_URL`                 |         |          | K8s API URL, for out-of-cluster usage only                                                            |
| `PSL_HC_DAEMONSET_ENABLED`        | true    |          | Enabled DaemonSets health check                                                                       |
| `PSL_HC_DAEMONSET_NAMESPACE`      |         |          | Target K8s namespace where to perform DaemonSets healthcheck, leave blank for all namespaces          |
| `PSL_HC_DAEMONSET_HOST_NETWORK`   | false   |          | Check only DaemonSets bind to the `host network`                                                      |
| `PSL_HC_DAEMONSET_INCLUDE_LABELS` |         |          | DaemonSet labels to include in healthcheck, `label1:value1,label2:value2`                             |
| `PSL_HC_DAEMONSET_EXCLUDE_LABELS` |         |          | DaemonSet labels to exclude from healthcheck, `label1:value1,label2:value2`                           |
| `PSL_HC_DAEMONSET_PERIOD_FAIL`    | 10s     |          | Period of health checks if previous failed                                                            |
| `PSL_HC_DAEMONSET_PERIOD_PASS`    | 60s     |          | Period of health checks if previous succeeded                                                         |
| `PSL_HC_NODELOAD_ENABLED`         | false   |          | Enabled Node load health check                                                                        |
| `PSL_HC_NODELOAD_CPU_THRESHOLD`   | 80      |          | Node CPU utilisation in percent above which it is treated as unhealthy                                |
| `PSL_HC_NODELOAD_PERIOD`          | 10s     |          | Period of health checks                                                                               |
| `PSL_LOG`                         | info    |          | Log level                                                                                             |

## How to run locally

Examples with some options:

```bash
kubectl proxy -p 57585
PSL_K8S_API_URL="http://127.0.0.1:57585" \
  NODE_NAME="10.11.10.11" \
  PSL_HC_DAEMONSET_INCLUDE_LABELS="app:test,version:1.1" \
  PSL_HC_DAEMONSET_HOST_NETWORK=true \
  bin/health
```

```bash
NODE_NAME="k3s-agent-2" \
  PSL_HC_NODELOAD_ENABLED=true \
  PSL_HC_NODELOAD_CPU_THRESHOLD=60 \
  bin/health
```


## How to deploy to Kubernetes

The preferable way is to deploy as a DaemonSet.
You can find example deployment in [.k8s/k8s-health](../.k8s/k8s-health) directory.
