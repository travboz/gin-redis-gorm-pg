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
	@docker compose down -v


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
