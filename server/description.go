package server

//Description is the help text for the `serer` command
const Description = `Start a small HTTP web server for receiving build requests.

Configure through the environment:

Global:
DOCKER_BUILDER_LOGLEVEL         =>     --log-level (global)
DOCKER_BUILDER_LOGFORMAT        =>     --log-format (global)

Server:
DOCKER_BUILDER_PORT             =>     --port
DOCKER_BUILDER_APITOKEN         =>     --api-token
DOCKER_BUILDER_SKIPPUSH         =>     --skip-push

Basic Auth:
DOCKER_BUILDER_USERNAME         =>     --username
DOCKER_BUILDER_PASSWORD         =>     --password

Travis Auth:
DOCKER_BUILDER_TRAVISTOKEN      =>     --travis-token
DOCKER_BUILDER_NOTRAVIS         =>     --no-travis

GitHub Auth:
DOCKER_BUILDER_GITHUBSECRET     =>     --github-secret
DOCKER_BUILDER_NOGITHUB         =>     --no-github

Docker Registry Auth:
DOCKER_BUILDER_CFGUN            =>     --dockercfg-un
DOCKER_BUILDER_CFGPASS          =>     --dockercfg-pass
DOCKER_BUILDER_CFGEMAIL         =>     --dockercfg-email

NOTE: If username and password are both empty (i.e. not provided), basic auth will not be used.
`
