.PHONY: build run test clean install-deps templ css css-watch dev

# Build the application
build: templ css
	go build -o bin/server cmd/server/main.go

# Run the application
run: build
	./bin/server

# Run tests
test:
	go test -v ./...

# Clean build artifacts
clean:
	rm -rf bin/
	rm -f static/css/output.css
	find . -name "*_templ.go" -type f -delete

# Install dependencies
install-deps:
	go mod download
	go install github.com/a-h/templ/cmd/templ@latest
	npm install

# Generate templ files
templ:
	~/go/bin/templ generate

# Build CSS
css:
	npx tailwindcss -i ./static/css/input.css -o ./static/css/output.css --minify

# Watch CSS for development
css-watch:
	npx tailwindcss -i ./static/css/input.css -o ./static/css/output.css --watch

# Development mode (run templ generate and server with live reload)
dev:
	@echo "Starting development mode..."
	@echo "Run 'make css-watch' in another terminal for CSS hot reload"
	~/go/bin/templ generate --watch --proxy="http://localhost:8080" --cmd="go run cmd/server/main.go"