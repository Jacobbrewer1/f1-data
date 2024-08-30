#!/bin/bash

buildApp=$1
toPush=$2

function fail() {
  echo "Error: $1"
  exit 1
}

function handle_subdirectory() {
  hash=$(git rev-parse --short HEAD)
  DATE=$(date -u +'%Y-%m-%dT%H:%M:%SZ')

  echo "$hash"
  echo "$DATE"

  # Build binary for the subdirectory
  echo "Building binary for $1"
  cd ./"$1" || fail "Unable to cd to cmd/$1"

  GOOS=linux GOARCH=amd64 go build -o bin/app -ldflags "-X main.Commit=$hash -X main.Date=$DATE" . || fail "Unable to build binary for $1"

  # Move app to bin directory from the root of the repo. Create the bin directory if it doesn't exist
  mkdir -p ../../bin || fail "Unable to create bin directory"
  mv bin/app ../../bin/app || fail "Unable to move binary to bin/app"

  rm -rf bin

  cd ../../ || fail "Unable to cd to root of repo"

  # Remove the "./cmd/" prefix from the directory name
  appName=${1#"./cmd/"}
  appName=${appName#"cmd/"}

  # Build the docker image
  echo "Building docker image for $appName"
  registry="ghcr.io/jacobbrewer1/f1-data-$appName"

  docker build -t "$registry:$hash" -t "$registry:latest" . || fail "Unable to build docker image for $appName"

  # Check if the toPush variable is set to true. If yes, push the docker image to the GitHub Container Registry
  if [ "$toPush" != "true" ]; then
    echo "Not pushing docker image to GitHub Container Registry"

    # Cleanup the binary
    rm -rf bin
    rm -rf ./cmd/"$appName"/bin

    return 0
  fi

  # Push the docker image to the GitHub Container Registry
  echo "Pushing docker image to GitHub Container Registry"

  docker push "$registry:$hash" || fail "Unable to push docker image to GitHub Container Registry"
  echo "Docker image pushed to $registry:$hash"

  docker push "$registry:latest" || fail "Unable to push docker image to GitHub Container Registry"
  echo "Docker image pushed to $registry:latest"

  echo "Done building $appName"

  # Cleanup the binary
  rm -rf bin
  rm -rf ./cmd/"$appName"/bin
}

# If the toPush variable is not set, set it to false
if [ -z "$toPush" ]; then
  toPush="false"
fi

if [ "$buildApp" == "build-all" ]; then
  echo "Building all apps"
elif [ -z "$buildApp" ]; then
  echo "No app specified. Building all apps"
  buildApp="build-all"
else
  echo "Building $buildApp"
  handle_subdirectory "cmd/$buildApp"
  exit 0
fi

# Get all the subdirectories in the cmd directory
subdirectories=$(ls -d ./cmd/*)

# For each subdirectory of the cmd directory and run the subdirectory function
for dir in $subdirectories; do
  if [ -d "$dir" ]; then
    echo "Running $dir"
    handle_subdirectory "$dir"
  fi
done