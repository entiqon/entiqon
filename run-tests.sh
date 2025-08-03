#!/bin/bash
set -e

modules=(common db)

WITH_COVERAGE=false

# Check if the first argument is '-with-coverage'
if [[ "$1" == "-with-coverage" ]]; then
  WITH_COVERAGE=true
  shift  # Remove first argument so it doesn’t interfere with modules
fi

for mod in "${modules[@]}"
do
  echo "Running tests in $mod"
  if $WITH_COVERAGE; then
    if [ ! -f coverage.out ]; then
      go test -coverprofile=coverage.out ./$mod/...
    else
      go test -coverprofile=profile.tmp ./$mod/...
      tail -n +2 profile.tmp >> coverage.out
      rm profile.tmp
    fi
  else
    go test ./$mod/...
  fi
done

if $WITH_COVERAGE; then
  echo "Coverage report generated at coverage.out"
fi
