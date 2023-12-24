.PHONY: build run clean

SRCS = $(shell find . -name '*.go') Dockerfile Makefile go.mod go.sum _tools/go.mod _tools/go.sum
PORT = 8080

FAKE_SHADER ?= 0

run: .build
	docker run --rm -e FAKE_SHADER=${FAKE_SHADER} -p '${PORT}:8080' -it -v '${PWD}:/app' $(shell cat $<)

build: .build
.build: ${SRCS}
	docker build -q --network=none . > $@

clean:
	rm -f .build

.DELETE_ON_ERROR:
