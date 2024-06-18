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
You may specify additional command line options to override defaults:

| Option      | Default | Description                                                            |
|-------------|---------|------------------------------------------------------------------------|
| `--host`    | 0.0.0.0 | Address to bind                                                        |
| `--port`    | 8888    | Port to bind                                                           |
| `--locks`   | 1       | Number of locks allowed to be acquired at the same time                |
| `--timeout` | 10      | Default time until the acquired lock is released, seconds              |
| `--check`   | *none*  | List of endpoints to check before allow locking, see example below     |
| `--failHc`  | 10      | Pause between endpoint checks if the previous check failed, seconds    |
| `--passHc`  | 60      | Pause between endpoint checks if the previous check succeeded, seconds |

## How to run locally
Example with some command line options:
```bash
go run lock/main.go --port 9000 --locks 2 --check http://myelasticsearch:9200 --check tcp://mymongodb:27017
```

## How to deploy to Kubernetes
The preferable way is to deploy as a DaemonSet.
You can find example deployment in [.k8s/lock](../.k8s/lock) directory.
