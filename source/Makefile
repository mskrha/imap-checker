BUILD		= imap-checker
OUTPUT		= $(BUILD)
VERSION		?= 0.0

all: clean format build

clean:
	rm -f $(OUTPUT)

format:
	go fmt

build:
	CGO_ENABLED=0 go build -ldflags "-X main.version=$(VERSION)" -o $(OUTPUT) *.go
