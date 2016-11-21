TARGETS_NOVENDOR=$(shell glide novendor)

all: bin/primus-server bin/primus-client

bundle:
	glide install

bin/primus-server: cmd/primus/primus-server.go *.go server/*.go
	GO15VENDOREXPERIMENT=1 go build $(GOFLAGS) -o bin/primus-server cmd/primus/primus-server.go

bin/primus-client: cmd/primus/primus-client.go *.go client/*.go
	GO15VENDOREXPERIMENT=1 go build $(GOFLAGS) -o bin/primus-client cmd/primus/primus-client.go

fmt:
	@echo $(TARGETS_NOVENDOR) | xargs go fmt

clean:
	rm -rf bin/primus*
