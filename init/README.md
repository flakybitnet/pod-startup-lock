# Init container for the Lock service

Application repeatedly queries `lock` service endpoint until it gets `200 OK` response.
Then exits letting main application to start.

**Designed to be deployed as an Init Container**.

## Configuration

You may specify environment variables to override defaults:

| Option                   | Default | Required | Description                       |
|--------------------------|---------|----------|-----------------------------------|
| `PSL_LOCK_HOST`          | *none*  | +        | Lock Service's hostname           |
| `PSL_LOCK_PORT`          | 8080    |          | Lock Service's HTTP port          |
| `PSL_LOCK_DURATION`      | *none*  |          | Custom lock duration to request   |
| `PSL_LOCK_CHECK_PERIOD`  | 3s      |          | Period of Lock acquiring attempts |
| `PSL_LOCK_CHECK_TIMEOUT` | 1s      |          | Timeout of Lock acquiring request |
| `PSL_LOG`                | info    |          | Log level                         |

## How to run locally

Example with some options:

```bash
PSL_LOCK_HOST=localhost PSL_LOCK_DURATION=30s bin/init
```

## How to deploy to Kubernetes

Should be deployed as an Init Container.
You can find example deployment of Postgres database in [.k8s/init](../.k8s/init) directory.
