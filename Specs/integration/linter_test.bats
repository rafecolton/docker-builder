#!/usr/bin/env bats

@test "docker-builder correct lints a valid Builderfile" {
  run $GOPATH/bin/docker-builder -q lint Specs/fixtures/bob.toml
  [ "$status" -eq 0  ]
}

@test "docker-builder exits 5 when asked to lint an invalid file" {
  run $GOPATH/bin/docker-builder -q lint README.md
  [ "$status" -eq 5  ]
}

@test "docker-builder exits 17 when asked to lint a file that does not exist" {
  run $GOPATH/bin/docker-builder -q lint foo
  [ "$status" -eq 17  ]
}

#vim:ft=bats
