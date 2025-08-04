#!/bin/bash
set -e

WITH_COVERAGE=false
OPEN_COVERAGE=false

print_usage() {
  cat << EOF
Usage: $0 [options]

Options:
  --with-coverage     Run tests with coverage report generation (coverage.out)
  --open-coverage     Open coverage report in browser (only valid with --with-coverage on macOS)
  -h, --help          Show this help message and exit

Examples:
  $0                         Run tests normally
  $0 --with-coverage          Run tests with coverage, no auto-open
  $0 --with-coverage --open-coverage   Run tests with coverage and open report (macOS only)
EOF
}

# Show help and exit if requested
for arg in "$@"; do
  case $arg in
    -h|--help)
      print_usage
      exit 0
      ;;
  esac
done

# Parse flags
for arg in "$@"; do
  case $arg in
    --with-coverage)
      WITH_COVERAGE=true
      ;;
    --open-coverage)
      OPEN_COVERAGE=true
      ;;
  esac
done

# Only allow --open-coverage on macOS
if $OPEN_COVERAGE; then
  if [[ "$(uname)" != "Darwin" ]]; then
    echo "Error: --open-coverage flag is only supported on macOS (Darwin)."
    exit 1
  fi
fi

if $WITH_COVERAGE; then
  echo "Running tests with coverage across all packages"
  go test -coverprofile=coverage.out -covermode=atomic ./...
else
  echo "Running tests normally across all packages"
  go test ./...
fi

if $WITH_COVERAGE; then
  echo "Coverage report saved as coverage.out"
  echo "Run 'go tool cover -html=coverage.out' to view it"

  if $OPEN_COVERAGE; then
    echo "Opening coverage report in browser..."
    go tool cover -html=coverage.out
  fi
fi
