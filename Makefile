tidy:
	go mod tidy
	go mod vendor


run:
	go run ./cmd/cli/

build:
	go build -o smart-water-heater ./cmd/cli/ 