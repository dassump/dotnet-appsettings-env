.DEFAULT_GOAL := default

APP     := dotnet-appsettings-env
VERSION := $(shell git describe --tags)
VERSION := $(if $(VERSION:-=),$(VERSION),unknown)

GOCMD   := $(shell which go)
GOROOT  := $(shell $(GOCMD) env GOROOT)
GOPATH  := $(shell $(GOCMD) env GOPATH)
GOCGO   := 0

LDFLAGS   := -ldflags "-s -w -X main.app=$(APP) -X main.version=$(VERSION)"
MAKEFLAGS += --silent

clean:
	$(GOCMD) clean -cache
	rm -rf build/$(APP)-*

fmt:
	$(GOCMD) fmt ./...

vet:
	$(GOCMD) vet ./...

compile:
	CGO_ENABLED=$(GOCGO) GOOS=linux   GOARCH=amd64 $(GOCMD) build $(LDFLAGS) -o build/$(APP)-linux-amd64 .
	CGO_ENABLED=$(GOCGO) GOOS=linux   GOARCH=arm64 $(GOCMD) build $(LDFLAGS) -o build/$(APP)-linux-arm64 .
	CGO_ENABLED=$(GOCGO) GOOS=windows GOARCH=386   $(GOCMD) build $(LDFLAGS) -o build/$(APP)-windows-386.exe .
	CGO_ENABLED=$(GOCGO) GOOS=windows GOARCH=amd64 $(GOCMD) build $(LDFLAGS) -o build/$(APP)-windows-amd64.exe .
	CGO_ENABLED=$(GOCGO) GOOS=windows GOARCH=arm64 $(GOCMD) build $(LDFLAGS) -o build/$(APP)-windows-arm64.exe .
	CGO_ENABLED=$(GOCGO) GOOS=darwin  GOARCH=amd64 $(GOCMD) build $(LDFLAGS) -o build/$(APP)-darwin-amd64 .
	CGO_ENABLED=$(GOCGO) GOOS=darwin  GOARCH=arm64 $(GOCMD) build $(LDFLAGS) -o build/$(APP)-darwin-arm64 .

default: clean fmt vet compile;
