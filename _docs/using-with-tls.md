# Using with TLS

If you are using a version of `docker` with TLS enabled (supported in
`docker` `v1.3.0` and up, enabled by default with `boot2docker`), you
will need to use `docker-builder` `v0.9.2` or greater.

Additionally, you must set the following environment variables:

```bash
# all values are the boot2docker defaults
export DOCKER_CERT_PATH="$HOME/.boot2docker/certs/boot2docker-vm"
export DOCKER_TLS_VERIFY=1
export DOCKER_HOST="tcp://127.0.0.1:2376"
```

**NOTE:** `docker-builder` will automatically set the correct url scheme
for TLS if you are using port 2376.  If you are using another port and
wish to enable TLS, you must set the following additional environment
variable:

```bash
export DOCKER_HOST_SCHEME="https"
```
