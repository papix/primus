TARGETS_NOVENDOR=$(shell glide novendor)

all: bin/primus-server bin/primus-client

bundle:
	glide install

bin/primus-server: cmd/primus/server/main.go server/*.go
	go build $(GOFLAGS) -o bin/primus-server cmd/primus/server/main.go

bin/primus-client: cmd/primus/client/main.go client/*.go
	go build $(GOFLAGS) -o bin/primus-client cmd/primus/client/main.go

fmt:
	@echo $(TARGETS_NOVENDOR) | xargs go fmt

clean:
	rm -rf bin/primus*
