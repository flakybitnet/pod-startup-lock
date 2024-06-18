# Init container for the Lock service

Application repeatedly queries `lock` service endpoint until it gets `200 OK` response.
Then exits letting main application to start.

**Designed to be deployed as an Init Container**.

## Additional Configuration
You may specify additional command line options to override defaults:

| Option       | Default   | Description                                    |
|--------------|-----------|------------------------------------------------|
| `--port`     | 8888      | Lock Service's HTTP port                       |
| `--host`     | localhost | Lock Service's hostname                        |
| `--pause`    | 1         | Pause between Lock acquiring attempts, seconds |
| `--duration` | *none*    | Custom lock duration to request, seconds       |

## How to run locally
Example with some command line options:
```bash
go run init/main.go --port 9000 --duration 15
```

## How to deploy to Kubernetes
Should be deployed as an Init Container.
You can find example deployment of Postgres database in [.k8s/init](../.k8s/init) directory.
