BINDIR := ./bin
BIN    := auth-service

GOBIN    := $(shell go env GOPATH)/bin
GOSRC    := $(shell find . -type f -name '*.go' -print) go.mod go.sum
PROTOSRC := $(shell find . -type f -name '*.proto' -print)
SQLSRC   := $(shell find . -type f -name '*.sql' -print)
GOMOD    := $(shell ./scripts/module go.mod)

GOGEN     := $(GOBIN)/protoc-gen-go
GOGRPC    := $(GOBIN)/protoc-gen-go-grpc
GOIMPORTS := $(GOBIN)/goimports
SQLC      := $(GOBIN)/sqlc
MIGRATE   := $(GOBIN)/migrate

PROTODIR := ./api
PROTOGEN := $(PROTODIR)/pb
PROTODEF := $(patsubst $(PROTODIR)/%,%,$(PROTOSRC))

SQLDIR := ./build/package/auth-service
SQLGEN := ./internal/db/sqlc

LDFLAGS := -w -s

COUNT ?= 1

# -----------------------------------------------------------------
#  build

.PHONY: all
all: build

.PHONY: build
build: $(BINDIR)/$(BIN)

$(BINDIR)/$(BIN): $(GOSRC)
	go build -trimpath -ldflags '$(LDFLAGS)' -o $(BINDIR)/$(BIN) ./cmd/$(BIN)

# -----------------------------------------------------------------
#  test

.PHONY: test
test:
	go test -race -v -count=$(COUNT) ./...

# -----------------------------------------------------------------
#  generate

.PHONY: generate
generate: $(GOGEN) $(GOGRPC) $(SQLC) $(PROTODIR)/pb/.protogen $(SQLGEN)/.sqlgen

$(PROTOGEN)/.protogen: $(PROTOSRC)
	protoc -I=$(PROTODIR) \
		--go_opt=default_api_level=API_OPAQUE \
		--go_out=.\
		--go-grpc_out=.\
		--plugin protoc-gen-go=$(GOGEN) \
		--plugin protoc-gen-go-grpc=$(GOGRPC) \
		$(PROTODEF)
	@touch $(PROTOGEN)/.protogen

$(SQLGEN)/.sqlgen: $(SQLSRC)
	$(SQLC) -f $(SQLDIR)/sqlc.yaml generate
	@touch $(SQLGEN)/.sqlgen

# -----------------------------------------------------------------
#  dependencies

$(GOGEN):
	( cd /; go install google.golang.org/protobuf/cmd/protoc-gen-go@latest)

$(GOGRPC):
	( cd /; go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest)

$(GOIMPORTS):
	(cd /; go install golang.org/x/tools/cmd/goimports@latest)

$(SQLC):
	(cd /; go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest)

$(MIGRATE):
	(cd /; go install -tags 'sqlite' github.com/golang-migrate/migrate/v4/cmd/migrate@latest)

# -----------------------------------------------------------------
#  misc

.PHONY: format
format: $(GOIMPORTS)
	GO111MODULE=on go list -f '{{.Dir}}' ./... | xargs $(GOIMPORTS) -w -local $(GOMOD)

.PHONY: migration
migration: $(MIGRATE)
	$(MIGRATE) create -ext sql -dir build/package/auth-service/migrations/ -seq -digits 4 $(NAME)

.PHONY: clean
clean:
	rm -rf $(PROTOGEN)
	rm -rf $(SQLGEN)
	rm -rf $(BINDIR)
