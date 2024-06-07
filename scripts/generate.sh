#!/usr/bin/env bash

set -euo pipefail

apply() {
  original="${1}"
  clone="must/${original}"
  cp "${original}" "${clone}"
  sed -i.bak "s|package test|package must|g" "${clone}"
  sed -i.bak -e "1s|^|// Code generated via scripts/generate.sh. DO NOT EDIT.\n\n|g" "${clone}"
}

apply invocations.go
apply invocations_test.go
apply settings.go
apply settings_test.go
apply scripts.go
apply scripts_test.go
apply test.go
apply test_test.go
apply examples_test.go

# rename core test files
mv must/test.go must/must.go
mv must/test_test.go must/must_test.go

# cleanup *.bak files
find . -name *.bak | xargs rm
