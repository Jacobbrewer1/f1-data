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
	go generate ./...
deps:
	sudo apt-get install dos2unix
	dos2unix ./pkg/models/generate.sh
	chmod +x ./pkg/models/generate.sh
	go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest
	go install github.com/charmbracelet/gum@latest
	go install golang.org/x/tools/cmd/goimports@latest
	go install github.com/jacobbrewer1/goschema@latest
models:
	chmod +x ./pkg/models/generate.sh
	go generate ./pkg/models
apis:
	go generate ./pkg/codegen/...
