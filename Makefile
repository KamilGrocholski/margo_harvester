harvest-dev:
	go run cmd/harvest/main.go

test:
	go test -v -count=1 ./...
