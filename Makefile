NAME      := brazier
PACKAGES  := $(shell glide novendor)

.PHONY: all build $(NAME) deps restore test

all: build

build: $(NAME)

$(NAME):
	go install ./cmd/$@

deps:
	glide up

install:
	glide install
	go get github.com/favadi/protoc-go-inject-tag

test:
	go test -v -cover $(PACKAGES)

testrace:
	go test -v -race -cover $(PACKAGES)

bench:
	go test -run=NONE -bench=. -benchmem $(PACKAGES)

gen:
	go generate $(PACKAGES)
