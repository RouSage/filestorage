# ==================================== #
# HELPERS
# ==================================== #

# help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

# ==================================== #
# DEVELOPMENT
# ==================================== #

.PHONY: run
run:
	go run main.go

# ==================================== #
# QUALITY CHECK
# ==================================== #

# tidy: format all .go files and tidy module dependencies
.PHONY: tidy
tidy:
	@echo "Formatting .go files..."
	go fmt ./...
	@echo "Tidying module dependencies..."
	go mod tidy
	@echo "Verifying module dependencies..."
	go mod verify

# audit: run quality control checks
.PHONY: audit
audit:
	@echo "Checking module dependencies..."
	go mod tidy -diff
	go mod verify
	@echo "Vetting code..."
	go vet ./...
	go tool staticcheck ./...
	go test -race -vet=off ./...

# test: run tests
.PHONY: test
test:
	go test -v ./...

# ==================================== #
# BUILD
# ==================================== #

# build: build the application
.PHONY: build
build:
	@echo "Building..."
	go build -ldflags="-s" -o=./bin/fs
