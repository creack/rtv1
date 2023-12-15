.PHONY: build run clean

SRCS = $(shell find . -name '*.go') $(shell find maps -type f) Dockerfile Makefile go.mod go.sum _tools/go.mod _tools/go.sum
PORT = 8080

run: .build
	docker run --rm -p '${PORT}:8080' -it -v '${PWD}:/app' $(shell cat $<)

build: .build
.build: ${SRCS}
	docker build -q --network=none . > $@

clean:
	rm -f .build
