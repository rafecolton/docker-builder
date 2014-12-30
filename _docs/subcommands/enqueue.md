# enqueue

Use `docker-builder enqueue` to push a build for your *current working
directory* to your Docker build server.  To use on localhost, run the server
in one tab and enqueue in another.  For example:

```bash
docker-builder serve &
docker-builder enqueue
```

Or, you may push directly to your build server by setting the
docker-build-server host:

```bash
# via the environment
export DOCKER_BUILDER_HOST="http://localhost:5000"
docker-builder enqueue
```

or

```bash
# via the command line
docker-builder enqueue --host "http://localhost:5000"
```
