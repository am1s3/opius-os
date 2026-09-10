BINARY  := opius
VERSION := 1.0.0
LDFLAGS := -X github.com/opius-os/opius/internal/build.Version=$(VERSION)

.PHONY: build run tidy clean install web doctor

build:
	go build -ldflags "$(LDFLAGS)" -o bin/$(BINARY) .

run:
	go run .

tidy:
	go mod tidy

clean:
	rm -rf bin

install:
	go install .

web:
	go run . web

doctor:
	go run . doctor