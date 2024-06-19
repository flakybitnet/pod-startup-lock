# Simple time-based lock service with HTTP interface.

### Designed at [Oath](https://www.oath.com) to solve the [Thundering herd problem](https://en.wikipedia.org/wiki/Thundering_herd_problem) during multiple applications startup in the [Kubernetes](https://kubernetes.io) cluster. 

## The Problem

Starting multiple applications simultaneously on the same host may cause a performance bottleneck.
In Kubernetes this usually happens when applications are automatically deployed to a newly added Node.
In the worst-case scenario, application startup may be slowed down so dramatically that they fail to pass the healthcheck. 
They are then restarted by Kubernetes just to start fighting for shared resources again, in an endless loop.

See also [Containers startup throttling](https://github.com/kubernetes/kubernetes/issues/3312) Kubernetes issue.

## The Solution

Kubernetes allows a Pod to have additional, [Init container](https://kubernetes.io/docs/concepts/workloads/pods/init-containers/#examples),
and postpone application startup until Init container finishes execution.
The solution is to deploy Lock service as a DaemonSet on a Pod, and each init container will sequentially acquire this lock.
So moments of application container starts will be distributed in time.

## Components

See Readmes in subfolders for details.

* [Lock](lock/README.md)
  
  HTTP service to be deployed one instance per Node (as a DaemonSet).
  Returns code `200 OK` as a response to the first request.
  Returns `423 Locked` to the subsequent requests until timeout exceeded.
  May depend on additional endpoint check.


* [Init](init/README.md)

  Lightweight client for the Lock service. To be deployed as Init Container alongside the main application container.
  Periodically tries to acquire the lock. Once succeeded, terminates, allowing the main container to start running.
  

* [K8s-health](k8s-health/README.md)
  
  Optional component. Performs healthcheck of Kubernetes DaemonSets.
  May be used by Lock service to postpone lock acquiring until all DaemonSets on the Node are up and running.

## Images

The OCI images are available at 
* **Quay**
  * [K8s health](https://quay.io/repository/flakybitnet/psl-k8s-health)
  * [Lock](https://quay.io/repository/flakybitnet/psl-lock)
  * [Init](https://quay.io/repository/flakybitnet/psl-init)
* **GHCR**
  * [K8s health](https://github.com/flakybitnet/pod-startup-lock/pkgs/container/psl-k8s-health)
  * [Lock](https://github.com/flakybitnet/pod-startup-lock/pkgs/container/psl-lock)
  * [Init](https://github.com/flakybitnet/pod-startup-lock/pkgs/container/psl-init)
* **AWS ECR Public**
  * [K8s health](https://gallery.ecr.aws/flakybitnet/psl/k8s-health)
  * [Lock](https://gallery.ecr.aws/flakybitnet/psl/lock)
  * [Init](https://gallery.ecr.aws/flakybitnet/psl/init)
* **FlakyBit's Harbor**

## How to build locally

1.  Set target platform for Go binaries.
    
    Optional step, default is `linux`. You can use Make task as shown below.

    ```bash
    export GOOS=darwin
    ```

2.  Change directory.

    Open directory with a service you want to build.

    ```bash
    # cd <service directory>
    cd init
    # cd lock
    # cd k8s-health
    ```

3.  Build a binary.

    ```bash
    # go build -v -a -o <output binary file path>
    go build -v -a -o ../bin/init
    # go build -v -a -o ../bin/lock
    # go build -v -a -o ../bin/health
    ```

4.  Obtain the binaries.

    The binaries will be located in `bin` directory of project's root:

   * `pod-startup-lock/bin/init`
   * `pod-startup-lock/bin/health`
   * `pod-startup-lock/bin/lock`

## Release Notes

* `1.0.1`
    - Added connection timeouts for http and tcp connections
    - Added keep-alive for http connections
* `1.0.0`
    - Initial version
    
## Contributing

Please feel free to submit issues, fork the repository and send pull requests!

## Source

Source code is available at [Gitea](https://gitea.flakybit.net/flakybit/pod-startup-lock)
and mirrored to [GitHub](https://github.com/flakybitnet/pod-startup-lock).