fmt:
	go fmt ./...

mock:
	go generate -v ./...

test: mock
	go-acc ./...

cover:
	go tool cover -html coverage.txt
