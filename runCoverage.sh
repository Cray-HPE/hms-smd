#!/usr/bin/env bash

# Fail on error and print executions
set -ex

# Build the build base image (if it's not already)
docker build -t cray/smd-base -f Dockerfile.smd --target base .

# Run the tests.
docker build -t cray/smd-unit-testing -f Dockerfile.coverage --no-cache .
