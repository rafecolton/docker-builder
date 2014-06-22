## Advanced Usage

```bash
docker-builder help serve

# NAME:
#    serve - serve <options> - start a small HTTP web server for receiving build requests
#
# USAGE:
#    command serve [command options] [arguments...]
#
# DESCRIPTION:
#    Start a small HTTP web server for receiving build requests.
#
# Configure through the environment:
#
# DOCKER_BUILDER_LOGLEVEL     =>     --log-level (global)
# DOCKER_BUILDER_LOGFORMAT    =>     --log-format (global)
# DOCKER_BUILDER_PORT         =>     --port
# DOCKER_BUILDER_APITOKEN     =>     --api-token
# DOCKER_BUILDER_SKIPPUSH     =>     --skip-push
#
#
# OPTIONS:
#    --port, -P '5000'  port on which to serve
#    --api-token, -T  GitHub API token
#    --skip-push    override Bobfile behavior and do not push any images (useful for testing)
```
