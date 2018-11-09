GO ?= go
BUILD := $(CURDIR)/static/libgo.wasm
TARGET := $(CURDIR)/lib/main.go

all: wasm

wasm:
	GOOS=js GOARCH=wasm $(GO) build -o $(BUILD) $(TARGET)

test:
	$(GO) run server.go