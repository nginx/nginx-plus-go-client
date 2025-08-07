# renovate: datasource=github-tags depName=golangci/golangci-lint
GOLANGCI_LINT_VERSION = v2.3.1

test: unit-test test-integration test-integration-no-stream-block clean

lint:
	go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION) run --fix

unit-test:
	go test -v -shuffle=on -race client/*.go

test-integration:
	docker compose up -d --build test
	docker compose logs -f test

test-integration-no-stream-block:
	docker compose up -d --build test-no-stream
	docker compose logs -f test-no-stream

clean:
	docker compose down --remove-orphans
