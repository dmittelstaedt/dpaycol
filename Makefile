CONTAINERCMD=docker
VERSION = $(shell git tag --list | tail -1 | cut -c 2-)
GIT_COMMIT = $(shell git rev-parse --short HEAD)
BUILD_DATE=$(shell date +"%Y-%m-%d %T")
GO_VERSION=$(shell go version | awk '{print $$3}' | cut -c 3-)

all:
	echo "ToDo: Implement all $(CONTAINERCMD)"

install:
	echo "ToDo: Implement install"
	echo $(GO_VERSION)

build:
	go build -ldflags "-X github.com/dmittelstaedt/dpaycol/cmd.versionNumber=$(VERSION) -X github.com/dmittelstaedt/dpaycol/cmd.gitCommit=$(GIT_COMMIT) -X 'github.com/dmittelstaedt/dpaycol/cmd.buildDate=$(BUILD_DATE)' -X github.com/dmittelstaedt/dpaycol/cmd.goVersionNumber=$(GO_VERSION)" -o dpaycol main.go

build-container:
	echo "Building image"
	$(CONTAINERCMD) build -t dataport.de/dpaycol --no-cache --build-arg http_proxy=$(http_proxy) --build-arg https_proxy=$(https_proxy) .
	echo "Creating container"
	$(CONTAINERCMD) create -it --name dpaycol dataport.de/dpaycol
	echo "Copying executable"
	$(CONTAINERCMD) cp dpaycol:/data/dpaycol/dpaycol .
	echo "Removing container"
	$(CONTAINERCMD) container rm dpaycol

build-container-s390x:
	echo "Building image"
	$(CONTAINERCMD) build -t dataport.de/dpaycol --no-cache --build-arg http_proxy=$(http_proxy) --build-arg https_proxy=$(https_proxy) -f Dockerfile.s390x .
	echo "Creating container"
	$(CONTAINERCMD) create -it --name dpaycol dataport.de/dpaycol
	echo "Copying executable"
	$(CONTAINERCMD) cp dpaycol:/root/dpaycol .
	echo "Removing container"
	$(CONTAINERCMD) container rm dpaycol

clean:
	rm -rvf dpaycol

run:
	go run main.go
