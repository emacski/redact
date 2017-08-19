REPO=golang
TAG=1.8-alpine3.6
BIN=redact
VERSION?=dev
GOOS=linux
GOARCH=amd64

.PHONY: build pull shell

build: pull
	@docker run --rm \
		-v $(PWD):/go/src/github.com/emacski/$(BIN) \
		-w /go/src/github.com/emacski/$(BIN) \
		$(REPO):$(TAG) sh -c \
			'apk --no-cache add git \
				&& cd $(BIN) && go get -d && cd .. \
				&& CGO_ENABLED=0 go test \
				&& cd $(BIN) \
				&& GOOS=$(GOOS) GOARCH=$(GOARCH) CGO_ENABLED=0 go build -ldflags "-s -w -X main.version=$(VERSION)" -v'

shell: pull
	@docker run --rm -ti --init \
		-v $(PWD):/go/src/github.com/emacski/$(BIN) \
		-w /go/src/github.com/emacski/$(BIN) \
		-e CGO_ENABLED=0 \
		$(REPO):$(TAG) /bin/sh

pull:
	@docker pull $(REPO):$(TAG)
