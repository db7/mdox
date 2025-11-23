GO ?= go
BINARY ?= mdox

.PHONY: mdox clean test

mdox:
	$(GO) build -mod=vendor -o $(BINARY) .

test:
	$(GO) test ./...

clean:
	rm -rf $(BINARY)
