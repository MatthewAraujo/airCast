# Variáveis
MIGRATION_DIR=internal/cmd/migrate/migrations

build:
	@go build -o bin/airCast internal/cmd/main.go

test:
	@echo "Testing..."
	@go test -v ./...
	
run: build
	@echo "Building..."
	@./bin/airCast

migration:
	@migrate create -ext sql -dir $(MIGRATION_DIR) $(NAME)

migrate-up:
	@go run internal/cmd/migrate/main.go up 

migrate-down:
	@go run internal/cmd/migrate/main.go down

%:
	@:

watch:
	@if command -v air > /dev/null; then \
            air; \
            echo "Watching...";\
        else \
            read -p "Go's 'air' is not installed on your machine. Do you want to install it? [Y/n] " choice; \
            if [ "$$choice" != "n" ] && [ "$$choice" != "N" ]; then \
                go install github.com/air-verse/air@latest; \
                air; \
                echo "Watching...";\
            else \
                echo "You chose not to install air. Exiting..."; \
                exit 1; \
            fi; \
        fi
