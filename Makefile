fmt:
	go fmt ./...

mock:
	go generate -v ./...

test: mock
	go-acc --covermode=set -o coverage.txt ./...
	grep -v -E "colibri.go|_mock.go" coverage.txt > filtered_coverage.txt
	mv filtered_coverage.txt coverage.txt

cover:
	go tool cover -html coverage.txt
