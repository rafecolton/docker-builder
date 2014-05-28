#!/usr/bin/env bats

@test "short version is set by compile" {
  run $GOPATH/bin/builder -v
  status=
  if [[ "$output" =~ "unknown" ]] ; then status="fail" ; else status="pass" ; fi
  [ "$status" = "pass" ]
}

@test "long version is set by compile" {
  run $GOPATH/bin/builder --version
  status=
  if [[ "$output" =~ "unknown" ]] ; then status="fail" ; else status="pass" ; fi
  [ "$status" = "pass" ]
}

@test "branch is set by compile" {
  run $GOPATH/bin/builder --branch
  status=
  if [[ "$output" =~ "unknown" ]] ; then status="fail" ; else status="pass" ; fi
  [ "$status" = "pass" ]
}

@test "rev is set by compile" {
  run $GOPATH/bin/builder --rev
  status=
  if [[ "$output" =~ "unknown" ]] ; then status="fail" ; else status="pass" ; fi
  [ "$status" = "pass" ]
}

#vim:ft=bats
