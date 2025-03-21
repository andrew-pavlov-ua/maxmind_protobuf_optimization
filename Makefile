run: 
	go run cmd/main.go
bench:
	go test -bench=. -benchmem ./... -benchtime=1000000x
test:
	go test -v ./...