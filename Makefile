lint:
	@golangci-lint run -v ./...

test:
	@go test -v ./...

bench:
	@go test -bench=. -benchmem ./...

gen:
	./scripts/fetch-data.sh
	go run ./cmd/gen
