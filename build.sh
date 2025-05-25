#!/bin/bash
set -e

echo "Installing dependencies..."
go install github.com/a-h/templ/cmd/templ@latest
npm install

echo "Generating Templ files..."
/opt/buildhome/go/bin/templ generate

echo "Building CSS..."
npx tailwindcss -i ./static/css/input.css -o ./static/css/output.css --minify

echo "Building Go binary..."
GOOS=linux GOARCH=amd64 go build -o netlify/functions/server ./netlify/functions/server/main.go

echo "Build complete!"