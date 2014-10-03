#!/usr/bin/env bats

@test "docker-builder server can be started successfully without basic auth" {
  run $GOPATH/bin/docker-builder serve --integration-test-mode
  [ "$status" -eq 166  ]
}

@test "docker-builder server can be started successfully with basic auth" {
  run $GOPATH/bin/docker-builder serve --username foo --password bar --integration-test-mode
  [ "$status" -eq 166  ]
}

#vim:ft=bats
