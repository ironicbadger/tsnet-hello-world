# Makefile for ts-hello-world

# Default target
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  run       - Run the application with Docker Compose"
	@echo "  build     - Build Docker image"
	@echo "  clean     - Stop and remove containers"
	@echo "  logs      - Show logs from running container"

# Run application
.PHONY: run
run:
	docker compose up --build

# Build image
.PHONY: build
build:
	docker compose build

# Clean up
.PHONY: clean
clean:
	docker compose down -v

# Show logs
.PHONY: logs
logs:
	docker compose logs -f