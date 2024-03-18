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

RANDOM=$(shell date +%s)
RND=$(shell echo "("$RANDOM" % 2039) + 63490" | bc)
SERVER_PORT=$(RND)
ADDRESS=localhost:$(SERVER_PORT)
TEMP_FILE=$(shell mktemp)

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

build: go-build-agent go-build-server

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

.PHONY: go-update-deps
go-update-deps:
	@echo ">> updating Go dependencies"
	@for m in $$(go list -mod=readonly -m -f '{{ if and (not .Indirect) (not .Main)}}{{.Path}}{{end}}' all); do \
		go get $$m; \
	done
	go mod tidy
ifneq (,$(wildcard vendor))
	go mod vendor
endif

go-install:
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go install $(GOFILES)

go-clean:
	@echo "  >  Cleaning build cache"
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go clean

test10:
	@echo "  > Test Iteration 10 ..."
	cd bin && ./metricstest -test.v -test.run=^TestIteration10$$ -agent-binary-path=./agent -binary-path=./server -database-dsn='postgres://postgres:postgres@postgres:5432/praktikum?sslmode=disable' -server-port=$(SERVER_PORT) -file-storage-path=$(TEMP_FILE) -source-path=../.

test9: test8
	@echo "  > Test Iteration 9 ..."
	cd bin && ./metricstest -test.v -test.run=^TestIteration9$$ -agent-binary-path=./agent -binary-path=./server -server-port=$(SERVER_PORT) -file-storage-path=$(TEMP_FILE) -source-path=../.

test8: test7
	@echo "  > Test Iteration 8 ..."
	cd bin && ./metricstest -test.v -test.run=^TestIteration8$$ -agent-binary-path=./agent -binary-path=./server -server-port=$(SERVER_PORT) -source-path=../.

test7: test6
	@echo "  > Test Iteration 7 ..."
	cd bin && ./metricstest -test.v -test.run=^TestIteration7$$ -agent-binary-path=./agent -binary-path=./server -server-port=$(SERVER_PORT) -source-path=../.

test6: test5
	@echo "  > Test Iteration 6 ..."
	cd bin && ./metricstest -test.v -test.run=^TestIteration6$$ -agent-binary-path=./agent -binary-path=./server -server-port=$(SERVER_PORT) -source-path=../.

test5: test4
	@echo "  > Test Iteration 5 ..."
	cd bin && ./metricstest -test.v -test.run=^TestIteration5$$ -agent-binary-path=./agent -binary-path=./server -server-port=$(SERVER_PORT) -source-path=../.

test4: test3
	@echo "  > Test Iteration 4 ..."
	cd bin && ./metricstest -test.v -test.run=^TestIteration4$$ -agent-binary-path=./agent -binary-path=./server -server-port=$(SERVER_PORT) -source-path=../.

test3: test2
	@echo "  > Test Iteration 3 ..."
	cd bin && ./metricstest -test.v -test.run=^TestIteration3[AB]*$$ -source-path=../. -agent-binary-path=./agent -binary-path=./server

test2: test1
	@echo "  > Test Iteration 1 ..."
	cd bin && ./metricstest -test.v -test.run=^TestIteration2[AB]*$$ -source-path=../. -agent-binary-path=./agent

test1:
	@echo "  > Test Iteration 1 ..."
	cd bin && ./metricstest -test.v -test.run=^TestIteration1$$ -binary-path=./server

.PHONY: help
all: help
help: Makefile
	@echo
	@echo " Choose a command run in "$(PROJECTNAME)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo
