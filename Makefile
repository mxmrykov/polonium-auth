#!make
include ./migrations/.env

LOCAL_BIN:=$(CURDIR)/bin
EASYJSON := ~/go/bin/easyjson

.PHONY: install-requirements
install-requirements:
	GOBIN=$(LOCAL_BIN) go install github.com/pressly/goose/v3/cmd/goose@v3.24.3

.PHONY: migrate-up
migrate-up: install-requirements
	$(LOCAL_BIN)/goose -dir ./migrations/sql postgres postgres://${POSTGRES_RW_USER}:${POSTGRES_RW_PASS}@localhost:5491/${POSTGRES_DB} up

easyjson: internal/model/models_easyjson.go

internal/model/models_easyjson.go: $(filter-out internal/model/models_easyjson.go, $(wildcard internal/model/*.go))
	@echo "⏳ Running easyjson for all models..."
	GOFLAGS=-mod=mod $(EASYJSON) -all internal/model
	@echo "✅ easyjson generation completed!"

easyjson_install:
	@echo "⏳ Installing easyjson dependencies..."
	go get github.com/mailru/easyjson && go install github.com/mailru/easyjson/...@latest
	@echo "✅ easyjson dependencies installed!"

up_builds:
	docker compose -f ./migrations/compose.yaml up -d