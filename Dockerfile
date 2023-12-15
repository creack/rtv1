ARG GO_VERSION=1.21
FROM golang:${GO_VERSION}

ENV CGO_ENABLED=0 GOFLAGS='-mod=vendor'

# Precompile stdlib for js/wasm and native os/arch.
RUN go build -a std && \
    GOOS=js GOARCH=wasm go build -a std

# Tell git our directory is safe.
RUN git config --global --add safe.directory /app

WORKDIR /app

# Add the vendored tools directory.
ADD _tools/ _tools/

# Install reflex and wasmserve.
RUN cd _tools && \
    go install ./vendor/github.com/cespare/reflex && \
    go install ./vendor/github.com/hajimehoshi/wasmserve

# Add the rest of the code.
ADD . .

# Build the wasm binary.
RUN GOOS=js GOARCH=wasm go build -o /tmp/out .

# Run the dev server as command.
EXPOSE 8080
CMD reflex curl -v http://localhost:8080/_notify& wasmserve .
