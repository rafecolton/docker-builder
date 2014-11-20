#!/usr/bin/env bats

@test "docker-builder correct lints a valid Builderfile" {
  run $GOPATH/bin/docker-builder -q lint Specs/fixtures/bob.toml
  [ "$status" -eq 0  ]
}

@test "docker-builder exits nonzero when asked to lint an invalid file" {
  run $GOPATH/bin/docker-builder -q lint README.md
  [ "$status" -ne 0  ]
}

@test "docker-builder exits nonzero when asked to lint a file that does not exist" {
  run $GOPATH/bin/docker-builder -q lint foo
  [ "$status" -ne 17  ]
}

#vim:ft=bats
