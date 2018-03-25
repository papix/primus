DEP ?= dep
BUILD_DIR = ./build
PRIMUS_CLIENT = $(BUILD_DIR)/primus-client
PRIMUS_SERVER = $(BUILD_DIR)/primus-server

all: deps $(PRIMUS_CLIENT) $(PRIMUS_SERVER)

deps:
	$(DEP) ensure -vendor-only

$(PRIMUS_CLIENT):
	go build -o $(PRIMUS_CLIENT) cmd/primus/client/main.go

$(PRIMUS_SERVER):
	go build -o $(PRIMUS_SERVER) cmd/primus/server/main.go

clean:
	rm -rf $(BUILD_DIR)
