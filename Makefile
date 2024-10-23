linter:
	@golangci-lint -c .golangci.yml run

run:
	go run cmd/main.go