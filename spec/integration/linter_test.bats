#!/usr/bin/env bats

@test "builder correct lints a valid Builderfile" {
  run $GOBIN/builder -q --lint spec/fixtures/Builderfile
  [ "$status" -eq 0  ]
}

@test "builder exits 5 when asked to lint an invalid file" {
  run $GOBIN/builder -q --lint README.md
  [ "$status" -eq 5  ]
}

@test "builder exits 17 when asked to lint an invalid file" {
  run $GOBIN/builder -q --lint foo
  [ "$status" -eq 17  ]
}

#vim:ft=bats
