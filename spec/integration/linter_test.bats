#!/usr/bin/env bats

@test "builder correct lints a valid Builderfile" {
  run $GOBIN/builder --quiet --lint spec/integration/Builderfile
  [ "$status" -ne 0  ]
}

@test "builder correct lints an invalid Builderfile" {
  run $GOBIN/builder --quiet --lint README.md
  [ "$status" -ne 0  ]
}

#vim:ft=bats
