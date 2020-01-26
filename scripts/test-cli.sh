#!/bin/bash

# =================================================================
#
# Copyright (C) 2020 Spatial Current, Inc. - All Rights Reserved
# Released as open source under the MIT License.  See LICENSE file.
#
# =================================================================

set -euo pipefail

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

testZero() {
  local expected='hello world'
  local output=$(echo 'hello world' | goslicer -i 0)
  assertEquals "unexpected output" "$(echo -e "${expected}")" "${output}"
}

testOne() {
  local expected='world'
  local output=$(echo 'hello world' | goslicer -i 6)
  assertEquals "unexpected output" "$(echo -e "${expected}")" "${output}"
}

testTwo() {
  local expected='hello'
  local output=$(echo 'hello world' | goslicer -i 0:5)
  assertEquals "unexpected output" "$(echo -e "${expected}")" "${output}"
}

testTwoNegative() {
  local expected='hello'
  local output=$(echo 'hello world' | goslicer --lines -i 0:-6)
  assertEquals "unexpected output" "$(echo -e "${expected}")" "${output}"
}

oneTimeSetUp() {
  echo "Setting up"
  echo "Using temporary directory at ${SHUNIT_TMPDIR}"
}

oneTimeTearDown() {
  echo "Tearing Down"
}

# Load shUnit2.
. "${DIR}/shunit2"
