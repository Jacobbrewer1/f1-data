# Define variables
hash = $(shell git rev-parse --short HEAD)
DATE = $(shell date -u +'%Y-%m-%dT%H:%M:%SZ')

pr-approval:
	@echo "Running PR CI"
	go build ./...
	go vet ./...
	go test ./...
codegen:
	@echo "Generating code"
	go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest

	cd ./pkg/codegen/apis

	# Loop through each directory. Inside the directory we want to use the directory name as the package name.
	# Ignore the README.md file
	for d in ./pkg/codegen/apis/*; do \
		if [ -d $$d ]; then \
			(cd $$d && go generate); \
		fi \
	done
mocks:
	@echo "Generating mocks"
	# Loop through each directory inside the pkg/repositories directory. Inside the directory we want to use the directory name as the package name.
	# Ignore the README.md file
	for d in ./pkg/repositories/*; do \
		if [ -d $$d ]; then \
			(cd $$d && go generate); \
		fi \
	done
	cd ./pkg/vault && go generate
