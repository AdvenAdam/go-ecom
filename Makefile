
build:
	@go build -o bin/go-ecom cmd/main.go

test:
	@go test -v ./...

clean:
	ifeq ($(OS),Windows_NT)
		del bin\go-ecom.exe
	else
		rm -f bin/go-ecom
	endif
run: build
	./bin/go-ecom

watch:
	@echo "Watching for changes..."
	@if [ "$(OS)" = "Windows_NT" ]; then \
		/reflex -r '\.go$$' -s -- sh -c "make run" \
	else \
		find . -name '*.go' | entr -r make run || (echo "ERROR: entr not found. Install entr by running: sudo apt-get install entr" && exit 1); \
	fi

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
