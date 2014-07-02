## Travis and GitHub Webhooks

While you can
[enqueue a job yourself via HTTP](enqueueing-a-build.md),
docker-builder also supports jobs being triggered by Github pushes or
successful Travis builds via their webhook notifications.

Travis
---
Travis can issue a webhook after a build, and you can limit it to
successful or failed builds. For docker-builder, you will need to
start the server with the `--travis-token` flag or
`DOCKER_BUILDER_TRAVISTOKEN` environment variable.

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
builds are currently supported. Additionally, docker-builder will not
enqueue a job if the Travis build failed, so the
`on-failure`/`on-success` settings above are reccomended.

The next successful build that Travis completes from a push should
trigger a new build.

Github
---
The only Github event type currently supported is the "push"
event. First you will need to start your server with a Github secret,
which you can specify with either the `--github-secret` flag or the
`DOCKER_BUILDER_GITHUBSECRET` environment variable.

```bash
$ docker-builder serve --github-secret <api-secret>
```

You can then add a Github webhook to your repository by accessing the
settings page at
https://github.com/USERNAME/REPOSITORY/settings/hooks. Make sure the
webhook is set to trigger on "Just the `push` event".

The next time you push to your repository you should see the
docker-builder server start a build.

Disabling Authentication and Endpoints
---
If a Github secret or Travis token is not supplied, then requests to
that endpoint will be not verified. While this is useful for
debugging, it is not reccomended to leave the webhook endpoints
unsecured in production.

If you wish to disable these endpoints entirely, you can do so with
the `--no-travis` and `--no-github` endpoints. Theere are also
correspdonig environment variables for these settings
(`DOCKER_BUILDER_NOTRAVIS` and `DOCKER_BUILDER_NOGITHUB`).
