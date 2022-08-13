#!/usr/bin/env bash

set -euo pipefail

files=$(git ls-files -m)
if [ -n "$files" ]; then
  echo "Files have been modified ..."
  echo "$files"
  exit 1
fi
