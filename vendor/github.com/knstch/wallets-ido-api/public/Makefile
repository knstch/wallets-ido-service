PROTO_FILES := $(wildcard *.proto)
GOPATH := $(shell go env GOPATH)

GOOGLEAPIS := $(GOPATH)/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@v1.16.0/third_party/googleapis

.PHONY: all proto openapi clean

all: proto

proto:
	@for file in $(PROTO_FILES); do \
		protoc \
		  -I . \
		  -I $(GOOGLEAPIS) \
		  --go_out=$(dir $$file) \
		  --go-grpc_out=$(dir $$file) \
		  $$file; \
	done

openapi:
	@for file in $(PROTO_FILES); do \
		protoc \
		  -I . \
		  -I $(GOOGLEAPIS) \
		  --openapiv2_out=logtostderr=true,json_names_for_fields=false,output_format=json:. \
		  $$file; \
	done

clean:
	@find . -name "*.pb.go" -type f -delete
