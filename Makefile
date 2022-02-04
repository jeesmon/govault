GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
BINARY_NAME=govault
VERSION=v0.0.2

all: darwin linux windows

version:
	@echo $(VERSION)
darwin: clean
	GOOS=darwin $(GOBUILD) -o release/$(BINARY_NAME)-darwin-$(VERSION) -v govault.go
linux: clean
	GOOS=linux $(GOBUILD) -o release/$(BINARY_NAME)-linux-$(VERSION) -v govault.go
windows: clean
	GOOS=windows $(GOBUILD) -o release/$(BINARY_NAME)-windows-$(VERSION).exe -v govault.go
clean:
	$(GOCLEAN)
	rm -rf release
