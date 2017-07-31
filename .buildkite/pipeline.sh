#!/usr/bin/env bash

set -eu

cat << STEP
steps:
  - command: 'docker-compose run --rm make clean test'
    label: ':pray: Testing'
    agents:
      queue: 'central-dev'
STEP

if [[ -n "${BUILDKITE_TAG:-}" ]]; then
  cat << STEP

  - wait

  - command: 'docker-compose run --rm create-release'
    label: ':airplane: Releasing'
    agents:
      queue: 'central-prod'

  - wait
STEP

  for GOOS in darwin linux windows; do
    for GOARCH in 386 amd64; do
      cat << STEP

  - command: |-
      docker-compose run --rm make clean install build
      zip -j build/konfigurator-$GOOS-$GOARCH.zip build/konfigurator
      docker-compose run --rm upload-release
    label: ':$GOOS: ($GOARCH) [${BUILDKITE_TAG}]'
    agents:
      queue: 'central-prod'
    env:
      GOOS: $GOOS
      GOARCH: $GOARCH
      BUILDKITE_TAG: $BUILDKITE_TAG

STEP
    done
  done
fi
