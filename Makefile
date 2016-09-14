NAME			:= brazier
GOFILES		:= $(shell find . -type f -name '*.go')
PACKAGES	:= $(shell glide novendor)

.PHONY: all build $(NAME) deps restore test

all: build

build: $(NAME)

$(NAME): $(GOFILES)
	go install ./cmd/$@

deps:
	glide up

install:
	glide install

test:
	go test -v -cover $(PACKAGES)

testrace:
	go test -v -race -cover $(PACKAGES)

bench:
	go test -v -run=NONE -bench=. -benchmem $(PACKAGES)

gen:
	go generate $(PACKAGES)
