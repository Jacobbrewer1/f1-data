#!/bin/bash

# Check that the IBM OpenAPI Linter is installed
if ! command -v lint-openapi >/dev/null; then
  echo "Error: IBM OpenAPI Linter is not installed. Please install the linter by following the instructions at https://github.com/IBM/openapi-validator"
  exit 1
fi

# Find all routes.yaml files in teh ./pkg/codegen/apis directory
routesFiles=$(find ./pkg/codegen/apis -name "routes.yaml")

touch ./pr-report.md
echo "## Linting Report" >./pr-report.md

totalErrors=0
totalWarnings=0
totalInfos=0
totalHints=0

# Lint each routes.yaml file
for file in $routesFiles; do
  rm -rf ./lint-output.json

  lint-openapi -c ./openapi-lint-config.yaml -s "$file" >./lint-output.json

  # Make ./pkg/codegen/apis/api/routes.yaml -> api/routes.yaml
  prettyName=$(echo $file | sed 's/\.\/pkg\/codegen\/apis\///' | sed 's/\/routes.yaml//')

  echo "Report for $prettyName"

  # Print the lint output
  cat ./lint-output.json

  # Put the header on the PR report
  cat <<EOF >>./pr-report.md
### Linting $prettyName
\`\`\`
EOF

  # Get the total number of errors, warnings, infos, and hints
  errors=$(cat ./lint-output.json | jq .error.summary.total)
  warnings=$(cat ./lint-output.json | jq .warning.summary.total)
  infos=$(cat ./lint-output.json | jq .info.summary.total)
  hints=$(cat ./lint-output.json | jq .hint.summary.total)

  # Put the lint output in the PR report
  cat <<EOF >>./pr-report.md
$errors errors, $warnings warnings, $infos infos, $hints hints
EOF

  # Put the footer on the PR report
  cat <<EOF >>./pr-report.md
\`\`\`

EOF

  if [[ $errors -gt 0 ]]; then
    cat <<EOF >>./pr-report.md
#### Error Messages
\`\`\`
EOF

    cat ./lint-output.json | jq -r '.error.summary.entries[].generalizedMessage' >> ./pr-report.md

    cat <<EOF >>./pr-report.md
\`\`\`

EOF
  fi

    if [[ $warnings -gt 0 ]]; then
      cat <<EOF >>./pr-report.md
#### Warning Messages
\`\`\`
EOF

    cat ./lint-output.json | jq -r '.warning.summary.entries[].generalizedMessage' >> ./pr-report.md

    cat <<EOF >>./pr-report.md
\`\`\`

EOF
    fi

    if [[ $infos -gt 0 ]]; then
      cat <<EOF >>./pr-report.md
#### Info Messages
\`\`\`
EOF

    cat ./lint-output.json | jq -r '.info.summary.entries[].generalizedMessage' >> ./pr-report.md

    cat <<EOF >>./pr-report.md
\`\`\`

EOF
    fi

    if [[ $hints -gt 0 ]]; then
      cat <<EOF >>./pr-report.md
#### Hint Messages
\`\`\`
EOF

    cat ./lint-output.json | jq -r '.hint.summary.entries[].generalizedMessage' >> ./pr-report.md

    cat <<EOF >>./pr-report.md
\`\`\`

EOF
    fi

  # Add the errors, warnings, infos, and hints to the total
  totalErrors=$((totalErrors + errors))
  totalWarnings=$((totalWarnings + warnings))
  totalInfos=$((totalInfos + infos))
  totalHints=$((totalHints + hints))
done

if [[ $totalErrors -gt 0 ]]; then
  echo "FAIL: Linting failed with $totalErrors errors, $totalWarnings warnings, $totalInfos infos, and $totalHints hints"
  exit 1
elif [[ $totalWarnings -gt 0 ]]; then
  echo "FAIL: Linting failed with $totalErrors errors, $totalWarnings warnings, $totalInfos infos, and $totalHints hints"
  exit 1
else
  echo "PASS: Linting passed with $totalErrors errors, $totalWarnings warnings, $totalInfos infos, and $totalHints hints"
fi
