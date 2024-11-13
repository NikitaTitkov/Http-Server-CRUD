LOCAL_BIN:=$(CURDIR)/bin

install-golangci-lint:
	GOBIN=$(LOCAL_BIN) go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.61.0

lint:
	$(LOCAL_BIN)/golangci-lint run ./... --config .golangci.pipeline.yaml
	
local-migration-up:
	migrate -path ./schema -database 'postgres://postgres:qwerty@localhost:5436/postgres?sslmode=disable' up

local-migration-down:
	migrate -path ./schema -database 'postgres://postgres:qwerty@localhost:5436/postgres?sslmode=disable' up