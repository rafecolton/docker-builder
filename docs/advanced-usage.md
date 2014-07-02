## Advanced Usage

Topics:

0. [Running the Server](#running-the-server)
0. [Enqueueing a Build](enqueueing-a-build.md)
0. [Travis and GitHub Webhooks](travis-and-github-webhooks.md)

### Running the Server

#### Usage

```bash
> docker-builder help serve

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
# DOCKER_BUILDER_LOGLEVEL             =>     --log-level (global)
# DOCKER_BUILDER_LOGFORMAT            =>     --log-format (global)
# DOCKER_BUILDER_PORT                 =>     --port
# DOCKER_BUILDER_APITOKEN             =>     --api-token
# DOCKER_BUILDER_SKIPPUSH             =>     --skip-push
# DOCKER_BUILDER_USERNAME             =>     --username
# DOCKER_BUILDER_PASSWORD             =>     --password
# DOCKER_BUILDER_TRAVISTOKEN          =>     --travis-token
# DOCKER_BUILDER_NOTRAVIS             =>     --no-travis
# DOCKER_BUILDER_GITHUBSECRET         =>     --github-secret
# DOCKER_BUILDER_NOGITHUB             =>     --no-github
# 
# NOTE: If username and password are both empty (i.e. not provided), basic
# auth will not be used.
# 
# 
# OPTIONS:
#    --port, -p '5000'    port on which to serve
#    --api-token, -t      GitHub API token
#    --skip-push          override Bobfile behavior and do not push any images (useful for testing)
#    --username           username for basic auth
#    --password           password for basic auth
#    --travis-token       Travis API token for webhooks
#    --github-secret      GitHub secret for webhooks
#    --no-travis          do not include route for Travis CI webhook
#    --no-github          do not include route for GitHub webhook
```

#### Healthcheck

The `docker-builder` server has a healthcheck route available at
`/health`.  As long as the server is running, an HTTP request to
`/health` will return 200/OK.
