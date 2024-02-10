include .env

PROJECTNAME=$(shell basename "$(PWD)")

# Go related variables.
GOBASE=$(shell pwd)
GOPATH="$(GOBASE)/vendor:$(GOBASE)"
GOBIN=$(GOBASE)/bin
GOFILES=$(wildcard *.go)

# Redirect error output to a file, so we can show it in development mode.
STDERR=/tmp/.$(PROJECTNAME)-stderr.txt

# PID file will keep the process id of the server
PID=/tmp/.$(PROJECTNAME).pid
PID_AGENT=/tmp/.$(PROJECTNAME)-agent.pid
PID_SERVER=/tmp/.$(PROJECTNAME)-server.pid

# Make is verbose in Linux. Make it silent.
MAKEFLAGS += --silent

## run: Compile and run server and agent
run: go-compile start

## start: Start in development mode. Auto-starts when code changes.
start: start-server start-agent

## stop: Stop development mode.
stop: stop-agent stop-server

start-server: stop-server
	@echo "  >  $(PROJECTNAME) is available at $(ADDRESS)"
	@-$(GOBIN)/server & echo $$! > $(PID_SERVER)
	@cat $(PID_SERVER) | sed "/^/s/^/  \>  PID: /"

start-agent: stop-agent
	@echo "  >  $(PROJECTNAME) is available at $(ADDRESS)"
	@-$(GOBIN)/agent 2>&1 & echo $$! > $(PID_AGENT)
	@cat $(PID_AGENT) | sed "/^/s/^/  \>  PID: /"

stop-server:
	@-touch $(PID_SERVER)
	@-kill `cat $(PID_SERVER)` 2> /dev/null || true
	@-rm $(PID_SERVER)

stop-agent:
	@-touch $(PID_AGENT)
	@-kill `cat $(PID_SERVER)` 2> /dev/null || true
	@-rm $(PID_AGENT)


restart-server: stop-server start-server


## clean: Clean build files. Runs `go clean` internally.
clean:
	@(MAKEFILE) go-clean

## compile: Compile the binary.
go-compile: go-build-agent go-build-server

go-build-agent:
	@echo "  >  Building agent binary..."
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) cd ./cmd/agent && go build -o $(GOBIN)/agent $(GOFILES)

go-build-server:
	@echo "  >  Building server binary..."
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) cd ./cmd/server && go build -o $(GOBIN)/server $(GOFILES)

go-generate:
	@echo "  >  Generating dependency files..."
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go generate $(generate)

go-get:
	@echo "  >  Checking if there is any missing dependencies..."
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go get $(get)

go-install:
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go install $(GOFILES)

go-clean:
	@echo "  >  Cleaning build cache"
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go clean

.PHONY: help
all: help
help: Makefile
	@echo
	@echo " Choose a command run in "$(PROJECTNAME)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo
