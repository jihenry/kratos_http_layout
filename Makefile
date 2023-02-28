# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BIN=$(PWD)/bin
BIN_NAME=gaas-`date +"%Y%m%d%H%M"`
VERSION=$(shell git describe --tags --always)
BRANCH=$(shell git rev-parse --abbrev-ref HEAD)
# Proto parameters
PROTOC_PATH=$(GOPATH)/bin
PROTO_PATH=./proto
GAAS_PROTO_PATH=./gaas_proto
PROTOC=$(PROTOC_PATH)/protoc
GAAS_PPAOTO_PLUGIN_NAME=gaasproto
GAAS_PROTO_PLUGIN_PATH=$(GAAS_PROTO_PATH)
PBDIRS= $(shell find $(PROTO_PATH) -maxdepth 1 -type d | grep pb$)
GRPC_PROTO_FILES=$(shell find $(PROTO_PATH) -name \*.proto)
GAAS_PROTO_FILES=$(shell find $(GAAS_PROTO_PATH) -name \*.proto)
CMD=status -s

all: clean pb gaas_pb build

pb:
	@echo 'proto files generating...'
	$(PROTOC) -I=. --go_out=. --go-grpc_out=. --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative,require_unimplemented_servers=false $(GRPC_PROTO_FILES)
	@echo 'generate finished.'
gaas_pb:
	@echo 'gaas proto files generating...'
	$(PROTOC) -I=. --plugin $(GAAS_PROTO_PLUGIN_PATH) --$(GAAS_PPAOTO_PLUGIN_NAME)_out=./pkg/gaasproto/ --$(GAAS_PPAOTO_PLUGIN_NAME)_opt=outputFile=proto_map,type=server --go_out=. --go_opt=paths=source_relative $(GAAS_PROTO_FILES)
	@echo 'gaas generate finished.'

build:
	@echo 'go building...'
	mkdir -p $(BIN) && go build -ldflags "-X main.Version=$(VERSION)" -o $(BIN) ./cmd/...
	@echo 'go build finished.'

build_linux:
ifdef ONLINE
ifneq (${BRANCH},master)
	$(error the current git branch:${BRANCH} not master)
endif
endif
	mkdir -p $(BIN) && env GOOS=linux GOARCH=amd64 go build -ldflags "-X main.Version=$(VERSION)" -o $(BIN) ./cmd/...

build_all:
	@echo 'go building all..., BLIST=${BLIST}'
ifdef BLIST
ifeq ($(findstring task, $(BLIST)), task)
	cd ../task && make build_linux BIN=$(BIN)
endif
ifeq ($(findstring activity, $(BLIST)), activity)
	cd ../activity && make build_linux BIN=$(BIN)
endif
ifeq ($(findstring publish, $(BLIST)), publish)
	cd ../publish && make build_linux BIN=$(BIN)
endif
ifeq ($(findstring base, $(BLIST)), base)
	cd ../base && make build_linux BIN=$(BIN)
endif
ifeq ($(findstring account, $(BLIST)), account)
	cd ../account && make build_linux BIN=$(BIN)
endif
ifeq ($(findstring tact, $(BLIST)), tact)
	cd ../tact && make build_linux BIN=$(BIN)
endif
ifeq ($(findstring nprize, $(BLIST)), nprize)
	cd ../nprize && make build_linux BIN=$(BIN)
endif
ifeq ($(findstring thirdparty_agent, $(BLIST)), thirdparty_agent)
	cd ../thirdparty_agent && make build_linux BIN=$(BIN)
endif
ifeq ($(findstring gamespace, $(BLIST)), gamespace)
	make build_linux BIN=$(BIN)
endif
else
	cd ../task && make build_linux BIN=$(BIN)
	cd ../activity && make build_linux BIN=$(BIN)
	cd ../publish && make build_linux BIN=$(BIN)
	cd ../base && make build_linux BIN=$(BIN)
	cd ../account && make build_linux BIN=$(BIN)
	cd ../tact && make build_linux BIN=$(BIN)
	cd ../nprize && make build_linux BIN=$(BIN)
	cd ../thirdparty_agent && make build_linux BIN=$(BIN)
	make build_linux BIN=$(BIN)
endif
	tar -zcvf $(BIN_NAME).tgz ./bin
	@echo 'go build finished.'
online:
	@echo 'go building online..., BLIST=${BLIST}'
ifdef BLIST
ifeq ($(findstring task, $(BLIST)), task)
	cd ../task && make build_linux BIN=$(BIN) ONLINE=true
endif
ifeq ($(findstring activity, $(BLIST)), activity)
	cd ../activity && make build_linux BIN=$(BIN) ONLINE=true
endif
ifeq ($(findstring publish, $(BLIST)), publish)
	cd ../publish && make build_linux BIN=$(BIN) ONLINE=true
endif
ifeq ($(findstring base, $(BLIST)), base)
	cd ../base && make build_linux BIN=$(BIN) ONLINE=true
endif
ifeq ($(findstring account, $(BLIST)), account)
	cd ../account && make build_linux BIN=$(BIN) ONLINE=true
endif
ifeq ($(findstring tact, $(BLIST)), tact)
	cd ../tact && make build_linux BIN=$(BIN) ONLINE=true
endif
ifeq ($(findstring nprize, $(BLIST)), nprize)
	cd ../nprize && make build_linux BIN=$(BIN) ONLINE=true
endif
ifeq ($(findstring thirdparty_agent, $(BLIST)), thirdparty_agent)
	cd ../thirdparty_agent && make build_linux BIN=$(BIN) ONLINE=true
