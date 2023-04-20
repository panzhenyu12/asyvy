VERSION := $(cat version) #Strips the v prefix from the tag
LDFLAGS := -ldflags "-s -w -X=main.version=$(VERSION)"

GOPATH := $(firstword $(subst :, ,$(shell go env GOPATH)))
GOBIN := $(GOPATH)/bin
GOSRC := $(GOPATH)/src

TEST_MODULE_DIR := pkg/module/testdata
TEST_MODULE_SRCS := $(wildcard $(TEST_MODULE_DIR)/*/*.go)
TEST_MODULES := $(patsubst %.go,%.wasm,$(TEST_MODULE_SRCS))

EXAMPLE_MODULE_DIR := examples/module
EXAMPLE_MODULE_SRCS := $(wildcard $(EXAMPLE_MODULE_DIR)/*/*.go)
EXAMPLE_MODULES := $(patsubst %.go,%.wasm,$(EXAMPLE_MODULE_SRCS))


export CGO_ENABLED := 0 

u := $(if $(update),-u)

# Tools
$(GOBIN)/wire:
	go install github.com/google/wire/cmd/wire@v0.5.0

$(GOBIN)/crane:
	go install github.com/google/go-containerregistry/cmd/crane@v0.9.0

$(GOBIN)/golangci-lint:
	curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh| sh -s -- -b $(GOBIN) v1.49.0

$(GOBIN)/labeler:
	go install github.com/knqyf263/labeler@latest

$(GOBIN)/easyjson:
	go install github.com/mailru/easyjson/...@v0.7.7

$(GOBIN)/goyacc:
	go install golang.org/x/tools/cmd/goyacc@latest

.PHONY: wire
wire: $(GOBIN)/wire
	wire gen ./pkg/commands/... ./pkg/rpc/...

.PHONY: mock
mock: $(GOBIN)/mockery
	mockery -all -inpkg -case=snake -dir $(DIR)

.PHONY: deps
deps:
	go get ${u} -d
	go mod tidy

# Run unit tests
.PHONY: test
test: $(TEST_MODULES)
	go test -v -short -coverprofile=coverage.txt -covermode=atomic ./...


.PHONY: lint
lint: $(GOBIN)/golangci-lint
	$(GOBIN)/golangci-lint run --timeout 5m

.PHONY: fmt
fmt:
	find ./ -name "*.proto" | xargs clang-format -i

.PHONY: build
build:
	go build $(LDFLAGS) ./cmd/trivy

.PHONY: protoc
protoc:
	docker build -t trivy-protoc - < Dockerfile.protoc
	docker run --rm -it -v ${PWD}:/app -w /app trivy-protoc make _$@

_protoc:
	for path in `find ./rpc/ -name "*.proto" -type f`; do \
		protoc --twirp_out=. --twirp_opt=paths=source_relative --go_out=. --go_opt=paths=source_relative $${path} || exit; \
	done

.PHONY: install
install:
	go install $(LDFLAGS) ./cmd/trivy

.PHONY: clean
clean:
	rm -rf integration/testdata/fixtures/images

# Create labels on GitHub
.PHONY: label
label: $(GOBIN)/labeler
	labeler apply misc/triage/labels.yaml -r aquasecurity/trivy -l 5

# Run MkDocs development server to preview the documentation page
.PHONY: mkdocs-serve
mkdocs-serve:
	docker build -t $(MKDOCS_IMAGE) -f docs/build/Dockerfile docs/build
	docker run --name mkdocs-serve --rm -v $(PWD):/docs -p $(MKDOCS_PORT):8000 $(MKDOCS_IMAGE)

# Generate JSON marshaler/unmarshaler for TinyGo/WebAssembly as TinyGo doesn't support encoding/json.
.PHONY: easyjson
easyjson: $(GOBIN)/easyjson
	easyjson pkg/module/serialize/types.go

# Generate license parser with goyacc
.PHONY: yacc
yacc: $(GOBIN)/goyacc
	go generate ./pkg/licensing/expression/... 	
