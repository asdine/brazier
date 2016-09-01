NAME		:= brazier
GOFILES := $(shell find . -type f -name '*.go')

.PHONY: all build $(NAME) deps restore test

all: build

build: $(NAME)

$(NAME): $(GOFILES)
	go install ./cmd/$@

deps:
	godep save `go list ./... | grep -v /vendor/`

restore:
	godep restore

test:
	go test -v -cover `go list ./... | grep -v /vendor/`

testrace:
	go test -v -race -cover `go list ./... | grep -v /vendor/`

proto:
	protoc --proto_path=./store/boltdb/proto store/boltdb/proto/item.proto --go_out=store/boltdb/proto
	protoc --proto_path=./rpc/proto ./rpc/proto/*.proto --go_out=plugins=grpc:rpc/proto
