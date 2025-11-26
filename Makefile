APP_NAME=go-platform-starter-service
CMD_PATH=./cmd/server

.PHONY: run build test lint

run:
	go run $(CMD_PATH)

build:
	go build -o bin/$(APP_NAME) $(CMD_PATH)

test:
	go test ./...

lint:
	@echo "No linter configured yet. Add golangci-lint or staticcheck later."
