.PHONY: setup build run test clean

# Name of your database for local development
DB_NAME=dev.db
PORT=8085
HOST=http://localhost:$(PORT)

setup:
	@echo "Setting up development database..."
	@if [ ! -f $(DB_NAME) ]; then \
		cp internal/testdata/test_db $(DB_NAME); \
		echo "$(DB_NAME) created successfully."; \
	else \
		echo "$(DB_NAME) already exists."; \
	fi

build:
	@echo "Tidying go modules and building..."
	go mod tidy
	go build -o bookmarks_run_binary main.go

run: build setup
	@echo "Starting bookmarks server on $(HOST)..."
	./bookmarks_run_binary $(HOST) $(PORT) $(DB_NAME)

test:
	@echo "Running tests..."
	go test -v ./...

clean:
	@echo "Cleaning up..."
	rm -f bookmarks_run_binary
	rm -f $(DB_NAME)
