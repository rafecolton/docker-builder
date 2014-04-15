#!/usr/bin/env bats

@test "short version is set by compile" {
  run builder -v
  status=
  if [[ "$output" =~ "unknown" ]] ; then status="fail" ; else status="pass" ; fi
  [ "$status" = "pass" ]
}

@test "long version is set by compile" {
  run builder --version
  status=
  if [[ "$output" =~ "unknown" ]] ; then status="fail" ; else status="pass" ; fi
  [ "$status" = "pass" ]
}

@test "branch is set by compile" {
  run builder --branch
  status=
  if [[ "$output" =~ "unknown" ]] ; then status="fail" ; else status="pass" ; fi
  [ "$status" = "pass" ]
}

@test "rev is set by compile" {
  run builder --rev
  status=
  if [[ "$output" =~ "unknown" ]] ; then status="fail" ; else status="pass" ; fi
  [ "$status" = "pass" ]
}

#vim:ft=bats
