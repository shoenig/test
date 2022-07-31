#!/usr/bin/env bash

set -euo pipefail

cp scripts.go must/scripts.go
cp scripts_test.go must/scripts_test.go
cp test.go must/must.go
cp test_test.go must/must_test.go
sed -i "s|package test|package must|g" must/scripts.go
sed -i "s|package test|package must|g" must/scripts_test.go
sed -i "s|package test|package must|g" must/must.go
sed -i "s|package test|package must|g" must/must_test.go
sed -i -e "1s|^|// Code generated via scripts/generate.sh. DO NOT EDIT.\n\n|g" must/scripts.go
sed -i -e "1s|^|// Code generated via scripts/generate.sh. DO NOT EDIT.\n\n|g" must/scripts_test.go
sed -i -e "1s|^|// Code generated via scripts/generate.sh. DO NOT EDIT.\n\n|g" must/must.go
sed -i -e "1s|^|// Code generated via scripts/generate.sh. DO NOT EDIT.\n\n|g" must/must_test.go
