#/bin/bash

# This script is used to run the test suite for the project.
# It sets up the environment and runs the tests.
# Usage: ./scripts/test.sh
# Make sure to run this script from the root of the project
# Check if the script is being run from the root of the project
if [ ! -f go.mod ]; then
  echo "Please run this script from the root of the project."
  exit 1
fi
go test ./...
