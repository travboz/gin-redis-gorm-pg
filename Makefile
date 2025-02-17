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


ACCESS_CONTAINER=6379
INTERNAL_EXPOSED_PORT=6379
IMAGE_NAME="redis:7.4.2-alpine3.21"
CONTAINER_NAME="redis-cache"

.PHONY: run-redis
run-redis:
	@echo "Running container, access by using port $(ACCESS_CONTAINER)"
	@docker run -d --rm --name $(CONTAINER_NAME) -p $(ACCESS_CONTAINER):$(INTERNAL_EXPOSED_PORT) $(IMAGE_NAME)


stop-redis:
	@docker container stop $(CONTAINER_NAME)