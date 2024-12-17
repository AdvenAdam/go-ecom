
# build:
# 	@go build -o bin/go-ecom cmd/main.go

# test:
# 	@go test -v ./...

# clean:
# 	ifeq ($(OS),Windows_NT)
# 		del bin\go-ecom.exe
# 	else
# 		rm -f bin/go-ecom
# 	endif
# run: build
# 	./bin/go-ecom

# watch:
# 	find . -name '*.go' | entr -r make run 


# Define variables
APP_NAME := go-ecom
SRC_DIR := ./cmd
BIN_DIR := ./bin
WATCH_DIR := ./cmd

# Define the build command
build:
	@go build -o $(BIN_DIR)/$(APP_NAME) $(SRC_DIR)/main.go

# Define the run command
run: build
	@$(BIN_DIR)/$(APP_NAME)

# Force Stop
stop:
ifeq ($(OS),Windows_NT)
	@taskkill /F /IM $(APP_NAME).exe
else
	@pkill -f $(APP_NAME)
endif

# Running Test
test:
	@go test -v ./...

# Define the watch command using reflex for livereload
watch:
ifeq ($(OS),Windows_NT)
	@reflex -r '\\.go$$' -s -- sh -c 'make run'
else
	@reflex -r '\.go$$' -s -- make run
endif

# Clean up binary files
clean:
ifeq ($(OS),Windows_NT)
	@del $(BIN_DIR)\$(APP_NAME).exe
else
	@rm -f $(BIN_DIR)/$(APP_NAME)
endif

.PHONY: build run watch clean


migration:
	@migrate create -ext sql -dir cmd/migrate/migrations $(filter-out $@,$(MAKECMDGOALS))

migrate-up:
	@go run cmd/migrate/main.go up

migrate-down:
	@go run cmd/migrate/main.go down

help:
	@echo "Available commands:"
	@echo "  make build: Build the application"
	@echo "  make test: Run tests"
	@echo "  make clean: Clean build artifacts"
	@echo "  make run: Run the application with gowatch"
	@echo "  make migration <name>: Create a new migration"
	@echo "  make migrate-up: Run migrations up"
	@echo "  make migrate-down: Run migrations down"
