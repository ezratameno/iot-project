tidy:
	go mod vendor
	go mod tidy

run:
	go run ./cmd/cli/main.go