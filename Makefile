build:
	templ generate view
	go build -o bin/main ./cmd/app
templ:
	templ generate -watch -proxy=http://localhost:8080
