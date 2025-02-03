# Variáveis
MIGRATION_DIR=internal/cmd/migrate/migrations

build:
	@echo "Building..."
	@templ generate

	# @./tailwindcss -i static/css/input.css -o static/css/output.css
	@go build -o bin/airCast internal/cmd/main.go

test:
	@echo "Testing..."
	@go test -v ./...

templ-install:
	@if ! command -v templ > /dev/null; then \
		read -p "Go's 'templ' is not installed on your machine. Do you want to install it? [Y/n] " choice; \
		if [ "$$choice" != "n" ] && [ "$$choice" != "N" ]; then \
			go install github.com/a-h/templ/cmd/templ@latest; \
			if [ ! -x "$$(command -v templ)" ]; then \
				echo "templ installation failed. Exiting..."; \
				exit 1; \
			fi; \
		else \
			echo "You chose not to install templ. Exiting..."; \
			exit 1; \
		fi; \
	fi

tailwind:
	@if [ ! -f tailwindcss ]; then curl -sL https://github.com/tailwindlabs/tailwindcss/releases/latest/download/tailwindcss-linux-x64 -o tailwindcss; fi
	
	@chmod +x tailwindcss

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
