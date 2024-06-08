ifeq ($(POSTGRES_SETUP_TEST),)
	POSTGRES_SETUP_TEST := user=root password=secret dbname=TodoList host=localhost port=5432 sslmode=disable
endif
PKG := ./...
MAIN := /cmd/server/main.go
INTERNAL_PKG_PATH=$(CURDIR)/internal/
MIGRATION_FOLDER=$(INTERNAL_PKG_PATH)/migrations

.PHONY: migration-create
migration-create:
	goose -dir "$(MIGRATION_FOLDER)" create "$(name)" sql
.PHONY: test-migration-up
test-migration-up:
	goose -dir "$(MIGRATION_FOLDER)" postgres "$(POSTGRES_SETUP_TEST)" up

.PHONY: test-migration-down
test-migration-down:
	goose -dir "$(MIGRATION_FOLDER)" postgres "$(POSTGRES_SETUP_TEST)" down

.PHONY: test
test:
	go test $(PKG) -v -cover

.PHONY: compose-up
compose-up:
	docker-compose build
	docker-compose up -d db

.PHONY: compose-rm
compose-rm:
	docker-compose down

.PHONY: swag
swag:
	swag init -g $(MAIN)
