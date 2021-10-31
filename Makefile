CURRENT_VERSION:=v1
PLATFORM_INFO:=$$(uname)-x$$(getconf LONG_BIT)
PROTO_ENTRY:=aadhaar_service.proto
PROTO_ENTRY_FOLDER:=api/proto/$(CURRENT_VERSION)
FBS_FOLDER:=api/flatc
FBS_BUILDER_FOLDER=pkg/flatbuffservice/api
IGNORE:=api/proto/$(CURRENT_VERSION)/aadhaar_service.proto
LAST_COMMIT:=$$(git rev-list -1 HEAD)
TAG:=$$(git describe)
GO_VERSION:=$$(go version)
GO_OUT:=pkg/api/$(CURRENT_VERSION)
all: proto restgateway swagger releaseb
#  chnageloggenerate

run:
	go run cmd/server/main.go
chnageloggenerate:
	git-chglog -o CHANGELOG.md

proto:
	@ if ! which protoc > /dev/null; then \
		echo "error: protoc not installed" >&2; \
		exit 1; \
	fi
	@ echo "Compiling protobufs to go definitions..."
	@ for file in $(shell ls $(PROTO_ENTRY_FOLDER)/*.proto); do \
		echo "proto - $$file"; \
		protoc \
			--proto_path=$(PROTO_ENTRY_FOLDER)  --proto_path=third_party/ \
			--go-grpc_out=$(GO_OUT) --go_out=$(GO_OUT) $$file $$i || exit 1 ;\
		protoc \
			--proto_path=$(PROTO_ENTRY_FOLDER)  --proto_path=third_party/ \
			--go_out=":$(GO_OUT)" \
			--validate_out="lang=go:$(GO_OUT)" \
			$$file $$i || exit 1 ;\
	done
	@ echo "Compiling protobufs to go - success \n"

	


restgateway:
	@ if ! which protoc > /dev/null; then \
		echo "error: protoc not installed" >&2; \
		exit 1; \
	fi
	@ echo "Compiling protobufs to restgateway..."
	@ protoc \
 		--proto_path=$(PROTO_ENTRY_FOLDER) \
  		--proto_path=third_party/ --grpc-gateway_out=logtostderr=true:$(GO_OUT) \
		  $(PROTO_ENTRY);
	@ echo "Compiling protobufs to restgateway - success \n"


swagger:
	@ if ! which protoc > /dev/null; then \
		echo "error: protoc not installed" >&2; \
		exit 1; \
	fi
	@ echo "Compiling protobufs to swagger definitions..."
	@ protoc \
 		--proto_path=$(PROTO_ENTRY_FOLDER) --proto_path=third_party/ \
 		--openapiv2_out=logtostderr=true:api/swagger/$(CURRENT_VERSION) $(PROTO_ENTRY)
	@ echo "Compiling protobufs to swagger - success \n"

# Create release binary for current Platform
releaseb:
	@ echo "Compiling release binary"
	@ go build -ldflags "-s -w -X main.gitCommit=$(LAST_COMMIT) -X main.gitTag=$(TAG)" -o dist/server-$(PLATFORM_INFO) cmd/server/main.go;
	@ echo "Compiling binary success - output=dist/server-$(PLATFORM_INFO)"

# Create release binary for current Platform
releaseclientb:
	@ echo "Compiling client release binary"
	@ go build -ldflags "-s -w" -o dist/client-$(PLATFORM_INFO) cmd/client-grpc/main.go;
	@ echo "Compiling client binary success - output=dist/client-$(PLATFORM_INFO)"


# Create debug binary for current platform
debugb:
	go build -o debug cmd/server/main.go;
