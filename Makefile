# ---- config ---------------------------------------------------------------
DRIVER   := sqlite3          # goose driver name
DB       := localSQLite.db   # SQLite file
MIGR_DIR := ./migrations     # migrations live in project root

GOOSE    ?= goose            # assumes $(go env GOPATH)/bin is on $PATH

# ---- targets --------------------------------------------------------------
.PHONY: deps migrate up down status redo

deps:                            # one-time install (skips if goose already exists)
	@command -v $(GOOSE) >/dev/null || go install github.com/pressly/goose/v3/cmd/goose@latest

migrate: up                      # default shorthand

up: deps                         # apply all pending migrations
	$(GOOSE) -dir $(MIGR_DIR) $(DRIVER) $(DB) up

down: deps                       # roll back the last migration
	$(GOOSE) -dir $(MIGR_DIR) $(DRIVER) $(DB) down

status: deps                     # show applied / pending list
	$(GOOSE) -dir $(MIGR_DIR) $(DRIVER) $(DB) status

redo: deps                       # roll back & re-apply the most recent migration
	$(GOOSE) -dir $(MIGR_DIR) $(DRIVER) $(DB) redo
