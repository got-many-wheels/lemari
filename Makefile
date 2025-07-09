build:
	templ generate view
	go build -o bin/main ./app
templ:
	templ generate -watch -proxy=http://localhost:8080
