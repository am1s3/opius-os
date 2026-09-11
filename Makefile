BINARY  := opius
VERSION := 1.2.0
CHANNEL := alpha
LDFLAGS := -X github.com/opius-os/opius/internal/build.Version=$(VERSION) \
           -X github.com/opius-os/opius/internal/build.Channel=$(CHANNEL)

.PHONY: build run tidy clean install web doctor device-list pkg-list

build:
	go build -ldflags "$(LDFLAGS)" -o bin/$(BINARY) .

run:
	go run .

tidy:
	go mod tidy

clean:
	rm -rf bin

install: build
	mkdir -p $(HOME)/.local/bin
	cp bin/$(BINARY) $(HOME)/.local/bin/$(BINARY)
	@echo "Installed to $(HOME)/.local/bin/$(BINARY)"

web:
	go run . web

doctor:
	go run . doctor

device-list:
	go run . device list

pkg-list:
	go run . pkg list

release: build
	@echo "Building release $(VERSION)..."
	@mkdir -p release
	@tar -czf release/opius_$(VERSION)_darwin_amd64.tar.gz -C bin $(BINARY)
	@echo "Release archive: release/opius_$(VERSION)_darwin_amd64.tar.gz"

completion-zsh:
	@echo "Add to ~/.zshrc:"
	@echo 'source <(opius completion zsh)'