.PHONY: build
build:
	go build -v -o ./dist/gophermart ./cmd/gophermart/main.go

.PHONY: test
test:
	go fmt ./internal/app/...
	go vet ./internal/app/...
	go test -cover -v -timeout 30s ./internal/app/...

.PHONY: run
run: 
	go run ./cmd/gophermart/main.go
