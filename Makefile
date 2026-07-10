mod:
	go list -m --versions

test:
	go test -race -cover ./...

vet:
	go vet ./...

fmt:
	gofmt -w .

clean:
	go clean -cache
	go clean -modcache
	go clean -testcache