#!/bin/bash

# Check that the IBM OpenAPI Linter is installed
if ! command -v lint-openapi >/dev/null; then
  gum style --foreground 196 "Error: IBM OpenAPI Linter is not installed. Please install the linter by following the instructions at https://github.com/IBM/openapi-validator"
  exit 1
fi

score=0

lint-openapi -c ./openapi-lint-config.yaml -s ./common/common.yaml >./lint-output.json
score=$(jq '.impactScore.categorizedSummary.overall' ./lint-output.json)

if [ "$score" -lt 100 ]; then
  echo "IBM OpenAPI Linter found issues with the OpenAPI specification. Please fix the issues and try again."
  exit 1
fi
