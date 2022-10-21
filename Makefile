
fmt:
	go fmt ./...

mock:
	go generate -v ./...

test: mock
	go test ./... --coverprofile coverage.out

cover:
	go tool cover -html coverage.out
