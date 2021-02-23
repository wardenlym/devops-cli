.PHONY: all build

all:
	rm -rf ~/.devops-cli
	go build

build:
	go build