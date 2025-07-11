# ---- config ---------------------------------------------------------------
DRIVER       := sqlite3
DB           := localSQLite.db
MIGR_DIR     := ./migrations
GOOSE_VER    := v3.24.3
SQLC_VER     := v1.29.0

GOOSE ?= goose
SQLC  ?= sqlc

# ---- helpers --------------------------------------------------------------
define maybe-install
	@command -v $(1) >/dev/null || \
		go install $(2)@$(3)
endef

# ---- targets --------------------------------------------------------------
.PHONY: deps migrate up down status redo create generate vet tidy test

deps:                                  ## one-time tool install
	$(call maybe-install,$(GOOSE),github.com/pressly/goose/v3/cmd/goose,$(GOOSE_VER))
	$(call maybe-install,$(SQLC),github.com/sqlc-dev/sqlc/cmd/sqlc,$(SQLC_VER))

migrate: up                            ## default

up: deps                               ## run all pending migrations
	$(GOOSE) -dir $(MIGR_DIR) $(DRIVER) $(DB) up

down: deps                             ## rollback last migration
	$(GOOSE) -dir $(MIGR_DIR) $(DRIVER) $(DB) down

status: deps                           ## show migration status
	$(GOOSE) -dir $(MIGR_DIR) $(DRIVER) $(DB) status

redo: deps                             ## down & up last migration
	$(GOOSE) -dir $(MIGR_DIR) $(DRIVER) $(DB) redo

create: deps                           ## make a new timestamped migration
	@read -p "name: " name && \
	$(GOOSE) -dir $(MIGR_DIR) $(DRIVER) $(DB) create $$name sql

generate: deps                         ## sqlc code-gen
	$(SQLC) generate

vet: deps                              ## govet + sqlc vet
	go vet ./...
	$(SQLC) vet

tidy:                                  ## keep go.mod clean
	go mod tidy

test:                                  ## run unit tests
	go test ./...
