#!/bin/bash

# This script is used to run the test suite for the project.
# It sets up the environment and runs the tests.
# Usage: ./scripts/test.sh
# Make sure to run this script from the root of the project
# Check if the script is being run from the root of the project
if [ ! -f go.mod ]; then
  echo "Please run this script from the root of the project."
  exit 1
fi

declare -a test_list=()
while read pkg; do 
  pkg_dir=$(go list -f '{{.Dir}}' $pkg)
  if [ ! -f "$pkg_dir/skip_test.go" ]; then 
    test_list+=("$pkg")
  else
    echo "Skipping tests in $pkg"
  fi
done < <(go list ./...)

# echo "Running tests in the following packages:"
# for pkg in "${test_list[@]}"; do
#   echo " - $pkg"
# done

go test -cover -covermode=atomic "${test_list[@]}" || exit 1
