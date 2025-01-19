GO := go
BIN := app

build:
	$(GO) build -o $(BIN) ./cmd

test:
	$(GO) test ./...

fmt:
	$(GO) fmt ./...

clean:
	rm -f $(BIN)