## Travis and GitHub Webhooks

While you can [enqueue a job yourself via HTTP](enqueueing-a-build.md),
`docker-builder` also supports jobs being triggered by both GitHub pushes and
successful Travis builds via webhook notifications.

### Travis

Travis can issue a webhook after any build, and can be limited to
successful or failed builds only. To provide authentication for Travis
requests, `docker-builder`, you will need to start the server with the
`--travis-token` flag or `DOCKER_BUILDER_TRAVISTOKEN` environment
variable.

```bash
$ docker-builder serve --travis-token <token>
```

Next, add the following to your `.travis.yml` file:

```yaml
notifications:
  webhooks:
    urls:
      - http://BUILD_SERVER_LOCATION/docker-build/travis
    on_failure: never
    on_success: always
```

Note that while Travis supports multiple build types, only `push`
builds are currently supported.  Builds from pull requests will not be
enqueued to the build server. Additionally, `docker-builder` will not
enqueue a job if the Travis build failed, so the
`on-failure`/`on-success` settings above are reccomended.

### GitHub

The only Github event type currently supported is the "push" event.  In
order to authenticate requests from GitHub, you will need to start your
server with the GitHub secret you used when creating the webhook.  This
can be specified to the server with either the `--github-secret` flag or
the `DOCKER_BUILDER_GITHUBSECRET` environment variable.

```bash
$ docker-builder serve --github-secret <webhook-secret>
```

You can add a Github webhook to your repository by accessing the
settings page at https://github.com/USERNAME/REPOSITORY/settings/hooks.
Make sure the webhook is set to trigger on "Just the `push` event".

Note that the route for GitHub hooks is `/docker-build/github`

### Disabling Authentication and Endpoints

If a GitHub secret is not supplied, requests to the GitHub endpoint will
not be authenticated.  The same applies to Travis.  While this is useful
for debugging, it is not reccomended to leave the webhook endpoints
unsecured in production.

If you wish to disable these endpoints entirely, you can do so with the
`--no-travis` and `--no-github` flags. Theere are also corresponding
environment variables for these settings (`DOCKER_BUILDER_NOTRAVIS` and
`DOCKER_BUILDER_NOGITHUB`).
