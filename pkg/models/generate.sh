#!/bin/bash

if ! command -v goschema &> /dev/null; then
  gum style --foreground 196 "goschema is required to generate models. Please install it"
  exit 1
fi

if ! command -v gum &> /dev/null; then
  gum style --foreground 196 "gum is required to generate models. Please install it by running 'make deps'"
  exit 1
fi

if ! command -v goimports &> /dev/null; then
  gum style --foreground 196 "goimports is required to generate models. Please install it by running 'go get golang.org/x/tools/cmd/goimports'"
  exit 1
fi

all=false
clean=false
forced=false
silent=false

# Get the flags passed to the script and set the variables accordingly
while getopts "acf" flag; do
  case $flag in
    a)
      all=true
      ;;
    c)
      clean=true
      ;;
    f)
      forced=true
      ;;
    *)
      gum style --foreground 196 "Invalid flag $flag"
      exit 1
      ;;
  esac
done

# If the silent environment variable is set, set the silent flag to true
if [ "$GS_SILENT" = true ]; then
  silent=true
fi

# If the -c flag is passed, remove all generated models
if [ "$clean" = true ]; then
  if [ "$forced" = false ]; then
    gum style "Are you sure you want to remove all generated models?"
    gum confirm || exit 0
  fi

  rm -f ./*.xo.go
fi

# If the -a flag is passed, generate all models
if [ "$all" = true ]; then
  if [ "$silent" = true ]; then
    echo "Generating all models"
    goschema generate --out=./ --sql=./schemas/*.sql --extension=xo
    go fmt ./*.xo.go
    goimports -w ./*.xo.go
    exit 0
  fi

  gum spin --spinner dot --title "Generating all models" -- goschema generate --out=./ --sql=./schemas/*.sql --extension=xo
  go fmt ./*.xo.go
  goimports -w ./*.xo.go
  exit 0
fi

# Allow the user too select the model/s to generate
schemas=$(find ./schemas -type f -name "*.sql" | sed 's/\.\/schemas\///' | sed 's/\.sql//')

togen=$(echo "$schemas" | gum choose --no-limit --header "Select the model(s) to generate")

gum style "Are you sure you want to generate the following models?"

for model in $togen; do
  gum style --foreground 222 "  - $model"
done

if [ "$forced" = false ]; then
  gum confirm || exit 0
fi

for model in $togen; do
  gum spin --spinner dot --title "Generating model $model" -- goschema generate --out=./ --sql=./schemas/"$model".sql --extension=xo
  go fmt ./"$model".xo.go
  goimports -w ./"$model".xo.go
done
