# Lock service

Service to manage locks. Releases the lock after configured timeout.
Uses `k8s-health` service to postpone acquiring the locks.
`init` container uses this service to acquire the lock.

## How it works

##### 1. Starts and listens for HTTP requests

Responds with `423 Locked` until initialized.
Binding `host` and `port` are configurable.

##### 2. Checks TCP/HTTP endpoints if configured

Waits till all dependent endpoints are accessible. Then allows acquiring of the lock. 
Endpoints are checked constantly. Locking is allowed only if all are OK.

##### 3. First request(s) acquires the lock

Client gets `200 OK`. Lock is acquired for a specific time. Multiple locks can be configured for parallel acquiring.
Custom lock duration may be specified in the request itself.

##### 4. Subsequent requests are denied to acquire the lock

Client gets `423 Locked` until lock timeout exceeds.

##### That's it. Steps 2 - 4 are constantly repeated

## Request custom lock duration

You can configure default lock timeout. But each client may request custom duration with `GET` parameter, for example, 
`http://lock.psl.svc.cluster.local:8888?duration=60`
To acquire a lock for 60 seconds.

## Dependent Endpoints check

This is useful when you need to wait for certain service(s) start before allowing starting of applications in the cluster.

* You may specify `http/https` endpoint, like `http://elastic.search.svc.cluster.local:9200`.
  Response with HTTP code `2XX` is considered as a success.
* You may specify `tcp` endpoint, like `tcp://mongodb.database.svc.cluster.local:27017`.
  Established TCP connection is considered as a success.

## Configuration

You may specify environment variables to override defaults:

| Option               | Default | Required | Description                                       |
|----------------------|---------|----------|---------------------------------------------------|
| `PSL_BIND_HOST`      | 0.0.0.0 |          | Address to bind                                   |
| `PSL_BIND_PORT`      | 8080    |          | Port to bind                                      |
| `PSL_PARALLEL_LOCKS` | 1       |          | Number of locks allowed to acquire simultaneously |
| `PSL_LOCK_DURATION`  | 10s     |          | Default lock duration                             |
| `PSL_HC_ENABLED`     | false   |          | Enabled health checks                             |
| `PSL_HC_ENDPOINTS`   | *none*  |          | List of endpoints to check before allow locking   |
| `PSL_HC_PERIOD_FAIL` | 10s     |          | Period of health checks if previous failed        |
| `PSL_HC_PERIOD_PASS` | 60s     |          | Period of health checks if previous succeeded     |
| `PSL_HC_TIMEOUT`     | 5s      |          | Timeout of health check requests                  |
| `PSL_LOG`            | info    |          | Log level                                         |

## How to run locally

Example with some options:

```bash
PSL_HC_ENABLED=true PSL_HC_ENDPOINTS=http://k8s-health:8080,http://myelasticsearch:9200,tcp://mymongodb:27017 bin/lock
```

## How to deploy to Kubernetes

The preferable way is to deploy as a DaemonSet.
You can find example deployment in [.k8s/lock](../.k8s/lock) directory.