endif
ifeq ($(findstring gamespace, $(BLIST)), gamespace)
	make build_linux BIN=$(BIN) ONLINE=true
endif
else
	cd ../task && make build_linux BIN=$(BIN) ONLINE=true
	cd ../activity && make build_linux BIN=$(BIN) ONLINE=true
	cd ../publish && make build_linux BIN=$(BIN) ONLINE=true
	cd ../base && make build_linux BIN=$(BIN) ONLINE=true
	cd ../account && make build_linux BIN=$(BIN)  ONLINE=true
	cd ../tact && make build_linux BIN=$(BIN)  ONLINE=true
	cd ../nprize && make build_linux BIN=$(BIN)  ONLINE=true
	cd ../thirdparty_agent && make build_linux BIN=$(BIN)  ONLINE=true
	make build_linux BIN=$(BIN)  ONLINE=true
endif
	tar -zcvf $(BIN_NAME).tgz ./bin
	@echo 'go build finished.'
online_proto:
	cd ../task && go get -u gitlab.yeahka.com/gaas/proto@master
	cd ../activity && go get -u gitlab.yeahka.com/gaas/proto@master
	cd ../publish && go get -u gitlab.yeahka.com/gaas/proto@master
	cd ../base && go get -u gitlab.yeahka.com/gaas/proto@master
	cd ../account && go get -u gitlab.yeahka.com/gaas/proto@master
	cd ../tact && go get -u gitlab.yeahka.com/gaas/proto@master
	cd ../nprize && go get -u gitlab.yeahka.com/gaas/proto@master
	cd ../thirdparty_agent && go get -u gitlab.yeahka.com/gaas/proto@master
	go get -u gitlab.yeahka.com/gaas/proto@master
test_proto:
	cd ../task && go get -u gitlab.yeahka.com/gaas/proto@test
	cd ../activity && go get -u gitlab.yeahka.com/gaas/proto@test
	cd ../publish && go get -u gitlab.yeahka.com/gaas/proto@test
	cd ../base && go get -u gitlab.yeahka.com/gaas/proto@test
	cd ../account && go get -u gitlab.yeahka.com/gaas/proto@test
	cd ../tact && go get -u gitlab.yeahka.com/gaas/proto@test
	cd ../nprize && go get -u gitlab.yeahka.com/gaas/proto@test
	cd ../thirdparty_agent && go get -u gitlab.yeahka.com/gaas/proto@test
	go get -u gitlab.yeahka.com/gaas/proto@test
online_pkg:
	cd ../task && go get -u gitlab.yeahka.com/gaas/pkg@master
	cd ../activity && go get -u gitlab.yeahka.com/gaas/pkg@master
	cd ../publish && go get -u gitlab.yeahka.com/gaas/pkg@master
	cd ../base && go get -u gitlab.yeahka.com/gaas/pkg@master
	cd ../account && go get -u gitlab.yeahka.com/gaas/pkg@master
	cd ../tact && go get -u gitlab.yeahka.com/gaas/pkg@master
	cd ../nprize && go get -u gitlab.yeahka.com/gaas/pkg@master
	cd ../thirdparty_agent && go get -u gitlab.yeahka.com/gaas/pkg@master
	go get -u gitlab.yeahka.com/gaas/pkg@master
test_pkg:
	cd ../task && go get -u gitlab.yeahka.com/gaas/pkg@test
	cd ../activity && go get -u gitlab.yeahka.com/gaas/pkg@test
	cd ../publish && go get -u gitlab.yeahka.com/gaas/pkg@test
	cd ../base && go get -u gitlab.yeahka.com/gaas/pkg@test
	cd ../account && go get -u gitlab.yeahka.com/gaas/pkg@test
	cd ../tact && go get -u gitlab.yeahka.com/gaas/pkg@test
	cd ../nprize && go get -u gitlab.yeahka.com/gaas/pkg@test
	cd ../thirdparty_agent && go get -u gitlab.yeahka.com/gaas/pkg@test
	go get -u gitlab.yeahka.com/gaas/pkg@test
tidy:
	cd ../task && rm -f go.sum && go mod tidy
	cd ../activity && rm -f go.sum && go mod tidy
	cd ../publish && rm -f go.sum && go mod tidy
	cd ../base && rm -f go.sum && go mod tidy
	cd ../account && rm -f go.sum && go mod tidy
	cd ../tact && rm -f go.sum && go mod tidy
	cd ../nprize && rm -f go.sum && go mod tidy
	cd ../thirdparty_agent && rm -f go.sum && go mod tidy
	rm -f go.sum && go mod tidy
git:
	cd ../task && git $(CMD)
	cd ../activity && git $(CMD)
	cd ../publish && git $(CMD)
	cd ../base && git $(CMD)
	cd ../account && git $(CMD)
	cd ../tact && git $(CMD)
	cd ../nprize && git $(CMD)
	cd ../thirdparty_agent && git $(CMD)
	git $(CMD)
tar:
	@echo 'tar begin...'
	tar -zcvf $(BIN_NAME).tgz ./bin
	@echo 'tar finished.'
clean:
	@echo 'cleaning...'
	@if [ "`ls $(BIN)`" = "" ];then \
        echo "$(BIN) is empty."; \
    else \
        rm $(BIN)/*; \
    fi
	@if [ "`ls $(BIN)/.DS_Store`" = "" ];then \
        echo "$(BIN) is empty."; \
    else \
        rm $(BIN)/.DS_Store; \
    fi
	@if [ "`ls *.tgz`" = "" ];then \
        echo "not find .tgz file"; \
    else \
        rm -f *.tgz; \
    fi
	@echo 'clean finished.'
