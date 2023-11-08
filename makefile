# Makefile

# Define the Go compiler command
GOCMD=go

# Define the command to build the binary
GOBUILD=$(GOCMD) build

# Define the command to run the Go program
GORUN=$(GOCMD) run

# Define the binary name
BINARY_NAME=myapp

# List all the Go files that need to be run
GOFILES=main.go serializer.go deserializer.go

# Default target executed when no arguments are given to make
default: run

# Target for building the binary
build:
	$(GOBUILD) -o $(BINARY_NAME) $(GOFILES)

# Target for running the Go program with all files
run:
	$(GORUN) $(GOFILES)

# Target to clean up the binary
clean:
	rm -f $(BINARY_NAME)

.PHONY: default build run clean
