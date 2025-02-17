OUTPUT_BINARY = gredis
OUTPUT_DIR = ./bin
ENTRY_DIR = ./
ADDR = ":4000" # Default address


.PHONY: build
build:
	@mkdir -p $(OUTPUT_DIR)
	@go build -o $(OUTPUT_DIR)/$(OUTPUT_BINARY) $(ENTRY_DIR)

.PHONY: run
run: build
	@$(OUTPUT_DIR)/$(OUTPUT_BINARY)

.PHONY: clean
clean:
	@rm -rf $(OUTPUT_DIR)

# Docker commands
.PHONY: up
up:	
	@echo "Starting containers..."
	@docker compose up -d

.PHONY: down
down:
	@echo "Stopping containers..."
	@docker compose down

list-containers:
	@echo "Listing containers..."
	@docker container ls


REDIS_ACCESS=6379
REDIS_INTERNAL=6379
REDIS_IMAGE="redis:7.4.2-alpine3.21"
REDIS_CONTAINER_NAME="redis-cache"

.PHONY: run-redis stop-redis
run-redis:
	@echo "Running Redis container, access by using port $(REDIS_ACCESS)"
	@docker run -d --rm --name $(REDIS_CONTAINER_NAME) -p $(REDIS_ACCESS):$(REDIS_INTERNAL) $(REDIS_IMAGE)

stop-redis:
	@docker container stop $(REDIS_CONTAINER_NAME)

PG_ACCESS=5432
PG_INTERNAL=5432
PG_CONTAINER_NAME="gredis-pg-db"

.PHONY: run-pg stop-pg
run-pg:
	@echo "Running Postgres container, access by using port $(PG_ACCESS)"
	@docker run -d --rm \
	--name gredis-pg-db \
	-v gredis-db:/var/lib/postgresql/data \
	-e POSTGRES_DB=gredis \
	-e POSTGRES_PASSWORD=adminpass \
	-e POSTGRES_USER=admin \
	-p $(PG_ACCESS):$(PG_INTERNAL) \
	postgres:16.3-alpine

stop-pg:
	@docker container stop $(PG_CONTAINER_NAME)


# goose migrations commands
DB_ADDR="postgres://admin:adminpass@localhost/gredis?sslmode=disable"
MIGRATIONS_DIR =./sql/migrations

.PHONY: goose-up goose-down goose-status

goose-up:
	@goose -dir $(MIGRATIONS_DIR) postgres "$(DB_ADDR)" up

goose-down:
	@goose -dir $(MIGRATIONS_DIR) postgres "$(DB_ADDR)" down

goose-status:
	goose -dir $(MIGRATIONS_DIR) postgres "$(DB_ADDR)" status
